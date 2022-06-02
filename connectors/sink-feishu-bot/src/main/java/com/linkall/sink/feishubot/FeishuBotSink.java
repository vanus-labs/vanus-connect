package com.linkall.sink.feishubot;

import com.linkall.common.env.EnvUtil;
import com.linkall.common.json.JsonMapper;
import com.linkall.core.Sink;
import com.linkall.core.http.HttpClient;
import com.linkall.core.http.HttpServer;
import io.vertx.core.json.JsonArray;
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
import java.text.SimpleDateFormat;
import java.util.Date;
import java.util.concurrent.atomic.AtomicInteger;


public class FeishuBotSink implements Sink {
    private static final Logger LOGGER = LoggerFactory.getLogger(FeishuBotSink.class);
    private static final AtomicInteger eventNum = new AtomicInteger(0);
    private static final String QUERY_WEATHER = "天气预报";
    private static final String QUERY_WEATHER_BASE = "http://apis.juhe.cn/simpleWeather/query";
    private static final String QUERY_DREAM = "周公解梦";
    private static final String QUERY_DREAM_BASE = "http://v.juhe.cn/dream/query";
    private static final String QUERY_IP_ADDRESS = "IP地址";
    private static final String QUERY_IP_ADDRESS_BASE = "http://apis.juhe.cn/ip/ipNewV3";
    private static final SimpleDateFormat SDF= new SimpleDateFormat("yyyy-MM-dd");
    private static final String WEATHER_TEMPLATE = "{ \"config\": { \"wide_screen_mode\": true }, \"elements\": [ { \"fields\": [ { \"is_short\": true, \"text\": { \"content\": \"**☀️ 温度：\", \"tag\": \"lark_md\" } }, { \"is_short\": true, \"text\": { \"content\": \"**\uD83C\uDF21 湿度：\", \"tag\": \"lark_md\" } } ], \"tag\": \"div\" }, { \"tag\": \"div\", \"fields\": [ { \"is_short\": true, \"text\": { \"tag\": \"lark_md\", \"content\": \"**☁️天气情况：\" } }, { \"is_short\": true, \"text\": { \"tag\": \"lark_md\", \"content\": \"**\uD83D\uDCDD空气质量：\" } } ] }, { \"tag\": \"div\", \"fields\": [ { \"is_short\": true, \"text\": { \"tag\": \"lark_md\", \"content\": \"**\uD83D\uDDF3风向：\" } }, { \"is_short\": true, \"text\": { \"tag\": \"lark_md\", \"content\": \"**\uD83C\uDF91风力：\" } } ] } ], \"header\": { \"template\": \"yellow\", \"title\": { \"content\": \" \", \"tag\": \"plain_text\" } } }";
    private static final String DREAM_TEMPLATE = "{ \"config\": { \"wide_screen_mode\": true }, \"elements\": [ { \"tag\": \"markdown\", \"content\": \"\" } ], \"header\": { \"template\": \"green\", \"title\": { \"content\": \" \", \"tag\": \"plain_text\" } } }";
    @Override
    public void start(){
        HttpServer server = HttpServer.createHttpServer();
        server.ceHandler(event -> {
            int num = eventNum.addAndGet(1);
            LOGGER.info("receive a new event, in total: "+num);
            //System.out.println("receive a new event, in total: "+num);
            JsonObject js = JsonMapper.wrapCloudEvent(event);
            String data = js.getJsonObject("data").getJsonObject("body").getString("data");
            if(null != data){
                LOGGER.info("Feishu-BOT receive a message: "+data);
                handleQuery(data);
            }else{
                LOGGER.error("invalid data format");
            }

        });
        server.listen();
    }

    private void handleQuery(String botMsg){
        if(!botMsg.contains(" ")) {
            LOGGER.info("Not a valid query msg"+botMsg+"\ndirectly send it to "+EnvUtil.getVanceSink());
            JsonObject deliverJson = getDeliverJsObject();
            sendTextToFeishuServer(deliverJson,botMsg);
            return;
        }
        String queryContent = botMsg.substring(botMsg.indexOf(" ")+1);
        LOGGER.info("queryContent: "+queryContent);
        String uri ;
        if(botMsg.startsWith(QUERY_WEATHER)){
            uri = QUERY_WEATHER_BASE+"?city="+queryContent+"&key=89c6e3a553dea9f345a1515582e38acc";
        }else if(botMsg.startsWith(QUERY_DREAM)){
            uri = QUERY_DREAM_BASE+"?q="+queryContent+"&cid=&full=&key=b1473e9ff2d06ce40f46227de236910c";
        }else if(botMsg.startsWith(QUERY_IP_ADDRESS)){
            uri =QUERY_IP_ADDRESS_BASE+"?ip="+queryContent+"&key=d8547e2918b29b4618a223ff86bbd99b";
        }else{
            LOGGER.info("Not a valid query msg"+botMsg+"\ndirectly send it to "+EnvUtil.getVanceSink());
            JsonObject deliverJson = getDeliverJsObject();
            sendTextToFeishuServer(deliverJson,botMsg);
            return;
        }
        LOGGER.info("request URI: "+uri);
        HttpClient.sendGetRequest(uri)
                .onSuccess(response -> {
                    JsonObject resp = new JsonObject(response.body());
                    LOGGER.info(resp.encodePrettily());
                    int errCode = resp.getInteger("error_code");
                    if(errCode !=0) {
                        LOGGER.error("get API result failed: "+resp.getString("reason"));
                    }else{
                        JsonObject deliverJson = getDeliverJsObject();
                        try{
                            handleResult(queryContent,resp.getJsonObject("result"),deliverJson);
                        }catch (Exception e){
                            handleJsonArray(queryContent,resp.getJsonArray("result"),deliverJson);
                        }
                    }
                })
                .onFailure(err ->
                        LOGGER.error("get response err: " + err.getMessage()));
    }

