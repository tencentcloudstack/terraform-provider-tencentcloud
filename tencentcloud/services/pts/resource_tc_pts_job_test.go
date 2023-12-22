package pts_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	svcpts "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/pts"
)

// go test -i; go test -test.run TestAccTencentCloudPtsJobResource_basic -v
func TestAccTencentCloudPtsJobResource_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckPtsJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPtsJob,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPtsJobExists("tencentcloud_pts_job.job"),
					resource.TestCheckResourceAttr("tencentcloud_pts_job.job", "scenario_id", tcacctest.DefaultScenarioIdJob),
					resource.TestCheckResourceAttr("tencentcloud_pts_job.job", "job_owner", "username"),
					resource.TestCheckResourceAttr("tencentcloud_pts_job.job", "project_id", tcacctest.DefaultPtsProjectId),
					resource.TestCheckResourceAttr("tencentcloud_pts_job.job", "note", "desc"),
				),
			},
			{
				ResourceName:      "tencentcloud_pts_job.job",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckPtsJobDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcpts.NewPtsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_pts_project" {
			continue
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		projectId := idSplit[0]
		scenarioId := idSplit[1]
		jobId := idSplit[2]

		job, err := service.DescribePtsJob(ctx, projectId, scenarioId, jobId)
		if job != nil {
			return fmt.Errorf("pts job %s still exists", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckPtsJobExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}

		idSplit := strings.Split(rs.Primary.ID, tccommon.FILED_SP)
		if len(idSplit) != 3 {
			return fmt.Errorf("id is broken,%s", rs.Primary.ID)
		}
		projectId := idSplit[0]
		scenarioId := idSplit[1]
		jobId := idSplit[2]

		service := svcpts.NewPtsService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		job, err := service.DescribePtsJob(ctx, projectId, scenarioId, jobId)
		if job == nil {
			return fmt.Errorf("pts job %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccPtsJobVar = `
variable "project_id" {
  default = "` + tcacctest.DefaultPtsProjectId + `"
}
variable "scenario_id" {
	default = "` + tcacctest.DefaultScenarioIdJob + `"
}
`

const testAccPtsJob = testAccPtsJobVar + `

resource "tencentcloud_pts_job" "job" {
	scenario_id = var.scenario_id
	job_owner = "username"
	project_id = var.project_id
	# debug = ""
	note = "desc"
  }

`
