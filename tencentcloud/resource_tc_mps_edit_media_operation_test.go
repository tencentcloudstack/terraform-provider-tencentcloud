package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsEditMediaOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsEditMediaOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_edit_media_operation.edit_media_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_edit_media_operation.edit_media_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsEditMediaOperation = `

resource "tencentcloud_mps_edit_media_operation" "edit_media_operation" {
  file_infos {
		input_info {
			type = ""
			cos_input_info {
				bucket = ""
				region = ""
				object = ""
			}
			url_input_info {
				url = ""
			}
			s3_input_info {
				s3_bucket = ""
				s3_region = ""
				s3_object = ""
				s3_secret_id = ""
				s3_secret_key = ""
			}
		}
		start_time_offset = 
		end_time_offset = 

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
  output_object_path = ""
  output_config {
		container = ""
		type = ""

  }
  task_notify_config {
		cmq_model = ""
		cmq_region = ""
		topic_name = ""
		queue_name = ""
		notify_mode = ""
		notify_type = ""
		notify_url = ""
		aws_s_q_s {
			s_q_s_region = ""
			s_q_s_queue_name = ""
			s3_secret_id = ""
			s3_secret_key = ""
		}

  }
  tasks_priority = 
  session_id = ""
  session_context = ""
}

`
