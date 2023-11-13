package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudVpcBandwidthPackageAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccVpcBandwidthPackageAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_vpc_bandwidth_package_attachment.bandwidth_package_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_vpc_bandwidth_package_attachment.bandwidth_package_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccVpcBandwidthPackageAttachment = `

resource "tencentcloud_vpc_bandwidth_package_attachment" "bandwidth_package_attachment" {
  resource_id = &lt;nil&gt;
  bandwidth_package_id = &lt;nil&gt;
  network_type = &lt;nil&gt;
  resource_type = &lt;nil&gt;
  protocol = &lt;nil&gt;
}

`
