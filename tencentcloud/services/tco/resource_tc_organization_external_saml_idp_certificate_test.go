package tco_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudOrganizationExternalSamlIdpCertificateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationExternalSamlIdpCertificate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_organization_external_saml_idp_certificate.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_external_saml_idp_certificate.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_organization_external_saml_idp_certificate.example", "x509_certificate"),
				),
			},
			{
				ResourceName:      "tencentcloud_organization_external_saml_idp_certificate.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccOrganizationExternalSamlIdpCertificate = `
resource "tencentcloud_organization_external_saml_idp_certificate" "example" {
  zone_id          = "z-dsj3ieme"
  x509_certificate = "MIIBtjCCAVugAwIBAgITBmyf1XSXNmY/Owua2eiedgPySjAKBggqhkj********"
}
`
