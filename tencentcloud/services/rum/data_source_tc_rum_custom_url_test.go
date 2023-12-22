package rum_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumCustomUrlDataSource_basic -v
func TestAccTencentCloudRumCustomUrlDataSource_basic(t *testing.T) {
	t.Parallel()

	startTime := time.Now().AddDate(0, 0, -29).Unix()
	endTime := time.Now().Unix()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccRumCustomUrlDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_rum_custom_url.custom_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_custom_url.custom_url", "result"),
				),
			},
		},
	})
}

const testAccRumCustomUrlDataSource = `

data "tencentcloud_rum_custom_url" "custom_url" {
  start_time = %v
  type       = "top"
  end_time   = %v
  project_id = 120000
}

`
