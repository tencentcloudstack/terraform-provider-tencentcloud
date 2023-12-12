Provides a resource to create a mps edit_media_operation

Example Usage

Operation through COS

```hcl
resource "tencentcloud_cos_bucket" "output" {
	bucket = "tf-bucket-mps-output-${local.app_id}"
  }

data "tencentcloud_cos_bucket_object" "object" {
	bucket = "keep-bucket-${local.app_id}"
	key    = "/mps-test/test.mov"
}

resource "tencentcloud_mps_edit_media_operation" "operation" {
  file_infos {
		input_info {
			type = "COS"
			cos_input_info {
				bucket = data.tencentcloud_cos_bucket_object.object.bucket
				region = "%s"
				object = data.tencentcloud_cos_bucket_object.object.key
			}
		}
		start_time_offset = 60
		end_time_offset = 120
  }
  output_storage {
		type = "COS"
		cos_output_storage {
			bucket = tencentcloud_cos_bucket.output.bucket
			region = "%s"
		}
  }
  output_object_path = "/output"
}
```