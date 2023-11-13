package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
				Config: testAccDbbrainDiagEventDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_diag_event.diag_event")),
			},
		},
	})
}

const testAccDbbrainDiagEventDataSource = `

data "tencentcloud_dbbrain_diag_event" "diag_event" {
  instance_id = ""
  event_id = 
  product = ""
                      }

`
