package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataAddressTemplateGroupsNameAll = "data.tencentcloud_address_template_groups.all_test"

func TestAccTencentCloudDataAddressTemplateGroups(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAddressTemplateGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataAddressTemplateGroupsBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckAddressTemplateGroupExists("tencentcloud_address_template_group.mygroup"),
					resource.TestCheckResourceAttrSet(testDataAddressTemplateGroupsNameAll, "group_list.0.id"),
					resource.TestCheckResourceAttr(testDataAddressTemplateGroupsNameAll, "group_list.0.name", "mygroup"),
					resource.TestCheckResourceAttr(testDataAddressTemplateGroupsNameAll, "group_list.0.template_ids.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudDataAddressTemplateGroupsBasic = `
resource "tencentcloud_address_template" "myaddress_template" {
  name        = "testacctcrtemplate"
  addresses = ["1.1.1.1"]
}

resource "tencentcloud_address_template_group" "mygroup" {
  name        = "mygroup"
  template_ids = [tencentcloud_address_template.myaddress_template.id]
}

data "tencentcloud_address_template_groups" "all_test" {
  name = tencentcloud_address_template_group.mygroup.name
  id = tencentcloud_address_template_group.mygroup.id
}

`
