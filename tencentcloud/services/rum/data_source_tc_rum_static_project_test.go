package rum_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumStaticProjectDataSource_basic -v
func TestAccTencentCloudRumStaticProjectDataSource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccRumStaticProjectDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_rum_static_project.static_project"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_static_project.static_project", "result"),
				),
			},
		},
	})
}

const testAccRumStaticProjectDataSource = `

data "tencentcloud_rum_static_project" "static_project" {
  start_time = %v
  type       = "allcount"
  end_time   = %v
  project_id = 120000
}

`
