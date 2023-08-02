package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCynosdbSwitchClusterZoneResource_basic -v
func TestAccTencentCloudCynosdbSwitchClusterZoneResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbSwitchClusterZone,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_switch_cluster_zone.switch_cluster_zone", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_switch_cluster_zone.switch_cluster_zone", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_switch_cluster_zone.switch_cluster_zone", "zone", "ap-guangzhou-6"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_switch_cluster_zone.switch_cluster_zone",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCynosdbSwitchClusterZoneUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_switch_cluster_zone.switch_cluster_zone", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_switch_cluster_zone.switch_cluster_zone", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_switch_cluster_zone.switch_cluster_zone", "zone", "ap-guangzhou-4"),
				),
			},
		},
	})
}

const testAccCynosdbSwitchClusterZoneVar = defaultSecurityGroupData + defaultVpcSubnets + `
variable "availability_zone" {
	default = "ap-guangzhou-4"
}

variable "new_availability_zone" {
  default = "ap-guangzhou-6"
}

variable "my_param_template" {
  default = "15765"
}

resource "tencentcloud_cynosdb_cluster" "instance" {
  available_zone               = var.availability_zone
  vpc_id                       = local.vpc_id
  subnet_id                    = local.subnet_id
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf_test_cynosdb_cluster_switch_zone"
  password                     = "cynos@123"
  instance_maintain_duration   = 3600
  instance_maintain_start_time = 10800
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  instance_cpu_core    = 1
  instance_memory_size = 2
  param_items {
    name          = "character_set_server"
    current_value = "utf8"
  }
  param_items {
    name          = "time_zone"
    current_value = "+09:00"
  }

  force_delete = true

  rw_group_sg = [
    local.sg_id
  ]
  ro_group_sg = [
    local.sg_id
  ]
  prarm_template_id = var.my_param_template
}

resource "tencentcloud_cynosdb_cluster_slave_zone" "cluster_slave_zone" {
	cluster_id = tencentcloud_cynosdb_cluster.instance.id
	slave_zone = var.new_availability_zone
  }
`

const testAccCynosdbSwitchClusterZone = testAccCynosdbSwitchClusterZoneVar + `

resource "tencentcloud_cynosdb_switch_cluster_zone" "switch_cluster_zone" {
  cluster_id = tencentcloud_cynosdb_cluster.instance.id
  zone = var.new_availability_zone

  depends_on = [ tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone ]
}

`

const testAccCynosdbSwitchClusterZoneUp = testAccCynosdbSwitchClusterZoneVar + `

resource "tencentcloud_cynosdb_switch_cluster_zone" "switch_cluster_zone" {
  cluster_id = tencentcloud_cynosdb_cluster.instance.id
  zone = var.availability_zone

  depends_on = [ tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone ]
}

`
