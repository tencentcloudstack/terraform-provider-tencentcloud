package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPrivatednsPrivateZoneVpcResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivatednsPrivateZoneVpc,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_privatedns_private_zone_vpc.private_zone_vpc", "id")),
			},
			{
				ResourceName:      "tencentcloud_privatedns_private_zone_vpc.private_zone_vpc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPrivatednsPrivateZoneVpc = `

resource "tencentcloud_privatedns_private_zone_vpc" "private_zone_vpc" {
  zone_id = "zone-xxxxxxx"
  vpc_set {
		uniq_vpc_id = "vpc-xadsafsdasd"
		region = "ap-guangzhou"

  }
  account_vpc_set {
		uniq_vpc_id = "vpc-xadsafsdasd"
		region = "ap-guangzhou"
		uin = "123456789"
		vpc_name = "testname"

  }
  tags = {
    "createdBy" = "terraform"
  }
}

`
