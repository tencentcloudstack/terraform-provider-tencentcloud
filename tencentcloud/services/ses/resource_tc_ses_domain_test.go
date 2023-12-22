package ses_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudSesDomain_basic -v
func TestAccTencentCloudSesDomain_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckBusiness(t, tcacctest.ACCOUNT_TYPE_SES) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_ses_domain.domain", "id"),
					resource.TestCheckResourceAttr("tencentcloud_ses_domain.domain", "email_identity", "iac.cloud"),
				),
			},
			{
				ResourceName:      "tencentcloud_ses_domain.domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSesDomain = `

resource "tencentcloud_ses_domain" "domain" {
  email_identity = "iac.cloud"
}

`
