package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixApigatewayUpstreamsDataSource_basic -v
func TestAccTencentCloudNeedFixApigatewayUpstreamsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayUpstreamDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_api_gateway_upstreams.example"),
				),
			},
		},
	})
}

const testAccApigatewayUpstreamDataSource = `
data "tencentcloud_api_gateway_upstreams" "example" {
  upstream_id = "upstream-4n5bfklc"
}
`
