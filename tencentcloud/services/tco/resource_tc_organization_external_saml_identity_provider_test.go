package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationExternalSamlIdentityProviderResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccOrganizationExternalSamlIdentityProvider,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_organization_external_saml_identity_provider.organization_external_saml_identity_provider", "id")),
		}, {
			ResourceName:      "tencentcloud_organization_external_saml_identity_provider.organization_external_saml_identity_provider",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccOrganizationExternalSamlIdentityProvider = `

resource "tencentcloud_organization_external_saml_identity_provider" "organization_external_saml_identity_provider" {
}
`
