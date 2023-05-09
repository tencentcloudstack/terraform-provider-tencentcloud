package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAPIGatewayAPIAppsResourceName = "data.tencentcloud_api_gateway_api_apps"

// go test -i; go test -test.run TestAccTencentAPIGatewayAPIAppsDataSource_basic -v
func TestAccTencentAPIGatewayAPIAppsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayAPIAppDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayAPIApps(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAPIGatewayAPIAppExists(testAPIGatewayAPIAppResourceName+".test"),
					resource.TestCheckResourceAttr(testAPIGatewayAPIAppsResourceName+".test", "api_app_list.#", "1"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIAppsResourceName+".test", "api_app_list.0.api_app_id"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIAppsResourceName+".test", "api_app_list.0.api_app_name"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIAppsResourceName+".test", "api_app_list.0.api_app_key"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIAppsResourceName+".test", "api_app_list.0.api_app_secret"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIAppsResourceName+".test", "api_app_list.0.created_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIAppsResourceName+".test", "api_app_list.0.modified_time"),
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIAppsResourceName+".test", "api_app_list.0.api_app_desc"),
				),
			},
		},
	})
}

func testAccTestAccTencentAPIGatewayAPIApps() string {
	return `
resource "tencentcloud_api_gateway_api_app" "test" {
  api_app_name = "app_test1"
  api_app_desc = "create app desc"
}

data "tencentcloud_api_gateway_api_apps" "test" {
  api_app_id   = tencentcloud_api_gateway_api_app.test.id
}
`
}
