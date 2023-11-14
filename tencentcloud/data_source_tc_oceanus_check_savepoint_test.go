package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudOceanusCheckSavepointDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOceanusCheckSavepointDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_oceanus_check_savepoint.check_savepoint")),
			},
		},
	})
}

const testAccOceanusCheckSavepointDataSource = `

data "tencentcloud_oceanus_check_savepoint" "check_savepoint" {
  job_id = "cql-52xkpymp"
    record_type = 1
  savepoint_path = "cosn://xxx/xxx"
  work_space_id = "space-1327"
  }

`
