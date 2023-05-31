package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlRoInstanceIpResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlRoInstanceIp,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_instance_ip.ro_instance_ip", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_ro_instance_ip.ro_instance_ip",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlRoInstanceIp = `

resource "tencentcloud_mysql_ro_instance_ip" "ro_instance_ip" {
  instance_id = ""
  uniq_subnet_id = ""
  uniq_vpc_id = ""
}

`
