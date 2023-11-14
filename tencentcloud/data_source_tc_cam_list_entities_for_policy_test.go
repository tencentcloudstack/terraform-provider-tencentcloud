package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamListEntitiesForPolicyDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamListEntitiesForPolicyDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cam_list_entities_for_policy.list_entities_for_policy")),
			},
		},
	})
}

const testAccCamListEntitiesForPolicyDataSource = `

data "tencentcloud_cam_list_entities_for_policy" "list_entities_for_policy" {
  policy_id = 
  rp = 
  entity_filter = ""
    }

`
