package ckafka_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaRouteResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaRouteBasic,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ckafka_route.route", "id")),
			},
		},
	})
}

func TestAccTencentCloudCkafkaRouteResource_withIp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaRouteBasicWithIp,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ckafka_route.route_with_ip", "id")),
			},
		},
	})
}

func TestAccTencentCloudCkafkaRouteResource_securityGroupIds(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaRouteSecurityGroupIds,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_route.route_with_sg", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_route.route_with_sg", "security_group_ids.#", "1"),
				),
			},
			{
				Config: testAccCkafkaRouteSecurityGroupIdsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_route.route_with_sg", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_route.route_with_sg", "security_group_ids.#", "2"),
				),
			},
			{
				Config: testAccCkafkaRouteSecurityGroupIdsRemove,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_route.route_with_sg", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_route.route_with_sg", "security_group_ids.#", "0"),
				),
			},
		},
	})
}

const testAccCkafkaRouteBasic = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_ckafka_route" "route" {
	instance_id = var.instance_id
	vip_type = 3
	vpc_id = "vpc-hyuojwa7"
	subnet_id = "subnet-96ufvn1s"
	access_type = 0
	public_network = 3
  }
`

const testAccCkafkaRouteBasicWithIp = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_ckafka_route" "route_with_ip" {
	instance_id = var.instance_id
	vip_type = 3
	vpc_id = "vpc-hyuojwa7"
	subnet_id = "subnet-96ufvn1s"
	access_type = 0
	public_network = 3
	ip = "10.0.0.55"
  }
`

const testAccCkafkaRouteSecurityGroupIds = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_ckafka_route" "route_with_sg" {
	instance_id = var.instance_id
	vip_type = 3
	vpc_id = "vpc-hyuojwa7"
	subnet_id = "subnet-96ufvn1s"
	access_type = 0
	public_network = 3
	security_group_ids = ["sg-xxxx"]
  }
`

const testAccCkafkaRouteSecurityGroupIdsUpdate = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_ckafka_route" "route_with_sg" {
	instance_id = var.instance_id
	vip_type = 3
	vpc_id = "vpc-hyuojwa7"
	subnet_id = "subnet-96ufvn1s"
	access_type = 0
	public_network = 3
	security_group_ids = ["sg-xxxx", "sg-yyyy"]
  }
`

const testAccCkafkaRouteSecurityGroupIdsRemove = tcacctest.DefaultKafkaVariable + `
resource "tencentcloud_ckafka_route" "route_with_sg" {
	instance_id = var.instance_id
	vip_type = 3
	vpc_id = "vpc-hyuojwa7"
	subnet_id = "subnet-96ufvn1s"
	access_type = 0
	public_network = 3
  }
`
