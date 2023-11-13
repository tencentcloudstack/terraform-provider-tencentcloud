package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmChcDeniedActionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcDeniedActionsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_chc_denied_actions.chc_denied_actions")),
			},
		},
	})
}

const testAccCvmChcDeniedActionsDataSource = `

data "tencentcloud_cvm_chc_denied_actions" "chc_denied_actions" {
  chc_ids = 
  }

`
