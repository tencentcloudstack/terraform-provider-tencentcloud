package dts_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcdts "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dts"

	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func init() {
	resource.AddTestSweepers("tencentcloud_dts_compare_task", &resource.Sweeper{
		Name: "tencentcloud_dts_compare_task",
		F:    testSweepDtsCompareTask,
	})
}

// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_dts_compare_task
func testSweepDtsCompareTask(r string) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	cli, _ := tcacctest.SharedClientForRegion(r)
	dtsService := svcdts.NewDtsService(cli.(tccommon.ProviderMeta).GetAPIV3Conn())
	// find migrate job first
	ret, err := dtsService.DescribeDtsMigrateJobsByFilter(ctx, map[string]interface{}{})
	if err != nil {
		return err
	}

	for _, v := range ret {
		if v.JobId == nil || v.CompareTask == nil || v.CompareTask.CompareTaskId == nil {
			continue
		}

		jobId := v.JobId
		compareTaskId := v.CompareTask.CompareTaskId

		ret, err := dtsService.DescribeDtsCompareTask(ctx, jobId, compareTaskId)
		if err != nil {
			return err
		}

		task := ret[0]

		if strings.HasPrefix(*task.TaskName, tcacctest.KeepResource) || strings.HasPrefix(*task.TaskName, tcacctest.DefaultResource) {
			continue
		}

		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			err := dtsService.DeleteDtsCompareTaskById(ctx, *task.JobId, *task.CompareTaskId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("[ERROR] sweeper tencentcloud_dts_compare_task:[%s#%s] failed! reason:[%s]", *jobId, *compareTaskId, err.Error())
		}
	}
	return nil
}

func TestAccTencentCloudDtsCompareTaskResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckDtsCompareTaskDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDtsCompareTask_basic, tcacctest.DefaultDTSJobId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDtsCompareTaskExists("tencentcloud_dts_compare_task.compare_task"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_compare_task.compare_task", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_compare_task.compare_task", "task_name", "tf_test_compare_task"),
				),
			},
			{
				Config: fmt.Sprintf(testAccDtsCompareTask_stop, tcacctest.DefaultDTSJobId, tcacctest.DefaultDTSJobId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDtsCompareTaskExists("tencentcloud_dts_compare_task.compare_task"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_compare_task.compare_task", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dts_compare_task.compare_task", "task_name", "tf_test_compare_task"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_compare_task_stop_operation.stop", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_compare_task_stop_operation.stop", "job_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_dts_compare_task_stop_operation.stop", "compare_task_id"),
				),
			},
		},
	})
}

func testAccCheckDtsCompareTaskDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	dtsService := svcdts.NewDtsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_dts_compare_task" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		jobId := helper.String(idSplit[0])
		compareTaskId := helper.String(idSplit[1])

		task, err := dtsService.DescribeDtsCompareTask(ctx, jobId, compareTaskId)
		if err != nil {
			return err
		}

		if task != nil {
			status := *task[0].Status
			if status != "canceled" {
				return fmt.Errorf("DTS compare task still exist, Id: %v, status:%s", rs.Primary.ID, status)
			}
		}
	}
	return nil
}

func testAccCheckDtsCompareTaskExists(re string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		dtsService := svcdts.NewDtsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())

		rs, ok := s.RootModule().Resources[re]
		if !ok {
			return fmt.Errorf("DTS compare task %s is not found", re)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("DTS compare task id is not set")
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 2 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		jobId := helper.String(idSplit[0])
		compareTaskId := helper.String(idSplit[1])

		task, err := dtsService.DescribeDtsCompareTask(ctx, jobId, compareTaskId)
		if err != nil {
			return err
		}

		if task == nil {
			return fmt.Errorf("DTS compare task not found, Id: %v", rs.Primary.ID)
		}
		return nil
	}
}

const testAccDtsCompareTask_basic = `

resource "tencentcloud_dts_compare_task" "compare_task" {
  job_id = "%s"
  task_name = "tf_test_compare_task"
  objects {
	object_mode = "partial"
  }
}

`

const testAccDtsCompareTask_stop = `

resource "tencentcloud_dts_compare_task" "compare_task" {
  job_id = "%s"
  task_name = "tf_test_compare_task"
  objects {
	object_mode = "partial"
  }
}

resource "tencentcloud_dts_compare_task_stop_operation" "stop" {
	job_id = "%s"
	compare_task_id = tencentcloud_dts_compare_task.compare_task.compare_task_id
  }

`
