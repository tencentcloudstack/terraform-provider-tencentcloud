package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveRecordTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveRecordTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_live_record_template.record_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_live_record_template.record_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLiveRecordTemplate = `

resource "tencentcloud_live_record_template" "record_template" {
  template_name = ""
  description = ""
  flv_param {
		record_interval = 
		storage_time = 
		enable = 
		vod_sub_app_id = 
		vod_file_name = ""
		procedure = ""
		storage_mode = ""
		class_id = 

  }
  hls_param {
		record_interval = 
		storage_time = 
		enable = 
		vod_sub_app_id = 
		vod_file_name = ""
		procedure = ""
		storage_mode = ""
		class_id = 

  }
  mp4_param {
		record_interval = 
		storage_time = 
		enable = 
		vod_sub_app_id = 
		vod_file_name = ""
		procedure = ""
		storage_mode = ""
		class_id = 

  }
  aac_param {
		record_interval = 
		storage_time = 
		enable = 
		vod_sub_app_id = 
		vod_file_name = ""
		procedure = ""
		storage_mode = ""
		class_id = 

  }
  is_delay_live = 
  hls_special_param {
		flow_continue_duration = 

  }
  mp3_param {
		record_interval = 
		storage_time = 
		enable = 
		vod_sub_app_id = 
		vod_file_name = ""
		procedure = ""
		storage_mode = ""
		class_id = 

  }
  remove_watermark = 
  flv_special_param {
		upload_in_recording = 

  }
}

`
