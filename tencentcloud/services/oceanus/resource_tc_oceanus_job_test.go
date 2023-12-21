package oceanus_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusJobResource_basic -v
func TestAccTencentCloudNeedFixOceanusJobResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusJob,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "job_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "cluster_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "cu_mem"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "folder_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "flink_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "work_space_id"),
				),
			},
			{
				Config: testAccOceanusJobUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "job_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "cluster_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "cu_mem"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "remark"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "folder_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "flink_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job.example", "work_space_id"),
				),
			},
		},
	})
}

const testAccOceanusJob = `
resource "tencentcloud_oceanus_job" "example" {
  name          = "tf_example_job"
  job_type      = 1
  cluster_type  = 2
  cluster_id    = "cluster-1kcd524h"
  cu_mem        = 4
  remark        = "remark."
  folder_id     = "folder-7ctl246z"
  flink_version = "Flink-1.16"
  work_space_id = "space-2idq8wbr"
}
`

const testAccOceanusJobUpdate = `
resource "tencentcloud_oceanus_job" "example" {
  name          = "tf_example_job_update"
  job_type      = 1
  cluster_type  = 2
  cluster_id    = "cluster-1kcd524h"
  cu_mem        = 4
  remark        = "remark update."
  folder_id     = "folder-7ctl246z"
  flink_version = "Flink-1.16"
  work_space_id = "space-2idq8wbr"
}
`
