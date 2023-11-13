package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesDomain,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ses_domain.domain", "id")),
			},
			{
				ResourceName:      "tencentcloud_ses_domain.domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSesDomain = `

resource "tencentcloud_ses_domain" "domain" {
  email_identity = "mail.qcloud.com"
}

`
