package apigateway_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudAPIGateWayUpdateApiAppKey_basic -v
func TestAccTencentCloudAPIGateWayUpdateApiAppKey_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccUpdateApiAppKey,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_update_api_app_key.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_update_api_app_key.example", "api_app_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_api_gateway_update_api_app_key.example", "api_app_key"),
				),
			},
		},
	})
}

const testAccUpdateApiAppKey = `
resource "tencentcloud_api_gateway_update_api_app_key" "example" {
  api_app_id  = "app-krljp4wn"
  api_app_key = "APID6JmG21yRCc03h4z16hlsTqj1wpO3dB3ZQcUP"
}
`
