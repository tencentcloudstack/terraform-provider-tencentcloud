package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTseCngwGatewayResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTseCngwGateway,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tse_cngw_gateway.cngw_gateway", "id")),
			},
			{
				ResourceName:      "tencentcloud_tse_cngw_gateway.cngw_gateway",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTseCngwGateway = `

resource "tencentcloud_tse_cngw_gateway" "cngw_gateway" {
  name = "test"
  type = "kong"
  gateway_version = "2.4.1"
  node_config {
		specification = ""
		number = 

  }
  vpc_config {
		vpc_id = "vpc-xxxxxx"
		subnet_id = "subnet-xxxxxx"

  }
  description = "for test"
  enable_cls = false
  feature_version = ""
  internet_max_bandwidth_out = 
  engine_region = "ap-guangzhou"
  ingress_class_name = ""
  trade_type = 
  internet_config {
		internet_address_version = ""
		internet_pay_mode = ""
		internet_max_bandwidth_out = 
		description = ""
		sla_type = ""
		multi_zone_flag = 
		master_zone_id = ""
		slave_zone_id = ""

  }
  tags = {
    "createdBy" = "terraform"
  }
}

`
