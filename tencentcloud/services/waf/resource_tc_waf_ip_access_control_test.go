package waf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafIpAccessControlResource_basic -v
func TestAccTencentCloudWafIpAccessControlResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafIpAccessControl,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_ip_access_control.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_ip_access_control.example", "instance_id", "waf_2kxtlbky00b3b4qz"),
					resource.TestCheckResourceAttr("tencentcloud_waf_ip_access_control.example", "domain", "keep.qcloudwaf.com"),
					resource.TestCheckResourceAttr("tencentcloud_waf_ip_access_control.example", "edition", "sparta-waf"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_ip_access_control.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccWafIpAccessControlUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_ip_access_control.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_waf_ip_access_control.example", "instance_id", "waf_2kxtlbky00b3b4qz"),
					resource.TestCheckResourceAttr("tencentcloud_waf_ip_access_control.example", "domain", "keep.qcloudwaf.com"),
					resource.TestCheckResourceAttr("tencentcloud_waf_ip_access_control.example", "edition", "sparta-waf"),
				),
			},
		},
	})
}

const testAccWafIpAccessControl = `
resource "tencentcloud_waf_ip_access_control" "example" {
  instance_id = "waf_2kxtlbky00b3b4qz"
  domain      = "keep.qcloudwaf.com"
  edition     = "sparta-waf"
  items {
    ip   = "1.1.1.1"
    note = "desc info."
    action = 42
    valid_ts = "2019571199"
  }

  items {
    ip   = "2.2.2.2"
    note = "desc info."
    action = 40
    valid_ts = "2019571199"
  }

  items {
    ip   = "3.3.3.3"
    note = "desc info."
    action = 42
    valid_ts = "1680570420"
  }
}
`

const testAccWafIpAccessControlUpdate = `
resource "tencentcloud_waf_ip_access_control" "example" {
  instance_id = "waf_2kxtlbky00b3b4qz"
  domain      = "keep.qcloudwaf.com"
  edition     = "sparta-waf"
  items {
    ip   = "1.1.1.1"
    note = "desc info."
    action = 42
    valid_ts = "2019571199"
  }

  items {
    ip   = "3.3.3.3"
    note = "desc info update."
    action = 42
    valid_ts = "2019571199"
  }

  items {
    ip   = "4.4.4.4"
    note = "desc info."
    action = 40
    valid_ts = "2019571199"
  }
}
`
