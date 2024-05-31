package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudInstanceSetDataSource_Basic -v
func TestAccTencentCloudInstanceSetDataSource_Basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstancesSetBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudInstanceExists("tencentcloud_instance.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances_set.example", "instance_list.0.instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances_set.example", "instance_list.0.instance_name", "tf_example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances_set.example", "instance_list.0.instance_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances_set.example", "instance_list.0.cpu"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances_set.example", "instance_list.0.memory"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances_set.example", "instance_list.0.availability_zone"),
					resource.TestCheckResourceAttr("data.tencentcloud_instances_set.example", "instance_list.0.project_id", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_instances_set.example", "instance_list.0.system_disk_type"),
				),
			},
		},
	})
}

const testAccTencentCloudInstancesSetBasic = `
# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create cvm
resource "tencentcloud_instance" "example" {
  instance_name     = "tf_example"
  availability_zone = "ap-guangzhou-6"
  image_id          = "img-9qrfy1xt"
  instance_type     = "SA3.MEDIUM4"
  system_disk_type  = "CLOUD_HSSD"
  system_disk_size  = 100
  hostname          = "example"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id

  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

data "tencentcloud_instances_set" "example" {
  instance_id       = tencentcloud_instance.example.id
  instance_name     = tencentcloud_instance.example.instance_name
  availability_zone = tencentcloud_instance.example.availability_zone
  project_id        = tencentcloud_instance.example.project_id
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id

  tags = {
    tagKey = "tagValue"
  }
}
`
