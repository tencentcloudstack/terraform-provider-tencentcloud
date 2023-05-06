package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAPIGatewayAPIDocsResourceName = "data.tencentcloud_api_gateway_api_docs"

// go test -i; go test -test.run TestAccTencentAPIGatewayAPIDocsDataSource_basic -v
func TestAccTencentAPIGatewayAPIDocsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAPIGatewayAPIDocDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTestAccTencentAPIGatewayAPIDocs,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testAPIGatewayAPIDocsResourceName+".test", "api_doc_list.#"),
				),
			},
		},
	})
}

const testAccTestAccTencentAPIGatewayAPIDocs = `
data "tencentcloud_api_gateway_api_docs" "test" {
  
}
`
