# Source Tiktok Proposal

## Description

The Tiktok Source connector is used to Receive events from TikTok when a user’s post is liked, shared, or commented on, and when a user is followed and put the data into a CloudEvent to send it to the target
This page describes the design of the Tiktok Source in detail.

## Programming Language

-[ ] Golang
-[x] Java

## Prerequisites

- A Tiktok WebHook suscription
- SSL signed certificate

## Connector Details

### Configuration

The Tiktok event have the following specification:

| Key               | Type     | Description                                                                  |
| :---------------- | :------- | :----------------------------------------------------------------------------|
| client_key        | string   | The unique identification key provisioned to the partner.                    |
| event             | string   | Event name.                                                                  |
| create_time       | int64    | The time in which the event occurred. UTC epoch time is in seconds.          |
| user_openid       | string   | The TikTok user's unique identifier; obtained through /oauth/access_token/.  |
| content           | string   | A serialized JSON string of event information.                               |


### Example Payload

```
> {
    "client_key": "bwo2m45353a6k85",
    "event": "video.publish.completed",
    "create_time": 1615338610,
    "user_openid": "act.acv4fasd234asd1c123124asda",
    "content":"{\"share_id\":\"video.6974245311675353080.VDCxrcMJ\"}"
}
```

### Connector Behavior


If an incoming data looks like:

```text
> {
    "client_key": "bwo2m45353a6k85",
    "event": "video.publish.completed",
    "create_time": 1615338610,
    "user_openid": "act.acv4fasd234asd1c123124asda",
    "content":"{\"share_id\":\"video.6974245311675353080.VDCxrcMJ\"}"
}
```
The CloudEvent will look like:

```JSON
{
  "id" : "4ad0b59fc-3e1f-484d-8925-bd78aab15123",
  "source" : "tiktok_event",
  "type" : "event_type",
  "datacontenttype" : "application/json",
  "time" : "2022-09-07T10:21:49.668Z",
  "data" : {
  }
}
```
### Security

To avoid third party inteference, An HMAC with the SHA256 hash function will be computed with the client_secret as the key and the signed_payload string as the message.
To verify the payload, Compare the signature in the header to the generated signature. In the case they are equal, compute the difference between the current timestamp and the received timestamp in the header. Use this to decide whether the difference is tolerable.


