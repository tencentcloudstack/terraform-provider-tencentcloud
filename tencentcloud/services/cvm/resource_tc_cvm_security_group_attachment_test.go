package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudCvmSecurityGroupAttachmentResource_basic -v
func TestAccTencentCloudCvmSecurityGroupAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmSecurityGroupAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_security_group_attachment.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_security_group_attachment.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_security_group_attachment.example", "security_group_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cvm_security_group_attachment.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmSecurityGroupAttachment = `
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

# attachment security group
resource "tencentcloud_cvm_security_group_attachment" "example" {
  instance_id       = tencentcloud_instance.example.id
  security_group_id = "sg-5275dorp"
}
`
