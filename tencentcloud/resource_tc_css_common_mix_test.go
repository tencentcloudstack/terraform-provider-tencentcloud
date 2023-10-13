package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCssCommonMixResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssCommonMix,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_css_common_mix.common_mix", "id")),
			},
			{
				ResourceName:      "tencentcloud_css_common_mix.common_mix",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssCommonMix = `

resource "tencentcloud_css_common_mix" "common_mix" {
  mix_stream_session_id = ""
  input_stream_list {
		input_stream_name = ""
		layout_params {
			image_layer = 
			input_type = 
			image_height = 
			image_width = 
			location_x = 
			location_y = 
			color = ""
			watermark_id = 
		}
		crop_params {
			crop_width = 
			crop_height = 
			crop_start_location_x = 
			crop_start_location_y = 
		}

  }
  output_params {
		output_stream_name = ""
		output_stream_type = 
		output_stream_bit_rate = 
		output_stream_gop = 
		output_stream_frame_rate = 
		output_audio_bit_rate = 
		output_audio_sample_rate = 
		output_audio_channels = 
		mix_sei = ""

  }
  mix_stream_template_id = 
  control_params {
		use_mix_crop_center = 
		allow_copy = 
		pass_input_sei = 

  }
}

`
