package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
		delName :=*v.StreamName
		delId:=*v.TaskId

		if strings.HasPrefix(delName, "test_") {
			err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
				err := cssService.DeleteCssPullStreamTaskById(ctx, delId)
				if err != nil {
					return retryError(err)
				}
				return nil
			})
			if err != nil {
				return fmt.Errorf("[ERROR] instance %s:%s failed! reason:[%s]", delName,delId, err.Error())
			}
		}
	}
	return nil
}

func TestAccTencentCloudCSSPullStreamTaskResource_basic(t *testing.T) {
	t.Parallel()
	startTime := time.Now().Add(2*time.Hour).Format("2006-01-02T15:04:05+08:00")
	endTime := time.Now().Add(4*time.Hour).Format("2006-01-02T15:04:05+08:00")

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: testAccCheckCssPullStreamTaskDestroy,
		Steps: []resource.TestStep{
			{
				PreventDiskCleanup: false,
				Config: fmt.Sprintf(testAccCssPullStreamTask, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCssPullStreamTaskExists("tencentcloud_css_pull_stream_task.pull_stream_task"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "id"),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "app_name", "live"),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "stream_name", "test_stream_name"),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "domain_name", "177154.push.tlivecloud.com"),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "comment", "This is a e2e test case."),
					resource.TestCheckResourceAttr("tencentcloud_css_pull_stream_task.pull_stream_task", "source_type", "PullLivePushLive"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "start_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "end_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "callback_events.#"),
					
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "source_urls.#"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "create_by"),
					resource.TestCheckResourceAttrSet("tencentcloud_css_pull_stream_task.pull_stream_task", "push_args"),

				),
			},
			{
				ResourceName:      "tencentcloud_css_pull_stream_task.pullStreamTask",
				ImportState:       true,
				ImportStateVerify: true,
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
  source_type = "PullLivePushLive"
  source_urls = ["rtmp://5000.liveplay.myqcloud.com/live/stream1"]
  domain_name = "177154.push.tlivecloud.com"
  app_name = "live"
  stream_name = "test_stream_name"
  start_time = "%s"
  end_time = "%s"
  operator = "tf_admin"
  comment = "This is a e2e test case."
}
`
