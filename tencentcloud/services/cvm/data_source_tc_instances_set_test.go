package cvm_test

import (
	"testing"

	resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudInstanceSetDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccInstanceSetDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_instances_set.foo"), resource.TestCheckResourceAttr("data.tencentcloud_instances_set.foo", "instance_list.#", "1")),
			},
		},
	})
}

const testAccInstanceSetDataSource_BasicCreate = `

data "tencentcloud_availability_zones" "default" {
}
data "tencentcloud_images" "testing" {
    image_type = ["PUBLIC_IMAGE"]
}
data "tencentcloud_images" "default" {
    image_type = ["PUBLIC_IMAGE"]
    image_name_regex = "Final"
}
data "tencentcloud_instance_types" "default" {
    memory_size = 2
    exclude_sold_out = true
    
    filter {
        values = ["ap-guangzhou-7"]
        name = "zone"
    }
    filter {
        values = ["S1","S2","S3","S4","S5"]
        name = "instance-family"
    }
    cpu_core_count = 2
}
data "tencentcloud_instances_set" "foo" {
    instance_id = tencentcloud_instance.instances_set.id
}
resource "tencentcloud_vpc" "vpc" {
    name = "cvm-basic-vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    name = "cvm-basic-subnet"
    cidr_block = "10.0.0.0/16"
    availability_zone = "ap-guangzhou-7"
    vpc_id = tencentcloud_vpc.vpc.id
}
resource "tencentcloud_instance" "instances_set" {
    availability_zone = "ap-guangzhou-7"
    image_id = data.tencentcloud_images.default.images.0.image_id
    vpc_id = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_PREMIUM"
    project_id = 0
    instance_name = "tf-ci-test"
}

`
