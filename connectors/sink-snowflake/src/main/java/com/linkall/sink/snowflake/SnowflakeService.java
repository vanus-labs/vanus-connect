package com.linkall.sink.snowflake;

import com.linkall.sink.snowflake.localfile.LocalWriter;
import io.cloudevents.CloudEvent;
import io.vertx.core.json.Json;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.SQLException;
import java.time.Instant;
import java.util.*;
import java.util.concurrent.*;
import java.util.concurrent.locks.Lock;
import java.util.concurrent.locks.ReentrantLock;
import java.util.stream.Collectors;

public class SnowflakeService {

    private Logger LOGGER = LoggerFactory.getLogger(SnowflakeService.class);

    private static final String CREATE_STAGE = "CREATE STAGE IF NOT EXISTS %s FILE_FORMAT = ( TYPE = JSON) COPY_OPTIONS = ( ON_ERROR = SKIP_FILE )";
    private static final String PUT_FILE = "PUT file://%s @%s/%s AUTO_COMPRESS=TRUE;";
    private static final String COPY_INTO = "COPY INTO %s.%s FROM '@%s/%s' "
            + "FILE_FORMAT = (TYPE = JSON) MATCH_BY_COLUMN_NAME = CASE_INSENSITIVE FILES= (%s)";
    private static final String DROP_STAGE = "DROP STAGE IF EXISTS %s;";
    private static final String REMOVE = "REMOVE @%s;";

    private static final long BUFFER_SIZE_BYTES = 100 * 1024 * 1024;
    private static final long BUFFER_FLUSH_TIME_SEC = 10;

    private ScheduledExecutorService scheduledExecutorService;
    private ExecutorService executorService;
    private final Lock writeLock;
    private final Lock loadLock;


    private String schemaName;
    private String tableName;
    private String stageName;
    private String stagePath;
    private SnowflakeDatabase database;
    private TableMetadata metadata;


    private boolean hasInitialized;
    private long flushTime;
    private long fileSize;

    private Writer writer;
    private List<Writer> needLoadFiles;
    private long previousFlush;


    public SnowflakeService(SnowflakeConfig config) {
        this.database = new SnowflakeDatabase(config.getSnowflake());
        needLoadFiles = new CopyOnWriteArrayList<>();

        schemaName = config.getSnowflake().getSchema();
        tableName = config.getSnowflake().getTable();
        stageName = Utils.getStageName(schemaName, tableName);
        stagePath = String.format("%s/%d/", UUID.randomUUID(), System.currentTimeMillis());
        LOGGER.info("stage name {} path {}", stageName, stagePath);

        previousFlush = Instant.now().getEpochSecond();
        scheduledExecutorService = Executors.newSingleThreadScheduledExecutor();
        executorService = Executors.newFixedThreadPool(2);
        writeLock = new ReentrantLock();
        loadLock = new ReentrantLock();
        if (config.getFlushTime() > 0) {
            flushTime = config.getFlushTime();
        } else {
            flushTime = BUFFER_FLUSH_TIME_SEC;
        }
        if (config.getSizeBytes() > 0) {
            fileSize = config.getSizeBytes();
        } else {
            fileSize = BUFFER_SIZE_BYTES;
        }
    }

    public void start() {
        scheduledExecutorService.scheduleAtFixedRate(() -> loadData(), 5, 2, TimeUnit.SECONDS);
        scheduledExecutorService.scheduleAtFixedRate(() -> flushWriter(false), 5, 1, TimeUnit.SECONDS);
    }

    public void stop() {
        LOGGER.info("stop start");
        scheduledExecutorService.shutdown();
        flushWriter(true);
        if (writer!=null && writer.size()==0) {
            writer.close();
        }
        loadData();
        dropStage();
        LOGGER.info("stop success");
    }

    public void addData(CloudEvent event) throws Exception {
        if (!hasInitialized) {
            initMetadata(event);
            createStageIfNotExists();
            createTableIfNotExists();
            writer = createWriter();
            hasInitialized = true;
        }
        byte[] data = event.getData().toBytes();
        Writer tmpWriter = null;
        writeLock.lock();
        try {
            writer.write(data);
            if (writer.size() > fileSize) {
                tmpWriter = writer;
                writer = createWriter();
            }
        } finally {
            writeLock.unlock();
        }
        if (tmpWriter!=null) {
            flush(tmpWriter);
        }
    }

