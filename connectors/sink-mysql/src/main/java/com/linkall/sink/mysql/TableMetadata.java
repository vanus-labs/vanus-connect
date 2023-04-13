package com.linkall.sink.mysql;

import java.util.ArrayList;
import java.util.List;
import java.util.Objects;
import java.util.Set;

public class TableMetadata {
    private String tableName;
    private List<String> columnNameList;
    private Set<String> columnNames;

    public TableMetadata(String tableName, Set<String> columnNames) {
        this.tableName = tableName;
        this.columnNames = columnNames;
        this.columnNameList = new ArrayList<>(columnNames);
    }

    public String getTableName() {
        return tableName;
    }

    public List<String> getColumnNames() {
        return columnNameList;
    }

    public boolean columnChange(Set<String> columnNames) {
        if (columnNames==null || columnNames.isEmpty()) {
            return false;
        }
        return !Objects.equals(this.columnNames, columnNames);
    }
}
