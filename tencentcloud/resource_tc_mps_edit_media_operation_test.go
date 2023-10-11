package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Config: fmt.Sprintf(testAccMpsEditMediaOperation, defaultRegion, defaultRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_edit_media_operation.operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_edit_media_operation.operation", "file_infos.0.input_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_edit_media_operation.operation", "file_infos.0.input_info.0.type", "COS"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_edit_media_operation.operation", "file_infos.0.input_info.0.cos_input_info.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_edit_media_operation.operation", "file_infos.0.input_info.0.cos_input_info.0.bucket"),
					resource.TestCheckResourceAttr("tencentcloud_mps_edit_media_operation.operation", "file_infos.0.input_info.0.cos_input_info.0.region", defaultRegion),
					resource.TestCheckResourceAttr("tencentcloud_mps_edit_media_operation.operation", "file_infos.0.input_info.0.cos_input_info.0.object", "/mps-test/test.mov"),
				),
			},
		},
	})
}

const testAccMpsEditMediaOperation = userInfoData + `
resource "tencentcloud_cos_bucket" "output" {
	bucket      = "tf-bucket-mps-edit-media-output-${local.app_id}"
	force_clean = true
	acl         = "public-read"
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
		end_time_offset   = 120
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

`
