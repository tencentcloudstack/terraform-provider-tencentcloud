package dnspod_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDnspodCustomLineResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodCustomLine,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dnspod_custom_line.custom_line", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_custom_line.custom_line", "name", "testline8"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_custom_line.custom_line", "area", "6.6.6.1-6.6.6.3"),
				),
			},
			{
				Config: testAccDnspodCustomLineUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_dnspod_custom_line.custom_line", "domain", "iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_custom_line.custom_line", "name", "testline9"),
					resource.TestCheckResourceAttr("tencentcloud_dnspod_custom_line.custom_line", "area", "6.6.6.1-6.6.6.2"),
				),
			},
			{
				ResourceName:      "tencentcloud_dnspod_custom_line.custom_line",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDnspodCustomLine = `

resource "tencentcloud_dnspod_custom_line" "custom_line" {
  domain = "iac-tf.cloud"
  name = "testline8"
  area = "6.6.6.1-6.6.6.3"
}

`

const testAccDnspodCustomLineUp = `

resource "tencentcloud_dnspod_custom_line" "custom_line" {
  domain = "iac-tf.cloud"
  name = "testline9"
  area = "6.6.6.1-6.6.6.2"
}

`
