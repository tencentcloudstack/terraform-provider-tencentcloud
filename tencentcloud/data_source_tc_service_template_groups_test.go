package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataServiceTemplateGroupsNameAll = "data.tencentcloud_service_template_groups.all_test"

func TestAccTencentCloudDataServiceTemplateGroups(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceTemplateGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataServiceTemplateGroupsBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckServiceTemplateGroupExists("tencentcloud_service_template_group.mygroup"),
					resource.TestCheckResourceAttrSet(testDataServiceTemplateGroupsNameAll, "group_list.0.id"),
					resource.TestCheckResourceAttr(testDataServiceTemplateGroupsNameAll, "group_list.0.name", "mygroup"),
					resource.TestCheckResourceAttr(testDataServiceTemplateGroupsNameAll, "group_list.0.template_ids.#", "1"),
				),
			},
		},
	})
}

const testAccTencentCloudDataServiceTemplateGroupsBasic = `
resource "tencentcloud_service_template" "myservice_template" {
  name        = "testacctcrtemplate"
  services = ["udp:all"]
}

resource "tencentcloud_service_template_group" "mygroup" {
  name        = "mygroup"
  template_ids = [tencentcloud_service_template.myservice_template.id]
}

data "tencentcloud_service_template_groups" "all_test" {
  name = tencentcloud_service_template_group.mygroup.name
  id = tencentcloud_service_template_group.mygroup.id
}

`