    private static void handleResult(String queryCon,JsonObject result, JsonObject deliverData){
        //handle QUERY_WEATHER data
        if(null!=result.getJsonObject("realtime")){
            JsonObject realTime = result.getJsonObject("realtime");
            JsonObject wtJS = new JsonObject(WEATHER_TEMPLATE);
            JsonObject firstLine = wtJS.getJsonArray("elements").getJsonObject(0);
            JsonObject secondLine = wtJS.getJsonArray("elements").getJsonObject(1);
            JsonObject thirdLine = wtJS.getJsonArray("elements").getJsonObject(2);
            wtJS.getJsonObject("header").getJsonObject("title").put("content","时间："+SDF.format(System.currentTimeMillis())+"  城市："+queryCon);

            JsonObject temperature = firstLine.getJsonArray("fields").getJsonObject(0);
            System.out.println(temperature.getJsonObject("text").put("content",temperature.getJsonObject("text").getString("content")+realTime.getString("temperature")+"℃**"));
            JsonObject humidity = firstLine.getJsonArray("fields").getJsonObject(1);
            System.out.println(humidity.getJsonObject("text").put("content",humidity.getJsonObject("text").getString("content")+realTime.getString("humidity")+"**"));

            JsonObject info = secondLine.getJsonArray("fields").getJsonObject(0);
            JsonObject aqi = secondLine.getJsonArray("fields").getJsonObject(1);
            System.out.println(info.getJsonObject("text").put("content",info.getJsonObject("text").getString("content")+realTime.getString("info")+"**"));
            System.out.println(aqi.getJsonObject("text").put("content",aqi.getJsonObject("text").getString("content")+realTime.getString("aqi")+"**"));

            JsonObject direct = thirdLine.getJsonArray("fields").getJsonObject(0);
            JsonObject power = thirdLine.getJsonArray("fields").getJsonObject(1);
            System.out.println(direct.getJsonObject("text").put("content",direct.getJsonObject("text").getString("content")+realTime.getString("direct")+"**"));
            System.out.println(power.getJsonObject("text").put("content",power.getJsonObject("text").getString("content")+realTime.getString("power")+"**"));

            deliver(deliverData,wtJS);
        }else{
            //deliverData.put("msg_type","text");
            String text = queryCon+"来自"+result.getString("City")+
                    result.getString("District")+
                    result.getString("Isp");
            /*JsonObject content = new JsonObject();
            content.put("text",text);
            deliverData.put("content",content);
            HttpClient.deliver(deliverData);*/
            sendTextToFeishuServer(deliverData, text);
        }
    }

    private static void sendTextToFeishuServer(JsonObject deliverData, String text){
        deliverData.put("msg_type","text");
        JsonObject content = new JsonObject();
        content.put("text",text);
        deliverData.put("content",content);
        HttpClient.deliver(deliverData);
    }

    private static void handleJsonArray(String queryCon,JsonArray result,JsonObject deliverData){
        //handle QUERY_Dream data
        JsonObject dtJS = new JsonObject(DREAM_TEMPLATE);
        JsonObject element = dtJS.getJsonArray("elements").
                getJsonObject(0);
        StringBuilder sb = new StringBuilder();
        result.forEach((e)->{
            String ret = "**"+( (JsonObject) e).getString("title")+":** "+  ( (JsonObject) e).getString("des")+"\n";
            sb.append(ret);
        });
        element.put("content",sb.toString());
        dtJS.getJsonObject("header").getJsonObject("title").put("content","周公解梦："+queryCon);
        System.out.println(element.put("content",sb.toString()));
        deliver(deliverData,dtJS);
    }
    private static void deliver(JsonObject deliverData, JsonObject card){
        deliverData.put("card",card);
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
             sign = GenSign(EnvUtil.getEnvOrConfig("feishu_secret"),timeStamp);
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
    public static void main(String[] args) throws NoSuchAlgorithmException, InvalidKeyException {
        JsonObject deliverJson = getDeliverJsObject();
        JsonObject out = new JsonObject("{ \"Country\" : \"中国\", \"Province\" : \"北京\", \"City\" : \"北京\", \"District\" : \"海淀区\", \"Isp\" : \"联通\" }");

        handleResult("123.113.106.149",out,deliverJson);

        /*JsonObject out = new JsonObject("{\"result\" : [ { \"id\" : \"34f6ef89acb022dd3686355fa93529f2\", \"title\" : \"金钱豹\", \"des\" : \"梦见金钱豹，会与强人为敌。\" }, { \"id\" : \"a88636469d6b9556085eb418d858fa5e\", \"title\" : \"金钱\", \"des\" : \"梦见赢钱，表明你成功和繁荣，钱是信任，自我价值，获得成功，价值。你有很多自己的信念。\" }, { \"id\" : \"ccfb2ac7a5a68ee09cdb0fa177cbcf7d\", \"title\" : \"钱财 金钱 发财\", \"des\" : \"梦见发了财是祥兆。\" } ]}");
        JsonArray oo = out.getJsonArray("result");
        handleJsonArray("金钱",oo,deliverJson);*/

    }
}