package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcBandwidthPackageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandwidthPackage,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_bandwidth_package.bandwidth_package", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_bandwidth_package.bandwidth_package",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcBandwidthPackage = `

resource "tencentcloud_vpc_bandwidth_package" "bandwidth_package" {
  network_type = &lt;nil&gt;
  charge_type = &lt;nil&gt;
  bandwidth_package_name = &lt;nil&gt;
  bandwidth_package_count = &lt;nil&gt;
  internet_max_bandwidth = &lt;nil&gt;
  protocol = &lt;nil&gt;
  tags = {
    "createdBy" = "terraform"
  }
}

`
