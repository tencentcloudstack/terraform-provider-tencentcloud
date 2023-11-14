package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmRenewHostsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmRenewHosts,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_renew_hosts.renew_hosts", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_renew_hosts.renew_hosts",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmRenewHosts = `

resource "tencentcloud_cvm_renew_hosts" "renew_hosts" {
  host_ids = 
  host_charge_prepaid {
		period = 1
		renew_flag = "NOTIFY_AND_MANUAL_RENEW"

  }
}

`
