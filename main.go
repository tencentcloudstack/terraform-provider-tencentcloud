package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()

	// primary is the SDKv2 implementation of the provider. Its code is left
	// untouched: tf5to6server wraps the protocol 5 server in a protocol 6
	// facade, so SDKv2 resources keep behaving exactly as before.
	primary := tencentcloud.Provider()

	// Why protocol 6: the framework half uses nested attributes
	// (schema.ListNestedAttribute), which do not exist in protocol 5 - Terraform
	// rejects them with "protocol version 5 cannot have Attributes set".
	// tfplugin6 is required, and a plugin process can only announce a single
	// protocol version, so the SDKv2 half is upgraded instead of keeping it at
	// v5. The upgrade server is a pure wire-level translator.
	upgradedSdkv2, err := tf5to6server.UpgradeServer(ctx, primary.GRPCProvider)
	if err != nil {
		log.Fatal(err.Error())
	}

	providers := []func() tfprotov6.ProviderServer{
		func() tfprotov6.ProviderServer { return upgradedSdkv2 },
		providerserver.NewProtocol6(framework.NewProvider(primary)), // framework provider
	}

	// use the muxer
	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)
	if err != nil {
		log.Fatal(err.Error())
	}

	var serveOpts []tf6server.ServeOpt

	if debugMode {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/tencentcloudstack/tencentcloud",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
