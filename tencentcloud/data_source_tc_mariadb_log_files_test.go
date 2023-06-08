package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbLogFilesDataSource_basic -v
func TestAccTencentCloudMariadbLogFilesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbLogFilesDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_log_files.log_files"),
				),
			},
		},
	})
}

const testAccMariadbLogFilesDataSource = `
data "tencentcloud_mariadb_log_files" "log_files" {
  instance_id = "tdsql-9vqvls95"
  type        = 1
}
`
