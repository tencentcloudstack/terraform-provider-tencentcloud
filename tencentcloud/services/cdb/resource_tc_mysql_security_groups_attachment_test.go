package cdb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMysqlSecurityGroupsAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlSecurityGroupsAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_security_groups_attachment.security_groups_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_security_groups_attachment.security_groups_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMysqlSecurityGroupsAttachment = `

resource "tencentcloud_mysql_security_groups_attachment" "security_groups_attachment" {
  security_group_id = "sg-baxfiao5"
  instance_id       = "cdb-fitq5t9h"
}

`
