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

public class InsertExecutor<T> implements SqlExecutor<T> {
    private static final Logger LOGGER = LoggerFactory.getLogger(InsertExecutor.class);

    private final List<T> batch;
    private PreparedStatement preparedStatement;
    private final TableMetadata metadata;
    private final Dialect dialect;
    private StatementBinder binder;

    public InsertExecutor(TableMetadata metadata, Dialect dialect) {
        batch = new ArrayList<>();
        this.metadata = metadata;
        this.dialect = dialect;
    }

    @Override
    public void prepareStatement(Connection connection) throws SQLException {
        String sql = dialect.generateInsertSql(metadata.getTableName(), metadata.getColumnNames());
        LOGGER.info("insert sql: {}", sql);
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
