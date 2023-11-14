package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcAddressInternetChargeTypeModifyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcAddressInternetChargeTypeModify,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_address_internet_charge_type_modify.address_internet_charge_type_modify", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_address_internet_charge_type_modify.address_internet_charge_type_modify",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcAddressInternetChargeTypeModify = `

resource "tencentcloud_vpc_address_internet_charge_type_modify" "address_internet_charge_type_modify" {
  address_id = "eip-3456tghy"
  internet_charge_type = "BANDWIDTH_PREPAID_BY_MONTH"
  internet_max_bandwidth_out = 10
  address_charge_prepaid {
		period = 1
		auto_renew_flag = 1

  }
}

`
