package tencentcloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDbbrainDiagEventDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainDiagEventDataSource, defaultDbBrainInstanceId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_diag_event.diag_event"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "diag_item"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "diag_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "explanation"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "outline"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "problem"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "severity"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "suggestions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "metric"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "result_output_file"),
				),
			},
		},
	})
}

const testAccDbbrainDiagEventDataSource = `

data "tencentcloud_dbbrain_diag_event" "diag_event" {
  instance_id = "%s"
  product = "mysql"
}

`
