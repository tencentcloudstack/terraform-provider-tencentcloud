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
					resource.TestCheckResourceAttr("tencentcloud_teo_security_api_resource.example", "api_resources.0.name", "tf-example-api-resource"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_api_resource.example", "api_resources.0.path", "/api/v1/orders"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_api_resource.example", "api_resources.0.id"),
				),
			},
			{
				Config: testAccTeoSecurityAPIResourceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_security_api_resource.example", "api_resources.0.name", "tf-example-api-resource-updated"),
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

const testAccTeoSecurityAPIResourceBase = `
variable "zone_id" {
  default = "zone-123sfakjf"
}

variable "api_service_id" {
  default = "apisrv-123233399"
}
`

const testAccTeoSecurityAPIResource = testAccTeoSecurityAPIResourceBase + `
resource "tencentcloud_teo_security_api_resource" "example" {
  zone_id = var.zone_id

  api_resources {
    name            = "tf-example-api-resource"
    path            = "/api/v1/orders"
    api_service_ids = [var.api_service_id]
    methods         = ["GET", "POST"]
  }
}
`

const testAccTeoSecurityAPIResourceUpdate = testAccTeoSecurityAPIResourceBase + `
resource "tencentcloud_teo_security_api_resource" "example" {
  zone_id = var.zone_id

  api_resources {
    name               = "tf-example-api-resource-updated"
    path               = "/api/v2/orders"
    api_service_ids    = [var.api_service_id]
    methods            = ["GET", "POST", "PUT"]
    request_constraint = "${http.request.body.form['type']} in ['query']"
  }
}
`
