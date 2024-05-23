package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudInstanceSetDataSource_Basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstancesSetBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_instances_set.foo", "instance_list.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudInstancesSetBasic = `
data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "default" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
}
data "tencentcloud_images" "testing" {
  image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_instance_types" "default" {

  filter {
    name   = "instance-family"
    values = ["S1", "S2", "S3", "S4", "S5"]
  }
  filter {
    name   = "zone"
    values = ["ap-guangzhou-7"]
  }
  cpu_core_count   = 2
  memory_size      = 2
  exclude_sold_out = true
}
resource "tencentcloud_vpc" "vpc" {
  name       = "cvm-basic-vpc"
  cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-7"
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "cvm-basic-subnet"
  cidr_block        = "10.0.0.0/16"
}
resource "tencentcloud_instance" "instances_set" {
  instance_name     = "tf-ci-test"
  availability_zone = "ap-guangzhou-7"
  image_id          = data.tencentcloud_images.default.images.0.image_id
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  instance_type     = data.tencentcloud_instance_types.default.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  project_id        = 0
}

data "tencentcloud_instances_set" "foo" {
  instance_id = tencentcloud_instance.instances_set.id
}
`
