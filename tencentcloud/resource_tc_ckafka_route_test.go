package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCkafkaRouteResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
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
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaRouteBasicWithIp,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ckafka_route.route_with_ip", "id")),
			},
		},
	})
}

const testAccCkafkaRouteBasic = defaultKafkaVariable + `
resource "tencentcloud_ckafka_route" "route" {
	instance_id = var.instance_id
	vip_type = 3
	vpc_id = "vpc-hyuojwa7"
	subnet_id = "subnet-96ufvn1s"
	access_type = 0
	public_network = 3
  }
`

const testAccCkafkaRouteBasicWithIp = defaultKafkaVariable + `
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
