package cvm_test

import (
	"testing"

	resource "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	acctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCvmChcDeniedActionsDataSource_Basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.AccPreCheck(t)
		},
		Providers: acctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcDeniedActionsDataSource_BasicCreate,
				Check:  resource.ComposeTestCheckFunc(acctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_cvm_chc_denied_actions.chc_denied_actions"), resource.TestCheckResourceAttr("data.tencentcloud_cvm_chc_denied_actions.chc_denied_actions", "chc_host_denied_action_set.#", "1"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_denied_actions.chc_denied_actions", "chc_host_denied_action_set.0.chc_id"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_denied_actions.chc_denied_actions", "chc_host_denied_action_set.0.state"), resource.TestCheckResourceAttrSet("data.tencentcloud_cvm_chc_denied_actions.chc_denied_actions", "chc_host_denied_action_set.0.deny_actions.#")),
			},
		},
	})
}

const testAccCvmChcDeniedActionsDataSource_BasicCreate = `

data "tencentcloud_cvm_chc_denied_actions" "chc_denied_actions" {
    chc_ids = ["chc-mn3l1qf5"]
}

`
