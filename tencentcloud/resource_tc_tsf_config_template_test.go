package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudNeedFixTsfConfigTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfConfigTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_config_template.config_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_config_template.config_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfConfigTemplate = `

resource "tencentcloud_tsf_config_template" "config_template" {
  config_template_name = ""
  config_template_type = ""
  config_template_value = ""
  config_template_desc = ""
  program_id_list = 
      }

`
