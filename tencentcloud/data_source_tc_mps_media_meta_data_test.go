package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsMediaMetaDataDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsMediaMetaDataDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_media_meta_data.media_meta_data")),
			},
		},
	})
}

const testAccMpsMediaMetaDataDataSource = `

data "tencentcloud_mps_media_meta_data" "media_meta_data" {
  input_info {
		type = "COS"
		cos_input_info {
			bucket = "TopRankVideo-125xxx88"
			region = "ap-chongqing"
			object = "/movie/201907/WildAnimal.mov"
		}
		url_input_info {
			url = &lt;nil&gt;
		}

  }
  meta_data {
		size = &lt;nil&gt;
		container = &lt;nil&gt;
		bitrate = &lt;nil&gt;
		height = &lt;nil&gt;
		width = &lt;nil&gt;
		duration = 
		rotate = &lt;nil&gt;
		video_stream_set {
			bitrate = &lt;nil&gt;
			height = &lt;nil&gt;
			width = &lt;nil&gt;
			codec = &lt;nil&gt;
			fps = &lt;nil&gt;
			color_primaries = &lt;nil&gt;
			color_space = &lt;nil&gt;
			color_transfer = &lt;nil&gt;
			hdr_type = &lt;nil&gt;
		}
		audio_stream_set {
			bitrate = &lt;nil&gt;
			sampling_rate = &lt;nil&gt;
			codec = &lt;nil&gt;
			channel = &lt;nil&gt;
		}
		video_duration = &lt;nil&gt;
		audio_duration = &lt;nil&gt;

  }
}

`
