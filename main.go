package main

import (
	"context"
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud"
	"log"
)

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debuggable", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/tencentcloudstack/tencentcloud",
			&plugin.ServeOpts{
				ProviderFunc: tencentcloud.Provider,
			})
		if err != nil {
			log.Println(err.Error())
		}
	} else {
		plugin.Serve(&plugin.ServeOpts{
			ProviderFunc: tencentcloud.Provider})
	}
}
