package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudScfFunctionEventInvokeConfigDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfFunctionEventInvokeConfigDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_function_event_invoke_config.function_event_invoke_config")),
			},
		},
	})
}

const testAccScfFunctionEventInvokeConfigDataSource = `

data "tencentcloud_scf_function_event_invoke_config" "function_event_invoke_config" {
  function_name = ""
  namespace = ""
  qualifier = ""
  }

`
