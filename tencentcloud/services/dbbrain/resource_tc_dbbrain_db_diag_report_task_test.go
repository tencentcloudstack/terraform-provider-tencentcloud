package dbbrain_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	svcdbbrain "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbbrain"
)

func init() {
	resource.AddTestSweepers("tencentcloud_dbbrain_db_diag_report_task", &resource.Sweeper{
		Name: "tencentcloud_dbbrain_db_diag_report_task",
		F:    testSweepDbbrainDbDiagReportTask,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_dbbrain_db_diag_report_task
func testSweepDbbrainDbDiagReportTask(r string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(r)
	dbbrainService := svcdbbrain.NewDbbrainService(cli.(tccommon.ProviderMeta).GetAPIV3Conn())
	products := []string{
		"cynosdb",
		"mysql",
	}

	for _, product := range products {
		ret, err := dbbrainService.DescribeDbbrainDbDiagReportTaskById(ctx, nil, tcacctest.DefaultDbBrainInstanceId, product)
		if err != nil {
			return err
		}
		if ret == nil {
			return fmt.Errorf("Dbbrain Db diag report tasks not exists.")
		}

		delId := *ret.AsyncRequestId
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			err := dbbrainService.DeleteDbbrainDbDiagReportTaskById(ctx, delId, tcacctest.DefaultDbBrainInstanceId, product)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("[ERROR] sweeper Dbbrain Db diag report tasks:[AsyncRequestId:%v, instanceId:%s, product:%s] failed! reason:[%s]", delId, tcacctest.DefaultDbBrainInstanceId, product, err.Error())
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
			tcacctest.AccPreCheck(t)
		},
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDbbrainDbDiagReportTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainDbDiagReportTask, tcacctest.DefaultDbBrainInstanceId, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDbbrainDbDiagReportTaskExists("tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task"),
					resource.TestCheckResourceAttrSet("tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dbbrain_db_diag_report_task.db_diag_report_task", "instance_id", tcacctest.DefaultDbBrainInstanceId),
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
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	dbbrainService := svcdbbrain.NewDbbrainService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dbbrain_sql_filter" {
			continue
		}

		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		dbbrainService := svcdbbrain.NewDbbrainService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("Dbbrain sql filter  %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Dbbrain sql filter id is not set")
		}

		ids := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
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
