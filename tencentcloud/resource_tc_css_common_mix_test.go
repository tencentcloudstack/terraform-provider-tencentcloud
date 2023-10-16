package tencentcloud

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCssCommonMixResource_basic(t *testing.T) {
	t.Parallel()
	randIns := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := randIns.Intn(1000)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCssCommonMix, randomNum),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_common_mix.common_mix", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_common_mix.common_mix", "mix_stream_session_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_common_mix.common_mix", "input_stream_list.#"),
					resource.TestCheckResourceAttr("tencentcloud_css_common_mix.common_mix", "input_stream_list.0.input_stream_name", "test_stream1"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_common_mix.common_mix", "input_stream_list.0.layout_params.#"),
					resource.TestCheckResourceAttr("tencentcloud_css_common_mix.common_mix", "input_stream_list.0.layout_params.0.image_layer", "1"),
					resource.TestCheckResourceAttr("tencentcloud_css_common_mix.common_mix", "input_stream_list.1.input_stream_name", "test_stream2"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_common_mix.common_mix", "input_stream_list.1.layout_params.#"),
					resource.TestCheckResourceAttr("tencentcloud_css_common_mix.common_mix", "input_stream_list.1.layout_params.0.image_layer", "2"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_common_mix.common_mix", "output_params.#"),
					resource.TestCheckResourceAttr("tencentcloud_css_common_mix.common_mix", "output_params.0output_stream_name", "test_output_stream1"),
					resource.TestCheckResourceAttr("tencentcloud_css_common_mix.common_mix", "mix_stream_template_id", "30"),
				),
			},
		},
	})
}

const testAccCssCommonMix = `

resource "tencentcloud_css_common_mix" "common_mix" {
  mix_stream_session_id = "test_room_%d"
  input_stream_list {
		input_stream_name = "test_stream1"
		layout_params {
			image_layer = 1
		}
  }
  input_stream_list {
	input_stream_name = "test_stream2"
	layout_params {
		image_layer = 2
	}
}
  output_params {
		output_stream_name = "test_stream1"
  }
  mix_stream_template_id = 20
}

`
