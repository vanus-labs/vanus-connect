# Sink APITable Proposal

## Description

The APITable Sink is used to push data (from incoming CloudEvents) to rows/records in APITable.

## Programming Language

-[x] Golang

## Prerequisites

- APITable account
- API Token, visit the workbench of APITable, click on the personal avatar in the lower left corner, and enter [My Setting > Developer]. Click generate token (binding email required for first use).

## Connector Details

### Configuration

| Name        | Required | Default | Description                                                                               |
| :---------- | :------: | :-----: | ----------------------------------------------------------------------------------------- |
| datasheetId | YES | | Datasheet ID |
| token   | YES | | The user authentication token |
| fields  | YES | | Columns in APITable |

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
    "records": [
      {
        "fields": {
          "field1": 18,
          "field2": "xyz",
          "field3": "Development Manager",
          "field4": "2022-07-06"
        }
      }
    ]
  }
}
```

The APITable Sink will create a new row with values (18, "xyz", "Development Manager", "2022-07-06") to each coresponding field names to the Datasheet specified in the configuration.

### Used Libraries/APIs

[APITable SDKs](https://developers.apitable.com/api/quick-start)
