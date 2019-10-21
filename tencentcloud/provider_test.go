package tencentcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider

// var testAccProvidersWithTLS map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]terraform.ResourceProvider{
		"tencentcloud": testAccProvider,
	}

	//testAccProvidersWithTLS = map[string]terraform.ResourceProvider{
	//	"tls": tls.Provider(),
	//}

	//for k, v := range testAccProviders {
	//	testAccProvidersWithTLS[k] = v
	//}
}
