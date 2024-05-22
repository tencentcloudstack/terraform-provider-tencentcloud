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
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain_dns.global_domain_dns", "domain_id", "dm-br5seuhh"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain_dns.global_domain_dns", "proxy_id_list.0", "link-m9t4yho9"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain_dns.global_domain_dns", "nation_country_inner_codes.0", "101001"),
				),
			},
			{
				Config: testAccGaapGlobalDomainDnsUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_global_domain_dns.global_domain_dns", "id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain_dns.global_domain_dns", "domain_id", "dm-br5seuhh"),
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
	domain_id = "dm-br5seuhh"
	proxy_id_list = ["link-m9t4yho9"]
	nation_country_inner_codes = ["101001"]
}
`

const testAccGaapGlobalDomainDnsUpdate = `
resource "tencentcloud_gaap_global_domain_dns" "global_domain_dns" {
	domain_id = "dm-br5seuhh"
	proxy_id_list = ["link-m9t4yho9", "link-2rk61jn5"]
	nation_country_inner_codes = ["101002"]
}
`
