package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudHaVipInstanceAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccHaVipInstanceAttachmentEni,
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("tencentcloud_ha_vip_instance_attachment.ha_vip_instance_attachment", "id"),
				resource.TestCheckResourceAttr("tencentcloud_ha_vip_instance_attachment.ha_vip_instance_attachment", "instance_id", "eni-mdsbxqeb"),
				resource.TestCheckResourceAttrSet("tencentcloud_ha_vip_instance_attachment.ha_vip_instance_attachment", "ha_vip_id"),
				resource.TestCheckResourceAttr("tencentcloud_ha_vip_instance_attachment.ha_vip_instance_attachment", "instance_type", "ENI"),
			),
		}, {
			ResourceName:      "tencentcloud_ha_vip_instance_attachment.ha_vip_instance_attachment",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccHaVipInstanceAttachmentEni = `

resource "tencentcloud_ha_vip" "foo" {
  name      = "terraform_test"
  vpc_id    = "vpc-axrsmmrv"
  subnet_id = "subnet-kxaxknmg"
  vip       = "172.16.144.10"
  check_associate = true
}
resource "tencentcloud_ha_vip_instance_attachment" "ha_vip_instance_attachment" {
  instance_id = "eni-mdsbxqeb"
  ha_vip_id = tencentcloud_ha_vip.foo.id
  instance_type = "ENI"
}
`
