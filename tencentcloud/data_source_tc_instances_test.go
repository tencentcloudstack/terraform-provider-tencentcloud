package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDataSourceInstancesBase(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceInstancesBase,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.default"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances.foo", "instance_list.0.instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.foo", "instance_list.0.instance_name", defaultInsName),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances.foo", "instance_list.0.instance_type"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.foo", "instance_list.0.cpu", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.foo", "instance_list.0.memory", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.foo", "instance_list.0.availability_zone", defaultAZone),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.foo", "instance_list.0.project_id", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.foo", "instance_list.0.system_disk_type", "CLOUD_PREMIUM"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceInstancesBase = instanceCommonTestCase + `
data "tencentcloud_instances" "foo" {
  instance_id = "${tencentcloud_instance.default.id}"
  instance_name = "${tencentcloud_instance.default.instance_name}"
}
`
