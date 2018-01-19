package tencentcloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-template/template"
)

var testAccProviders map[string]terraform.ResourceProvider

var testAccProvidersWithTLS map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

var testAccTemplateProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccTemplateProvider = template.Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"tencentcloud": testAccProvider,
		"template":     testAccTemplateProvider,
	}
	//testAccProvidersWithTLS = map[string]terraform.ResourceProvider{
	//	"tls": tls.Provider(),
	//}

	//for k, v := range testAccProviders {
	//	testAccProvidersWithTLS[k] = v
	//}
}
