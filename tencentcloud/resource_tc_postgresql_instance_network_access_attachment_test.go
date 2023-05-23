package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudPostgresqlInstanceNetworkAccessAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlInstanceNetworkAccessAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_postgresql_instance_network_access_attachment.instance_network_access_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_postgresql_instance_network_access_attachment.instance_network_access_attachment",
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
