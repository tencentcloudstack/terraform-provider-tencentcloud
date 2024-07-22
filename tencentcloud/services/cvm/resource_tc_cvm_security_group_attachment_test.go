package cvm_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmSecurityGroupAttachmentResource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmSecurityGroupAttachmentResource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_security_group_attachment.example", "id"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_security_group_attachment.example", "instance_id"), resource.TestCheckResourceAttrSet("tencentcloud_cvm_security_group_attachment.example", "security_group_id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_security_group_attachment.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmSecurityGroupAttachmentResource_BasicCreate = `

data "tencentcloud_instance_types" "default" {
    
    filter {
        name = "instance-family"
        values = ["S1","S2","S3","S4","S5"]
    }
    filter {
        name = "zone"
        values = ["ap-guangzhou-7"]
    }
    cpu_core_count = 2
    memory_size = 2
    exclude_sold_out = true
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
    availability_zone = "ap-guangzhou-6"
    hostname = "example"
    project_id = 0
    vpc_id = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
    
    data_disks {
        encrypt = true
        data_disk_type = "CLOUD_HSSD"
        data_disk_size = 50
    }
    
    tags = {
        tagKey = "tagValue"
    }
    instance_name = "tf_example"
    image_id = "img-9qrfy1xt"
    instance_type = data.tencentcloud_instance_types.default.instance_types.0.instance_type
    system_disk_type = "CLOUD_HSSD"
    system_disk_size = 100
}
resource "tencentcloud_cvm_security_group_attachment" "example" {
    instance_id = tencentcloud_instance.example.id
    security_group_id = "sg-5275dorp"
}

`
