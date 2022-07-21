package com.linkall.source.github;

import com.linkall.vance.common.config.ConfigUtil;
import com.linkall.vance.common.config.SecretUtil;
import com.linkall.vance.core.Adapter;
import com.linkall.vance.core.Adapter2;
import com.linkall.vance.core.http.HttpClient;
import com.linkall.vance.core.http.HttpResponseInfo;
import io.cloudevents.CloudEvent;
import io.vertx.core.Handler;
import io.vertx.core.Vertx;
import io.vertx.core.buffer.Buffer;
import io.vertx.core.http.HttpServer;
import io.vertx.core.http.HttpServerRequest;
import io.vertx.core.json.JsonObject;
import io.vertx.ext.web.Router;
import io.vertx.ext.web.RoutingContext;
import org.apache.commons.codec.digest.HmacAlgorithms;
import org.apache.commons.codec.digest.HmacUtils;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class GitHubHttpServer implements com.linkall.vance.core.http.HttpServer {

    private static final Logger LOGGER = LoggerFactory.getLogger(GitHubHttpServer.class);
    private static final Vertx vertx = Vertx.vertx();
    private static Router router;
    private final HttpServer httpServer;
    private HttpResponseInfo handlerRI;

    public GitHubHttpServer(){
        this.httpServer = vertx.createHttpServer();
        this.router = Router.router(vertx);
        this.handlerRI = new HttpResponseInfo(200, "Receive success, deliver CloudEvents to" + ConfigUtil.getVanceSink() + "success", 500, "Receive success, deliver CloudEvents to" + ConfigUtil.getVanceSink() + "failure");
    }

    public void init(){
        this.httpServer.requestHandler(this.router);
    }

    @Override
    public <T extends Handler<RoutingContext> & Adapter> void handler(T handler) {
        this.handler("/", handler);
    }

    @Override
    public <T extends Handler<RoutingContext> & Adapter> void handler(String path, T handler) {
        this.router.route(path).handler(handler);
    }

    @Override
    public void simpleHandler(Adapter2<HttpServerRequest, Buffer> adapter2) {
        this.simpleHandler(adapter2, (HttpResponseInfo) null);
    }

    @Override
    public void simpleHandler(Adapter2<HttpServerRequest, Buffer> adapter2, HttpResponseInfo httpResponseInfo) {
        this.router.route("/").handler(request-> {
            String signature = request.request().getHeader("X-Hub-Signature-256");
            request.request().bodyHandler(body -> {
                JsonObject jsonObject = body.toJsonObject();
                String bodyStr = jsonObject.toString();
                if (!verifySignature(signature, bodyStr)) {
                    HttpResponseInfo info = httpResponseInfo;
                    if (null == info) {
                        info = this.handlerRI;
                    }
                    request.response().setStatusCode(info.getFailureCode());
                    request.response().end(info.getFailureChunk());
                    return;
                }
                CloudEvent ce = adapter2.adapt(request.request(), body);
                boolean ret = HttpClient.deliver(ce);
                HttpResponseInfo info = httpResponseInfo;
                if (null == info) {
                    info = this.handlerRI;
                }
                request.response().setStatusCode(info.getSuccessCode());
                request.response().end(info.getSuccessChunk());
            });
            System.out.println("receive the request.");
            String contentType = request.request().getHeader("content-type");
            if (null != contentType && contentType.equals("application/json")) {
                request.request().bodyHandler(body -> {
                    JsonObject jsonObject = body.toJsonObject();
                    String bodyStr = jsonObject.toString();
                    System.out.println(bodyStr);
                });
            }
        });
    }

    @Override
    public void ceHandler(Handler<CloudEvent> handler) {
        this.ceHandler(handler, (HttpResponseInfo)null);
    }

    @Override
    public void ceHandler(Handler<CloudEvent> handler, HttpResponseInfo httpResponseInfo) {

    }

    @Override
    public void listen() {
        int port = Integer.parseInt(ConfigUtil.getPort());
        this.httpServer.listen(port, (server) -> {
            if (server.succeeded()) {
                LOGGER.info("HttpServer is listening on port: " + ((HttpServer)server.result()).actualPort());
            } else {
                LOGGER.error(server.cause().getMessage());
            }

        });
    }
    public boolean verifySignature(String signature, String bodyStr){
        String code = SecretUtil.getString("githubWebHookSecret");

        String hex = "sha256=" + new HmacUtils(HmacAlgorithms.HMAC_SHA_256,code).hmacHex(bodyStr);

        if(!hex.equals(signature)){
            return false;
        }
        return true;
    }


}
