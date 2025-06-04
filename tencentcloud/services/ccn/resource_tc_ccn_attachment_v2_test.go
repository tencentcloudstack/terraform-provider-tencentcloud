package ccn_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCcnAttachmentV2Resource_basic(t *testing.T) {
	keyName := "tencentcloud_ccn_attachment_v2.example"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCcnAttachmentV2Attach,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(keyName, "ccn_id"),
					resource.TestCheckResourceAttrSet(keyName, "instance_type"),
					resource.TestCheckResourceAttrSet(keyName, "instance_region"),
					resource.TestCheckResourceAttrSet(keyName, "instance_id"),
					resource.TestCheckResourceAttrSet(keyName, "state"),
					resource.TestCheckResourceAttrSet(keyName, "attached_time"),
					resource.TestCheckResourceAttrSet(keyName, "cidr_block.#"),
				),
			},
			{
				Config: testAccCcnAttachmentV2AttachUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(keyName, "ccn_id"),
					resource.TestCheckResourceAttrSet(keyName, "instance_type"),
					resource.TestCheckResourceAttrSet(keyName, "instance_region"),
					resource.TestCheckResourceAttrSet(keyName, "instance_id"),
					resource.TestCheckResourceAttrSet(keyName, "description"),
					resource.TestCheckResourceAttrSet(keyName, "state"),
					resource.TestCheckResourceAttrSet(keyName, "attached_time"),
					resource.TestCheckResourceAttrSet(keyName, "cidr_block.#"),
				),
			},
			{
				ResourceName:      keyName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCcnAttachmentV2Attach = `
variable "region" {
  default = "ap-guangzhou"
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
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

# create ccn
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
  tags = {
    createBy = "Terraform"
  }
}

# attachment instance
resource "tencentcloud_ccn_attachment_v2" "example" {
  ccn_id          = tencentcloud_ccn.example.id
  instance_id     = tencentcloud_vpc.vpc.id
  instance_type   = "VPC"
  instance_region = var.region
}
`

const testAccCcnAttachmentV2AttachUpdate = `
variable "region" {
  default = "ap-guangzhou"
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
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

# create ccn
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
  tags = {
    createBy = "Terraform"
  }
}

# attachment instance
resource "tencentcloud_ccn_attachment_v2" "example" {
  ccn_id          = tencentcloud_ccn.example.id
  instance_id     = tencentcloud_vpc.vpc.id
  instance_type   = "VPC"
  instance_region = var.region
  description     = "desc update."
}
`
