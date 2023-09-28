package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsWithdrawsWatermarkOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMpsWithdrawsWatermarkOperation, defaultRegion),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_withdraws_watermark_operation.operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_withdraws_watermark_operation.operation", "input_info.#"),
					resource.TestCheckResourceAttr("tencentcloud_mps_withdraws_watermark_operation.operation", "input_info.0.type", "COS"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_withdraws_watermark_operation.operation", "input_info.0.cos_input_info.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_withdraws_watermark_operation.operation", "input_info.0.cos_input_info.0.bucket"),
					resource.TestCheckResourceAttr("tencentcloud_mps_withdraws_watermark_operation.operation", "input_info.0.cos_input_info.0.region", defaultRegion),
					resource.TestCheckResourceAttr("tencentcloud_mps_withdraws_watermark_operation.operation", "input_info.0.cos_input_info.0.object", "/mps-test/test.mov"),
					resource.TestCheckResourceAttr("tencentcloud_mps_withdraws_watermark_operation.operation", "session_context", "this is a example session context"),
				),
			},
		},
	})
}

const testAccMpsWithdrawsWatermarkOperation = userInfoData + `
// resource "tencentcloud_cos_bucket" "example" {
// 	bucket = "tf-test-mps-wm-${local.app_id}"
// 	acl    = "public-read"
//   }

//   resource "tencentcloud_cos_bucket_object" "example" {
// 	bucket  = tencentcloud_cos_bucket.example.bucket
// 	key     = "/test-file/test.mov"
// 	source  = "/Users/luoyin/Downloads/file_example_MOV_480_700kB.mov"
//   }

data "tencentcloud_cos_bucket_object" "object" {
	bucket = "keep-bucket-${local.app_id}"
	key    = "/mps-test/test.mov"
  }

resource "tencentcloud_mps_withdraws_watermark_operation" "operation" {
  input_info {
		type = "COS"
		cos_input_info {
			bucket = data.tencentcloud_cos_bucket_object.object.bucket
			region = "%s"
			object = data.tencentcloud_cos_bucket_object.object.key
		}
  }
//   task_notify_config {
// 		cmq_model = ""
// 		cmq_region = ""
// 		topic_name = ""
// 		queue_name = ""
// 		notify_mode = ""
// 		notify_type = "TDMQ-CMQ"
// 		notify_url = ""
// 		aws_sqs {
// 			sqs_region = ""
// 			sqs_queue_name = ""
// 			s3_secret_id = ""
// 			s3_secret_key = ""
// 		}
//   }
  session_context = "this is a example session context"
}

`
