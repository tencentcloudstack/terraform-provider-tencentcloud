package monitor_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBindingObjects(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceBindingObjects(),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_monitor_binding_objects.objects"),
				),
			},
		},
	})
}

func testAccDataSourceBindingObjects() string {
	return `data "tencentcloud_monitor_policy_groups" "name" {
}

data "tencentcloud_monitor_binding_objects" "objects" {
  group_id = data.tencentcloud_monitor_policy_groups.name.list[0].group_id
}`
}
