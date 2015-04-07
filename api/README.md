# Micro API

This is a lightweight proxy for [Micro](https://github.com/asim/micro) based microservices. It conforms to the [API Gateway](http://microservices.io/patterns/apigateway.html) pattern and can be used in conjuction with [go-micro](https://github.com/asim/go-micro) based apps or any future language implementation of the [Micro](https://github.com/asim/micro) toolchain.

Currently a work in progress.

### Run API
```bash
$ go get github.com/asim/micro
$ micro api
I0308 18:55:50.196418   93745 rpc_server.go:156] Rpc handler /_rpc
I0308 18:55:50.196581   93745 server.go:97] API Rpc handler /rpc
I0308 18:55:50.196672   93745 server.go:108] Listening on [::]:8080
I0308 18:55:50.196707   93745 server.go:90] Starting server go.micro.api id go.micro.api-bcee5e02-c5c4-11e4-a534-68a86d0d36b6
I0308 18:55:50.196782   93745 rpc_server.go:187] Listening on [::]:64983
I0308 18:55:50.196816   93745 server.go:76] Registering go.micro.api-bcee5e02-c5c4-11e4-a534-68a86d0d36b6
```

### Testing API

Let's start the template [go-micro](https://github.com/asim/go-micro) based service.
```bash
$ go get github.com/asim/go-micro/template
$ $GOPATH/bin/template 
I0308 18:58:15.297623   93764 rpc_server.go:156] Rpc handler /_rpc
I0308 18:58:15.297759   93764 server.go:90] Starting server go.micro.service.template id go.micro.service.template-136b13f0-c5c5-11e4-a290-68a86d0d36b6
I0308 18:58:15.297863   93764 rpc_server.go:187] Listening on [::]:65013
I0308 18:58:15.297898   93764 server.go:76] Registering go.micro.service.template-136b13f0-c5c5-11e4-a290-68a86d0d36b6
```

The template service has a handler registered called Example with a method named Call. 
Now let's query this through the API. 
```bash
$ curl -d 'service=go.micro.service.template' -d 'method=Example.Call' -d 'request={"name": "Asim Aslam"}' http://localhost:8080/rpc
{"msg":"go.micro.service.template-e4fc9d93-c5c5-11e4-93bf-68a86d0d36b6: Hello Asim Aslam"}
```

Alternatively let's try 'Content-Type: application/json'
```bash
$ curl -H 'Content-Type: application/json' -d '{"service": "go.micro.service.template", "method": "Example.Call", "request": {"name": "Asim Aslam"}}' http://localhost:8080/rpc
{"msg":"go.micro.service.template-7752615b-c5c5-11e4-a90f-68a86d0d36b6: Hello Asim Aslam"}
```
