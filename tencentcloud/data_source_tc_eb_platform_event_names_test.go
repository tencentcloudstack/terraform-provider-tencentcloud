package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudEbPlatformEventNamesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbPlatformEventNamesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_eb_platform_event_names.platform_event_names")),
			},
		},
	})
}

const testAccEbPlatformEventNamesDataSource = `

data "tencentcloud_eb_platform_event_names" "platform_event_names" {
  product_type = ""
  }

`
