package com.linkall.source.github;

import com.linkall.vance.common.config.ConfigUtil;
import com.linkall.vance.common.config.SecretUtil;
import com.linkall.vance.core.Adapter;
import com.linkall.vance.core.Adapter2;
import com.linkall.vance.core.Source;
import com.linkall.vance.core.http.HttpClient;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.http.HttpServerRequest;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.concurrent.atomic.AtomicInteger;

public class GitHubHttpSource implements Source {

    private static final Logger LOGGER = LoggerFactory.getLogger(GitHubHttpSource.class);
    private static final AtomicInteger eventNum = new AtomicInteger(0);
    @Override
    public Adapter getAdapter() {
        return new GitHubHttpAdapter();
    }

    @Override
    public void start() throws Exception {
        GitHubHttpServer server = new GitHubHttpServer();
        server.init();
        HttpClient.setDeliverSuccessHandler(resp->{
            int num = eventNum.addAndGet(1);
            LOGGER.info("send event in total: "+num);
        });
        server.simpleHandler((Adapter2<HttpServerRequest, Buffer>) getAdapter());
        server.listen();
    }

}
