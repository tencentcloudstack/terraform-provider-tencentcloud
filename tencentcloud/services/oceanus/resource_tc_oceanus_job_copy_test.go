package oceanus_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusJobCopyResource_basic -v
func TestAccTencentCloudNeedFixOceanusJobCopyResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusJobCopy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_copy.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_copy.example", "source_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_copy.example", "target_cluster_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_copy.example", "source_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_copy.example", "target_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_copy.example", "target_folder_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_copy.example", "job_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_oceanus_job_copy.example", "work_space_id"),
				),
			},
		},
	})
}

const testAccOceanusJobCopy = `
resource "tencentcloud_oceanus_job_copy" "example" {
  source_id         = "cql-0nob2hx8"
  target_cluster_id = "cluster-1kcd524h"
  source_name       = "keep_jar"
  target_name       = "tf_copy_example"
  target_folder_id  = "folder-7ctl246z"
  job_type          = 2
  work_space_id     = "space-2idq8wbr"
}
`
