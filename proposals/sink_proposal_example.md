# Sink XXX Proposal

## Description

The Google Sheets Sink is used to push data (from incoming CloudEvents) to a single Google Spreadsheet.
This page describes the design of the Google Sheets Sink in detail.

## Programming Language

-[x] Golang
-[ ] Java

## Prerequisites

- Google Cloud Account
- A created Google Spreadsheet with pre-set headers

## Connector Details

### Configuration

The Google Sheets Sink needs following configurations to work properly.

| Name        | Required | Default | Description                                                |
|:------------|:--------:|:-------:|------------------------------------------------------------|
| credentials |   YES    |         | The Google Account credentials which authorize the connector to write in the spreadsheet. |
| sheet_url   |   YES    |         | The Google SpreadSheet URL                                 |
| work_sheet  |   YES    |         | The name of the WorkSheet the connector will write data into.  |

### Required CloudEvents Data Format

The Google Sheets Sink requires following JSON data format in CloudEvent's `data` field.

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
    "id": 18,
    "name": "xyz",
    "email": "Development Manager",
    "date": "2022-07-06"
  }
}
```

In order to insert the data to the sheet successfully, the sheet must have **4 headers ("id", "name", "email", "data")**.

The Google Sheets Sink will create a new row with values (18, "xyz", "Development Manager", "2022-07-06") to the SpreadSheet
specified in the configuration. These Values will be placed under the header with the same name as its corresponding key.

### Used Libraries/APIs

[Google Sheets API-Golang](https://developers.google.com/sheets/api/quickstart/go)