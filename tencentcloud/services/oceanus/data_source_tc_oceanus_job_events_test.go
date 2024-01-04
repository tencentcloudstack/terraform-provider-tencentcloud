package oceanus_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusJobEventsDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusJobEventsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusJobEventsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_job_events.example"),
					resource.TestCheckResourceAttr("data.tencentcloud_oceanus_job_events.example", "job_id", "cql-6w8eab6f"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_job_events.example", "start_timestamp"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_job_events.example", "end_timestamp"),
					resource.TestCheckResourceAttr("data.tencentcloud_oceanus_job_events.example", "work_space_id", "space-6w8eab6f"),
				),
			},
		},
	})
}

const testAccOceanusJobEventsDataSource = `
data "tencentcloud_oceanus_job_events" "example" {
  job_id          = "cql-6w8eab6f"
  start_timestamp = 1630932161
  end_timestamp   = 1631232466
  types           = ["1", "2"]
  work_space_id   = "space-6w8eab6f"
}
`
