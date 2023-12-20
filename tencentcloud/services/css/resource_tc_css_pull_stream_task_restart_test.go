package css_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudCssPullStreamTaskRestartResource_basic -v
func TestAccTencentCloudNeedFixCssPullStreamTaskRestartResource_basic(t *testing.T) {
	t.Parallel()
	baseTime := time.Now().UTC().Add(10 * time.Hour)
	startTime := baseTime.Format(time.RFC3339)
	endTime := baseTime.Add(1 * time.Hour).Format(time.RFC3339)
	liveUrl := "rtmp://5000.liveplay.myqcloud.com/live/stream1"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCssRestartPushTask, tcacctest.DefaultCSSLiveType, liveUrl, tcacctest.DefaultCSSDomainName, tcacctest.DefaultCSSAppName, tcacctest.DefaultCSSStreamName, startTime, endTime, tcacctest.DefaultCSSOperator),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task_restart.restart_push_task", "id"),
				),
			},
		},
	})
}

const testAccCssRestartPushTask = testAccCssPullStreamTask + `

resource "tencentcloud_css_pull_stream_task_restart" "restart_push_task" {
  task_id = tencentcloud_css_pull_stream_task.pull_stream_task.id
  operator = "tf-test"
}

`
