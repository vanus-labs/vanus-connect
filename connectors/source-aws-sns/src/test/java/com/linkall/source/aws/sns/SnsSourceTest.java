package com.linkall.source.aws.sns;

import com.amazonaws.SdkClientException;
import com.amazonaws.services.sns.message.SnsMessageManager;
import com.linkall.source.aws.utils.SNSUtil;
import com.linkall.vance.common.config.ConfigUtil;
import io.vertx.core.json.JsonObject;
import org.junit.BeforeClass;
import org.junit.Test;
import org.mockito.Mockito;
import software.amazon.awssdk.awscore.exception.AwsErrorDetails;
import software.amazon.awssdk.http.SdkHttpResponse;
import software.amazon.awssdk.services.sns.SnsClient;
import software.amazon.awssdk.services.sns.model.*;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.UUID;

import static org.junit.Assert.*;

public class SnsSourceTest{

    SnsSource snsSource = new SnsSource();
    private static SnsClient snsClient;

    String snsTopicArn = ConfigUtil.getString("topic_arn");
    String endpoint = ConfigUtil.getString("endpoint");
    String protocol = ConfigUtil.getString("protocol");
    String region = SNSUtil.getRegion(snsTopicArn);
    private String subscriptionArn = snsTopicArn+":abc";

    @BeforeClass
    public static void beforeInit(){
        snsClient = Mockito.mock(SnsClient.class);
    }

    @Test
    public void testInit() {
        snsSource.init();
        assertEquals(snsTopicArn, snsSource.getTopicArn());
        assertEquals(endpoint, snsSource.getEndPoint());
        assertEquals(protocol, snsSource.getProtocol());
        assertEquals(region, snsSource.getRegion());
    }

    @Test
    public void start(){
        snsSource.start();
    }

    @Test
    public void testSubscribe() {
        SubscribeRequest request = SubscribeRequest.builder()
                .protocol(protocol)
                .endpoint(endpoint)
                .returnSubscriptionArn(true)
                .topicArn(snsTopicArn)
                .build();
        SubscribeResponse response = (SubscribeResponse) SubscribeResponse.builder()
                .subscriptionArn(subscriptionArn)
                .sdkHttpResponse(new SdkHttpResponse() {
                    @Override
                    public Optional<String> statusText() {
                        return Optional.empty();
                    }

                    @Override
                    public int statusCode() {
                        return 200;
                    }

                    @Override
                    public Map<String, List<String>> headers() {
                        return null;
                    }

                    @Override
                    public Builder toBuilder() {
                        return null;
                    }
                })
                .build();
        Mockito.when(snsClient.subscribe(request)).thenReturn(response);
        String ret = snsSource.subscribe(snsClient, snsTopicArn, endpoint, protocol);
        assertEquals(subscriptionArn, ret);

        AwsErrorDetails awsErrorDetails = AwsErrorDetails.builder()
                .errorMessage("An error occurred during subscription.")
                .build();
        SnsException snsException = (SnsException) SnsException.builder()
                .awsErrorDetails(awsErrorDetails)
                .build();
        Mockito.when(snsClient.subscribe(request)).thenThrow(snsException).thenReturn(null);
        String ret1 = snsSource.subscribe(snsClient, snsTopicArn, endpoint, protocol);
        assertEquals("", ret1);
    }

