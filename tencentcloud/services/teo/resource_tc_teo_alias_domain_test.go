package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoAliasDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoAliasDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_alias_domain.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_alias_domain.example", "zone_id", "zone-3fkff38fyw8s"),
					resource.TestCheckResourceAttr("tencentcloud_teo_alias_domain.example", "alias_name", "test1.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_alias_domain.example", "target_name", "yyyy.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_alias_domain.example", "cert_type", "none"),
					resource.TestCheckResourceAttr("tencentcloud_teo_alias_domain.example", "paused", "false"),
				),
			},
			{
				Config: testAccTeoAliasDomainUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_alias_domain.example", "target_name", "yyyy.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_alias_domain.example", "paused", "true"),
				),
			},
		},
	})
}

const testAccTeoAliasDomain = `
resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "test1.makn.cn"
  target_name = "yyyy.makn.cn"
  cert_type   = "none"
  paused      = false
}
`

const testAccTeoAliasDomainUpdate = `
resource "tencentcloud_teo_alias_domain" "example" {
  zone_id     = "zone-3fkff38fyw8s"
  alias_name  = "test1.makn.cn"
  target_name = "yyyy.makn.cn"
  cert_type   = "none"
  paused      = true
}
`
