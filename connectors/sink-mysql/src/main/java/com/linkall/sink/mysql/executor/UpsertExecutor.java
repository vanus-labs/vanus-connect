package com.linkall.sink.mysql.executor;

import com.linkall.sink.mysql.TableMetadata;
import com.linkall.sink.mysql.binder.PreparedStatementBinder;
import com.linkall.sink.mysql.binder.StatementBinder;
import com.linkall.sink.mysql.dialect.Dialect;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.SQLException;
import java.util.ArrayList;
import java.util.List;

public class UpsertExecutor<T> implements SqlExecutor<T> {
    private static final Logger LOGGER = LoggerFactory.getLogger(UpsertExecutor.class);

    private final List<T> batch;
    private PreparedStatement preparedStatement;
    private final TableMetadata metadata;
    private final Dialect dialect;
    private StatementBinder binder;

    public UpsertExecutor(TableMetadata metadata, Dialect dialect) {
        batch = new ArrayList<>();
        this.metadata = metadata;
        this.dialect = dialect;
    }

    @Override
    public void prepareStatement(Connection connection) throws SQLException {
        String sql = dialect.generateUpsertSql(metadata.getTableName(), metadata.getColumnNames());
        LOGGER.info("table {} upsert sql: {}", metadata.getTableName(), sql);
        preparedStatement = connection.prepareStatement(sql);
        this.binder = new PreparedStatementBinder(metadata, preparedStatement);
    }

    @Override
    public void addToBatch(T data) {
        batch.add(data);
    }

    @Override
    public void executeBatch() throws SQLException {
        if (batch.isEmpty()) {
            return;
        }
        for (T t : batch) {
            binder.bindData(t);
            preparedStatement.addBatch();
        }
        preparedStatement.executeBatch();
        LOGGER.debug("table {} flush size {}", metadata.getTableName(), batch.size());
        batch.clear();
    }

    @Override
    public void close() throws SQLException {
        if (preparedStatement==null) {
            return;
        }
        preparedStatement.close();
        preparedStatement = null;
    }
}
