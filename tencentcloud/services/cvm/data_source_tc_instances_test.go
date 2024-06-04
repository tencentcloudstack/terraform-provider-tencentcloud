package cvm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmInstancesDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstancesDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_instances.example"), resource.TestCheckResourceAttr("data.tencentcloud_instances.example", "instance_list.0.project_id", "0"), resource.TestCheckResourceAttr("data.tencentcloud_instances.example", "instance_list.0.instance_name", "tf_example")),
			},
		},
	})
}

const testAccCvmInstancesDataSource_BasicCreate = `

data "tencentcloud_instances" "example" {
    project_id = tencentcloud_instance.example.project_id
    vpc_id = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
    
    tags = {
        tagKey = "tagValue"
    }
    instance_id = tencentcloud_instance.example.id
    instance_name = tencentcloud_instance.example.instance_name
    availability_zone = tencentcloud_instance.example.availability_zone
}
resource "tencentcloud_vpc" "vpc" {
    name = "vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    name = "subnet"
    vpc_id = tencentcloud_vpc.vpc.id
    availability_zone = "ap-guangzhou-6"
    cidr_block = "10.0.20.0/28"
    is_multicast = false
}
resource "tencentcloud_instance" "example" {
    instance_type = "SA3.MEDIUM4"
    hostname = "example"
    project_id = 0
    
    data_disks {
        data_disk_size = 50
        encrypt = false
        data_disk_type = "CLOUD_HSSD"
    }
    availability_zone = "ap-guangzhou-6"
    image_id = "img-9qrfy1xt"
    system_disk_type = "CLOUD_HSSD"
    system_disk_size = 100
    vpc_id = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
    
    tags = {
        tagKey = "tagValue"
    }
    instance_name = "tf_example"
}

`
