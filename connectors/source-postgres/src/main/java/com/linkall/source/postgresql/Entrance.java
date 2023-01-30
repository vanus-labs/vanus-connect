package com.linkall.source.postgresql;

import com.linkall.cdk.Application;

public class Entrance {
    public static void main(String[] args) throws Exception {
        Application.run(PostgreSQLSource.class);
    }
}
