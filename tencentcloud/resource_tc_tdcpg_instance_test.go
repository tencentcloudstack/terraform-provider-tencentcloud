package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTdcpgInstance_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdcpgInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_tdcpg_instance.instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_tdcpg_instance.instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTdcpgInstance = `

resource "tencentcloud_tdcpg_instance" "instance" {
  cluster_id = ""
  c_p_u = ""
  memory = ""
  instance_name = ""
  instance_count = ""
  operation_timing = ""
}

`
