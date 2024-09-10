package cdwdoris_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCdwdorisInstanceResource_basic -v
func TestAccTencentCloudNeedFixCdwdorisInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdwdorisInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_instance.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_instance.example", "zone"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_instance.example", "user_vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_instance.example", "user_subnet_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_instance.example", "product_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_instance.example", "instance_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_instance.example", "doris_user_pwd"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_instance.example", "ha_flag"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_instance.example", "ha_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_instance.example", "case_sensitive"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_instance.example", "enable_multi_zones"),
				),
			},
		},
	})
}

const testAccCdwdorisInstance = `
# availability zone
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "172.16.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  is_multicast      = false
}

# create POSTPAID instance
resource "tencentcloud_cdwdoris_instance" "example" {
  zone               = var.availability_zone
  user_vpc_id        = tencentcloud_vpc.vpc.id
  user_subnet_id     = tencentcloud_subnet.subnet.id
  product_version    = "2.1"
  instance_name      = "tf-example"
  doris_user_pwd     = "Password@test"
  ha_flag            = true
  ha_type            = 1
  case_sensitive     = 0
  enable_multi_zones = false

  charge_properties {
    charge_type = "POSTPAID_BY_HOUR"
  }

  fe_spec {
    spec_name = "S_4_16_P"
    count     = 3
    disk_size = 200
  }

  be_spec {
    spec_name = "S_4_16_P"
    count     = 3
    disk_size = 200
  }

  tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}
`
