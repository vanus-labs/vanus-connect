package com.linkall.sink.mysql;

import com.linkall.sink.mysql.connection.ConnectionProvider;
import com.linkall.sink.mysql.connection.HikariConnectionProvider;
import com.linkall.sink.mysql.dialect.Dialect;
import com.linkall.sink.mysql.dialect.MySqlDialect;
import com.linkall.sink.mysql.executor.InsertExecutor;
import com.linkall.sink.mysql.executor.SqlExecutor;
import com.linkall.sink.mysql.executor.UpsertExecutor;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.SQLException;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;

public class SqlWriter {
    private static final Logger LOGGER = LoggerFactory.getLogger(SqlWriter.class);
    private final SqlExecutor<JsonObject> sqlExecutor;
    private final TableMetadata metadata;
    private int batchSize;
    private final int commitSize;
    private final ScheduledExecutorService executorService;
    public static final int MAX_BATCH_SIZE = 2000;
    private final ConnectionProvider connectionProvider;
    private MySQLConfig.InsertMode insertMode;
    private long commitInterval;

    public SqlWriter(MySQLConfig config, TableMetadata metadata) {
        if (config.getCommitSize()!=null && config.getCommitSize() > 0) {
            this.commitSize = config.getCommitSize();
        } else {
            this.commitSize = MAX_BATCH_SIZE;
        }
        if (config.getCommitInterval()!=null && config.getCommitInterval() > 0) {
            this.commitInterval = config.getCommitInterval();
        } else {
            this.commitInterval = 1000;
        }
        insertMode = config.getInsertMode();
        if (insertMode==null) {
            insertMode = MySQLConfig.InsertMode.INSERT;
        }
        this.metadata = metadata;
        this.connectionProvider = new HikariConnectionProvider(config.getDbConfig());
        this.executorService = Executors.newSingleThreadScheduledExecutor();
        this.sqlExecutor = getSqlExecutor();

    }

    private SqlExecutor getSqlExecutor() {
        Dialect dialect = new MySqlDialect();
        switch (insertMode) {
            case INSERT:
                return new InsertExecutor<>(metadata, dialect);
            case UPSERT:
                return new UpsertExecutor<>(metadata, dialect);
            default:
                throw new RuntimeException("invalid insert mode");
        }
    }

    public void init() throws SQLException {
        sqlExecutor.prepareStatement(connectionProvider.getConnection());
        executorService.scheduleAtFixedRate(() -> {
            try {
                flush();
            } catch (SQLException e) {
                LOGGER.warn("sql writer flush has error", e);
            }
        }, 2 * 1000, commitInterval, TimeUnit.MILLISECONDS);
    }

    public synchronized void add(JsonObject data) throws SQLException {
        sqlExecutor.addToBatch(data);
        batchSize++;
        if (batchSize >= commitSize) {
            flush();
        }
    }

    public synchronized void flush() throws SQLException {
        sqlExecutor.executeBatch();
        batchSize = 0;
    }

    public void close() throws Exception {
        executorService.shutdown();
        flush();
        sqlExecutor.close();
        connectionProvider.close();
    }
}
