package com.linkall.sink.aws;

import com.linkall.cdk.Application;

public class Entrance {
    public static void main(String[] args) {
        Application.run(S3Sink.class);
    }
}
