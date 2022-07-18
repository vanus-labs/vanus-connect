package com.linkall.source.aws.s3;

import com.linkall.vance.core.VanceApplication;

public class Entrance {
    public static void main(String[] args) {
        VanceApplication.run(S3Source.class);
    }
}
