package mps_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsProcessMediaOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMpsProcessMediaOperation, tcacctest.DefaultRegion, tcacctest.DefaultRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_media_operation.operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_media_operation.operation", "input_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_process_media_operation.operation", "input_info.0.type", "COS"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_media_operation.operation", "input_info.0.cos_input_info.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_media_operation.operation", "input_info.0.cos_input_info.0.bucket"),
					resource.TestCheckResourceAttr("tencentcloud_mps_process_media_operation.operation", "input_info.0.cos_input_info.0.region", tcacctest.DefaultRegion),
					resource.TestCheckResourceAttr("tencentcloud_mps_process_media_operation.operation", "input_info.0.cos_input_info.0.object", "/mps-test/test.mov"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_media_operation.operation", "output_storage.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_process_media_operation.operation", "output_storage.0.type", "COS"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_media_operation.operation", "output_storage.0.cos_output_storage.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_media_operation.operation", "output_storage.0.cos_output_storage.0.bucket"),
					resource.TestCheckResourceAttr("tencentcloud_mps_process_media_operation.operation", "output_storage.0.cos_output_storage.0.region", tcacctest.DefaultRegion),
				),
			},
		},
	})
}

const testAccMpsProcessMediaOperation = tcacctest.UserInfoData + `
resource "tencentcloud_cos_bucket" "output" {
  bucket      = "tf-bucket-mps-edit-media-output-${local.app_id}"
  force_clean = true
  acl         = "public-read"
}

data "tencentcloud_cos_bucket_object" "object" {
  bucket = "keep-bucket-${local.app_id}"
  key    = "/mps-test/test.mov"
}

resource "tencentcloud_mps_process_media_operation" "operation" {
  input_info {
    type = "COS"
    cos_input_info {
      bucket = data.tencentcloud_cos_bucket_object.object.bucket
      region = "%s"
      object = data.tencentcloud_cos_bucket_object.object.key
    }
  }
  output_storage {
    type = "COS"
    cos_output_storage {
      bucket = tencentcloud_cos_bucket.output.bucket
      region = "%s"
    }
  }
  output_dir = "output/"

  ai_content_review_task {
    definition = 10
  }

  ai_recognition_task {
    definition = 10
  }

  task_notify_config {
    cmq_model   = "Queue"
    cmq_region  = "gz"
    queue_name  = "test"
    topic_name  = "test"
    notify_type = "CMQ"
  }
}


`
