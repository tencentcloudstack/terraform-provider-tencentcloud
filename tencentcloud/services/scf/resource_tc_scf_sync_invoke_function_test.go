package scf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudScfSyncInvokeFunctionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfSyncInvokeFunction,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_scf_sync_invoke_function.sync_invoke_function", "id")),
			},
		},
	})
}

const testAccScfSyncInvokeFunction = `

resource "tencentcloud_scf_sync_invoke_function" "invoke_function" {
  function_name = "keep-1676351130"
  qualifier     = "2"
  namespace     = "default"
}

`
