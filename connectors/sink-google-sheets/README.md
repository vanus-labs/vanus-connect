---
title: Google Sheet
---

# Google Sheet Sink

## Introduction

The Google Sheet Sink is a [Vanus Connector][vc] which aims to send incoming CloudEvents in a JSON format to a Google Sheet.

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

The Google Sheet Sink will extract _data_ field and write it to a Google Sheet.

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
