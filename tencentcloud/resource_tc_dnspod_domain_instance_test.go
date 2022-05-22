package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDnspodDoamin(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDnspodDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dnspod_domain_instance.domain", "domain", "terraformer.com"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_domain_instance.domain", "is_mark", "no"),
				),
			},
		},
	})
}

const testAccTencentCloudDnspodDomain = `
resource "tencentcloud_dnspod_domain_instance" "domain" {
  domain      = "terraformer.com"
  is_mark     = "no"
}
`
