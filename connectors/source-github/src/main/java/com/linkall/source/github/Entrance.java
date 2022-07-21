package com.linkall.source.github;

import com.linkall.vance.core.VanceApplication;

public class Entrance {
    public static void main(String[] args) {
        VanceApplication.run(GitHubHttpSource.class);
    }
}
