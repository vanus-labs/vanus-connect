# Connector Generator

This document provides a brief introduction of the Connector Generator. It's also designed to guide you through 
the process of **generate a connector template** using Connector Generator.

## Introduction

Connector Generator is an executable program designed to generate a connector template. Then developer can 
develop the connector on the basis of the connector template, which will improve development efficiency of 
a connector. The connector template will also give a clear logical framework of the development of a connector
and the directory structure of the connector project. When using Connector Generator, you should set parameters
including cdk, name and type of your connector. And then Connector Generator will build a template for helping 
you develop your connector.

## Quick Start

This quick start will guide you through the process of running a connector generator.

Connector Generator is an executable program built by golang. At first, you should get the executable program
from [vance repository][vance-repo]. 

### Parameters

- `cdk`: the connector's language environment during development, can only be specified as java or go.
- `name`: the connector's name.
- `type`: the connector's type, can only be specified as source or sink.

Run Connector Generator using those parameters like following command.

```shell
./generator -cdk "java" -name "Github" -type "Source"
```

`generator` is the name of Connector Generator. This command should be executed in the peer directory of the generator.

Then you can see the printed logs on the terminal:

```shell
[root@iZ8vbhcwtixrzhsn023sghZ generator]# ./generator -cdk "java" -name "github" -type "source"
Writing file:  ./source-github/README.md
Writing file:  ./source-github/pom.xml
Writing file:  ./source-github/src/main/java/com/vance/source/github/Entrance.java
Writing file:  ./source-github/src/main/java/com/vance/source/github/githubSource.java
Writing file:  ./source-github/src/main/java/com/vance/source/github/githubConfig.java
Writing file:  ./source-github/src/main/java/com/vance/source/github/githubAdapter.java
Writing file:  ./source-github/src/main/java/com/vance/source/github/githubAdaptedSample.java
```

If you use `java` as the `cdk` parameter, after execute the command, you can get a maven project for your connector with
a template to help your development. 

[vance-repo]:https://github.com/linkall-labs/vance

