package tencentcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDomainsDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDomainsDataSourceBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_domains.domains", "list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_domains.domains", "list.0.auto_renew"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_domains.domains", "list.0.is_premium"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_domains.domains", "list.0.domain_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_domains.domains", "list.0.expiration_date"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_domains.domains", "list.0.domain_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_domains.domains", "list.0.code_tld"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_domains.domains", "list.0.creation_date"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_domains.domains", "list.0.tld"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_domains.domains", "list.0.buy_status"),
					resource.TestMatchOutput("domain", regexp.MustCompile(`\\w+\\.\\w+`)),
				),
			},
		},
	})
}

const testAccDomainsDataSourceBasic = `
data "tencentcloud_domains" "domains" {}

locals {
  domain1 = data.tencentcloud_domains.domains.list.0.domain_name
}

output "domain" {
  value = local.domain1
}
`
