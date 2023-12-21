package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapGlobalDomainDnsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapGlobalDomainDns,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_global_domain_dns.global_domain_dns", "id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain_dns.global_domain_dns", "domain_id", "dm-du60lmhj"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain_dns.global_domain_dns", "proxy_id_list.0", "link-m9p9fae3"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain_dns.global_domain_dns", "nation_country_inner_codes.0", "101001"),
				),
			},
			{
				Config: testAccGaapGlobalDomainDnsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_global_domain_dns.global_domain_dns", "id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain_dns.global_domain_dns", "domain_id", "dm-du60lmhj"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain_dns.global_domain_dns", "proxy_id_list.#", "2"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain_dns.global_domain_dns", "nation_country_inner_codes.0", "101002"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_global_domain_dns.global_domain_dns",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGaapGlobalDomainDns = `
resource "tencentcloud_gaap_global_domain_dns" "global_domain_dns" {
	domain_id = "dm-du60lmhj"
	proxy_id_list = ["link-m9p9fae3"]
	nation_country_inner_codes = ["101001"]
}
`

const testAccGaapGlobalDomainDnsUpdate = `
resource "tencentcloud_gaap_global_domain_dns" "global_domain_dns" {
	domain_id = "dm-du60lmhj"
	proxy_id_list = ["link-m9p9fae3", "link-8lpyo88p"]
	nation_country_inner_codes = ["101002"]
}
`
