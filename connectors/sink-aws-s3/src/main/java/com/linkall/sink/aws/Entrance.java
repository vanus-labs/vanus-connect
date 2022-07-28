package com.linkall.sink.aws;

import com.linkall.vance.core.VanceApplication;

public class Entrance {
    public static void main(String[] args) {
        VanceApplication.run(S3Sink.class);
    }
}
