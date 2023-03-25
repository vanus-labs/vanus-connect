# Source Facebook Proposal

## Description
The Facebook Source is used to pull data in form of events from a user's Facebook page. The events include:
* Follow events: When a user's Facebook page is followed
* Like events: When a post on a Facebook page is liked
* Message events: When a user's Facebook page receives a message
* Share events: When a post on a Facebook page is shared.
* When a post on a Facebook page is commented on

## Programming Language
- [ ] Golang
- [x] Java

## Prerequisites
- Create an endpoint on a secure server that can process HTTPS requests.
  Your endpoint must be able to process two types of HTTPS requests: Verification Requests and Event Notifications.
Also, your server must have a valid TLS or SSL certificate correctly configured and installed.

- Configure the Webhooks product in your app's App Dashboard

## Connector Details
### Configuration


| Property | Type   | Description                                                                                                                                      |
|:--------|:-------|:-------------------------------------------------------------------------------------------------------------------------------------------------|
| Object  | String | The object's type (e.g., user, page, etc.)                                                                                                       |
| Entry   | array  | An array containing an object describing the changes. Multiple changes from different objects that are of the same type may be batched together. |
|         |        |                                                                                                                                                  |
| id      | String | The object's id                                                                                                                                  |
| Time    | int    | A UNIX timestamp indicating when the Event Notification was sent                                                                                 |

Most payloads will contain the following common properties, but the contents and structure of each payload varies depending on the object fields you are subscribed to.
### Connector Behavior
When a post on the Facebook page is commented on, you will receive the following notification.
```JSON
{
  "object": "page",
  "entry": [
    {
      "id": "PAGE_ID",
      "time": 1458692752478,
      "messaging": [
        {
          "sender": {
            "id": "USER_ID"
          },
          "recipient": {
            "id": "PAGE_ID"
          },
          "timestamp": 1458692752478,
          "message": {
            "mid": "mid.1457764197618:41d102a3e1ae206a38",
            "text": "hello, world!"
          }
        }
      ]
    }
  ]
}
```
### Used Libraries/APIs
[FaceBook SDK](https://mvnrepository.com/artifact/com.facebook.business.sdk/facebook-java-business-sdk)




