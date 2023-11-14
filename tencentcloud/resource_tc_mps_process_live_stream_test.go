package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsProcessLiveStreamResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsProcessLiveStream,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_process_live_stream.process_live_stream", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_process_live_stream.process_live_stream",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsProcessLiveStream = `

resource "tencentcloud_mps_process_live_stream" "process_live_stream" {
  url = ""
  task_notify_config {
		cmq_model = ""
		cmq_region = ""
		queue_name = ""
		topic_name = ""
		notify_type = ""
		notify_url = ""

  }
  output_storage {
		type = ""
		cos_output_storage {
			bucket = ""
			region = ""
		}
		s3_output_storage {
			s3_bucket = ""
			s3_region = ""
			s3_secret_id = ""
			s3_secret_key = ""
		}

  }
  output_dir = ""
  ai_content_review_task {
		definition = 

  }
  ai_recognition_task {
		definition = 

  }
  ai_analysis_task {
		definition = 
		extended_parameter = ""

  }
  ai_quality_control_task {
		definition = 
		channel_ext_para = ""

  }
  session_id = ""
  session_context = ""
  schedule_id = 
}

`
