package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCdhInstancesDataSource(t *testing.T) {
	dataSourceName := "data.tencentcloud_cdh_instances.list"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccTencentCloudCdhInstancesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID(dataSourceName),
					resource.TestCheckResourceAttrSet(dataSourceName, "cdh_instance_list.0.host_id"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cdh_instance_list.0.create_time"),
					resource.TestCheckResourceAttrSet(dataSourceName, "cdh_instance_list.0.expired_time"),
					resource.TestCheckResourceAttr(dataSourceName, "cdh_instance_list.0.host_name", "unit-test"),
					resource.TestCheckResourceAttr(dataSourceName, "cdh_instance_list.0.host_type", "HC20"),
					resource.TestCheckResourceAttr(dataSourceName, "cdh_instance_list.0.host_state", "RUNNING"),
					resource.TestCheckResourceAttr(dataSourceName, "cdh_instance_list.0.charge_type", "PREPAID"),
					resource.TestCheckResourceAttr(dataSourceName, "cdh_instance_list.0.availability_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr(dataSourceName, "cdh_instance_list.0.project_id", "0"),
					resource.TestCheckResourceAttr(dataSourceName, "cdh_instance_list.0.prepaid_renew_flag", "DISABLE_NOTIFY_AND_MANUAL_RENEW"),
					resource.TestCheckResourceAttr(dataSourceName, "cdh_instance_list.0.host_resource.#", "1"),
				),
			},
		},
	})
}

const TestAccTencentCloudCdhInstancesDataSourceConfig = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_cdh_instance" "foo" {
  availability_zone = var.availability_zone
  host_type = "HC20"
  charge_type = "PREPAID"
  host_name = "unit-test"
  prepaid_renew_flag = "DISABLE_NOTIFY_AND_MANUAL_RENEW"
}

data "tencentcloud_cdh_instances" "list" {
  availability_zone = var.availability_zone
  host_id = tencentcloud_cdh_instance.foo.id
  host_name = tencentcloud_cdh_instance.foo.host_name
  host_state = "RUNNING"
}

`
