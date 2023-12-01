package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixClbReplaceCertForLbsResource_basic(t *testing.T) {
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
		},
	})
}

var testAccClbReplaceCertForLbs = fmt.Sprintf(`

resource "tencentcloud_clb_replace_cert_for_lbs" "replace_cert_for_lbs" {
  old_certificate_id = "zjUMifFK"
  certificate {
    cert_ca_name = "test"
	cert_ca_content = "%s"
  }
}
`, testAccSslCertificateCA)
