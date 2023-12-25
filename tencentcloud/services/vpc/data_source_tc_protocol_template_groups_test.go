package vpc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataProtocolTemplateGroupsNameAll = "data.tencentcloud_protocol_template_groups.all_test"

func TestAccTencentCloudDataProtocolTemplateGroups(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckProtocolTemplateGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataProtocolTemplateGroupsBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckProtocolTemplateGroupExists("tencentcloud_protocol_template_group.mygroup"),
					resource.TestCheckResourceAttrSet(testDataProtocolTemplateGroupsNameAll, "group_list.0.id"),
					resource.TestCheckResourceAttr(testDataProtocolTemplateGroupsNameAll, "group_list.0.name", "mygroup"),
					resource.TestCheckResourceAttr(testDataProtocolTemplateGroupsNameAll, "group_list.0.template_ids.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudDataProtocolTemplateGroupsBasic = `
resource "tencentcloud_protocol_template" "myprotocol_template" {
  name        = "testacctcrtemplate"
  protocols = ["udp:all"]
}

resource "tencentcloud_protocol_template_group" "mygroup" {
  name        = "mygroup"
  template_ids = [tencentcloud_protocol_template.myprotocol_template.id]
}

data "tencentcloud_protocol_template_groups" "all_test" {
  name = tencentcloud_protocol_template_group.mygroup.name
  id = tencentcloud_protocol_template_group.mygroup.id
}

`
