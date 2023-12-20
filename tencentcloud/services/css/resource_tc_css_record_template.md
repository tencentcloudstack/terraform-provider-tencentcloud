Provides a resource to create a css record_template

Example Usage

```hcl
resource "tencentcloud_css_record_template" "record_template" {
  template_name = "demo"
  description = "this is demo"
  flv_param {
		record_interval = 30
		storage_time = 2
		enable = 0
		vod_sub_app_id = 123
		vod_file_name = "demo"
		procedure = ""
		storage_mode = "Normal"
		class_id = 123

  }
  hls_param {
		record_interval = 40
		storage_time = 2
		enable = 1
		vod_sub_app_id = 123
		vod_file_name = "test"
		procedure = ""
		storage_mode = "Cold"
		class_id = 123

  }
  mp4_param {
		record_interval = 45
		storage_time = 56
		enable = 0
		vod_sub_app_id = 234
		vod_file_name = "test"
		procedure = ""
		storage_mode = "Cold"
		class_id = 123

  }
  aac_param {
		record_interval = 5678
		storage_time = 1234
		enable = 1
		vod_sub_app_id = 123
		vod_file_name = "test"
		procedure = ""
		storage_mode = "Normal"
		class_id = 123

  }
  is_delay_live = 1
  hls_special_param {
		flow_continue_duration = 1200

  }
  mp3_param {
		record_interval = 100
		storage_time = 100
		enable = 1
		vod_sub_app_id = 123
		vod_file_name = "test"
		procedure = ""
		storage_mode = "Normal"
		class_id = 123

  }
  remove_watermark = true
  flv_special_param {
		upload_in_recording = true

  }
}
```

Import

css record_template can be imported using the id, e.g.

```
terraform import tencentcloud_css_record_template.record_template templateId
```
