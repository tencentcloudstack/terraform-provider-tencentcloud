package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudHaVipEipAttachmentsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckHaVipEipAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccHaVipEipAttachmentsDataSource_basic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckHaVipEipAttachmentExists("tencentcloud_ha_vip_eip_attachment.ha_vip_eip_attachment"),
					resource.TestCheckResourceAttr("data.tencentcloud_ha_vip_eip_attachments.ha_vip_eip_attachments", "ha_vip_eip_attachment_list.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ha_vip_eip_attachments.ha_vip_eip_attachments", "ha_vip_eip_attachment_list.0.havip_id"),
				),
			},
		},
	})
}

const testAccHaVipEipAttachmentsDataSource_basic = tcacctest.DefaultVpcVariable + `
#Create EIP
resource "tencentcloud_eip" "eip" {
  name = "havip_eip"
}
resource "tencentcloud_ha_vip" "havip" {
  name       = "terraform_test"
  vpc_id     = var.vpc_id
  subnet_id  = var.subnet_id
}
resource "tencentcloud_ha_vip_eip_attachment" "ha_vip_eip_attachment" {
  havip_id   = tencentcloud_ha_vip.havip.id
  address_ip = tencentcloud_eip.eip.public_ip
}

data "tencentcloud_ha_vip_eip_attachments" "ha_vip_eip_attachments" {
  havip_id = tencentcloud_ha_vip_eip_attachment.ha_vip_eip_attachment.havip_id
}
`
