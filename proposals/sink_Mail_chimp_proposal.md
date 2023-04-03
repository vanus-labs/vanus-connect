# Sink Mail-chimp Proposal

## Description

The Mail Chimp Sink is used to 
- Allow users to Create a new campaign:
  - You can create a new campaign using the Mailchimp API. This includes specifying the campaign type (regular, automated, plain-text, A/B testing, etc.), the list or segment to which the campaign will be sent, and the content of the campaign.

- Allow users to Create tags:
  - Tags are a way to categorize and organize your subscribers within your mailing list.

- Allow users to Add or update subscribers:
   - With Mailchimp's API, you can add or update subscribers to your mailing list. This includes specifying the email address, name, and any additional fields you've created in your list. You can also update existing subscriber data, such as their name, email address, or list membership.

- Delete subscribers:
   - You can use the Mailchimp API to delete subscribers from your list.

- Remove subscriber from tags:
   - Tags are a way to categorize and organize your subscribers within your mailing list.

- Send to an existing campaign:
   - This will send an email to all subscribers of that campaign.

- Add subscriber to tag
   - add an email address to a tag.

## Programming Language

- [x] Golang

## Prerequisites

- Mail Chimp Account

## Connector Details

### Configuration

The Mail Chimp Sink needs following configurations to work properly.

| Name        | Required | Default | Description                                                                               |
| :---------- | :------: | :-----: | ----------------------------------------------------------------------------------------- |
| api_key |   YES    |         | The api key is needed for the unique identifier of the mail chimp account|
| server_prefix |   YES    |         | The server prefix is needed to tell mail chimp from which domain the campaign will run|

### Required CloudEvents Data Format

The Mail Chimp Sink requires following JSON data format in CloudEvent's `data` field to create a new campaign.

```json
{
  "type": "regular",
  "recipients": {
    "list_id": 1,
    "segment_opts" : {
        ...
    }
  },
  "settings": {
    ...
  }
}
```

### Connector Behavior

If an incoming CloudEvents looks like:

```json
{
  "id": "88767821-92c2-477d-9a6f-bfdfbed19c6a",
  "source": "source-connector",
  "specversion": "1.0",
  "type": " create-campaign",
  "time": "2022-07-08T03:17:03.139Z",
  "datacontenttype": "application/json",
  "data": {
    "type": "regular",
    "recipients": {
        "list_id": 1,
        "segment_opts" : {
            ...
        }
    },
    "settings": {
        ...
    }
  }
}
```
This sink connector will create a campaign with recipients and proper content according to the campaign type successfully.

### Used Libraries/APIs

[Mail Chimp API Docs](https://mailchimp.com/developer/marketing/guides/quick-start/)
