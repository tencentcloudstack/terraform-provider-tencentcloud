package oceanus_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusJobSubmissionLogDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusJobSubmissionLogDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusJobSubmissionLogDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_job_submission_log.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_job_submission_log.example", "job_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_job_submission_log.example", "start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_job_submission_log.example", "end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_job_submission_log.example", "running_order_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_job_submission_log.example", "order_type"),
				),
			},
		},
	})
}

const testAccOceanusJobSubmissionLogDataSource = `
data "tencentcloud_oceanus_job_submission_log" "example" {
  job_id           = "cql-314rw6w0"
  start_time       = 1696130964345
  end_time         = 1698118169241
  running_order_id = 0
  order_type       = "desc"
}
`
