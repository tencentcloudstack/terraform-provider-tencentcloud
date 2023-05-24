package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const testAccPgInstanceNetworkAccessAttachmentObject = "tencentcloud_postgresql_instance_network_access_attachment.instance_network_access_attachment"

func TestAccTencentCloudPostgresqlInstanceNetworkAccessAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		CheckDestroy: testAccCheckPgInstanceNetworkAccessAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlInstanceNetworkAccessAttachment,
				Check:  resource.ComposeTestCheckFunc(
					testAccCheckTCRVPCAttachmentExists(testAccPgInstanceNetworkAccessAttachmentObject)
					resource.TestCheckResourceAttrSet(testAccPgInstanceNetworkAccessAttachmentObject, "id")),
			},
			{
				ResourceName:      testAccPgInstanceNetworkAccessAttachmentObject,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresqlInstanceNetworkAccessAttachment = `

resource "tencentcloud_postgresql_instance_network_access_attachment" "instance_network_access_attachment" {
  d_b_instance_id = ""
  vpc_id = "vpc-xxx"
  subnet_id = "subnet-xxx"
  is_assign_vip = false
  vip = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
