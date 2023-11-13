package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbRollbackTimeValidityDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbRollbackTimeValidityDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_rollback_time_validity.rollback_time_validity")),
			},
		},
	})
}

const testAccCynosdbRollbackTimeValidityDataSource = `

data "tencentcloud_cynosdb_rollback_time_validity" "rollback_time_validity" {
  cluster_id = ""
  expect_time = ""
  expect_time_thresh = 
        }

`
