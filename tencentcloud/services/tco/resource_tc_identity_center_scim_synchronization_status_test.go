package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterScimSynchronizationStatusResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCenterScimSynchronizationStatus,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_scim_synchronization_status.identity_center_scim_synchronization_status", "id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_scim_synchronization_status.identity_center_scim_synchronization_status", "scim_synchronization_status", "Enabled"),
				),
			},
			{
				Config: testAccIdentityCenterScimSynchronizationStatusUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_scim_synchronization_status.identity_center_scim_synchronization_status", "id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_scim_synchronization_status.identity_center_scim_synchronization_status", "scim_synchronization_status", "Disabled"),
				),
			},
			{
				ResourceName:      "tencentcloud_identity_center_scim_synchronization_status.identity_center_scim_synchronization_status",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccIdentityCenterScimSynchronizationStatus = `
resource "tencentcloud_identity_center_scim_synchronization_status" "identity_center_scim_synchronization_status" {
  zone_id = "z-s64jh54hbcra"
  scim_synchronization_status = "Enabled"
}
`

const testAccIdentityCenterScimSynchronizationStatusUpdate = `
resource "tencentcloud_identity_center_scim_synchronization_status" "identity_center_scim_synchronization_status" {
  zone_id = "z-s64jh54hbcra"
  scim_synchronization_status = "Disabled"
}
`
