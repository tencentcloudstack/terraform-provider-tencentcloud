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

const testAccTencentCloudDataTCRVPCAttachmentsBasic = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance"
  instance_type = "basic"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "guagua-ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "guagua-ci-temp-test"
  vpc_id            = tencentcloud_vpc.foo.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_tcr_vpc_attachment" "mytcr_vpc_attachment" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  vpc_id = tencentcloud_vpc.foo.id
  subnet_id = tencentcloud_subnet.subnet.id
}

data "tencentcloud_tcr_vpc_attachments" "id_test" {
  instance_id = tencentcloud_tcr_vpc_attachment.mytcr_vpc_attachment.instance_id
}
`
