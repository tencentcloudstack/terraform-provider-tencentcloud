package cynosdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCynosdbClusterSlaveZoneResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterSlaveZone,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone", "slave_zone", "ap-guangzhou-6"),
				),
			},
			{
				Config: testAccCynosdbClusterSlaveZone_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone", "cluster_id"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone", "slave_zone", "ap-guangzhou-4"),
				),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_cluster_slave_zone.cluster_slave_zone",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbClusterSlaveZone_instance = tcacctest.DefaultSecurityGroupData + tcacctest.DefaultVpcSubnets + `
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
  cluster_name                 = "tf_test_cynosdb_cluster_slave_zone"
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
`

const testAccCynosdbClusterSlaveZone = testAccCynosdbClusterSlaveZone_instance + `

resource "tencentcloud_cynosdb_cluster_slave_zone" "cluster_slave_zone" {
  cluster_id = tencentcloud_cynosdb_cluster.instance.id
  slave_zone = var.new_availability_zone
}

`

const testAccCynosdbClusterSlaveZone_update = testAccCynosdbClusterSlaveZone_instance + `

resource "tencentcloud_cynosdb_cluster_slave_zone" "cluster_slave_zone" {
  cluster_id = tencentcloud_cynosdb_cluster.instance.id
  slave_zone = var.availability_zone
}

`
