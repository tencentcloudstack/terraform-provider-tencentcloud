package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoSecurityAPIResourceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoSecurityAPIResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_api_resource.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_api_resource.example", "api_resources.0.name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_api_resource.example", "api_resources.0.path", "/api/v1/orders"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_api_resource.example", "api_resources.0.id"),
				),
			},
			{
				Config: testAccTeoSecurityAPIResourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_security_api_resource.example", "api_resources.0.name", "tf-example-update"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_api_resource.example", "api_resources.0.path", "/api/v2/orders"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_security_api_resource.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoSecurityAPIResource = `
resource "tencentcloud_teo_security_api_resource" "example" {
  zone_id = "zone-3fkff38fyw8s"

  api_resources {
    name               = "tf-example"
    path               = "/api/v1/orders"
    api_service_ids    = [tencentcloud_teo_security_api_service.example.api_services[0].id]
    methods            = ["GET", "POST"]
    request_constraint = "$${http.request.body.form['operationType']} in ['query', 'create']"
  }
}
`

const testAccTeoSecurityAPIResourceUpdate = `
resource "tencentcloud_teo_security_api_resource" "example" {
  zone_id = "zone-3fkff38fyw8s"

  api_resources {
    name               = "tf-example-update"
    path               = "/api/v2/orders"
    api_service_ids    = [tencentcloud_teo_security_api_service.example.api_services[0].id]
    methods            = ["GET", "POST"]
    request_constraint = "$${http.request.body.form['operationType']} in ['query', 'create']"
  }
}
`
