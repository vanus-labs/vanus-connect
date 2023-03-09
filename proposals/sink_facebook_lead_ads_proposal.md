# Sink Facebook Lead Ads Proposal

## Description

The Facebook Lead Ads Sink Connector send data (from incoming CloudEvents) for ad campaigns, specifying the questions and fields you want to collect from users to create a lead form.

This page describes the design of the Facebook Lead Ads Sink in detail.

## Programming Language

-[x] Golang
-[] Java

## Prerequisites

Before you start you will need the following:
- The ``ads management`` permission
- The ``pages manage ads`` permission
- The ``pages read engagement`` permission
- The ``pages show list`` permission
- A Page access token from a User who can perform the ``ADVERTISE`` task on the Page



## Connector Details

### Configuration

The Facebook Lead Ads Sink needs following configurations to work properly.

| Name        | Required | Default | Description                                                |
|:------------|:--------:|:-------:|------------------------------------------------------------|
| name |   YES    |         | This is the Form name for the Facebook Lead Ads  |
| follow_up_action   |   YES    |         | This is the URL that Facebook will redirect the user after submitting the form.                                  |
| questions   |   YES    |         | This is the question(s) that will be asked in the form.  |
| context_card_id   |   YES     |       | This is the ID of the context card that will be shown to users after they submit the form.
| legal_case_id |   YES  |          | This is the ID of the legal content that will be shown to users before they submit the form.
| access_token |    YES  |          | This is the access token that is required to access the Facebook Graph API.

### Required CloudEvents Data Format

The Facebook Lead Ads Sink requires following JSON data format in CloudEvent's `data` field.

```json
{
  "key1" : "value1",
  "key2" : "value2",
  "key3" : "value3",
  "key4" : "value4",
  ...
}
```

### Connector Behavior

If an incoming CloudEvents looks like:

```json
{
  "id": "88767821-92c2-477d-9a6f-bfdfbed19c6a",
  "source": "quickstart",
  "specversion": "1.0",
  "type": "test",
  "time": "2022-07-08T03:17:03.139Z",
  "datacontenttype": "application/json",
  "data": {
    "name": "<FORM_NAME>",
    "follow_up_action": "<URL>",
    "questions": {
        "type":"<EMAIL>",
        ...
    }
    "context_card_id": "<CONTEXT_CARD_ID>",
    "legal_case_id": "<LEGAL_CASE_ID>"
  }
}
```
In order to create the Lead Ads form successfully, the JSON ``data`` field must contain the following JSON sub fields(``name``, ``follow_up_action``, ``questions``, ``context_card_id`` and ``legal_case_id``). 

### Used Libraries/APIs

[Lead Ads from META Marketing API](https://developers.facebook.com/docs/marketing-api/guides/lead-ads/create).
