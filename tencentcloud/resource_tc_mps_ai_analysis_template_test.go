package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMpsAiAnalysisTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsAiAnalysisTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_ai_analysis_template.ai_analysis_template", "id")),
			},
			{
				Config: testAccMpsAiAnalysisTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_ai_analysis_template.ai_analysis_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_ai_analysis_template.ai_analysis_template", "name", "terraform-for-test"),
				),
			},
			{
				ResourceName:      "tencentcloud_mps_ai_analysis_template.ai_analysis_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsAiAnalysisTemplate = `

resource "tencentcloud_mps_ai_analysis_template" "ai_analysis_template" {
  name = "terraform-test"

  classification_configure {
    switch = "OFF"
  }

  cover_configure {
    switch = "ON"
  }

  frame_tag_configure {
    switch = "ON"
  }

  tag_configure {
    switch = "ON"
  }
}


`

const testAccMpsAiAnalysisTemplateUpdate = `

resource "tencentcloud_mps_ai_analysis_template" "ai_analysis_template" {
  name = "terraform-for-test"

  classification_configure {
    switch = "OFF"
  }

  cover_configure {
    switch = "ON"
  }

  frame_tag_configure {
    switch = "ON"
  }

  tag_configure {
    switch = "ON"
  }
}


`
