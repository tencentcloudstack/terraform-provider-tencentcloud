package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcCcnInstancesRejectAttachResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcCcnInstancesRejectAttach,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_ccn_instances_reject_attach.ccn_instances_reject_attach", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_ccn_instances_reject_attach.ccn_instances_reject_attach",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcCcnInstancesRejectAttach = `

resource "tencentcloud_vpc_ccn_instances_reject_attach" "ccn_instances_reject_attach" {
  ccn_id = "ccn-gree226l"
  instances {
		instance_id = "vpc-r1x53ogh"
		instance_region = "ap-guangzhou"
		instance_type = ""
		description = ""
		route_table_id = ""

  }
}

`
