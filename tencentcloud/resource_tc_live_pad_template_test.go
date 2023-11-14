package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLivePadTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLivePadTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_pad_template.pad_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_pad_template.pad_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLivePadTemplate = `

resource "tencentcloud_live_pad_template" "pad_template" {
  template_name = ""
  url = ""
  description = ""
  wait_duration = 
  max_duration = 
  type = 
}

`
