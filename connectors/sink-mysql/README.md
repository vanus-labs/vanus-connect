---
title: MySQL (JDBC)
---

# MySQL Sink (JDBC)

## Introduction

The MySQL Sink is a [Vanus Connector][vc] that aims to handle incoming CloudEvents in a way that extracts the data part
of the original event and delivers these extracted data to a MySQL database using JDBC. 

For example, if the incoming CloudEvent looks like this:

```json
{
  "id": "88767821-92c2-477d-9a6f-bfdfbed19c6a",
  "source": "quickstart",
  "specversion": "1.0",
  "type": "quickstart",
  "time": "2022-07-08T03:17:03.139Z",
  "datacontenttype": "application/json",
  "data": {
    "id": 18,
    "name": "xdl",
    "email": "Development Manager",
    "date": "2022-07-06"
  }
}
```

The MySQL Sink will extract the data fields and write to the database table in the following way:

```text
+----+---------+---------------------+------------+
| id | name    | description         | date       |
+----+---------+---------------------+------------+
| 18 | xdl     | Development Manager | 2022-07-06 |
+----+---------+---------------------+------------+
```

## Quick Start

This quick start will guide you through the process of running an MySQL Sink Connector.

### Prerequisites

- Have a container runtime (i.e., docker).
- Have a running [MySQL][mysql] server.
- Have a database and table created.

### Prepare for db (Optional)

Connect MySQL and Create database and table

```sql
create database vanus_test;
CREATE TABLE IF NOT EXISTS vanus_test.user
(
  `id` int NOT NULL,
  `name` varchar(100) NOT NULL,
  `description` varchar(100) NOT NULL,
  `date` date NOT NULL,
  PRIMARY KEY (`id`)
);
 ```

### Create the config file

```shell
cat << EOF > config.yml
db:
  host: "localhost"
  port: 3306
  username: "vanus_test" 
  password: "123456"
  database: "vanus_test"
  table_name: "user"

insert_mode: UPSERT
EOF
```

| Name            | Required | Default | Description                                                |
|:----------------|:--------:|:-------:|------------------------------------------------------------|
| port            |    NO    |  8080   | the port which MySQL Sink listens on                       |
| db.host         |   YES    |         | IP address or host name of MySQL                           |
| db.port         |   YES    |         | integer port number of MySQL                               |
| db.username     |   YES    |         | username of MySQL                                          |
| db.password     |   YES    |         | password of MySQL                                          |
| db.database     |   YES    |         | database name of MySQL                                     |
| db.table_name   |   YES    |         | table name of MySQL                                        |
| insert_mode     |    NO    | INSERT  | MySQL insert data type: INSERT OR UPSERT                   |
| commit_interval |    NO    |  1000   | MySQL Sink batch data commit interval, unit is millisecond |
| commit_size     |    NO    |  2000   | MySQL Sink batch data commit event size                    |

The MYSQL Sink tries to find the config file at `/vanus-connect/config/config.yml` by default. You can specify the
position of config file by setting the environment variable `CONNECTOR_CONFIG` for your connector.

### Start with Docker

```shell
docker run -it --rm --network=host\
  -v ${PWD}:/vanus-connect/config \
  --name sink-mysql public.ecr.aws/vanus/connector/sink-mysql
```

### Test

Open a terminal and use the following command to send a CloudEvent to the Sink.
The data field must be according to your database structure.

```shell
curl --location --request POST 'localhost:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
  "id" : "88767821-92c2-477d-9a6f-bfdfbed19c6a",
  "source" : "quickstart",
  "specversion" : "1.0",
  "type" : "quickstart",
  "time" : "2022-07-08T03:17:03.139Z",
  "datacontenttype" : "application/json",
  "data" : {
    "id":18,
    "name":"xdl",
    "description":"Development Manager",
    "date": "2022-07-06"
  }
}'
```

Connect to MySQL and use the following command to make sure MySQL has the data

```sql
select * from vanus_test.user;
```

### Clean resource

```shell
docker stop sink-mysql
```

## Run in Kubernetes

```shell
kubectl apply -f sink-mysql.yaml
```

```yaml
apiVersion: v1
kind: Service
metadata:
  name: sink-mysql
  namespace: vanus
spec:
  selector:
    app: sink-mysql
  type: ClusterIP
  ports:
    - port: 8080
      name: sink-mysql
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: sink-mysql
  namespace: vanus
data:
  config.yml: |-
    port: 8080
    db:
      host: "localhost"
      port: 3306
      username: "vanus_test"
      password: "123456"
      database: "vanus_test"
      table_name: "user"

    insert_mode: UPSERT
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: sink-mysql
  namespace: vanus
  labels:
    app: sink-mysql
spec:
  selector:
    matchLabels:
      app: sink-mysql
  replicas: 1
  template:
    metadata:
      labels:
        app: sink-mysql
    spec:
      containers:
        - name: sink-mysql
          image: public.ecr.aws/vanus/connector/sink-mysql
          imagePullPolicy: Always
          resources:
            requests:
              memory: "128Mi"
              cpu: "100m"
            limits:
              memory: "128Mi"
              cpu: "100m"
          ports:
            - name: http
              containerPort: 8080
          volumeMounts:
            - name: config
              mountPath: /vanus-connect/config
      volumes:
        - name: config
          configMap:
            name: sink-mysql
```

## Integrate with Vanus

This section shows how a sink connector can receive CloudEvents from a
running [Vanus cluster](https://github.com/vanus-labs/vanus).

1. Run the sink-mysql.yaml

```shell
kubectl apply -f sink-mysql.yaml
```

2. Create an eventbus

```shell
vsctl eventbus create --name quick-start
```

3. Create a subscription (the sink should be specified as the sink service address or the host name with its port)

```shell
vsctl subscription create \
  --name quick-start \
  --eventbus quick-start \
  --sink 'http://sink-mysql:8080'
```

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect

[mysql]: https://www.mysql.com
