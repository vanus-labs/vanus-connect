package com.linkall.sink.mysql;

import com.linkall.cdk.config.Config;
import com.linkall.cdk.connector.Result;
import com.linkall.cdk.connector.Sink;
import com.linkall.sink.mysql.dialect.MySqlDialect;
import io.cloudevents.CloudEvent;
import io.vertx.core.json.JsonObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.sql.SQLException;
import java.util.concurrent.atomic.AtomicLong;

public class MySQLSink implements Sink {
    private static final Logger LOGGER = LoggerFactory.getLogger(MySQLSink.class);
    private DbWriter dbWriter;
    private JdbcConfig config;
    private AtomicLong total = new AtomicLong();

    public String getTableName(CloudEvent event) {
        Object obj = event.getExtension(Constants.ATTRIBUTE_TABLE_NAME);
        String tableName = config.getDbConfig().getTableName();
        if (obj!=null) {
            tableName = obj.toString();
        }
        return tableName;
    }

    public String getSplitColumnName(CloudEvent event) {
        Object obj = event.getExtension(Constants.ATTRIBUTE_SPLIT_COLUMN_NAME);
        if (obj==null) {
            return null;
        }
        return obj.toString();
    }

    @Override
    public Result Arrived(CloudEvent... events) {
        for (CloudEvent event : events) {
            String tableName = getTableName(event);
            String splitColumnName = getSplitColumnName(event);
            JsonObject data = new JsonObject(new String(event.getData().toBytes()));
            try {
                dbWriter.add(tableName, splitColumnName, data);
                LOGGER.info("total receive event:{}", total.incrementAndGet());
            } catch (SQLException e) {
                LOGGER.error("table {} write data has error", tableName, e);
            }
        }
        return Result.SUCCESS;
    }

    @Override
    public Class<? extends Config> configClass() {
        return JdbcConfig.class;
    }

    @Override
    public void initialize(Config config) {
        this.config = (JdbcConfig) config;
        this.dbWriter = new DbWriter(this.config, new MySqlDialect());
    }

    @Override
    public String name() {
        return "MySQL Sink";
    }

    @Override
    public void destroy() throws Exception {
        try {
            dbWriter.close();
        } catch (Exception e) {
            LOGGER.error("sql writer close error", e);
        }
    }
}
