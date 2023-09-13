package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixSesVerifyDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesVerifyDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ses_verify_domain.verify_domain", "id"),
				),
			},
		},
	})
}

const testAccSesVerifyDomain = `

resource "tencentcloud_ses_verify_domain" "verify_domain" {
  email_identity = "example.com"
}

`
