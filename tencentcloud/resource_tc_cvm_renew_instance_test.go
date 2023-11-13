package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmRenewInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmRenewInstance,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_renew_instance.renew_instance", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_renew_instance.renew_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmRenewInstance = `

resource "tencentcloud_cvm_renew_instance" "renew_instance" {
  instance_ids = 
  instance_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_AUTO_RENEW"

  }
  renew_portable_data_disk = true
}

`
