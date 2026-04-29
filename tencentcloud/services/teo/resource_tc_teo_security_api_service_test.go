package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoSecurityAPIServiceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoSecurityAPIService,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_api_service.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_api_service.example", "api_services.0.name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_api_service.example", "api_services.0.base_path", "/api/v1"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_security_api_service.example", "api_services.0.id"),
				),
			},
			{
				Config: testAccTeoSecurityAPIServiceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_security_api_service.example", "api_services.0.name", "tf-example-api"),
					resource.TestCheckResourceAttr("tencentcloud_teo_security_api_service.example", "api_services.0.base_path", "/api/v2"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_security_api_service.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoSecurityAPIService = `
resource "tencentcloud_teo_security_api_service" "example" {
  zone_id = "zone-3fkff38fyw8s"

  api_services {
    name      = "tf-example"
    base_path = "/api/v1"
  }
}
`

const testAccTeoSecurityAPIServiceUpdate = `
resource "tencentcloud_teo_security_api_service" "example" {
  zone_id = "zone-3fkff38fyw8s"

  api_services {
    name      = "tf-example"
    base_path = "/api/v2"
  }
}
`
