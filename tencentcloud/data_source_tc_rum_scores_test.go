package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRumScoresDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumScoresDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_scores.scores")),
			},
		},
	})
}

const testAccRumScoresDataSource = `

data "tencentcloud_rum_scores" "scores" {
  end_time = "2023082215"
  start_time = "2023082214"
  i_d = 1
  is_demo = 1
  }

`
