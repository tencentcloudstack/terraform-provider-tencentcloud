package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCssAuthenticateDomainOwnerOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssAuthenticateDomainOwnerOperation_dnscheck,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_authenticate_domain_owner_operation.authenticate", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_authenticate_domain_owner_operation.authenticate", "verify_type", "dnsCheck"),
				),
			},
			{
				Config: testAccCssAuthenticateDomainOwnerOperation_filecheck,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_authenticate_domain_owner_operation.authenticate", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_authenticate_domain_owner_operation.authenticate", "verify_type", "fileCheck"),
				),
			},
		},
	})
}

const testAccCssAuthenticateDomainOwnerOperation_dnscheck = `

resource "tencentcloud_css_authenticate_domain_owner_operation" "authenticate" {
  domain_name = tencentcloud_css_domain.domain.id
  verify_type = "dnsCheck"
}

resource "tencentcloud_css_domain" "domain" {
	domain_name = "iac-tf.cloud"
	domain_type = 0
	play_type = 1
  }

`

const testAccCssAuthenticateDomainOwnerOperation_filecheck = `

resource "tencentcloud_css_authenticate_domain_owner_operation" "authenticate" {
  domain_name = tencentcloud_css_domain.domain.id
  verify_type = "fileCheck"
}

resource "tencentcloud_css_domain" "domain" {
	domain_name = "iac-tf.cloud"
	domain_type = 0
	play_type = 1
  }

`
