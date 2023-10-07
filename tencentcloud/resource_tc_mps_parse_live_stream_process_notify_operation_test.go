package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixMpsParseLiveStreamProcessNotifyOperationResource_basic(t *testing.T) {
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
  content = "{\"EventType\":\"WorkflowTask\",\n    \"WorkflowTaskEvent\":{\n        \"TaskId\":\"245****654-WorkflowTask-f46dac7fe2436c47******d71946986t0\",\n        \"Status\":\"FINISH\",\n        \"ErrCode\":0,\n        \"Message\":\"\",\n        \"InputInfo\":{\n            \"Type\":\"COS\",\n            \"CosInputInfo\":{\n                \"Bucket\":\"macgzptest-125****654\",\n                \"Region\":\"ap-guangzhou\",\n                \"Object\":\"/dianping2.mp4\"\n            }\n        }\n    }\n}"
}

`
