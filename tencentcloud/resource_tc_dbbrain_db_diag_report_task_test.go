package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func init() {
	resource.AddTestSweepers("tencentcloud_dbbrain_db_diag_report_task", &resource.Sweeper{
		Name: "tencentcloud_dbbrain_db_diag_report_task",
		F:    testSweepDbbrainDbDiagReportTask,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_dbbrain_db_diag_report_task
func testSweepDbbrainDbDiagReportTask(r string) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cli, _ := sharedClientForRegion(r)
	dbbrainService := DbbrainService{client: cli.(*TencentCloudClient).apiV3Conn}
	products := []string{
		"cynosdb",
		"mysql",
	}

	for _, product := range products {
		ret, err := dbbrainService.DescribeDbbrainDbDiagReportTaskById(ctx, nil, defaultDbBrainInstanceId, product)
		if err != nil {
			return err
		}
		if ret == nil {
			return fmt.Errorf("Dbbrain Db diag report tasks not exists.")
		}

		delId := *ret.AsyncRequestId
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			err := dbbrainService.DeleteDbbrainDbDiagReportTaskById(ctx, delId, defaultDbBrainInstanceId, product)
			if err != nil {
				return retryError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("[ERROR] sweeper Dbbrain Db diag report tasks:[AsyncRequestId:%v, instanceId:%s, product:%s] failed! reason:[%s]", delId, defaultDbBrainInstanceId, product, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudDbbrainDbDiagReportTaskResource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -1).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().In(loc).Format("2006-01-02 15:04:05")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDbbrainDbDiagReportTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainDbDiagReportTask, defaultDbBrainInstanceId, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbbrainDbDiagReportTaskExists("tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task", "instance_id", defaultDbBrainInstanceId),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task", "start_time", startTime),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task", "end_time", endTime),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task", "send_mail_flag", "0"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task", "product", "mysql"),
				),
			},
		},
	})
}

func testAccCheckDbbrainDbDiagReportTaskDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	dbbrainService := DbbrainService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dbbrain_sql_filter" {
			continue
		}

		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		asyncRequestId := ids[0]
		instanceId := ids[1]
		product := ids[2]

		filter, err := dbbrainService.DescribeDbbrainDbDiagReportTaskById(ctx, helper.StrToInt64Point(asyncRequestId), instanceId, product)
		if err != nil {
			return err
		}

		if filter != nil {
			return fmt.Errorf("Dbbrain Db diag report task, Id: %v", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckDbbrainDbDiagReportTaskExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		dbbrainService := DbbrainService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Dbbrain sql filter  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Dbbrain sql filter id is not set")
		}

		ids := strings.Split(rs.Primary.ID, FILED_SP)
		if len(ids) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		asyncRequestId := ids[0]
		instanceId := ids[1]
		product := ids[2]

		filter, err := dbbrainService.DescribeDbbrainDbDiagReportTaskById(ctx, helper.StrToInt64Point(asyncRequestId), instanceId, product)
		if err != nil {
			return err
		}

		if filter == nil {
			return fmt.Errorf("Dbbrain Db diag report task, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccDbbrainDbDiagReportTask = `

resource "tencentcloud_dbbrain_db_diag_report_task" "db_diag_report_task" {
  instance_id = "%s"
  start_time = "%s"
  end_time = "%s"
  send_mail_flag = 0
  product = "mysql"
}

`
