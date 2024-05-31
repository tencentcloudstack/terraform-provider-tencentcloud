package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmChcDeniedActionsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcDeniedActionsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_chc_denied_actions.chc_denied_actions"),
					resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_denied_actions.chc_denied_actions", "chc_host_denied_action_set.#", "1"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_denied_actions.chc_denied_actions", "chc_host_denied_action_set.0.chc_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_denied_actions.chc_denied_actions", "chc_host_denied_action_set.0.state"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_denied_actions.chc_denied_actions", "chc_host_denied_action_set.0.deny_actions.#"),
				),
			},
		},
	})
}

const testAccCvmChcDeniedActionsDataSource = `

data "tencentcloud_cvm_chc_denied_actions" "chc_denied_actions" {
  chc_ids = ["chc-mn3l1qf5"]
}
`
