package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapGlobalDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapGlobalDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_global_domain.global_domain", "id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain.global_domain", "alias", "demo"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain.global_domain", "default_value", "mikatong.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain.global_domain", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain.global_domain", "status", "open"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain.global_domain", "tags.key", "value"),
				),
			},
			{
				Config: testAccGaapGlobalDomainUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_gaap_global_domain.global_domain", "id"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain.global_domain", "alias", "demo1"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain.global_domain", "default_value", "mikatong1.com"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain.global_domain", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain.global_domain", "status", "close"),
					resource.TestCheckResourceAttr("tencentcloud_gaap_global_domain.global_domain", "tags.key", "value1"),
				),
			},
			{
				ResourceName:      "tencentcloud_gaap_global_domain.global_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccGaapGlobalDomain = `
resource "tencentcloud_gaap_global_domain" "global_domain" {
	project_id = 0
	default_value = "mikatong.com"
	alias = "demo"
	tags={
		key = "value"
	}
}
`

const testAccGaapGlobalDomainUpdate = `
resource "tencentcloud_gaap_global_domain" "global_domain" {
	project_id = 0
	default_value = "mikatong1.com"
	alias = "demo1"
	tags={
		key = "value1"
	}
	status = "close"
}
`
