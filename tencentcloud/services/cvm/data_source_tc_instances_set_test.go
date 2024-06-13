package cvm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmInstancesSetDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmInstancesSetDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_instances_set.example"), resource.TestCheckResourceAttr("data.tencentcloud_instances_set.example", "instance_list.0.project_id", "0"), resource.TestCheckResourceAttr("data.tencentcloud_instances_set.example", "instance_list.0.instance_name", "tf_example")),
			},
		},
	})
}

const testAccCvmInstancesSetDataSource_BasicCreate = `

data "tencentcloud_instances_set" "example" {
    availability_zone = tencentcloud_instance.example.availability_zone
    project_id = tencentcloud_instance.example.project_id
    vpc_id = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
    
    tags = {
        tagKey = "tagValue"
    }
    instance_id = tencentcloud_instance.example.id
    instance_name = tencentcloud_instance.example.instance_name
}
resource "tencentcloud_vpc" "vpc" {
    name = "vpc"
    cidr_block = "10.0.0.0/16"
}
resource "tencentcloud_subnet" "subnet" {
    vpc_id = tencentcloud_vpc.vpc.id
    availability_zone = "ap-guangzhou-6"
    cidr_block = "10.0.20.0/28"
    is_multicast = false
    name = "subnet"
}
resource "tencentcloud_instance" "example" {
    instance_name = "tf_example"
    availability_zone = "ap-guangzhou-6"
    hostname = "example"
    project_id = 0
    vpc_id = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
    image_id = "img-9qrfy1xt"
    instance_type = "SA3.MEDIUM4"
    system_disk_type = "CLOUD_HSSD"
    system_disk_size = 100
    
    data_disks {
        data_disk_type = "CLOUD_HSSD"
        data_disk_size = 50
        encrypt = false
    }
    
    tags = {
        tagKey = "tagValue"
    }
}

`
