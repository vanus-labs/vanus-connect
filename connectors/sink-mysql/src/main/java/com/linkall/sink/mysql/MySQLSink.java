package com.linkall.sink.mysql;

import com.linkall.cdk.config.Config;
import com.linkall.cdk.connector.Result;
import com.linkall.cdk.connector.Sink;
import io.cloudevents.CloudEvent;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.SQLException;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.atomic.AtomicLong;

public class MySQLSink implements Sink {
    private static final Logger LOGGER = LoggerFactory.getLogger(MySQLSink.class);
    private SqlWriter sqlWriter;
    private MySQLConfig config;
    private AtomicLong total = new AtomicLong();

    public SqlWriter initSqlWriter(JsonObject data) throws SQLException {
        List<String> columnNames = new ArrayList<>(data.fieldNames());
        TableMetadata meta = new TableMetadata(config.getDbConfig().getTableName(), columnNames);
        SqlWriter sqlWriter = new SqlWriter(this.config, meta);
        sqlWriter.init();
        return sqlWriter;
    }

    @Override
    public Result Arrived(CloudEvent... events) {
        for (CloudEvent event : events) {
            JsonObject data = new JsonObject(new String(event.getData().toBytes()));
            if (sqlWriter==null) {
                try {
                    this.sqlWriter = initSqlWriter(data);
                } catch (SQLException e) {
                    LOGGER.error("get sql writer fail", e);
                    throw new RuntimeException("init sql writer error", e);
                }
            }
            try {
                sqlWriter.add(data);
            } catch (SQLException e) {
                LOGGER.error("write data has error", e);
            }
        }
        return Result.SUCCESS;
    }

    @Override
    public Class<? extends Config> configClass() {
        return MySQLConfig.class;
    }

    @Override
    public void initialize(Config config) throws Exception {
        this.config = (MySQLConfig) config;
    }

    @Override
    public String name() {
        return "MySQL Sink";
    }

    @Override
    public void destroy() throws Exception {
        if (sqlWriter==null) {
            return;
        }
        try {
            sqlWriter.close();
        } catch (Exception e) {
            LOGGER.error("sql writer close error", e);
        }
    }
}
