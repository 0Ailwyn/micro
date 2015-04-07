package main

import (
	"code.google.com/p/go.net/context"
	"code.google.com/p/goprotobuf/proto"

	"github.com/asim/go-micro/cmd"
	"github.com/asim/go-micro/server"
	hello "github.com/asim/micro/examples/greeter/server/proto/hello"
	log "github.com/golang/glog"
)

type Say struct{}

func (s *Say) Hello(ctx context.Context, req *hello.Request, rsp *hello.Response) error {
	log.Info("Received Say.Hello request")

	rsp.Msg = proto.String(server.Id + ": Hello " + req.GetName())

	return nil
}

func main() {
	// optionally setup command line usage
	cmd.Init()

	server.Name = "go.micro.srv.greeter"

	// Initialise Server
	server.Init()

	// Register Handlers
	server.Register(
		server.NewReceiver(
			new(Say),
		),
	)

	// Run server
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
