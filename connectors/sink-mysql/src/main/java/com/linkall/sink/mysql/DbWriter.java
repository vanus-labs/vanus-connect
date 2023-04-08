package com.linkall.sink.mysql;

import com.linkall.sink.mysql.connection.ConnectionProvider;
import com.linkall.sink.mysql.connection.SimpleConnectionProvider;
import com.linkall.sink.mysql.dialect.Dialect;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.Connection;
import java.sql.SQLException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class DbWriter {
    private static final Logger LOGGER = LoggerFactory.getLogger(DbWriter.class);
    public static final int MAX_BATCH_SIZE = 100;
    public static final int COMMIT_INTERVAL = 1000;

    private final int commitSize;
    private JdbcConfig.InsertMode insertMode;
    private long commitInterval;

    private final ScheduledExecutorService executorService;
    private final ConnectionProvider connectionProvider;
    private Dialect dialect;
    private Map<String, TableWriter> tableWriters;

    public DbWriter(JdbcConfig config, Dialect dialect) {
        if (config.getCommitSize()!=null && config.getCommitSize() > 0) {
            this.commitSize = config.getCommitSize();
        } else {
            this.commitSize = MAX_BATCH_SIZE;
        }
        if (config.getCommitInterval()!=null && config.getCommitInterval() > 0) {
            this.commitInterval = config.getCommitInterval();
        } else {
            this.commitInterval = COMMIT_INTERVAL;
        }
        insertMode = config.getInsertMode();
        if (insertMode==null) {
            insertMode = JdbcConfig.InsertMode.INSERT;
        }
        this.connectionProvider = new SimpleConnectionProvider(config.getDbConfig());
        this.executorService = Executors.newSingleThreadScheduledExecutor();
        this.tableWriters = new ConcurrentHashMap<>();
        this.dialect = dialect;
        init();
    }

    public void init() {
        executorService.scheduleAtFixedRate(
                () -> {
                    try {
                        flush();
                    } catch (SQLException e) {
                        LOGGER.warn("sql writer flush has error", e);
                    }
                },
                2 * 1000,
                commitInterval,
                TimeUnit.MILLISECONDS);
    }

    public synchronized void add(String tableName, JsonObject data) throws SQLException {
        TableWriter tableWriter = tableWriters.get(tableName);
        if (tableWriter==null) {
            List<String> columnNames = new ArrayList<>(data.fieldNames());
            TableMetadata metadata = new TableMetadata(tableName, columnNames);
            tableWriter = new TableWriter(metadata, dialect, insertMode, commitSize);
            tableWriter.updateConnection(connectionProvider.getConnection());
            tableWriters.put(tableName, tableWriter);
        }
        tableWriter.addToBatch(data);
    }

    public synchronized void flush() throws SQLException {
        if (tableWriters.size()==0) {
            return;
        }
        if (!connectionProvider.isConnectionValid()) {
            connectionProvider.close();
            Connection connection = connectionProvider.getConnection();
            for (TableWriter tableWriter : tableWriters.values()) {
                tableWriter.updateConnection(connection);
            }
        }
        for (TableWriter tableWriter : tableWriters.values()) {
            tableWriter.flush();
        }
    }

    public synchronized void close() throws Exception {
        executorService.shutdown();
        flush();

        connectionProvider.close();
    }
}
