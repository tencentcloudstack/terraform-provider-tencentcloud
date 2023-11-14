package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApigatewayAPIDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayAPIDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_apigateway_a_p_i.a_p_i")),
			},
		},
	})
}

const testAccApigatewayAPIDataSource = `

data "tencentcloud_apigateway_a_p_i" "a_p_i" {
  service_id = ""
    tags = {
    "createdBy" = "terraform"
  }
}

`
