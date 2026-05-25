package metadatasources_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tcfwacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework/acctest"
)

// TestAccTencentCloudProviderRuntimeDataSource_basic is the first
// acceptance test under the framework + SDKv2 mux architecture. It boots
// both stacks via the ProtoV5 factories, reads
// data.tencentcloud_provider_runtime, and asserts that the framework data
// source can plan/apply correctly inside a muxed environment.
//
// The test does not call any cloud API and therefore can run without real
// credentials (it only reads local provider runtime metadata). AccPreCheck
// is still executed so that the SDKv2 provider completes Configure and
// injects the shared client into sharedmeta.
func TestAccTencentCloudProviderRuntimeDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { tcacctest.AccPreCheck(t) },
		ProtoV5ProviderFactories: tcfwacctest.AccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderRuntimeDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"data.tencentcloud_provider_runtime.this", "id"),
					resource.TestCheckResourceAttrSet(
						"data.tencentcloud_provider_runtime.this", "region"),
					resource.TestCheckResourceAttrSet(
						"data.tencentcloud_provider_runtime.this", "client_version"),
					resource.TestCheckResourceAttr(
						"data.tencentcloud_provider_runtime.this", "stack_mode",
						"sdkv2+framework"),
					resource.TestCheckResourceAttrSet(
						"data.tencentcloud_provider_runtime.this", "protocol"),
					resource.TestCheckResourceAttrSet(
						"data.tencentcloud_provider_runtime.this", "domain"),
					resource.TestCheckResourceAttrSet(
						"data.tencentcloud_provider_runtime.this", "cos_domain"),
					resource.TestCheckResourceAttr(
						"data.tencentcloud_provider_runtime.this", "secret_id_present",
						"true"),
				),
			},
		},
	})
}

const testAccProviderRuntimeDataSourceConfigBasic = `
data "tencentcloud_provider_runtime" "this" {}
`
