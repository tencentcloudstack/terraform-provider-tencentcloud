package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsAiRecognitionTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsAiRecognitionTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_ai_recognition_templates.ai_recognition_templates")),
			},
		},
	})
}

const testAccMpsAiRecognitionTemplatesDataSource = `

data "tencentcloud_mps_ai_recognition_templates" "ai_recognition_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  a_i_recognition_template_set {
		definition = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		face_configure {
			switch = &lt;nil&gt;
			score = 
			default_library_label_set = &lt;nil&gt;
			user_define_library_label_set = &lt;nil&gt;
			face_library = "All"
		}
		ocr_full_text_configure {
			switch = &lt;nil&gt;
		}
		ocr_words_configure {
			switch = &lt;nil&gt;
			label_set = &lt;nil&gt;
		}
		asr_full_text_configure {
			switch = &lt;nil&gt;
			subtitle_format = &lt;nil&gt;
		}
		asr_words_configure {
			switch = &lt;nil&gt;
			label_set = &lt;nil&gt;
		}
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;
		type = &lt;nil&gt;

  }
}

`
