package com.linkall.source.github;

import com.fasterxml.jackson.annotation.JsonProperty;

public class SecretConfig {
    @JsonProperty("github_webhook_secret")
    private String githubWebHookSecret;

    public String getGithubWebHookSecret() {
        return githubWebHookSecret;
    }
}
