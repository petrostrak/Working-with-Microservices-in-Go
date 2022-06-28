### Working with Microservices in Go
Build highly available, scalable, resilient distributed applications using Go

For a long time, web applications were usually a single application that handled everything—in other words, a `monolithic` application. This monolith handled user authentication, logging, sending email, and everything else. While this is still a popular (and useful) approach, today, many larger scale applications tend to break things up into `microservices`. Today, most large organizations are focused on building web applications using this approach, and with good reason.

`Microservices`, also known as the microservice architecture, are an architectural style which structures an application as a loosely coupled collection of smaller applications. The microservice architecture allows for the rapid and reliable delivery of large, complex applications. Some of the most common features for a microservice are:

* it is maintainable and testable;

* it is loosely coupled with other parts of the application;

* it  can deployed by itself;

* it is organized around business capabilities;

* it is often owned by a small team.

In this repo, we've developed a number of small, self-contained, loosely coupled microservices that will will communicate with one another and a simple front-end application with a `REST` API, with `RPC`, over `gRPC`, and by sending and consuming messages using `AMQP`, the Advanced Message Queuing Protocol. The microservices we build will include the following functionality:

* A Front End service, that just displays web pages;

* An Authentication service, with a Postgres database;

* A Logging service, with a MongoDB database;

* A Listener service, which receives messages from RabbitMQ and acts upon them;

* A Broker service, which is an optional single point of entry into the microservice cluster;

* A Mail service, which takes a JSON payload, converts into a formatted email, and send it out.

All of these services are written in Go, commonly referred to as Golang, a language which is particularly well suited to building distributed web applications.

We deploy the distributed application to a `Docker Swarm` and `Kubernetes`, to scale up and down, as necessary, and to update individual microservices with little or no downtime.

### Makefile commands
To starts all containers in the background without forcing build:
```
make up
```

To stops docker-compose (if running), builds all projects and starts docker compose:
```
make up_build
```

To stop docker compose:
```
make down
```

To build the broker binary as a linux executable:
```
make build_broker
```

To build the auth binary as a linux executable:
```
make build_auth
```

To build the logger binary as a linux executable:
```
make build_logger
```

To build the front-end binary:
```
make build_front
```

To start the front-end:
```
make start
```

To stop the front-end:
```
make stop
```
