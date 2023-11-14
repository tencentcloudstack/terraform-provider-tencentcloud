package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsWithdrawsWatermarkResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsWithdrawsWatermark,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_withdraws_watermark.withdraws_watermark", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_withdraws_watermark.withdraws_watermark",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsWithdrawsWatermark = `

resource "tencentcloud_mps_withdraws_watermark" "withdraws_watermark" {
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
  session_context = ""
}

`
