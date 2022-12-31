# The Specification of Connector Development

## Project Layout
When contributing a new connector, you MUST follow the [templates](templates) directory 
structure to guarantee that all necessary files are added. 

For the consistent developer's experience, each connector's README.md SHOULD be created according to [README-template.md](templates/README.md)

## Developer Experience

This section gives some restrictions to make all connectors have a consistent style, which can ensure developers have a better experience in learning, understanding, and developing.

### Language
Because the connectors aim to be serverless application, So, each connector's programming SHOULD prefer to Golang in 
order to meet the minimum package size and maximum speed of starting application, unless there is the special consideration, 
like library and eco-system.


### Naming

#### Java
- package: `com.linkall.connector.{connector_name}.*`

#### Golang
- module: `github.com/linkall-labs/connector/{connector_name}`

### Error
TODO

### Log
TODO

### Testing
Because connectors really are stateless application, so the unit testing is pretty important to connectors. Thus, each 
connector's ut coverage should be greater than 80%.

### Observability
TODO


## Configuration
each connector will have 2 config files:
- **config.json**: including all properties of connector needs, except the secret information.
- **secret.json**: any sensitive property。

## Deploy
each connector should provide 3 methods to run:
- **docker**: how to run connector in a docker engine.
- **k8s**: how to run connector in k8s cluster.

These are already included in [templates](templates/README.md), whose 'how to use' section has been displayed it.

## How to create a new connector

### RDD

If you want to create a new connector, you MUST finish the `README.md` firstly. We call this the **RDD(README Drive Development)**.

The reason is that `README.md` is the first thing that the users will know what this connector is and how to use it.
There are many sections that users will care about in [README-template.md](templates/README.md), so we can think as a real
user and pay attention in details by writing readme doc.

When the readme doc is finished, the connector's design also almost be finished.

### Proposal

When you finished the `README.md`, you can create a PR to submit your file. if the README would be accepted, follow the
[developer instruction](#) to start developing if you want. If you don't want to implement it by yourself, the vance community
will make it implemented.
