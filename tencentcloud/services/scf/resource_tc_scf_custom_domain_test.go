package scf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudScfCustomDomainResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfCustomDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_scf_custom_domain.scf_custom_domain", "id"),
					resource.TestCheckResourceAttr("tencentcloud_scf_custom_domain.scf_custom_domain", "domain", "scf.iac-tf.cloud"),
					resource.TestCheckResourceAttr("tencentcloud_scf_custom_domain.scf_custom_domain", "protocol", "HTTP"),
					resource.TestCheckResourceAttr("tencentcloud_scf_custom_domain.scf_custom_domain", "endpoints_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_scf_custom_domain.scf_custom_domain", "endpoints_config.0.namespace", "default"),
					resource.TestCheckResourceAttrSet("tencentcloud_scf_custom_domain.scf_custom_domain", "endpoints_config.0.function_name"),
					resource.TestCheckResourceAttr("tencentcloud_scf_custom_domain.scf_custom_domain", "endpoints_config.0.qualifier", "$LATEST"),
					resource.TestCheckResourceAttr("tencentcloud_scf_custom_domain.scf_custom_domain", "endpoints_config.0.path_match", "/a/*"),
					resource.TestCheckResourceAttr("tencentcloud_scf_custom_domain.scf_custom_domain", "waf_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_scf_custom_domain.scf_custom_domain", "waf_config.0.waf_open", "CLOSE"),
				),
			},
			{
				Config: testAccScfCustomDomainUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_scf_custom_domain.scf_custom_domain", "id"),
					resource.TestCheckResourceAttr("tencentcloud_scf_custom_domain.scf_custom_domain", "endpoints_config.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_scf_custom_domain.scf_custom_domain", "endpoints_config.0.path_match", "/aa/*"),
				),
			},
			{
				ResourceName:      "tencentcloud_scf_custom_domain.scf_custom_domain",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccScfFunctionForCustomDomain = `
resource "tencentcloud_scf_function" "foo" {
  name              = "%s"
  handler           = "first.do_it_first"
  runtime           = "Python3.6"
  enable_public_net = true
  async_run_enable  = "FALSE"

  dns_cache = false
  intranet_config {
    ip_fixed = "DISABLE"
  }

  zip_file = "%s"
  triggers {
    name         = "url-trigger"
    type         = "http"
    trigger_desc = "{\"AuthType\":\"NONE\",\"NetConfig\":{\"EnableIntranet\":true,\"EnableExtranet\":false}}"
  }
}
`

var testAccScfCustomDomain = scfFunctionCodeEmbed("first.zip", testAccScfFunctionForCustomDomain) + `

resource "tencentcloud_scf_custom_domain" "scf_custom_domain" {
  domain ="scf.iac-tf.cloud"
  protocol = "HTTP"
  endpoints_config {
    namespace = "default"
    function_name = tencentcloud_scf_function.foo.name
    qualifier = "$LATEST"
    path_match = "/a/*"
  }
  waf_config {
    waf_open = "CLOSE"
  }
}
`

var testAccScfCustomDomainUpdate = scfFunctionCodeEmbed("first.zip", testAccScfFunctionForCustomDomain) + `

resource "tencentcloud_scf_custom_domain" "scf_custom_domain" {
  domain ="scf.iac-tf.cloud"
  protocol = "HTTP"
  endpoints_config {
    namespace = "default"
    function_name = tencentcloud_scf_function.foo.name
    qualifier = "$LATEST"
    path_match = "/aa/*"
  }
  waf_config {
    waf_open = "CLOSE"
  }
}
`
