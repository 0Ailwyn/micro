# Micro

Micro is a suite of libraries and tools for developing and running microservices.

# Overview
The goal of **Micro** is to provide a toolchain for microservice development and management. At the core, micro is simple and accessible enough that anyone can easily get started writing microservices. As you scale to hundreds of services, micro will provide the fundamental tools required to manage a microservice environment.

## Current Tools
- [go-micro](https://github.com/myodc/go-micro) - A microservices client/server library based on http/rpc protobuf
- [api](https://github.com/myodc/micro/tree/master/api) - A lightweight gateway/proxy for Micro based services
- [cli](https://github.com/myodc/micro/tree/master/cli) - A command line tool for micro
- [sidecar](https://github.com/myodc/micro/tree/master/sic) - Integrate any application into the Micro ecosystem

## Example Services
- [geo-srv](https://github.com/myodc/geo-srv) - A go-micro based geolocation tracking service using hailocab/go-geoindex
- [greeter](https://github.com/myodc/micro/tree/master/examples/greeter) - A greeter Go service

## About Microservices
Microservices is an architecture pattern used to decompose a single large application in to a smaller suite of services. Generally the goal is to create light weight services of 1000 lines of code or less. Each service alone provides a particular focused solution or set of solutions. These small services can be used as the foundational building blocks in the creation of a larger system.

The concept of microservices is not new, this is the reimagination of service orientied architecture but with an approach more holistically aligned with unix processes and pipes. For those of us with extensive experience in this field we're somewhat biased and feel this is an incredibly beneficial approach to system design at large and developer productivity.

Learn more about Microservices by watching Martin Fowler's presentation [here](https://www.youtube.com/watch?v=wgdBVIX9ifA) or his blog post [here](http://martinfowler.com/articles/microservices.html).

## Microservice Requirements

The foundation of a library enabling microservices is based around the following requirements:

- Server - an ability to define handlers and serve requests 
- Client - an ability to make requests to another service
- Discovery - a mechanism by which to discover other services

These 3 components form the minimum requirements for microservices development. An ecosystem of libraries and tools can be created around them to provide a feature rich system however at the foundation only these 3 things are required to write services and communicate between them.

### Server

The server is the core component which allows you to register request handlers and serve requests. Ideally it's transport agnostic so different transports such as http, rabbitmq, etc can be chosen. On start it should register itself with discovery system so other microservices know it exists and deregister when shutting down. The server should handle encoding/decoding incoming/outgoing requests, leaving the handlers to operate on the request/response types they expect.

Example interface:
```
server.New(name, options) - instantiate new server
server.Register(handler) - register a handler with the server
server.Start() - start
server.Stop() - stop
```
### Client

Where the server allows you to serve requests, the client lets you make them to other servers. The client should support request/response and pub/sub. Part of the microservices world is event driven programming, taking action based on events, which is why pub/sub is a requirement of the client. It should also make use of the discovery system so requests can be made by service name. 

Example interface:
```
client.Request(name, request) - Make a request to another server
client.Publish(topic, message) - Publish a message on a topic
client.Subscribe(topic, channel) - Subscribe to a topic
```
### Discovery

The discovery system is really vital to microservices development. Any sort of communication between servers will first require locating it and then making the request. Discovery should support registration and retrieval of servers. It should optionally support a keepalive mechanism to remove stale servers.

Example interface:
```
discovery.Register(name, hostname, ...) - Register a server
discovery.Deregister(name, hostname, ...) - Deregister a server
discovery.Get(name) - Get the details for a server
discovery.List() - List all servers
```

## Resources

- [A Journey into Microservices](https://sudo.hailoapp.com/services/2015/03/09/journey-into-a-microservice-world-part-1/)
- [A Journey into a Microservice World](https://speakerdeck.com/mattheath/a-journey-into-a-microservice-world) by Matt Heath (Slides)
- [Microservices](http://martinfowler.com/articles/microservices.html) by Martin Fowler
- [Microservices: Decomposing Applications for Deployability and Scalability](http://www.slideshare.net/chris.e.richardson/microservices-decomposing-applications-for-deployability-and-scalability-jax) by Chris Richardson (Slides)
- [4 reasons why microservices resonate](http://radar.oreilly.com/2015/04/4-reasons-why-microservices-resonate.html) by Neal Ford
