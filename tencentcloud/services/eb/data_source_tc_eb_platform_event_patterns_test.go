package eb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudEbPlatformEventPatternsDataSource_basic -v
func TestAccTencentCloudEbPlatformEventPatternsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-chongqing")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbPlatformEventPatternsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eb_platform_event_patterns.platform_event_patterns"),
				),
			},
		},
	})
}

const testAccEbPlatformEventPatternsDataSource = `

data "tencentcloud_eb_platform_event_patterns" "platform_event_patterns" {
  product_type = "eb_platform_test"
  }

`
