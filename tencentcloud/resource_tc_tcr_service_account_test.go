package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTcrServiceAccountResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTcrServiceAccount,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tcr_service_account.service_account", "id")),
			},
			{
				ResourceName:      "tencentcloud_tcr_service_account.service_account",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTcrServiceAccount = `

resource "tencentcloud_tcr_service_account" "service_account" {
  registry_id = "tcr-xxx"
  name = "robot"
  permissions {
		resource = "library"
		actions = 

  }
  description = "for namespace library"
  duration = 10
  expires_at = 1676897989000
  disable = false
  tags = {
    "createdBy" = "terraform"
  }
}

`
