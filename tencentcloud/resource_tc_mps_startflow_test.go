package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsStartflowResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsStartflow,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_startflow.startflow", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_startflow.startflow",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsStartflow = `

resource "tencentcloud_mps_startflow" "startflow" {
  flow_id = "your-flow-id"
}

`
