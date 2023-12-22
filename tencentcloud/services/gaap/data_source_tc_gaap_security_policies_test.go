package gaap_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudGaapSecurityPolices_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudGaapSecurityPolicesBasic,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_security_policies.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_policies.foo", "proxy_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_security_policies.foo", "status"),
					resource.TestCheckResourceAttr("data.tencentcloud_gaap_security_policies.foo", "action", "ACCEPT"),
				),
			},
		},
	})
}

var TestAccDataSourceTencentCloudGaapSecurityPolicesBasic = fmt.Sprintf(`

data tencentcloud_gaap_security_policies "foo" {
  id = "%s"
}
`, tcacctest.DefaultGaapSecurityPolicyId)
