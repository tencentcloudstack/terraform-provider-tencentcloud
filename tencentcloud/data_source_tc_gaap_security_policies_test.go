package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudGaapSecurityPolices_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityPolicesBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_policies.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_policies.foo", "proxy_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_policies.foo", "status"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_policies.foo", "action", "ACCEPT"),
				),
			},
		},
	})
}

var TestAccDataSourceTencentCloudGaapSecurityPolicesBasic = fmt.Sprintf(`
resource tencentcloud_gaap_security_policy "foo" {
  proxy_id = "%s"
  action   = "ACCEPT"
}

data tencentcloud_gaap_security_policies "foo" {
  id = tencentcloud_gaap_security_policy.foo.id
}
`, defaultGaapProxyId)