    @Test
    public void testVerifySignature()  {
        SnsMessageManager manager = Mockito.mock(SnsMessageManager.class);
        String subscriptionToken = UUID.randomUUID().toString();
        ConfirmSubscriptionRequest request = ConfirmSubscriptionRequest.builder()
                .token(subscriptionToken)
                .topicArn(snsTopicArn)
                .build();
        ConfirmSubscriptionResponse response = (ConfirmSubscriptionResponse) ConfirmSubscriptionResponse.builder()
                .subscriptionArn(subscriptionArn)
                .sdkHttpResponse(new SdkHttpResponse() {
                    @Override
                    public Optional<String> statusText() {
                        return Optional.empty();
                    }

                    @Override
                    public int statusCode() {
                        return 200;
                    }

                    @Override
                    public Map<String, List<String>> headers() {
                        return null;
                    }

                    @Override
                    public Builder toBuilder() {
                        return null;
                    }
                })
                .build();


        JsonObject jsonObject = new JsonObject();
        jsonObject.put("Type","SubscriptionConfirmation");
        jsonObject.put("MessageId","arn:aws:sns:us-west-2:123456789012:MyTopic");
        jsonObject.put("Subject","My First Message");
        jsonObject.put("Message","Hello world!");
        jsonObject.put("Timestamp","2012-05-02T00:54:06.655Z");
        jsonObject.put("Signature","EXAMPLEw6JRN...");
        jsonObject.put("Token", subscriptionToken);
        InputStream inputStream = new ByteArrayInputStream(jsonObject.toBuffer().getBytes());

        Mockito.when(snsClient.confirmSubscription(request)).thenReturn(response);
        Mockito.when(manager.parseMessage(inputStream)).thenReturn(null);
        boolean ret = snsSource.verifySignature(manager, snsClient, jsonObject, snsTopicArn);
        assertEquals(true, ret);

        Mockito.when(manager.parseMessage(Mockito.any()))
                .thenThrow(new SdkClientException("An error occurred while verifying the signature."));
        ret = snsSource.verifySignature(manager, snsClient, jsonObject, snsTopicArn);
        assertEquals(false, ret);

    }

    @Test
    public void testConfirmSubscriptionException(){
        SnsMessageManager manager = Mockito.mock(SnsMessageManager.class);
        String subscriptionToken = UUID.randomUUID().toString();
        ConfirmSubscriptionRequest request = ConfirmSubscriptionRequest.builder()
                .token(subscriptionToken)
                .topicArn(snsTopicArn)
                .build();

        ConfirmSubscriptionResponse response = (ConfirmSubscriptionResponse) ConfirmSubscriptionResponse.builder()
                .subscriptionArn(subscriptionArn)
                .sdkHttpResponse(new SdkHttpResponse() {
                    @Override
                    public Optional<String> statusText() {
                        return Optional.empty();
                    }

                    @Override
                    public int statusCode() {
                        return 200;
                    }

                    @Override
                    public Map<String, List<String>> headers() {
                        return null;
                    }

                    @Override
                    public Builder toBuilder() {
                        return null;
                    }
                })
                .build();

        JsonObject jsonObject = new JsonObject();
        jsonObject.put("Type","SubscriptionConfirmation");
        jsonObject.put("MessageId","arn:aws:sns:us-west-2:123456789012:MyTopic");
        jsonObject.put("Subject","My First Message");
        jsonObject.put("Message","Hello world!");
        jsonObject.put("Timestamp","2012-05-02T00:54:06.655Z");
        jsonObject.put("Signature","EXAMPLEw6JRN...");
        jsonObject.put("Token", subscriptionToken);
        InputStream inputStream = new ByteArrayInputStream(jsonObject.toBuffer().getBytes());

        Mockito.when(manager.parseMessage(inputStream)).thenReturn(null);

        AwsErrorDetails awsErrorDetails = AwsErrorDetails.builder()
                .errorMessage("An error occurred while confirming subscription.")
                .build();
        SnsException snsException = (SnsException) SnsException.builder()
                .awsErrorDetails(awsErrorDetails)
                .build();

        Mockito.when(snsClient.confirmSubscription(request))
                .thenThrow(snsException).thenReturn(response);
        boolean ret = snsSource.verifySignature(manager, snsClient, jsonObject, snsTopicArn);
        assertEquals(false, ret);
    }

    @Test
    public void testGetAdapter() {
        SnsAdapter snsAdapter = (SnsAdapter)snsSource.getAdapter();
        assertNotNull(snsAdapter);
    }
}