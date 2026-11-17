package trocket_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTrocketRocketmqInstanceResource_basic -v
func TestAccTencentCloudTrocketRocketmqInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTrocketRocketmqInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "instance_type", "BASIC"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "sku_code", "basic_2k"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "remark", "remark"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "tags.tag_key", "rocketmq"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "tags.tag_value", "basic_2k"),
				),
			},
			{
				Config: testAccTrocketRocketmqInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "name", "tf-example-update"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "instance_type", "BASIC"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "sku_code", "basic_4k"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "remark", "remark"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "tags.tag_key", "rocketmq"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "tags.tag_value", "basic_4k"),
				),
			},
			{
				ResourceName:      "tencentcloud_trocket_rocketmq_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// go test -i; go test -test.run TestAccTencentCloudTrocketRocketmqInstanceResource_enablePublic -v
func TestAccTencentCloudTrocketRocketmqInstanceResource_enablePublic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTrocketRocketmqInstancePublic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "instance_type", "BASIC"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "sku_code", "basic_4k"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "remark", "remark"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "enable_public", "true"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "bandwidth", "10"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "tags.tag_key", "rocketmq"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "tags.tag_value", "basic_4k"),
				),
			},
			{
				ResourceName:      "tencentcloud_trocket_rocketmq_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTrocketRocketmqInstance = `
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

# create rocketmq instance
resource "tencentcloud_trocket_rocketmq_instance" "example" {
  name          = "tf-example"
  instance_type = "BASIC"
  sku_code      = "basic_2k"
  remark        = "remark"
  vpc_id        = tencentcloud_vpc.vpc.id
  subnet_id     = tencentcloud_subnet.subnet.id
  tags = {
    tag_key   = "rocketmq"
    tag_value = "basic_2k"
  }
}
`

const testAccTrocketRocketmqInstanceUpdate = `
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

# create rocketmq instance
resource "tencentcloud_trocket_rocketmq_instance" "example" {
  name          = "tf-example-update"
  instance_type = "BASIC"
  sku_code      = "basic_4k"
  remark        = "remark"
  vpc_id        = tencentcloud_vpc.vpc.id
  subnet_id     = tencentcloud_subnet.subnet.id
  tags = {
    tag_key   = "rocketmq"
    tag_value = "basic_4k"
  }
}
`

const testAccTrocketRocketmqInstancePublic = `
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

# create rocketmq instance
resource "tencentcloud_trocket_rocketmq_instance" "example" {
  name          = "tf-example"
  instance_type = "BASIC"
  sku_code      = "basic_4k"
  remark        = "remark"
  vpc_id        = tencentcloud_vpc.vpc.id
  subnet_id     = tencentcloud_subnet.subnet.id
  enable_public = true
  bandwidth     = 10
  tags = {
    tag_key   = "rocketmq"
    tag_value = "basic_4k"
  }
}
`

// go test -i; go test -test.run TestAccTencentCloudTrocketRocketmqInstanceResource_params -v
func TestAccTencentCloudTrocketRocketmqInstanceResource_params(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTrocketRocketmqInstanceParams,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "name", "tf-params"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "instance_type", "BASIC"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "sku_code", "basic_2k"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "pay_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "max_topic_num", "500"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "zone_ids.#", "1"),
				),
			},
			{
				Config: testAccTrocketRocketmqInstanceParamsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_trocket_rocketmq_instance.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "name", "tf-params-update"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "sku_code", "basic_4k"),
					resource.TestCheckResourceAttr("tencentcloud_trocket_rocketmq_instance.example", "max_topic_num", "1000"),
				),
			},
			{
				ResourceName:      "tencentcloud_trocket_rocketmq_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTrocketRocketmqInstanceParams = `
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

# create rocketmq instance with new params
resource "tencentcloud_trocket_rocketmq_instance" "example" {
  name          = "tf-params"
  instance_type = "BASIC"
  sku_code      = "basic_2k"
  remark        = "remark"
  vpc_id        = tencentcloud_vpc.vpc.id
  subnet_id     = tencentcloud_subnet.subnet.id
  pay_mode      = 0
  renew_flag    = 0
  time_span     = 1
  max_topic_num = 500
  zone_ids      = [100003]
  tags = {
    tag_key   = "rocketmq"
    tag_value = "params"
  }
}
`

const testAccTrocketRocketmqInstanceParamsUpdate = `
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

# update rocketmq instance max_topic_num and zone_ids in-place
resource "tencentcloud_trocket_rocketmq_instance" "example" {
  name          = "tf-params-update"
  instance_type = "BASIC"
  sku_code      = "basic_4k"
  remark        = "remark"
  vpc_id        = tencentcloud_vpc.vpc.id
  subnet_id     = tencentcloud_subnet.subnet.id
  pay_mode      = 0
  renew_flag    = 0
  time_span     = 1
  max_topic_num = 1000
  zone_ids      = [100003]
  tags = {
    tag_key   = "rocketmq"
    tag_value = "params"
  }
}
`
