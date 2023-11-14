package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveCommonMixStreamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveCommonMixStream,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_common_mix_stream.common_mix_stream", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_common_mix_stream.common_mix_stream",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveCommonMixStream = `

resource "tencentcloud_live_common_mix_stream" "common_mix_stream" {
  mix_stream_session_id = "test_room"
  input_stream_list {
		input_stream_name = "demo"
		layout_params {
			image_layer = 1
			input_type = 1
			image_height = 
			image_width = 
			location_x = 
			location_y = 
			color = "0xcc0033"
			watermark_id = 123456
		}
		crop_params {
			crop_width = 
			crop_height = 
			crop_start_location_x = 
			crop_start_location_y = 
		}

  }
  output_params {
		output_stream_name = "demo"
		output_stream_type = 1
		output_stream_bit_rate = 20
		output_stream_gop = 5
		output_stream_frame_rate = 30
		output_audio_bit_rate = 20
		output_audio_sample_rate = 96000
		output_audio_channels = 1
		mix_sei = "demo_sei"

  }
  mix_stream_template_id = 123456
  control_params {
		use_mix_crop_center = 1
		allow_copy = 1
		pass_input_sei = 1

  }
}

`
