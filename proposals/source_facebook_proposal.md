# Source Facebook Proposal

## Description
The Facebook Source is used to pull data in form of events from a user's Facebook page. The events include:
* Follow events: When a user's Facebook page is followed
* Like events: When a user's Facebook page is liked
* Message events: When a user's Facebook page receives a message
* Share events: When a user's Facebook page is shared

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
If you successfully subscribe to the page object's field, here's a notification sent when a User posted to a Page.
```JSON
[
  {
    "entry": [
      {
        "changes": [
          {
            "field": "feed",
            "value": {
              "from": {
                "id": "{user-id}",
                "name": "Cinderella Hoover"
              },
              "item": "post",
              "post_id": "{page-post-id}",
              "verb": "add",
              "created_time": 1520544814,
              "is_hidden": false,
              "message": "It's Thursday and I want to eat cake."
            }
          }
        ],
        "id": "{page-id}",
        "time": 1520544816
      }
    ],
    "object": "page"
  }
]
```
### Used Libraries/APIs
[FaceBook SDK](https://mvnrepository.com/artifact/com.facebook.business.sdk/facebook-java-business-sdk)




