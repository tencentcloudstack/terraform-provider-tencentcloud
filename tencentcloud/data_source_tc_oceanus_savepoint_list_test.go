package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusSavepointListDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusSavepointListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusSavepointListDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_savepoint_list.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_savepoint_list.example", "job_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_savepoint_list.example", "work_space_id"),
				),
			},
		},
	})
}

const testAccOceanusSavepointListDataSource = `
data "tencentcloud_oceanus_savepoint_list" "example" {
  job_id        = "cql-314rw6w0"
  work_space_id = "space-2idq8wbr"
}
`
