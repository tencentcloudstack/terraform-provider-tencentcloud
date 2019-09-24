package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.instance"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances.data_instances", "instance_list.0.instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.data_instances", "instance_list.0.instance_name", "tf_test_instance"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances.data_instances", "instance_list.0.instance_type"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.data_instances", "instance_list.0.cpu", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.data_instances", "instance_list.0.memory", "2"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.data_instances", "instance_list.0.availability_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.data_instances", "instance_list.0.project_id", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances.data_instances", "instance_list.0.system_disk_type", "CLOUD_PREMIUM"),
				),
			},
		},
	})
}

const testAccInstancesDataSource = `
data "tencentcloud_image" "my_favorate_image" {
  os_name = "centos"
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

data "tencentcloud_instance_types" "my_favorate_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S2"]
  }
  cpu_core_count = 1
  memory_size    = 2
}

resource "tencentcloud_instance" "instance" {
  instance_name     = "tf_test_instance"
  availability_zone = "ap-guangzhou-3"
  image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
  instance_type     = "${data.tencentcloud_instance_types.my_favorate_instance_types.instance_types.0.instance_type}"
  system_disk_type = "CLOUD_PREMIUM"
}

data "tencentcloud_instances" "data_instances" {
  instance_id = "${tencentcloud_instance.instance.id}"
  instance_name = "${tencentcloud_instance.instance.instance_name}"
}
`
