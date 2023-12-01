package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumScoresDataSource_basic -v
func TestAccTencentCloudRumScoresDataSource_basic(t *testing.T) {
	t.Parallel()

	startTime := time.Now().AddDate(0, 0, -29).Unix()
	endTime := time.Now().Unix()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccRumScoresDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_scores.scores"),
				),
			},
		},
	})
}

const testAccRumScoresDataSource = `

data "tencentcloud_rum_scores" "scores" {
  start_time = %v
  end_time   = %v
  project_id = 120000
  is_demo = 1
}

`
