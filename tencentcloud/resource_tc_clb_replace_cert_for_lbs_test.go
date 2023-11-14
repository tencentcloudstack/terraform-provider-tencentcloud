package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudClbReplaceCertForLbsResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClbReplaceCertForLbs,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_clb_replace_cert_for_lbs.replace_cert_for_lbs", "id")),
			},
			{
				ResourceName:      "tencentcloud_clb_replace_cert_for_lbs.replace_cert_for_lbs",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClbReplaceCertForLbs = `

resource "tencentcloud_clb_replace_cert_for_lbs" "replace_cert_for_lbs" {
  old_certificate_id = "xxxxxxxx"
  certificate {
		s_s_l_mode = "UNIDIRECTIONAL"
		cert_id = "xxxxxxxx"
		cert_ca_id = "xxxxxxxx"
		cert_name = "test"
		cert_key = "xxxxxxxxxxxxxxxx"
		cert_content = ""
		cert_ca_name = "test"
		cert_ca_content = ""

  }
}

`
