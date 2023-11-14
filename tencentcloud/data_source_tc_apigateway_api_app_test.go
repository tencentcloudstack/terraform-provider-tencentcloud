package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApigatewayApiAppDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayApiAppDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_apigateway_api_app.api_app")),
			},
		},
	})
}

const testAccApigatewayApiAppDataSource = `

data "tencentcloud_apigateway_api_app" "api_app" {
  service_id = ""
  a_p_i_ids = 
  filters {
		name = ""
		values = 

  }
    tags = {
    "createdBy" = "terraform"
  }
}

`
