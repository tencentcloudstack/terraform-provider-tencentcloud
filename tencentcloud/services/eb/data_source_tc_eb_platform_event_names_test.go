package eb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudEbPlatformEventNamesDataSource_basic -v
func TestAccTencentCloudEbPlatformEventNamesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-chongqing")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbPlatformEventNamesDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eb_platform_event_names.platform_event_names"),
				),
			},
		},
	})
}

const testAccEbPlatformEventNamesDataSource = `

data "tencentcloud_eb_platform_event_names" "platform_event_names" {
  product_type = "eb_platform_test"
}

`
