package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainHealthScoreDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainHealthScoreDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_health_score.health_score")),
			},
		},
	})
}

const testAccDbbrainHealthScoreDataSource = `

data "tencentcloud_dbbrain_health_score" "health_score" {
  instance_id = ""
  time = ""
  product = ""
  }

`
