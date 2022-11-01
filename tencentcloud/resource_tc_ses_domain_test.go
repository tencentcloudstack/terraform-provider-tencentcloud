package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudSesDomain_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps:     []resource.TestStep{
			//{
			//	Config: testAccSesDomain,
			//	Check: resource.ComposeTestCheckFunc(
			//		resource.TestCheckResourceAttrSet("tencentcloud_ses_domain.domain", "id"),
			//	),
			//},
			//{
			//	ResourceName:      "tencentcloud_ses_domain.domain",
			//	ImportState:       true,
			//	ImportStateVerify: true,
			//},
		},
	})
}

const testAccSesDomain = `

resource "tencentcloud_ses_domain" "domain" {
  email_identity = ""
}

`
