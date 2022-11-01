package com.linkall.source.postgresql;

import io.vertx.core.json.JsonObject;

// {"schema":null,"payload":["vance_test",{"server":"vance_test"}]} ->
// {"last_snapshot_record":true,"lsn":24750664,"txId":746,"ts_usec":1667278571739000,"snapshot":true}
public class PostgreSQLOffset {
    private Long lsn;

    public PostgreSQLOffset(JsonObject config) {
        if (config.containsKey("offset_lsn")) {
            lsn = config.getLong("offset_lsn");
        }
    }

    public Long getLsn() {
        return lsn;
    }

    public void setLsn(Long lsn) {
        this.lsn = lsn;
    }
}
