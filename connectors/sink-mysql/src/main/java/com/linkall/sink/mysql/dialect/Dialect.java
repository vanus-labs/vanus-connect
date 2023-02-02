package com.linkall.sink.mysql.dialect;

import java.util.Collection;

public interface Dialect {

    String generateInsertSql(String tableName, Collection<String> fieldNames);

    String generateUpsertSql(String tableName, Collection<String> fieldNames);
}
