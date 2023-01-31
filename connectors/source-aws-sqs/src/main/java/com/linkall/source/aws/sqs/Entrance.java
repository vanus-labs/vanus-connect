package com.linkall.source.aws.sqs;

import com.linkall.cdk.Application;

public class Entrance {
    public static void main(String[] args) {
        Application.run(SqsSource.class);
    }
}
