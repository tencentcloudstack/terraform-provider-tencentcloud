package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseRegionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseRegionDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_region.region")),
			},
		},
	})
}

const testAccLighthouseRegionDataSource = `

data "tencentcloud_lighthouse_region" "region" {
}
`
