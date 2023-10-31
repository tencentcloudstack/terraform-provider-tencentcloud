package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsParseLiveStreamProcessNotificationDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsParseLiveStreamProcessNotificationDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_parse_live_stream_process_notification.notification"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_parse_live_stream_process_notification.notification", "content"),
				),
			},
		},
	})
}

const testAccMpsParseLiveStreamProcessNotificationDataSource = `

data "tencentcloud_mps_parse_live_stream_process_notification" "notification" {
	content = "{\"NotificationType\":\"ProcessEof\",\"TaskId\":\"2600010949-procedure-live-48a2680775c4d73651ca894aaa91052ctt7\",\"ProcessEofInfo\":{\"ErrCode\":0,\"Message\":\"Success\"},\"AiReviewResultInfo\":null,\"AiAnalysisResultInfo\":null,\"AiRecognitionResultInfo\":null,\"AiQualityControlResultInfo\":null,\"SessionId\":\"\",\"SessionContext\":\"\"}"
}

`
