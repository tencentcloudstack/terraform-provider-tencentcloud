package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudAsInstancesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsInstancesDataSource_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_as_instances.instances"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_as_instances.instances", "instance_list.#"),
				),
			},
		},
	})
}

func testAccAsInstancesDataSource_basic() string {
	return defaultAsVariable + `
resource "tencentcloud_vpc" "vpc" {
  name       = "tf-as-vpc"
  cidr_block = "10.2.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "tf-as-subnet"
  cidr_block        = "10.2.11.0/24"
  availability_zone = var.availability_zone
}

resource "tencentcloud_as_scaling_config" "launch_configuration" {
  configuration_name = "tf-as-configuration-ds-ins-basic"
  image_id           = "img-2lr9q49h"
  instance_types     = [data.tencentcloud_instance_types.default.instance_types.0.instance_type]
}

resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name = "tf-as-group-ds-ins-basic"
  configuration_id   = tencentcloud_as_scaling_config.launch_configuration.id
  max_size           = 1
  min_size           = 1
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_as_instances" "instances" {
	filters {
		name = "auto-scaling-group-id"
		values = [tencentcloud_as_scaling_group.scaling_group.id]
  }
}

`
}
