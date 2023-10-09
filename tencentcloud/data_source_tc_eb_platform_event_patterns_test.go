package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudEbPlatformEventPatternsDataSource_basic -v
func TestAccTencentCloudEbPlatformEventPatternsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-chongqing")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbPlatformEventPatternsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_eb_platform_event_patterns.platform_event_patterns"),
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
