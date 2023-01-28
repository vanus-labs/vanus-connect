package com.linkall.source.github;

import com.linkall.cdk.Application;

public class Entrance {
    public static void main(String[] args) {
        Application.run(GitHubHttpSource.class);
    }
}
