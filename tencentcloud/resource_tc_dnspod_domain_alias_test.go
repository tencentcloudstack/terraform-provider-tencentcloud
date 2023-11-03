package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodDomainAliasResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodDomainAlias,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dnspod_domain_alias.domain_alias", "id")),
			},
			{
				ResourceName:      "tencentcloud_dnspod_domain_alias.domain_alias",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDnspodDomainAlias = `

resource "tencentcloud_dnspod_domain_alias" "domain_alias" {
  domain_alias = "dnspod.com"
  domain = "dnspod.cn"
  domain_id = 123
  tags = {
    "createdBy" = "terraform"
  }
}

`
