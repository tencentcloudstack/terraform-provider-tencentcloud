package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafAntiFakeResource_basic -v
func TestAccTencentCloudWafAntiFakeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafAntiFake,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_anti_fake.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_anti_fake.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_anti_fake.example", "name"),
					resource.TestCheckResourceAttr("tencentcloud_waf_anti_fake.example", "uri", "/anti_fake_url.html"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_anti_fake.example", "status"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_anti_fake.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWafAntiFakeUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_anti_fake.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_anti_fake.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_anti_fake.example", "name"),
					resource.TestCheckResourceAttr("tencentcloud_waf_anti_fake.example", "uri", "/anti_fake_url_update.html"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_anti_fake.example", "status"),
				),
			},
		},
	})
}

const testAccWafAntiFake = `
resource "tencentcloud_waf_anti_fake" "example" {
  domain = "keep.qcloudwaf.com"
  name   = "tf_example"
  uri    = "/anti_fake_url.html"
  status = 0
}
`

const testAccWafAntiFakeUpdate = `
resource "tencentcloud_waf_anti_fake" "example" {
  domain = "keep.qcloudwaf.com"
  name   = "tf_example_update"
  uri    = "/anti_fake_url_update.html"
  status = 3
}
`
