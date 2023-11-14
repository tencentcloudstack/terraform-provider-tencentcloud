package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbRenewDBInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbRenewDBInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_renew_d_b_instance.renew_d_b_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_renew_d_b_instance.renew_d_b_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbRenewDBInstance = `

resource "tencentcloud_cdb_renew_d_b_instance" "renew_d_b_instance" {
  instance_id = ""
  time_span = 
  modify_pay_type = ""
}

`
