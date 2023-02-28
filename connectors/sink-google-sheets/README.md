---
title: Google Sheet
---

# Google Sheet Sink

## Introduction

The Google Sheet Sink is a [Vanus Connector][vc] which aims to send incoming CloudEvents in a JSON format to a Google
Sheet.

For example, if the incoming CloudEvent looks like::

```json
{
  "id": "88767821-92c2-477d-9a6f-bfdfbed19c6a",
  "source": "quickstart",
  "specversion": "1.0",
  "type": "quickstart",
  "time": "2022-07-08T03:17:03.139Z",
  "datacontenttype": "application/json",
  "data": {
    "id": "1",
    "name": "Ehis",
    "email": "ehis@gmail.com",
    "description": "Developer"
  }
}
```

The Google Sheet Sink will extract data field and write it to a Google Sheet.

## Pre-requisite

- A Google Sheet
- Service account on the google cloud platform for server authentication - Ensure you give Service Account Editor Access

**Note:** It’s necessary to share the spreadsheet with client_email of the service account to access it. Otherwise, you
will get 403 forbidden. You can find client_email in the downloaded key’s json file

### Config

| Name                | Required |   Default    | Description                                                              |
|:--------------------|:--------:|:------------:|--------------------------------------------------------------------------|
| port                |    NO    |     8080     | the port which Google Sheets Sink listens on                             |
| credentials         |   YES    |              | Google Service account                                                   |
| sheet_id            |   YES    |              | Google sheet ID, example: "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms" |
| sheet_name          |   YES    |              | Google sheet title name, example: "Sheet1"                               |

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
