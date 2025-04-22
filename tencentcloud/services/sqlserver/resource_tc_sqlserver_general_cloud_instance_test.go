package sqlserver_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSqlserverGeneralCloudInstanceResource_basic -v
func TestAccTencentCloudSqlserverGeneralCloudInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralCloudInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists("tencentcloud_sqlserver_general_cloud_instance.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_cloud_instance.example", "id"),
				),
			},
			{
				ResourceName:            "tencentcloud_sqlserver_general_cloud_instance.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
			{
				Config: testAccSqlserverGeneralCloudInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists("tencentcloud_sqlserver_general_cloud_instance.example"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_cloud_instance.example", "id"),
				),
			},
		},
	})
}

const testAccSqlserverGeneralCloudInstance = tcacctest.DefaultVpcSubnets + tcacctest.DefaultSecurityGroupData + `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_sqlserver_general_cloud_instance" "example" {
  name                 = "tf_example"
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  memory               = 4
  storage              = 100
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

const testAccSqlserverGeneralCloudInstanceUpdate = `
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "sqlserver"
}

resource "tencentcloud_sqlserver_general_cloud_instance" "example" {
  name                 = "tf_example_update"
  zone                 = data.tencentcloud_availability_zones_by_product.zones.zones.4.name
  memory               = 8
  storage              = 200
  cpu                  = 4
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

func TestAccTencentCloudSqlserverGeneralCloudInstanceResource_multiZonesAndMultiNodes(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		CheckDestroy: testAccCheckSqlserverInstanceDestroy,
		Providers:    tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSqlserverGeneralCloudInstance_multiZonesAndMultiNodes,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSqlserverInstanceExists("tencentcloud_sqlserver_general_cloud_instance.multi_zones_multi_nodes"),
					resource.TestCheckResourceAttrSet("tencentcloud_sqlserver_general_cloud_instance.multi_zones_multi_nodes", "id"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_general_cloud_instance.multi_zones_multi_nodes", "multi_zones", "true"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_general_cloud_instance.multi_zones_multi_nodes", "multi_nodes", "true"),
					resource.TestCheckResourceAttr("tencentcloud_sqlserver_general_cloud_instance.multi_zones_multi_nodes", "dr_zones.#", "2"),
				),
			},
			{
				ResourceName:            "tencentcloud_sqlserver_general_cloud_instance.multi_zones_multi_nodes",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"period"},
			},
		},
	})
}

const testAccSqlserverGeneralCloudInstance_multiZonesAndMultiNodes = tcacctest.DefaultVpcSubnets + `
resource "tencentcloud_sqlserver_general_cloud_instance" "multi_zones_multi_nodes" {
  name                 = "multi_zones_multi_nodes"
  zone                 = "ap-guangzhou-3"
  memory               = 4
  storage              = 20
  cpu                  = 2
  machine_type         = "CLOUD_BSSD"
  instance_charge_type = "POSTPAID"
  project_id           = 0
  subnet_id            = local.subnet_id
  vpc_id               = local.vpc_id
  db_version           = "2017"
  security_group_list  = ["sg-kensue7b"]
  weekly               = [1, 2, 3, 5, 6, 7]
  start_time           = "00:00"
  span                 = 6
  resource_tags {
    tag_key   = "test"
    tag_value = "test"
  }
  collation = "Chinese_PRC_CI_AS"
  time_zone = "China Standard Time"
  multi_zones = true
  multi_nodes = true
  dr_zones = ["ap-guangzhou-6", "ap-guangzhou-7"]
}
`
