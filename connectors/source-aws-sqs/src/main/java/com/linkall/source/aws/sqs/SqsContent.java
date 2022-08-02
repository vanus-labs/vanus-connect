package com.linkall.source.aws.sqs;

public class SqsContent {
    private String msgId;
    private String body;
    private String region;
    private String queueName;

    public String getRegion() {
        return region;
    }

    public void setRegion(String region) {
        this.region = region;
    }

    public String getQueueName() {
        return queueName;
    }

    public void setQueueName(String queueName) {
        this.queueName = queueName;
    }

    public SqsContent(String msgId, String body, String region, String queueName) {
        this.msgId = msgId;
        this.body = body;
        this.region = region;
        this.queueName = queueName;
    }

    public String getMsgId() {
        return msgId;
    }

    public void setMsgId(String msgId) {
        this.msgId = msgId;
    }

    public String getBody() {
        return body;
    }

    public void setBody(String body) {
        this.body = body;
    }
}
