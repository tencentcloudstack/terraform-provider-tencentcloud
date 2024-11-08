package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterScimCredentialStatusResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCenterScimCredentialStatus,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_scim_credential_status.identity_center_scim_credential_status", "id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_scim_credential_status.identity_center_scim_credential_status", "status", "Enabled"),
				),
			},
			{
				Config: testAccIdentityCenterScimCredentialStatusUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_scim_credential_status.identity_center_scim_credential_status", "id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_scim_credential_status.identity_center_scim_credential_status", "status", "Disabled"),
				),
			},
			{
				ResourceName:      "tencentcloud_identity_center_scim_credential_status.identity_center_scim_credential_status",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIdentityCenterScimCredentialStatus = `
resource "tencentcloud_identity_center_scim_synchronization_status" "identity_center_scim_synchronization_status" {
  zone_id = "z-s64jh54hbcra"
  scim_synchronization_status = "Enabled"
}

resource "tencentcloud_identity_center_scim_credential" "identity_center_scim_credential" {
  zone_id = "z-s64jh54hbcra"
  depends_on = [tencentcloud_identity_center_scim_synchronization_status.identity_center_scim_synchronization_status]
}

resource "tencentcloud_identity_center_scim_credential_status" "identity_center_scim_credential_status" {
  zone_id = "z-s64jh54hbcra"
  credential_id = tencentcloud_identity_center_scim_credential.identity_center_scim_credential.credential_id
  status = "Enabled"
}
`

const testAccIdentityCenterScimCredentialStatusUpdate = `
resource "tencentcloud_identity_center_scim_synchronization_status" "identity_center_scim_synchronization_status" {
  zone_id = "z-s64jh54hbcra"
  scim_synchronization_status = "Enabled"
}

resource "tencentcloud_identity_center_scim_credential" "identity_center_scim_credential" {
  zone_id = "z-s64jh54hbcra"
  depends_on = [tencentcloud_identity_center_scim_synchronization_status.identity_center_scim_synchronization_status]
}

resource "tencentcloud_identity_center_scim_credential_status" "identity_center_scim_credential_status" {
  zone_id = "z-s64jh54hbcra"
  credential_id = tencentcloud_identity_center_scim_credential.identity_center_scim_credential.credential_id
  status = "Disabled"
}
`
