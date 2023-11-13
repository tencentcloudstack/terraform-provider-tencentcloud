package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbRoInstanceIpResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbRoInstanceIp,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_ro_instance_ip.ro_instance_ip", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_ro_instance_ip.ro_instance_ip",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbRoInstanceIp = `

resource "tencentcloud_cdb_ro_instance_ip" "ro_instance_ip" {
  instance_id = ""
  uniq_subnet_id = ""
  uniq_vpc_id = ""
}

`
