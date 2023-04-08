package com.linkall.sink.mysql.dialect;

import com.linkall.sink.mysql.Constants;

import java.util.Collection;
import java.util.stream.Collectors;

public class MySqlDialect implements Dialect {

    @Override
    public String generateInsertSql(String tableName, Collection<String> columnNames) {
        StringBuilder builder = new StringBuilder();
        builder.append("INSERT INTO ");
        builder.append(quoteIdentifier(tableName));
        builder.append("(");
        builder.append(
                columnNames.stream()
                        .map(this::quoteIdentifier)
                        .collect(Collectors.joining(Constants.COMMA)));
        builder.append(")");
        builder.append(" VALUES(");
        builder.append(columnNames.stream().map(x -> "?").collect(Collectors.joining(Constants.COMMA)));
        builder.append(")");
        return builder.toString();
    }

    @Override
    public String generateUpsertSql(String tableName, Collection<String> columnNames) {
        StringBuilder builder = new StringBuilder();
        builder.append(generateInsertSql(tableName, columnNames));
        builder.append(" ON DUPLICATE KEY UPDATE ");
        builder.append(
                columnNames.stream()
                        .map(x -> quoteIdentifier(x) + "=VALUES(" + quoteIdentifier(x) + ")")
                        .collect(Collectors.joining(Constants.COMMA)));
        return builder.toString();
    }

    public String quoteIdentifier(String identifier) {
        return Constants.BACK_QUOTE + identifier + Constants.BACK_QUOTE;
    }
}
