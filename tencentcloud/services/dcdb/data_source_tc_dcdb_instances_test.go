package dcdb_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDCDBInstancesDataSource(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDataSourceDcdbInstances_basic, tcacctest.DefaultDcdbInstanceId, tcacctest.DefaultDcdbInstanceName, tcacctest.DefaultDcdbInsVpcId, tcacctest.DefaultDcdbInsIdSubnetId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_instances.instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_instances.instances", "list.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_instances.instances", "list.0.instance_id", tcacctest.DefaultDcdbInstanceId),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_instances.instances", "list.0.instance_name", tcacctest.DefaultDcdbInstanceName),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.app_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_instances.instances", "list.0.project_id", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_instances.instances", "list.0.region", "ap-guangzhou"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.vpc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.subnet_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.status_desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.vip"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.vport"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.auto_renew_flag"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.storage"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_instances.instances", "list.0.shard_count", "2"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.period_end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.isolated_timestamp"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.uin"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.shard_detail.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_instances.instances", "list.0.shard_detail.#", "2"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.node_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.is_tmp"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_instances.instances", "list.0.wan_domain", ""),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_instances.instances", "list.0.wan_vip", ""),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.wan_port"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.update_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.db_engine"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.db_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.paymode"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.wan_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.is_audit_supported"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.instance_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.0.resource_tags.#"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceDcdbInstances_id, tcacctest.DefaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_instances.instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.#"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceDcdbInstances_name, tcacctest.DefaultDcdbInstanceName),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_instances.instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.#"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceDcdbInstances_disable_excluster, tcacctest.DefaultDcdbInstanceName),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_instances.instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.#"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDataSourceDcdbInstances_disable_vpc_filter, tcacctest.DefaultDcdbInstanceId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_instances.instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_instances.instances", "list.#"),
				),
			},
		},
	})
}

const testAccDataSourceDcdbInstances_basic = `
data "tencentcloud_dcdb_instances" "instances" {
  instance_ids = ["%s"]
  search_name = "instancename"
  search_key = "%s" 
  project_ids = [0]
  is_filter_excluster = true
  excluster_type = 0
  is_filter_vpc = true
  vpc_id = "%s"
  subnet_id = "%s"
}
`

const testAccDataSourceDcdbInstances_id = `
data "tencentcloud_dcdb_instances" "instances" {
  instance_ids = ["%s"]
}
`

const testAccDataSourceDcdbInstances_name = `
data "tencentcloud_dcdb_instances" "instances" {
  search_name = "instancename"
  search_key = "%s" 
}
`

const testAccDataSourceDcdbInstances_disable_excluster = `
data "tencentcloud_dcdb_instances" "instances" {
  search_name = "instancename"
  search_key = "%s" 
  is_filter_excluster = false
  excluster_type = 2
}
`

const testAccDataSourceDcdbInstances_disable_vpc_filter = `
data "tencentcloud_dcdb_instances" "instances" {
  instance_ids = ["%s"]
  is_filter_vpc = false
  vpc_id = "This is a invalid string to verify vpc filter was disabled."
  subnet_id = "This is a invalid string to verify vpc filter was disabled."
}
`
