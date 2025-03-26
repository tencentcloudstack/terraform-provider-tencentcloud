package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttTopicResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMqttTopic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_topic.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_topic.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_topic.example", "topic"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_topic.example", "remark"),
				),
			},
			{
				Config: testAccMqttTopicUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_topic.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_topic.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_topic.example", "topic"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_topic.example", "remark"),
				),
			},
			{
				ResourceName:      "tencentcloud_mqtt_topic.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMqttTopic = `
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
  is_multicast      = false
}

// create mqtt instance
resource "tencentcloud_mqtt_instance" "example" {
  instance_type = "BASIC"
  name          = "tf-example"
  sku_code      = "basic_2k"
  remark        = "remarks."
  vpc_list {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
  }
  pay_mode                          = 0
  device_certificate_provision_type = "API"
  automatic_activation              = false
  tags = {
    createBy = "Terraform"
  }
}

// create topic
resource "tencentcloud_mqtt_topic" "example" {
  instance_id = tencentcloud_mqtt_instance.example.id
  topic       = "tf-example"
  remark      = "Remark."
}
`

const testAccMqttTopicUpdate = `
variable "availability_zone" {
  default = "ap-guangzhou-6"
}

// create vpc
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

// create subnet
resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
  is_multicast      = false
}

// create mqtt instance
resource "tencentcloud_mqtt_instance" "example" {
  instance_type = "BASIC"
  name          = "tf-example"
  sku_code      = "basic_2k"
  remark        = "remarks."
  vpc_list {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
  }
  pay_mode                          = 0
  device_certificate_provision_type = "API"
  automatic_activation              = false
  tags = {
    createBy = "Terraform"
  }
}

// create topic
resource "tencentcloud_mqtt_topic" "example" {
  instance_id = tencentcloud_mqtt_instance.example.id
  topic       = "tf-example"
  remark      = "Remark update."
}
`
