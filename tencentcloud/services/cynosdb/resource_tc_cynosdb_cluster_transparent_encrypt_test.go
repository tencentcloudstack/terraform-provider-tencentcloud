package cynosdb_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCynosdbClusterTransparentEncryptResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterTransparentEncrypt,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "key_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "key_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "key_type"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "is_open_global_encryption", "false"),
				),
			},
			{
				ResourceName: "tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt",
				ImportState:  true,
			},
			{
				Config: testAccCynosdbClusterTransparentEncryptUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "key_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "key_region"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "key_type"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_cluster_transparent_encrypt.cynosdb_cluster_transparent_encrypt", "is_open_global_encryption", "true"),
				),
			},
		},
	})
}

const testAccCynosdbClusterTransparentEncrypt = `
resource "tencentcloud_cynosdb_cluster" "foo" {
  auto_renew_flag              = 0
  available_zone               = "ap-guangzhou-3"
  charge_type                  = "POSTPAID_BY_HOUR"
  cluster_name                 = "cynosdbmysql-8iy4bgap"
  db_type                      = "MYSQL"
  db_version                   = "8.0"
  instance_cpu_core            = 2
  instance_maintain_duration   = 3600
  instance_maintain_start_time = 10800
  force_delete                 = true
  # cynos_version                = "3.1.15.006"
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Tue",
    "Wed",
  ]
  instance_memory_size = 4
  port                 = 3306
  project_id           = 0
  password             = "iac@123456"
  rw_group_sg = [
    "sg-5275dorp",
  ]
  serverless_status = null
  storage_pay_mode  = 0
  subnet_id         = "subnet-oi7ya2j6"
  tags              = {}
  vpc_id            = "vpc-axrsmmrv"
  param_items {
    name          = "lower_case_table_names"
    current_value = "0"
  }
  lifecycle {
    ignore_changes = [
      param_items["lower_case_table_names"]
    ]
  }
}
resource "tencentcloud_cynosdb_cluster_transparent_encrypt" "cynosdb_cluster_transparent_encrypt" {
  cluster_id                = tencentcloud_cynosdb_cluster.foo.id
  is_open_global_encryption = false
  key_id                    = "f063c18b-654b-11ef-9d9f-525400d3a886"
  key_region                = "ap-guangzhou"
  key_type                  = "custom"
}
`

const testAccCynosdbClusterTransparentEncryptUp = `
resource "tencentcloud_cynosdb_cluster" "foo" {
  auto_renew_flag              = 0
  available_zone               = "ap-guangzhou-3"
  charge_type                  = "POSTPAID_BY_HOUR"
  cluster_name                 = "cynosdbmysql-8iy4bgap"
  db_type                      = "MYSQL"
  db_version                   = "8.0"
  instance_cpu_core            = 2
  instance_maintain_duration   = 3600
  instance_maintain_start_time = 10800
  force_delete                 = true
  # cynos_version                = "3.1.15.006"
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Tue",
    "Wed",
  ]
  instance_memory_size = 4
  port                 = 3306
  project_id           = 0
  password             = "iac@123456"
  rw_group_sg = [
    "sg-5275dorp",
  ]
  serverless_status = null
  storage_pay_mode  = 0
  subnet_id         = "subnet-oi7ya2j6"
  tags              = {}
  vpc_id            = "vpc-axrsmmrv"
  param_items {
    name          = "lower_case_table_names"
    current_value = "0"
  }
  lifecycle {
    ignore_changes = [
      param_items["lower_case_table_names"]
    ]
  }
}
resource "tencentcloud_cynosdb_cluster_transparent_encrypt" "cynosdb_cluster_transparent_encrypt" {
  cluster_id                = tencentcloud_cynosdb_cluster.foo.id
  is_open_global_encryption = true
  key_id                    = "f063c18b-654b-11ef-9d9f-525400d3a886"
  key_region                = "ap-guangzhou"
  key_type                  = "custom"
}
`
