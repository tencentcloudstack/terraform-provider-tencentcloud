package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCdbSecurityGroupsAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbSecurityGroupsAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_security_groups_attachment.security_groups_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_security_groups_attachment.security_groups_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbSecurityGroupsAttachment = `

resource "tencentcloud_cdb_security_groups_attachment" "security_groups_attachment" {
  security_group_id = &lt;nil&gt;
  instance_ids = &lt;nil&gt;
  for_readonly_instance = false
}

`
