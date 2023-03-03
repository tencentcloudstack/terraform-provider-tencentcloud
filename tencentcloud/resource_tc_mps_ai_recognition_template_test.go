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
				Config: testAccMpsAiRecognitionTemplateUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_ai_recognition_template.ai_recognition_template", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_ai_recognition_template.ai_recognition_template", "name", "terraform-for-test"),
				),
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
  name = "terraform-test"

  asr_full_text_configure {
    switch = "OFF"
  }

  asr_words_configure {
    label_set = []
    switch    = "OFF"
  }

  face_configure {
    default_library_label_set     = [
      "entertainment",
      "sport",
    ]
    face_library                  = "All"
    score                         = 85
    switch                        = "ON"
    user_define_library_label_set = []
  }

  ocr_full_text_configure {
    switch = "OFF"
  }

  ocr_words_configure {
    label_set = []
    switch    = "OFF"
  }
}

`

const testAccMpsAiRecognitionTemplateUpdate = `

resource "tencentcloud_mps_ai_recognition_template" "ai_recognition_template" {
  name = "terraform-for-test"

  asr_full_text_configure {
    switch = "OFF"
  }

  asr_words_configure {
    label_set = []
    switch    = "OFF"
  }

  face_configure {
    default_library_label_set     = [
      "entertainment",
      "sport",
    ]
    face_library                  = "All"
    score                         = 85
    switch                        = "ON"
    user_define_library_label_set = []
  }

  ocr_full_text_configure {
    switch = "OFF"
  }

  ocr_words_configure {
    label_set = []
    switch    = "OFF"
  }
}

`
