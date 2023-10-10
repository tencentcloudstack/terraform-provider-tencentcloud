package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsStartFlowOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsStartFlowOperation_start,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_start_flow_operation.operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_start_flow_operation.operation", "id"),
				),
			},
		},
	})
}

const testAccMpsStartFlowOperation_start = `

resource "tencentcloud_mps_start_flow_operation" "operation" {
  flow_id = "your-flow-id"
  start   = true
}

`
