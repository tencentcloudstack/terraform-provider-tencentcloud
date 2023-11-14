package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsStopflowResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsStopflow,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_stopflow.stopflow", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_stopflow.stopflow",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsStopflow = `

resource "tencentcloud_mps_stopflow" "stopflow" {
  flow_id = "your flow id"
}

`
