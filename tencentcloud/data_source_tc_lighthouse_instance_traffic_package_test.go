package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseInstanceTrafficPackageDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseInstanceTrafficPackageDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_instance_traffic_package.instance_traffic_package")),
			},
		},
	})
}

const testAccLighthouseInstanceTrafficPackageDataSource = `
data "tencentcloud_lighthouse_instance_traffic_package" "instance_traffic_package" {
}
`
