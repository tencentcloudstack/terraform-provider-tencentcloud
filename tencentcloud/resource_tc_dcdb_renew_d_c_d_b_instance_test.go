package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbRenewDCDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbRenewDCDBInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dcdb_renew_d_c_d_b_instance.renew_d_c_d_b_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_dcdb_renew_d_c_d_b_instance.renew_d_c_d_b_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDcdbRenewDCDBInstance = `

resource "tencentcloud_dcdb_renew_d_c_d_b_instance" "renew_d_c_d_b_instance" {
  instance_id = ""
  period = 
  auto_voucher = 
  voucher_ids = 
}

`
