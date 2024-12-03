package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterScimCredentialResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCenterScimCredential,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_scim_credential.identity_center_scim_credential", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_scim_credential.identity_center_scim_credential", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_scim_credential.identity_center_scim_credential", "credential_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_scim_credential.identity_center_scim_credential", "credential_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_scim_credential.identity_center_scim_credential", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_scim_credential.identity_center_scim_credential", "expire_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_scim_credential.identity_center_scim_credential", "credential_secret"),
				),
			},
			{
				ResourceName:            "tencentcloud_identity_center_scim_credential.identity_center_scim_credential",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"credential_secret"},
			},
		},
	})
}

const testAccIdentityCenterScimCredential = `
resource "tencentcloud_identity_center_scim_synchronization_status" "identity_center_scim_synchronization_status" {
  zone_id = "z-s64jh54hbcra"
  scim_synchronization_status = "Enabled"
}

resource "tencentcloud_identity_center_scim_credential" "identity_center_scim_credential" {
  zone_id = "z-s64jh54hbcra"
  depends_on = [tencentcloud_identity_center_scim_synchronization_status.identity_center_scim_synchronization_status]
}
`
