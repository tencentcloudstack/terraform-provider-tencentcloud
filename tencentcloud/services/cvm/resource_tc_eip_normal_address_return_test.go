package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixEipNormalAddressReturnResource_basic -v
func TestAccTencentCloudNeedFixEipNormalAddressReturnResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEipNormalAddressReturn,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_eip_normal_address_return.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_eip_normal_address_return.example", "address_ips.#"),
				),
			},
		},
	})
}

const testAccEipNormalAddressReturn = `
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
  instance_name              = "tf_example"
  availability_zone          = "ap-guangzhou-6"
  image_id                   = "img-9qrfy1xt"
  instance_type              = "SA3.MEDIUM4"
  system_disk_type           = "CLOUD_HSSD"
  system_disk_size           = 100
  hostname                   = "example"
  project_id                 = 0
  vpc_id                     = tencentcloud_vpc.vpc.id
  subnet_id                  = tencentcloud_subnet.subnet.id
  allocate_public_ip         = true
  internet_max_bandwidth_out = 10

  data_disks {
    data_disk_type = "CLOUD_HSSD"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

resource "tencentcloud_eip_normal_address_return" "example" {
  address_ips = [tencentcloud_instance.example.public_ip]
}
`
