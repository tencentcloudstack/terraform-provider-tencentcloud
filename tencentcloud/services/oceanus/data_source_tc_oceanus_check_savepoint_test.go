package oceanus_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixOceanusCheckSavepointDataSource_basic -v
func TestAccTencentCloudNeedFixOceanusCheckSavepointDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusCheckSavepointDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_check_savepoint.example"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_check_savepoint.example", "job_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_check_savepoint.example", "serial_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_check_savepoint.example", "record_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_check_savepoint.example", "savepoint_path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_oceanus_check_savepoint.example", "work_space_id"),
				),
			},
		},
	})
}

const testAccOceanusCheckSavepointDataSource = `
data "tencentcloud_oceanus_check_savepoint" "example" {
  job_id         = "cql-314rw6w0"
  serial_id      = "svp-52xkpymp"
  record_type    = 1
  savepoint_path = "cosn://52xkpymp-12345/12345/10000/cql-12345/2/flink-savepoints/savepoint-000000-12334"
  work_space_id  = "space-2idq8wbr"
}
`
