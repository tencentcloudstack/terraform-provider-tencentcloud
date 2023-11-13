package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbRenewInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbRenewInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mariadb_renew_instance.renew_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_mariadb_renew_instance.renew_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMariadbRenewInstance = `

resource "tencentcloud_mariadb_renew_instance" "renew_instance" {
  instance_id = ""
  period = 
  auto_voucher = 
  voucher_ids = 
}

`
