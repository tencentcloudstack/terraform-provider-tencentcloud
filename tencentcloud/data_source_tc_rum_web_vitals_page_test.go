package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumWebVitalsPageDataSource_basic -v
func TestAccTencentCloudRumWebVitalsPageDataSource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccRumWebVitalsPageDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_web_vitals_page.web_vitals_page"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_web_vitals_page.web_vitals_page", "result"),
				),
			},
		},
	})
}

const testAccRumWebVitalsPageDataSource = `

data "tencentcloud_rum_web_vitals_page" "web_vitals_page" {
  start_time = %v
  type       = "from"
  end_time   = %v
  project_id = 120000
}

`
