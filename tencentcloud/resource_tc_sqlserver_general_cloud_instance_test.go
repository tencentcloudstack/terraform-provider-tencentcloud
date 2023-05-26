package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverGeneralCloudInstanceResource_basic -v
func TestAccTencentCloudSqlserverGeneralCloudInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralCloudInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists("tencentcloud_sqlserver_general_cloud_instance.general_cloud_instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_cloud_instance.general_cloud_instance", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_sqlserver_general_cloud_instance.general_cloud_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSqlserverGeneralCloudInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists("tencentcloud_sqlserver_general_cloud_instance.general_cloud_instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_cloud_instance.general_cloud_instance", "id"),
				),
			},
		},
	})
}

const testAccSqlserverGeneralCloudInstance = testAccSqlserverBasicInstanceNetwork + `
resource "tencentcloud_sqlserver_general_cloud_instance" "general_cloud_instance" {
  name                 = "create_db_instance"
  zone                 = "ap-guangzhou-6"
  memory               = 4
  storage              = 20
  cpu                  = 2
  machine_type         = "CLOUD_HSSD"
  instance_charge_type = "POSTPAID"
  project_id           = 0
  subnet_id            = local.subnet_id
  vpc_id               = local.vpc_id
  db_version           = "2008R2"
  security_group_list  = [local.sg_id]
  weekly               = [1, 2, 3, 5, 6, 7]
  start_time           = "00:00"
  span                 = 6
  resource_tags {
    tag_key   = "test"
    tag_value = "test"
  }
  collation = "Chinese_PRC_CI_AS"
  time_zone = "China Standard Time"
}
`

const testAccSqlserverGeneralCloudInstanceUpdate = testAccSqlserverBasicInstanceNetwork + `
resource "tencentcloud_sqlserver_general_cloud_instance" "general_cloud_instance" {
  name                 = "update_db_instance"
  zone                 = "ap-guangzhou-6"
  memory               = 8
  storage              = 40
  cpu                  = 4
  machine_type         = "CLOUD_HSSD"
  instance_charge_type = "POSTPAID"
  project_id           = 0
  subnet_id            = local.subnet_id
  vpc_id               = local.vpc_id
  db_version           = "2012SP3"
  security_group_list  = [local.sg_id]
  weekly               = [1, 2, 3, 5, 6, 7]
  start_time           = "00:00"
  span                 = 6
  resource_tags {
    tag_key   = "test"
    tag_value = "test"
  }
  collation = "Chinese_PRC_CI_AS"
  time_zone = "China Standard Time"
}
`
