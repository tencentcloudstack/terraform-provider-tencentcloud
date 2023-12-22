package dbbrain_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainDbSpaceStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainDbSpaceStatusDataSource, tcacctest.DefaultDbBrainInstanceId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_db_space_status.db_space_status"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_db_space_status.db_space_status", "instance_id", tcacctest.DefaultDbBrainInstanceId),
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
