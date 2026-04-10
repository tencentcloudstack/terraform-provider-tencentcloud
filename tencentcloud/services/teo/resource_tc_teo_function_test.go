package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoFunctionResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunction,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "name", "aaa-zone-2qtuhspy7cr6-1310708577"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "remark", "test"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "content", `addEventListener('fetch', e => {
  const response = new Response('Hello World!!');
  e.respondWith(response);
});
`),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_function.teo_function",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccTeoFunctionUp,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "name", "aaa-zone-2qtuhspy7cr6-1310708577"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "remark", "test-update"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "content", `addEventListener('fetch', e => {
  const response = new Response('Hello World');
  e.respondWith(response);
});
`),
				),
			},
		},
	})
}

const testAccTeoFunction = `

resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World!!');
          e.respondWith(response);
        });
    EOT
    name        = "aaa-zone-2qtuhspy7cr6-1310708577"
    remark      = "test"
    zone_id     = "zone-2qtuhspy7cr6"
}
`
const testAccTeoFunctionUp = `

resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World');
          e.respondWith(response);
        });
    EOT
    name        = "aaa-zone-2qtuhspy7cr6-1310708577"
    remark      = "test-update"
    zone_id     = "zone-2qtuhspy7cr6"
}
`

func TestAccTencentCloudTeoFunctionResource_withEnvironmentVariables(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunctionWithEnvVars,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "name", "test-func-env-vars"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "environment_variables.#", "2"),
				),
			},
		},
	})
}

func TestAccTencentCloudTeoFunctionResource_withRules(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoFunctionWithRules,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_function.teo_function", "id"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "name", "test-func-rules"),
					resource.TestCheckResourceAttr("tencentcloud_teo_function.teo_function", "rules.#", "1"),
				),
			},
		},
	})
}

const testAccTeoFunctionWithEnvVars = `

resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World with Env Vars!');
          e.respondWith(response);
        });
    EOT
    name        = "test-func-env-vars"
    remark      = "test with environment variables"
    zone_id     = "zone-2qtuhspy7cr6"

    environment_variables {
        key   = "ENV_VAR_1"
        value = "value1"
        type  = "string"
    }
    environment_variables {
        key   = "ENV_VAR_2"
        value = "{\"key\":\"value\"}"
        type  = "json"
    }
}
`

const testAccTeoFunctionWithRules = `

resource "tencentcloud_teo_function" "teo_function" {
    content     = <<-EOT
        addEventListener('fetch', e => {
          const response = new Response('Hello World with Rules!');
          e.respondWith(response);
        });
    EOT
    name        = "test-func-rules"
    remark      = "test with rules"
    zone_id     = "zone-2qtuhspy7cr6"

    rules {
        function_rule_conditions {
            rule_conditions {
                target   = "url"
                operator = "equals"
                values   = ["/test"]
            }
        }
    }
}
`
