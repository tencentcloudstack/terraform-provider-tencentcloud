package tencentcloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccBindingObjects(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBindingObjects(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_binding_objects.objects"),
				),
			},
		},
	})
}

func testAccDataSourceBindingObjects() string {
	return fmt.Sprintf(`data "tencentcloud_monitor_policy_groups" "name" {
}

data "tencentcloud_monitor_binding_objects" "objects" {
  group_id = data.tencentcloud_monitor_policy_groups.name.list[0].group_id
}`)
}
