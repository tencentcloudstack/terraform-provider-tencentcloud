package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoBindSecurityTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoBindSecurityTemplate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_teo_bind_security_template.teo_bind_security_template", "zone_id", "zone-39quuimqg8r6"),
					resource.TestCheckResourceAttr("tencentcloud_teo_bind_security_template.teo_bind_security_template", "template_id", "temp-7dr7dm78"),
					resource.TestCheckResourceAttr("tencentcloud_teo_bind_security_template.teo_bind_security_template", "entity", "aaa.makn.cn"),
					resource.TestCheckResourceAttr("tencentcloud_teo_bind_security_template.teo_bind_security_template", "operate", "unbind-use-default"),
					resource.TestCheckResourceAttr("tencentcloud_teo_bind_security_template.teo_bind_security_template", "status", "online"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_bind_security_template.teo_bind_security_template",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"operate",
				},
			},
		},
	})
}

const testAccTeoBindSecurityTemplate = `

resource "tencentcloud_teo_bind_security_template" "teo_bind_security_template" {
  operate     = "unbind-use-default"
  template_id = "temp-7dr7dm78"
  zone_id     = "zone-39quuimqg8r6"
  entity 	  = "aaa.makn.cn"
  over_write  = false
}

`
