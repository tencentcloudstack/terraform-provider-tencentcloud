package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsParseLiveStreamProcessNotifyOperationResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				// Config: fmt.Sprintf(testAccMpsParseLiveStreamProcessNotifyOperation,content),
				Config: testAccMpsParseLiveStreamProcessNotifyOperation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_mps_parse_live_stream_process_notify_operation.operation", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_mps_parse_live_stream_process_notify_operation.operation", "content"),
				),
			},
		},
	})
}

const testAccMpsParseLiveStreamProcessNotifyOperation = `

resource "tencentcloud_mps_parse_live_stream_process_notify_operation" "operation" {
    content = "{\"NotificationType\":\"ProcessEof\",\"TaskId\":\"2600010949-procedure-live-48a2680775c4d73651ca894aaa91052ctt7\",\"ProcessEofInfo\":{\"ErrCode\":0,\"Message\":\"Success\"},\"AiReviewResultInfo\":null,\"AiAnalysisResultInfo\":null,\"AiRecognitionResultInfo\":null,\"AiQualityControlResultInfo\":null,\"SessionId\":\"\",\"SessionContext\":\"\"}"
}

`
