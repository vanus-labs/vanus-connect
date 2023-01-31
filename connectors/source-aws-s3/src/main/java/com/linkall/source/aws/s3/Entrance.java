package com.linkall.source.aws.s3;


import com.linkall.cdk.Application;

public class Entrance {
    public static void main(String[] args) {
        Application.run(S3Source.class);
    }
}
