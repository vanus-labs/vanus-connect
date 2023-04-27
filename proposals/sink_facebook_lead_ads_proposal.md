# Sink Facebook Lead Ads Proposal

## Description

The Facebook Lead Ads Sink Connector send data (from incoming CloudEvents) for ad campaigns, specifying the questions and fields you want to collect from users to create a lead form.

This page describes the design of the Facebook Lead Ads Sink in detail.

## Programming Language

- [x] Golang
- [ ] Java

## Prerequisites

Before you start you will need the following:
- A facebook app with
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
| access_token |    YES  |          | This is the access token that is required to access the Facebook Graph API.
| page_id   |       YES   |          | This is the page ID that is going to run the Facebook Ads.

### Required CloudEvents Data Format

The Facebook Lead Ads Sink requires following JSON data format in CloudEvent's `data` field.

```json
{
    "name": "<FORM_NAME>",
    "follow_up_action_url": "<URL>",
    "questions": {
        {
            "type": "FIRST_NAME"
        },
        {
            "type": "LAST_NAME"
        },
        {
            "type": "EMAIL"
        },
        ...
    }
    "privacy_policy": {
        "url": "<Policy_URL>"
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
    "follow_up_action_url": "<URL>",
    "questions": {
        {
            "type": "FIRST_NAME"
        },
        {
            "type": "LAST_NAME"
        },
        {
            "type": "EMAIL"
        },
    }
    "privacy_policy": {
        "url": "<Policy_URL>"
    }
}
```
In order to create the Lead Ads form successfully, the JSON ``data`` field must contain the following JSON sub fields(``name``, ``follow_up_action_url`` and ``questions``).

The Facebook lead ads sink will create a form ``<FORM_NAME>`` with the following fields: ``FIRST_NAME``, ``LAST_NAME`` and ``EMAIL``. The redirection link after the user fill the form is ``<URL>``.
### Used Libraries/APIs

[Lead Ads from Meta Marketing API](https://developers.facebook.com/docs/marketing-api/guides/lead-ads/create).
