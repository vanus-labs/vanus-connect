package com.linkall.sink.feishubot;

import com.linkall.vance.common.config.ConfigUtil;
import com.linkall.vance.common.json.JsonMapper;
import com.linkall.vance.core.Sink;
import com.linkall.vance.core.http.HttpClient;
import com.linkall.vance.core.http.HttpServer;
import io.vertx.core.json.JsonObject;
import org.apache.commons.codec.binary.Base64;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.sql.Timestamp;
import java.util.concurrent.atomic.AtomicInteger;


public class FeishuBotSink implements Sink {
    private static final Logger LOGGER = LoggerFactory.getLogger(FeishuBotSink.class);
    private static final AtomicInteger eventNum = new AtomicInteger(0);
    @Override
    public void start(){
        HttpServer server = HttpServer.createHttpServer();
        server.ceHandler(event -> {
            int num = eventNum.addAndGet(1);
            LOGGER.info("receive a new event, in total: "+num);
            //System.out.println("receive a new event, in total: "+num);
            JsonObject js = JsonMapper.wrapCloudEvent(event);
            String data = js.getString("data");
            if(null != data){
                LOGGER.info("Feishu-BOT receive a message: "+data);
                JsonObject deliverJson = getDeliverJsObject();
                sendTextToFeishuServer(deliverJson,data);
            }else{
                LOGGER.error("invalid data format");
            }

        });
        server.listen();
    }

    private static void sendTextToFeishuServer(JsonObject deliverData, String text){
        deliverData.put("msg_type","text");
        JsonObject content = new JsonObject();
        content.put("text",text);
        deliverData.put("content",content);
        HttpClient.deliver(deliverData);
    }

    private static String GenSign(String secret, long timestamp) throws NoSuchAlgorithmException, InvalidKeyException {
        //Use timestamp + "\n" + key as the signature string
        String stringToSign = timestamp + "\n" + secret;

        //Calculate signature using the HmacSHA256 algorithm
        Mac mac = Mac.getInstance("HmacSHA256");
        mac.init(new SecretKeySpec(stringToSign.getBytes(StandardCharsets.UTF_8), "HmacSHA256"));
        byte[] signData = mac.doFinal(new byte[]{});

        return new String(Base64.encodeBase64(signData));
    }
    private static JsonObject getDeliverJsObject(){
        Timestamp ts = new Timestamp(System.currentTimeMillis());
        long timeStamp =  ts.getTime()/1000;
        String sign = null;
        try {
             sign = GenSign(ConfigUtil.getString("feishu_secret"),timeStamp);
        } catch (NoSuchAlgorithmException e) {
            e.printStackTrace();
        } catch (InvalidKeyException e) {
            e.printStackTrace();
        }
        JsonObject data = new JsonObject();
        data.put("timestamp",String.valueOf(timeStamp));
        data.put("sign",sign);
        data.put("msg_type","interactive");
        return data;
    }

}