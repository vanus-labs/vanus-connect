package com.linkall.source.postgresql;

import io.debezium.connector.postgresql.connection.Lsn;

import java.util.HashMap;
import java.util.Map;

// ["quick_start",{"server":"quick_start"}] ->
// {"last_snapshot_record":true,"lsn":24750664,"txId":746,"ts_usec":1667278571739000,"snapshot":true}
public class PostgreSQLOffset {
    private String lsn;

    public Map<String, Object> getOffset() {
        if (lsn==null) {
            return null;
        }
        Lsn lsn = Lsn.valueOf(this.lsn);
        if (lsn==Lsn.INVALID_LSN) {
            throw new IllegalArgumentException(String.format("offset lsn %s is invalid", this.lsn));
        }
        Map<String, Object> map = new HashMap<>();
        map.put("ts_usec", 0L);
        map.put("lsn", lsn.asLong());
        return map;
    }
}
