package com.linkall.sink.mysql;

import com.linkall.cdk.Application;

public class Entrance {
    public static void main(String[] args) {
        Application.run(MySQLSink.class);
    }
}
