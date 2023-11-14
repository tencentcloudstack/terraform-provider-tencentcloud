package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsAiAnalysisTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsAiAnalysisTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_ai_analysis_templates.ai_analysis_templates")),
			},
		},
	})
}

const testAccMpsAiAnalysisTemplatesDataSource = `

data "tencentcloud_mps_ai_analysis_templates" "ai_analysis_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  a_i_analysis_template_set {
		definition = &lt;nil&gt;
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
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;
		type = &lt;nil&gt;

  }
}

`
