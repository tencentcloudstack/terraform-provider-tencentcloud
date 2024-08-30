package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudIdentityCenterExternalSamlIdentityProviderResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityCenterExternalSamlIdentityProvider,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_external_saml_identity_provider.identity_center_external_saml_identity_provider", "id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_external_saml_identity_provider.identity_center_external_saml_identity_provider", "sso_status", "Enabled"),
				),
			},
			{
				Config: testAccIdentityCenterExternalSamlIdentityProviderUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_identity_center_external_saml_identity_provider.identity_center_external_saml_identity_provider", "id"),
					resource.TestCheckResourceAttr("tencentcloud_identity_center_external_saml_identity_provider.identity_center_external_saml_identity_provider", "sso_status", "Disabled"),
				),
			},
			{
				ResourceName:      "tencentcloud_identity_center_external_saml_identity_provider.identity_center_external_saml_identity_provider",
				ImportState:       true,
				ImportStateVerify: true,
			}},
	})
}

const testAccIdentityCenterExternalSamlIdentityProvider = `
resource "tencentcloud_identity_center_external_saml_identity_provider" "identity_center_external_saml_identity_provider" {
    zone_id = "z-s64jh54hbcra"
    sso_status = "Enabled"
}
`

const testAccIdentityCenterExternalSamlIdentityProviderUpdate = `
resource "tencentcloud_identity_center_external_saml_identity_provider" "identity_center_external_saml_identity_provider" {
    zone_id = "z-s64jh54hbcra"
    sso_status = "Disabled"
}
`
