# MongoDB Connector

## Support Version

## How to use

### Docker

```yaml
docker run -it --rm public.ecr.aws/vanus/connector/mongodb:latest /run/start.sh \
  --volume /xxx/secret.json /var/mongodb/secret.json \
  --env MONGODB_HOSTS=xxx \
  --env MONGODB_NAME=xxx \
  --env MONGODB_AUTHSOURCE=xxx
```

### k8s

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: mongodb-connector
  labels:
    app: mongodb-connector
spec:
  selector:
    matchLabels:
      app: mongodb-connector
  replicas: 1
  template:
    metadata:
      labels:
        app: mongodb-connector
    spec:
      containers:
        - name: mongodb-connector
          image: public.ecr.aws/vanus/connector/mongodb:latest
          imagePullPolicy: Always
          command: ["sh", "-c", "/var/mongodb/start.sh"]
          resources:
            requests:
              cpu: 100m
              memory: 1000Mi
          env:
            - name: MONGODB_HOSTS
              value: "localhost:27017"
            - name: MONGODB_NAME
              value: "admin"
            - name: MONGODB_AUTHSOURCE
              value: "admin"
            - name: DB_INCLUDE_LIST
              value: "test"  
          volumeMounts:    
          - name: secret     
            mountPath: "/var/mongodb/secret.json"
            readOnly: true
        volumes:
        - name: secret
          secret:
            secretName: mongodb-secret 
---
apiVersion: v1
kind: Secret
metadata:
  name: mongodb-secret
type: Opaque
data:
  user: admin
  password: admin
```