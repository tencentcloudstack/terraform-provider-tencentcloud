package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	// primary is the SDKv2 implementation of the provider
	primary := tencentcloud.Provider()

	providers := []func() tfprotov5.ProviderServer{
		primary.GRPCProvider, // sdk provider
		providerserver.NewProtocol5(framework.NewProvider(primary)), // framework provider
	}

	// use the muxer
	muxServer, err := tf5muxserver.NewMuxServer(context.Background(), providers...)
	if err != nil {
		log.Fatal(err.Error())
	}

	var serveOpts []tf5server.ServeOpt

	if debugMode {
		serveOpts = append(serveOpts, tf5server.WithManagedDebug())
	}

	err = tf5server.Serve(
		"registry.terraform.io/tencentcloudstack/tencentcloud",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
