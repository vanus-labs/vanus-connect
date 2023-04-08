package com.linkall.sink.mysql;

import com.linkall.sink.mysql.dialect.Dialect;
import com.linkall.sink.mysql.executor.InsertExecutor;
import com.linkall.sink.mysql.executor.SqlExecutor;
import com.linkall.sink.mysql.executor.UpsertExecutor;
import io.vertx.core.json.JsonObject;

import java.sql.Connection;
import java.sql.SQLException;

public class TableWriter {
    private SqlExecutor<JsonObject> sqlExecutor;
    private TableMetadata metadata;
    private int batchSize;
    private int commitSize;

    public TableWriter(
            TableMetadata metadata,
            Dialect dialect,
            JdbcConfig.InsertMode insertMode,
            int commitSize) {
        this.commitSize = commitSize;
        this.metadata = metadata;
        this.sqlExecutor = getSqlExecutor(dialect, insertMode);
    }

    private SqlExecutor<JsonObject> getSqlExecutor(
            Dialect dialect, JdbcConfig.InsertMode insertMode) {
        switch (insertMode) {
            case INSERT:
                return new InsertExecutor<>(metadata, dialect);
            case UPSERT:
                return new UpsertExecutor<>(metadata, dialect);
            default:
                throw new RuntimeException("invalid insert mode");
        }
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
        sqlExecutor.executeBatch();
        batchSize = 0;
    }

    public synchronized void updateConnection(Connection connection) throws SQLException {
        sqlExecutor.close();
        sqlExecutor.prepareStatement(connection);
    }
}
