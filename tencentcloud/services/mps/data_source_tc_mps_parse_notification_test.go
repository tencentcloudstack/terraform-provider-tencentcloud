package mps_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsParseNotificationDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsParseNotificationDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mps_parse_notification.notification"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mps_parse_notification.notification", "content"),
				),
			},
		},
	})
}

const testAccMpsParseNotificationDataSource = `

data "tencentcloud_mps_parse_notification" "notification" {
	content = "{\"EventType\":\"WorkflowTask\",\"WorkflowTaskEvent\":{\"TaskId\":\"2600010949-WorkflowTask-764a7f3ca6f6eb13cda2048cfb034025tt7\",\"Status\":\"FINISH\",\"ErrCode\":0,\"Message\":\"SUCCESS\",\"InputInfo\":{\"Type\":\"COS\",\"CosInputInfo\":{\"Bucket\":\"keep-bucket-1308919341\",\"Region\":\"ap-guangzhou\",\"Object\":\"/mps-test/test.mov\"},\"UrlInputInfo\":null,\"S3InputInfo\":null},\"MetaData\":{\"AudioDuration\":30.528,\"AudioStreamSet\":[{\"Bitrate\":139635,\"Codec\":\"aac\",\"SamplingRate\":48000,\"Channel\":2,\"Codecs\":\"\",\"Loudness\":0,\"Duration\":30.528}],\"Bitrate\":185735,\"Container\":\"mov\",\"Duration\":30.571,\"Height\":270,\"Rotate\":0,\"Size\":709764,\"VideoDuration\":30.033,\"VideoStreamSet\":[{\"Bitrate\":37717,\"Codec\":\"h264\",\"Fps\":30,\"Height\":270,\"Width\":480,\"ColorPrimaries\":\"\",\"ColorSpace\":\"\",\"ColorTransfer\":\"\",\"HdrType\":\"\",\"Codecs\":\"\",\"Duration\":30.033}],\"Width\":480},\"MediaProcessResultSet\":[],\"AiRecognitionResultSet\":[],\"AiContentReviewResultSet\":[],\"AiAnalysisResultSet\":[],\"AiQualityControlTaskResult\":null},\"ScheduleTaskEvent\":null,\"LiveScheduleTaskEvent\":null,\"EditMediaTaskEvent\":null,\"GetMediaAttributesTaskEvent\":null,\"SessionId\":\"\",\"SessionContext\":\"\",\"ExtInfo\":\"\"}"
}

`
