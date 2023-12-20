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
