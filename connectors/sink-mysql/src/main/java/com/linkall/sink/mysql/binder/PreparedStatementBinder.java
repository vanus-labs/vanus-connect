package com.linkall.sink.mysql.binder;

import com.linkall.sink.mysql.TableMetadata;
import io.vertx.core.json.JsonObject;

import java.sql.PreparedStatement;
import java.sql.SQLException;

public class PreparedStatementBinder implements StatementBinder<JsonObject> {

    private final PreparedStatement preparedStatement;
    private final TableMetadata metadata;

    public PreparedStatementBinder(TableMetadata metadata, PreparedStatement preparedStatement) {
        this.metadata = metadata;
        this.preparedStatement = preparedStatement;
    }

    @Override
    public void bindData(JsonObject data) throws SQLException {
        int index = 1;
        for (String columnName : metadata.getColumnNames()) {
            preparedStatement.setObject(index++, data.getValue(columnName));
        }
    }
}
