package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudCiMediaSpeechRecognitionTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaSpeechRecognitionTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_speech_recognition_template.media_speech_recognition_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_speech_recognition_template.media_speech_recognition_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaSpeechRecognitionTemplate = `

resource "tencentcloud_ci_media_speech_recognition_template" "media_speech_recognition_template" {
  name = &lt;nil&gt;
  speech_recognition {
		engine_model_type = &lt;nil&gt;
		channel_num = &lt;nil&gt;
		res_text_format = &lt;nil&gt;
		filter_dirty = &lt;nil&gt;
		filter_modal = &lt;nil&gt;
		convert_num_mode = &lt;nil&gt;
		speaker_diarization = &lt;nil&gt;
		speaker_number = &lt;nil&gt;
		filter_punc = &lt;nil&gt;
		output_file_type = &lt;nil&gt;

  }
}

`
