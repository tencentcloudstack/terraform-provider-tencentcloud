package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudApigatewayUpstreamDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayUpstreamDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_apigateway_upstream.upstream")),
			},
		},
	})
}

const testAccApigatewayUpstreamDataSource = `

data "tencentcloud_apigateway_upstream" "upstream" {
  upstream_id = ""
  filters {
		name = ""
		values = 

  }
    tags = {
    "createdBy" = "terraform"
  }
}

`
