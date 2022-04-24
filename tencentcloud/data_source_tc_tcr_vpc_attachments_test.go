package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTCRVPCAttachmentsNameAll = "data.tencentcloud_tcr_vpc_attachments.id_test"

func TestAccTencentCloudDataTCRVPCAttachments(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTCRVPCAttachmentsBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRVPCAttachmentExists("tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment"),
					resource.TestCheckResourceAttr(testDataTCRVPCAttachmentsNameAll, "vpc_attachment_list.#", "1"),
					resource.TestCheckResourceAttrSet(testDataTCRVPCAttachmentsNameAll, "vpc_attachment_list.0.status"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTCRVPCAttachmentsBasic = defaultVpcSubnets + `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "test-tcr-attach"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_vpc_attachment" "mytcr_vpc_attachment" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  vpc_id = local.vpc_id
  subnet_id = local.subnet_id
}

data "tencentcloud_tcr_vpc_attachments" "id_test" {
  instance_id = tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment.instance_id
}
`
