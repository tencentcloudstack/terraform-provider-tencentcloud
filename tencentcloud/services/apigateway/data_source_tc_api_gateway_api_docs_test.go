package apigateway_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testAPIGatewayAPIDocsResourceName = "data.tencentcloud_api_gateway_api_docs"

// go test -i; go test -test.run TestAccTencentAPIGatewayAPIDocsDataSource_basic -v
func TestAccTencentAPIGatewayAPIDocsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
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
resource "tencentcloud_api_gateway_api_doc" "test" {
  api_doc_name = "doc_test1"
  service_id   = "service-7lybgojo"
  environment  = "release"
  api_ids      = ["api-2bntitvw"]
}

data "tencentcloud_api_gateway_api_docs" "test" {
  depends_on = [tencentcloud_api_gateway_api_doc.test]
}
`
