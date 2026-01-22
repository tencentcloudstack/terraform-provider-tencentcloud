package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudNeedFixMqttInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMqtt,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "instance_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "sku_code"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "pay_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "device_certificate_provision_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "automatic_activation"),
				),
			},
			{
				Config: testAccMqttUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "instance_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "sku_code"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "pay_mode"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "device_certificate_provision_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance.example", "automatic_activation"),
				),
			},
		},
	})
}

const testAccMqtt = `
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
`

const testAccMqttUpdate = `
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
  name          = "tf-example-update"
  sku_code      = "basic_2k"
  remark        = "remarks update."
  vpc_list {
    vpc_id    = tencentcloud_vpc.vpc.id
    subnet_id = tencentcloud_subnet.subnet.id
  }
  pay_mode                          = 0
  device_certificate_provision_type = "JITP"
  automatic_activation              = true
  tags = {
    createBy = "Terraform"
  }
}
`
