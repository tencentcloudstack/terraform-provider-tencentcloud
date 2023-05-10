package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainHealthScoresDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainHealthScoresDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_health_scores.health_scores")),
			},
		},
	})
}

const testAccDbbrainHealthScoresDataSource = `

data "tencentcloud_dbbrain_health_scores" "health_scores" {
  instance_id = ""
  time = ""
  product = ""
  }

`
