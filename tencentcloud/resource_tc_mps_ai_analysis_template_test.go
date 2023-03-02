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
				ResourceName:      "tencentcloud_mps_ai_analysis_template.ai_analysis_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsAiAnalysisTemplate = `

resource "tencentcloud_mps_ai_analysis_template" "ai_analysis_template" {
  name = &lt;nil&gt;
  comment = &lt;nil&gt;
  classification_configure {
		switch = &lt;nil&gt;

  }
  tag_configure {
		switch = &lt;nil&gt;

  }
  cover_configure {
		switch = &lt;nil&gt;

  }
  frame_tag_configure {
		switch = &lt;nil&gt;

  }
}

`
