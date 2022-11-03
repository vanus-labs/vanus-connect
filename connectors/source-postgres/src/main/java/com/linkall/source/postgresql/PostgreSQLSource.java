package com.linkall.source.postgresql;

import com.linkall.source.debezium.DebeziumSource;
import com.linkall.vance.core.Adapter;
import com.linkall.vance.core.Source;
import io.debezium.connector.postgresql.PostgresConnector;

import java.util.*;

public class PostgreSQLSource extends DebeziumSource implements Source {

    private PostgreSQLOffset offset;

    public PostgreSQLSource() {
        offset = new PostgreSQLOffset(config);
    }

    @Override
    public String getConnectorClass() {
        return PostgresConnector.class.getCanonicalName();
    }

    @Override
    public Map<String, Object> getConfigOffset() {
        if (offset.getLsn()!=null) {
            Map<String, Object> map = new HashMap<>();
            map.put("offset", offset.getLsn());
            return map;
        }
        return null;
    }

    @Override
    public Properties getDebeziumProperties() {
        final Properties props = new Properties();

        // convert
//    props.setProperty("converters", "datetime");
//    props.setProperty("datetime.type", PostgresConverter.class.getName());
        props.setProperty("include.unknown.datatypes", "true");
        if (config.containsKey("plugin_name")) {
            props.setProperty("plugin.name", config.getString("plugin_name"));
        } else {
            props.setProperty("plugin.name", "pgoutput");
        }
        if (config.containsKey("snapshot_mode")) {
            props.setProperty("snapshot.mode", config.getString("snapshot_mode"));
        }
        if (config.containsKey("slot_name")) {
            props.setProperty("slot.name", config.getString("slot_name"));
        } else {
            props.setProperty("slot.name", "vance_slot");
        }
        if (config.containsKey("publication_name")) {
            props.setProperty("publication.name", config.getString("publication_name"));
        } else {
            props.setProperty("publication.name", "vance_publication");
        }
        String schemaName = "public";
        if (config.containsKey("schema_name")) {
            schemaName = config.getString("schema_name");
        }
        props.setProperty("schema.include.list", schemaName);
        if (config.containsKey("include_table")) {
            String includeTable = config.getString("include_table");
            props.setProperty("table.include.list", tableFormat(schemaName, Arrays.stream(includeTable.split(","))));
        } else {
            Set<String> excludeTables = new HashSet<>();
            if (config.containsKey("exclude_table")) {
                String excludeTable = config.getString("exclude_table");
                Arrays.stream(excludeTable.split(",")).forEach(v -> {
                    excludeTables.add(v);
                });
            }
            props.setProperty("table.exclude.list", tableFormat(schemaName, getExcludedTables(excludeTables).stream()));
        }

//        props.setProperty("publication.autocreate.mode", "disabled");
        return props;
    }

    public Set<String> getSystemExcludedTables() {
        return new HashSet<>(Arrays.asList("information_schema", "pg_catalog", "pg_internal", "catalog_history"));
    }

    @Override
    public Adapter getAdapter() {
        return new PostgreSQLAdapter();
    }
}
