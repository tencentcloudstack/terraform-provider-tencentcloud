package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudCssPlayDomainCertAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssPlayDomainCertAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_css_play_domain_cert_attachment.play_domain_cert_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_css_play_domain_cert_attachment.play_domain_cert_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssPlayDomainCertAttachment = `

resource "tencentcloud_css_play_domain_cert_attachment" "play_domain_cert_attachment" {
  cloud_cert_id = &lt;nil&gt;
  domain_name = ""
  status = 
  }

`
