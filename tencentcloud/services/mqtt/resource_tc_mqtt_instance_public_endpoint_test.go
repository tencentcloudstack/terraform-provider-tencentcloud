package mqtt_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudMqttInstancePublicEndpointResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMqttInstancePublicEndpoint,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance_public_endpoint.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance_public_endpoint.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance_public_endpoint.example", "bandwidth"),
				),
			},
			{
				Config: testAccMqttInstancePublicEndpointUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance_public_endpoint.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance_public_endpoint.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mqtt_instance_public_endpoint.example", "bandwidth"),
				),
			},
			{
				ResourceName:      "tencentcloud_mqtt_instance_public_endpoint.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMqttInstancePublicEndpoint = `
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

// create public endpoint
resource "tencentcloud_mqtt_instance_public_endpoint" "example" {
  instance_id = tencentcloud_mqtt_instance.example.id
  bandwidth   = 100
  rules {
    ip_rule = "192.168.1.0/24"
    remark  = "Remark."
  }

  rules {
    ip_rule = "172.16.1.0/24"
    remark  = "Remark."
  }
}
`

const testAccMqttInstancePublicEndpointUpdate = `
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

// create public endpoint
resource "tencentcloud_mqtt_instance_public_endpoint" "example" {
  instance_id = tencentcloud_mqtt_instance.example.id
  bandwidth   = 10
  rules {
    ip_rule = "192.168.1.0/24"
    remark  = "Remark."
  }
}
`
