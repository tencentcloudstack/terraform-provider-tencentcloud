package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainDiagHistoryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainDiagHistoryDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_diag_history.diag_history")),
			},
		},
	})
}

const testAccDbbrainDiagHistoryDataSource = `

data "tencentcloud_dbbrain_diag_history" "diag_history" {
  instance_id = ""
  start_time = ""
  end_time = ""
  product = ""
  }

`
