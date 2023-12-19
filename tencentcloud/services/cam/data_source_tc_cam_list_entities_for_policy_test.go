package cam_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamListEntitiesForPolicyDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamListEntitiesForPolicyDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cam_list_entities_for_policy.list_entities_for_policy"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_list_entities_for_policy.list_entities_for_policy", "policy_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cam_list_entities_for_policy.list_entities_for_policy", "entity_filter"),
				),
			},
		},
	})
}

const testAccCamListEntitiesForPolicyDataSource = `

data "tencentcloud_cam_list_entities_for_policy" "list_entities_for_policy" {
  policy_id = 1
  entity_filter = "All"
    }

`