    private Writer createWriter() throws Exception {
        return new LocalWriter();
    }

    private void initMetadata(CloudEvent event) {
        metadata = new TableMetadata(tableName);
        Map<String, Object> data = Json.CODEC.fromString(new String(event.getData().toBytes()), Map.class);
        for (Map.Entry<String, Object> entry : data.entrySet()) {
            metadata.addColumn(entry.getKey(), entry.getValue());
        }
    }

    private void createStageIfNotExists() throws SQLException {
        database.execute(String.format(CREATE_STAGE, stageName));
    }

    private void createTableIfNotExists() throws SQLException {
        StringJoiner joiner = new StringJoiner(",\n");
        metadata.getColumns().forEach(column -> joiner.add(column.getName() + " " + column.getType().name()));

        String createTable = String.format("CREATE TABLE IF NOT EXISTS %s.%s (\n" +
                "%s" +
                "\n" +
                ")", schemaName, tableName, joiner);
        LOGGER.info("create table sql \n{}", createTable);
        database.execute(createTable);
    }

    private void flushWriter(boolean force) {
        if (!force && Instant.now().getEpochSecond() - previousFlush < flushTime) {
            return;
        }
        if (writer==null || writer.size()==0) {
            return;
        }
        Writer tmpWriter;
        writeLock.lock();
        try {
            tmpWriter = writer;
            this.writer = createWriter();
        } catch (Exception e) {
            LOGGER.error("flush writer new writer error", e);
            return;
        } finally {
            writeLock.unlock();
        }
        flush(tmpWriter);
    }

    private void flush(Writer writer) {
        long size = writer.size();
        if (size==0) {
            return;
        }
        previousFlush = Instant.now().getEpochSecond();
        writer.flush();
        needLoadFiles.add(writer);
        LOGGER.info("a new file {} wait to load", writer.getFilepath());
    }

    private void loadData() {
        loadLock.lock();
        try {
            if (needLoadFiles.isEmpty()) {
                return;
            }
            List<Writer> writers = new ArrayList<>();
            List<CompletableFuture<Void>> futures = needLoadFiles.stream()
                    .map(writer -> CompletableFuture.runAsync(() -> {
                        boolean b = putFile(writer.getFilepath());
                        if (b) {
                            writers.add(writer);
                        }
                    }, executorService))
                    .collect(Collectors.toList());
            CompletableFuture.allOf(futures.toArray(new CompletableFuture[0])).join();
            if (writers.isEmpty()) {
                LOGGER.info("all put file failed");
                return;
            }
            List<String> filenames = new ArrayList<>();
            writers.forEach(writer -> filenames.add(writer.getFilename()));
            if (copyInto(filenames)) {
                for (Writer writer : writers) {
                    writer.close();
                    needLoadFiles.remove(writer);
                }
            }
            cleanStage();
        } finally {
            loadLock.unlock();
        }
    }

    private boolean putFile(String filepath) {
        try {
            database.execute(String.format(PUT_FILE, filepath, stageName, stagePath));
        } catch (SQLException e) {
            LOGGER.error("put file {} to db error", filepath, e);
            return false;
        }
        LOGGER.info("put file {} success", filepath);
        return true;
    }

    private boolean copyInto(List<String> filenames) {
        StringJoiner joiner = new StringJoiner(",");
        filenames.forEach(filename -> joiner.add("'" + filename + ".gz" + "'"));
        try {
            database.execute(String.format(COPY_INTO, schemaName, tableName, stageName, stagePath, joiner));
        } catch (SQLException e) {
            LOGGER.error("copy into files {} to db error", joiner, e);
            return false;
        }
        LOGGER.info("copy into files {} success", joiner);
        return true;
    }

    private boolean cleanStage() {
        try {
            database.execute(String.format(REMOVE, stageName));
        } catch (SQLException e) {
            LOGGER.error("remove stage {} file error", stageName, e);
            return false;
        }
        LOGGER.info("remove stage {} file success", stageName);
        return true;
    }

    private boolean dropStage() {
        try {
            database.execute(String.format(DROP_STAGE, stageName));
        } catch (SQLException e) {
            LOGGER.error("drop stage error", stageName, e);
            return false;
        }
        LOGGER.info("drop stage {} success", stageName);
        return true;
    }
}
