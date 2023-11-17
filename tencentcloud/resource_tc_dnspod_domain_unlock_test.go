package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudNeedFixDnspodDomainUnlockResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodDomainUnlock,
				Check:  resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dnspod_domain_unlock.domain_unlock", "id"),
				),
			},
		},
	})
}

const testAccDnspodDomainUnlock = `

resource "tencentcloud_dnspod_domain_unlock" "domain_unlock" {
  domain = "iac-tf.cloud"
  lock_code = "MTAyMTIxOTd8aWFjLXRmLmNsb3VkfDE3MDAyMTEyMjF8YzE4NGEyNzI5ZDI4OGFjNxxxxx"
}

`
