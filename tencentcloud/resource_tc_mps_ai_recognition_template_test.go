package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMpsAiRecognitionTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsAiRecognitionTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_ai_recognition_template.ai_recognition_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_ai_recognition_template.ai_recognition_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsAiRecognitionTemplate = `

resource "tencentcloud_mps_ai_recognition_template" "ai_recognition_template" {
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
}

`
