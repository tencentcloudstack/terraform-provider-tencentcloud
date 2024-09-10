package cdwdoris_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCdwdorisWorkloadGroupResource_basic -v
func TestAccTencentCloudNeedFixCdwdorisWorkloadGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdwdorisWorkloadGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_workload_group.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_workload_group.example", "instance_id"),
				),
			},
			{
				Config: testAccCdwdorisWorkloadGroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_workload_group.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cdwdoris_workload_group.example", "instance_id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cdwdoris_workload_group.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdwdorisWorkloadGroup = `
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

# create instance
resource "tencentcloud_cdwdoris_instance" "example" {
  zone                  = var.availability_zone
  user_vpc_id           = tencentcloud_vpc.vpc.id
  user_subnet_id        = tencentcloud_subnet.subnet.id
  product_version       = "2.1"
  instance_name         = "tf-example"
  doris_user_pwd        = "Password@test"
  ha_flag               = false
  case_sensitive        = 0
  enable_multi_zones    = false
  workload_group_status = true

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

# create workload group
resource "tencentcloud_cdwdoris_workload_group" "example" {
  instance_id = tencentcloud_cdwdoris_instance.example.id
  workload_group {
    workload_group_name       = "example"
    cpu_share                 = 1024
    memory_limit              = 20
    enable_memory_over_commit = true
    cpu_hard_limit            = "30%"
  }
}
`

const testAccCdwdorisWorkloadGroupUpdate = `
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

# create instance
resource "tencentcloud_cdwdoris_instance" "example" {
  zone                  = var.availability_zone
  user_vpc_id           = tencentcloud_vpc.vpc.id
  user_subnet_id        = tencentcloud_subnet.subnet.id
  product_version       = "2.1"
  instance_name         = "tf-example"
  doris_user_pwd        = "Password@test"
  ha_flag               = false
  case_sensitive        = 0
  enable_multi_zones    = false
  workload_group_status = true

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

# create workload group
resource "tencentcloud_cdwdoris_workload_group" "example" {
  instance_id = tencentcloud_cdwdoris_instance.example.id
  workload_group {
    workload_group_name       = "example"
    cpu_share                 = 1024
    memory_limit              = 30
    enable_memory_over_commit = false
    cpu_hard_limit            = "20%"
  }
}
`
