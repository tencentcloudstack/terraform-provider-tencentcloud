package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseModifyDiskRenewFlagResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseModifyDiskRenewFlag,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_modify_disk_renew_flag.modify_disk_renew_flag", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_modify_disk_renew_flag.modify_disk_renew_flag",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseModifyDiskRenewFlag = `

resource "tencentcloud_lighthouse_modify_disk_renew_flag" "modify_disk_renew_flag" {
  disk_ids = 
  renew_flag = ""
}

`
