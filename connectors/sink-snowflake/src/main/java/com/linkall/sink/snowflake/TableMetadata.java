package com.linkall.sink.snowflake;

import java.util.HashSet;
import java.util.Set;

public class TableMetadata {
    private String tableName;

    private Set<ColumnMetadata> columns;

    public TableMetadata(String tableName) {
        this.tableName = tableName;
        this.columns = new HashSet<>();
    }

    public ColumnType convertColumnType(Object value) {
        if (value==null) {
            return ColumnType.VARIANT;
        }
        if (value instanceof String) {
            return ColumnType.VARCHAR;
        }
        if (value instanceof Number) {
            return ColumnType.NUMBER;
        }
        if (value instanceof Boolean) {
            return ColumnType.BOOLEAN;
        }
        return ColumnType.VARIANT;
    }

    public void addColumn(String name, Object value) {
        ColumnType type = convertColumnType(value);
        columns.add(new ColumnMetadata(name, type));
    }

//    public void addColumn(String name, ColumnType type) {
//        columns.add(new ColumnMetadata(name, type));
//    }

    public String getTableName() {
        return tableName;
    }

    public Set<ColumnMetadata> getColumns() {
        return columns;
    }

    class ColumnMetadata {
        private String name;

        private ColumnType type;

        public ColumnMetadata(String name, ColumnType type) {
            this.name = name;
            this.type = type;
        }

        public String getName() {
            return name;
        }

        public ColumnType getType() {
            return type;
        }
    }

    enum ColumnType {
        BOOLEAN, NUMBER, VARCHAR, VARIANT
    }
}
