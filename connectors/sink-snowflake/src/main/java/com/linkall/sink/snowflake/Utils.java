package com.linkall.sink.snowflake;


public class Utils {


    public static boolean isValidSnowflakeTableName(String tableName) {
        return tableName.matches("^([_a-zA-Z]{1}[_$a-zA-Z0-9]+\\.){0,2}[_a-zA-Z]{1}[_$a-zA-Z0-9]+$");
    }

    public static String getStageName(String schemaName, String tableName) {
        return "VANUS_CONNECTOR_" + schemaName + "_STAGE_" + tableName;
    }

}
