package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudPostgresqlReadonlyGroupNetworkAccessAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlReadonlyGroupNetworkAccessAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_postgresql_readonly_group_network_access_attachment.readonly_group_network_access_attachment", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_postgresql_readonly_group_network_access_attachment.readonly_group_network_access_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccPostgresqlReadonlyGroupNetworkAccessAttachment = `

resource "tencentcloud_postgresql_readonly_group_network_access_attachment" "readonly_group_network_access_attachment" {
  readonly_group_id = "pgro-xxxx"
  vpc_id = "vpc-xxx"
  subnet_id = "subnet-xxx"
  is_assign_vip = false
  vip = ""
  tags = {
    "createdBy" = "terraform"
  }
}

`
