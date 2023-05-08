package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func init() {
	resource.AddTestSweepers("tencentcloud_css_pull_stream_task", &resource.Sweeper{
		Name: "tencentcloud_css_pull_stream_task",
		F:    testSweepCSSPullStreamTask,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_css_pull_stream_task
func testSweepCSSPullStreamTask(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	cssService := CssService{client: cli.(*TencentCloudClient).apiV3Conn}

	info, err := cssService.DescribeCssPullStreamTask(ctx, "")
	if err != nil {
		return err
	}
	if info == nil {
		return fmt.Errorf("instance not exists.")
	}

	for _, v := range info {
		delName := v.StreamName
		delId := v.TaskId

		if strings.HasPrefix(*delName, defaultCSSPrefix) {
			err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
				err := cssService.DeleteCssPullStreamTaskById(ctx, delId, helper.String(defaultCSSOperator))
				if err != nil {
					return retryError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[ERROR] sweeper instance %s:%s failed! reason:[%s]", *delName, *delId, err.Error())
			}
		}
	}
	return nil
}

func TestAccTencentCloudCssPullStreamTaskResource_basic(t *testing.T) {
	t.Parallel()
	baseTime := time.Now().UTC().Add(10 * time.Hour)
	startTime := baseTime.Format(time.RFC3339)
	endTime := baseTime.Add(1 * time.Hour).Format(time.RFC3339)
	startTimeNew := baseTime.Add(30 * time.Minute).Format(time.RFC3339)
	endTimeNew := baseTime.Add(2 * time.Hour).Format(time.RFC3339)
	liveUrl := "rtmp://5000.liveplay.myqcloud.com/live/stream1"
	// vodUrl := "https://main.qcloudimg.com/video/TVP_HOME.mp4"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCssPullStreamTaskDestroy,
		Steps: []resource.TestStep{
			{
				// PreventDiskCleanup: false,
				Config: fmt.Sprintf(testAccCssPullStreamTask, defaultCSSLiveType, liveUrl, defaultCSSDomainName, defaultCSSAppName, defaultCSSStreamName, startTime, endTime, defaultCSSOperator),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCssPullStreamTaskExists("tencentcloud_css_pull_stream_task.pull_stream_task"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "app_name", defaultCSSAppName),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "stream_name", defaultCSSStreamName),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "domain_name", defaultCSSDomainName),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "comment", "This is a e2e test case."),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "source_type", defaultCSSLiveType),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "start_time", startTime),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "end_time", endTime),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "update_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "callback_events.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "source_urls.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "create_by"),
				),
			},
			{
				// update
				Config: fmt.Sprintf(testAccCssPullStreamTask_update, defaultCSSLiveType, liveUrl, defaultCSSDomainName, defaultCSSAppName, defaultCSSStreamName, startTimeNew, endTimeNew, defaultCSSOperator),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCssPullStreamTaskExists("tencentcloud_css_pull_stream_task.pull_stream_task"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "app_name", defaultCSSAppName),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "stream_name", defaultCSSStreamName),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "domain_name", defaultCSSDomainName),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "comment", "This is a e2e test case_changed."),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "source_type", defaultCSSLiveType),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "start_time", startTimeNew),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "end_time", endTimeNew),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "update_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "callback_events.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "source_urls.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "create_by"),
				),
			},
			{
				ResourceName:            "tencentcloud_css_pull_stream_task.pull_stream_task",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"operator"},
			},
		},
	})
}

func testAccCheckCssPullStreamTaskDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	cssService := CssService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_css_pull_stream_task" {
			continue
		}

		tasks, err := cssService.DescribeCssPullStreamTask(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(tasks) > 0 {
			return fmt.Errorf("css pull stream task still exist, taskId: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckCssPullStreamTaskExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("css pull stream task %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("css pull stream task id is not set")
		}

		cssService := CssService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		tasks, err := cssService.DescribeCssPullStreamTask(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			return fmt.Errorf("css pull stream task not found, taskId: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccCssPullStreamTask = `

resource "tencentcloud_css_pull_stream_task" "pull_stream_task" {
  source_type = "%s"
  source_urls = ["%s"]
  domain_name = "%s"
  app_name = "%s"
  stream_name = "%s"
  start_time = "%s"
  end_time = "%s"
  operator = "%s"
  comment = "This is a e2e test case."
}
`

const testAccCssPullStreamTask_update = `

resource "tencentcloud_css_pull_stream_task" "pull_stream_task" {
  source_type = "%s"
  source_urls = ["%s"]
  domain_name = "%s"
  app_name = "%s"
  stream_name = "%s"
  start_time = "%s"
  end_time = "%s"
  operator = "%s_changed"
  comment = "This is a e2e test case_changed."
}
`
