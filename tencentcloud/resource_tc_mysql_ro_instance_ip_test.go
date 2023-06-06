package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlRoInstanceIpResource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_instance_ip.ro_instance_ip", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_instance_ip.ro_instance_ip", "uniq_subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_instance_ip.ro_instance_ip", "uniq_vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_instance_ip.ro_instance_ip", "ro_vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_ro_instance_ip.ro_instance_ip", "ro_vport"),
				),
			},
		},
	})
}

const testAccMysqlRoInstanceIp = `

resource "tencentcloud_mysql_ro_instance_ip" "ro_instance_ip" {
	instance_id = "cdbro-f49t0gnj"
	uniq_subnet_id = "subnet-dwj7ipnc"
	uniq_vpc_id = "vpc-4owdpnwr"
  }

`
