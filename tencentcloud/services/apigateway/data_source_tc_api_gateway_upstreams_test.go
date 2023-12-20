package apigateway_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixApigatewayUpstreamsDataSource_basic -v
func TestAccTencentCloudNeedFixApigatewayUpstreamsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApigatewayUpstreamDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_api_gateway_upstreams.example"),
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
