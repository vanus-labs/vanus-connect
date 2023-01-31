package com.linkall.source.aws.sns;

import com.linkall.cdk.Application;

public class Entrance {
    public static void main(String[] args) {
        Application.run(SnsSource.class);
    }
}
