package mps_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsProcessLiveStreamOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMpsProcessLiveStreamOperation, tcacctest.DefaultRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_live_stream_operation.operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_mps_process_live_stream_operation.operation", "url", "http://www.abc.com/abc.m3u8"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_live_stream_operation.operation", "task_notify_config.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_live_stream_operation.operation", "output_storage.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_process_live_stream_operation.operation", "output_storage.0.type", "COS"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_live_stream_operation.operation", "output_storage.0.cos_output_storage.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_live_stream_operation.operation", "output_storage.0.cos_output_storage.0.bucket"),
					resource.TestCheckResourceAttr("tencentcloud_mps_process_live_stream_operation.operation", "output_storage.0.cos_output_storage.0.region", tcacctest.DefaultRegion),
					resource.TestCheckResourceAttr("tencentcloud_mps_process_live_stream_operation.operation", "output_dir", "/output/"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_live_stream_operation.operation", "ai_content_review_task.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_process_live_stream_operation.operation", "ai_content_review_task.0.definition", "10"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_process_live_stream_operation.operation", "ai_recognition_task.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_process_live_stream_operation.operation", "ai_recognition_task.0.definition", "10"),
				),
			},
		},
	})
}

const testAccMpsProcessLiveStreamOperation = tcacctest.UserInfoData + `
resource "tencentcloud_cos_bucket" "output" {
  bucket      = "tf-bucket-mps-process-live-stream-output-${local.app_id}"
  force_clean = true
  acl         = "public-read"
}

resource "tencentcloud_mps_process_live_stream_operation" "operation" {
  url = "http://www.abc.com/abc.m3u8"
  task_notify_config {
    cmq_model   = "Queue"
    cmq_region  = "gz"
    queue_name  = "test"
    topic_name  = "test"
    notify_type = "CMQ"
  }

  output_storage {
    type = "COS"
    cos_output_storage {
      bucket = tencentcloud_cos_bucket.output.bucket
      region = "%s"
    }
  }

  output_dir = "/output/"

  ai_content_review_task {
    definition = 10
  }

  ai_recognition_task {
    definition = 10
  }
}


`
