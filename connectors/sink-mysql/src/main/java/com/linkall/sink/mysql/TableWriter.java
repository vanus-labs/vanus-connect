package com.linkall.sink.mysql;

import com.linkall.sink.mysql.dialect.Dialect;
import com.linkall.sink.mysql.executor.InsertExecutor;
import com.linkall.sink.mysql.executor.SqlExecutor;
import com.linkall.sink.mysql.executor.UpsertExecutor;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.Connection;
import java.sql.SQLException;

public class TableWriter {
    private Logger LOGGER = LoggerFactory.getLogger(TableWriter.class);

    private SqlExecutor<JsonObject> sqlExecutor;
    private TableMetadata metadata;
    private int batchSize;
    private int commitSize;
    private Connection connection;
    private Dialect dialect;
    private JdbcConfig.InsertMode insertMode;
    private int exceptionCount;

    public TableWriter(
            Dialect dialect,
            JdbcConfig.InsertMode insertMode,
            int commitSize, Connection connection, TableMetadata metadata) {
        this.commitSize = commitSize;
        this.dialect = dialect;
        this.insertMode = insertMode;
        this.connection = connection;
        this.metadata = metadata;
    }

    private SqlExecutor<JsonObject> getSqlExecutor() {
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
        exceptionCount = 0;
        this.sqlExecutor = getSqlExecutor();
        sqlExecutor.prepareStatement(connection);
    }

    public synchronized void addToBatch(JsonObject data) throws SQLException {
        sqlExecutor.addToBatch(data);
        batchSize++;
        if (batchSize >= commitSize) {
            flush();
        }
    }

    public synchronized void flush() throws SQLException {
        if (batchSize==0) {
            return;
        }
        try {
            sqlExecutor.executeBatch();
        } catch (SQLException e) {
            if (++exceptionCount < 3) {
                throw e;
            }
            exceptionCount = 0;
            LOGGER.error("table {} flush failed max times, throw it", metadata.getTableName(), e);
        }
        batchSize = 0;
    }

    public synchronized void updateConnection(Connection connection) throws SQLException {
        this.connection = connection;
        sqlExecutor.close();
        sqlExecutor.prepareStatement(connection);
    }

    public TableMetadata getTableMetadata() {
        return this.metadata;
    }

    public void updateTableMetadata(TableMetadata tableMetadata) throws SQLException {
        try {
            flush();
            sqlExecutor.close();
        } catch (Throwable t) {
            LOGGER.warn("table {} change flush buffer data error", metadata.getTableName(), t);
        }
        this.metadata = tableMetadata;
        init();
    }
}
