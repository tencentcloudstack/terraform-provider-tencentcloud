package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainDbSpaceStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainDbSpaceStatusDataSource, defaultDbBrainInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_db_space_status.db_space_status"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_db_space_status.db_space_status", "instance_id", defaultDbBrainInstanceId),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_db_space_status.db_space_status", "range_days", "7"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_db_space_status.db_space_status", "product", "mysql"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_db_space_status.db_space_status", "growth"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_db_space_status.db_space_status", "remain"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_db_space_status.db_space_status", "total"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_db_space_status.db_space_status", "available_days"),
				),
			},
		},
	})
}

const testAccDbbrainDbSpaceStatusDataSource = `

data "tencentcloud_dbbrain_db_space_status" "db_space_status" {
  instance_id = "%s"
  range_days = 7
  product = "mysql"
}

`
