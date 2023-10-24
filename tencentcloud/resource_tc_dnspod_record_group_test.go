package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDnspodRecordGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDnspodRecordGroup,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dnspod_record_group.record_group", "id")),
			},
			{
				ResourceName:      "tencentcloud_dnspod_record_group.record_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDnspodRecordGroup = `

resource "tencentcloud_dnspod_record_group" "record_group" {
  domain = "dnspod.cn"
  group_name = "group_name_demo"
  domain_id = 123
  tags = {
    "createdBy" = "terraform"
  }
}

`
