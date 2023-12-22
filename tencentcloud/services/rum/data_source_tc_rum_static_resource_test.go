package rum_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumStaticResourceDataSource_basic -v
func TestAccTencentCloudRumStaticResourceDataSource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccRumStaticResourceDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_rum_static_resource.static_resource"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_static_resource.static_resource", "result"),
				),
			},
		},
	})
}

const testAccRumStaticResourceDataSource = `

data "tencentcloud_rum_static_resource" "static_resource" {
  start_time = %v
  type       = "top"
  end_time   = %v
  project_id = 120000
}

`
