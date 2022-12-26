---
title: AWS S3
---

# AWS S3 Sink

## Overview
The AWS S3 Sink is a Vance Connector which converts the received CloudEvents to JSON format, and uploads them to AWS Object Storage.
The received CloudEvents will be time-based partitioned into chunks. The size of each chunk is determined by the number of CloudEvents
and the scheduled interval of upload a chunk of CloudEvents data. 