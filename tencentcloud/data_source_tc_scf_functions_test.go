package tencentcloud

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudScfFunctions_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: scfFunctionCodeEmbed("main.py", TestAccDataSourceTencentCloudScfFunctions),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_functions.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_scf_functions.foo", "functions.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestMatchResourceAttr("data.tencentcloud_scf_functions.foo", "functions.0.name", regexp.MustCompile(`ci-test-function-basic`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.handler"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.mem_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.timeout"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.runtime"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.namespace"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.modify_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.code_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.err_no"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.install_dependency"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.eip_fixed"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.eips.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.l5_enable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.trigger_info.#"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudScfFunctions_namespace(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: scfFunctionCodeEmbed("main.py", TestAccDataSourceTencentCloudScfFunctionsNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_functions.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_functions.foo", "functions.#", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_functions.foo", "functions.0.name", "ci-test-function-ns"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_functions.foo", "functions.0.handler", "main.do_it"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_functions.foo", "functions.0.mem_size", "128"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_functions.foo", "functions.0.timeout", "3"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_functions.foo", "functions.0.runtime", "Python3.6"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.namespace"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.modify_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.code_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.err_no"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.install_dependency"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.eip_fixed"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.eips.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.l5_enable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.trigger_info.#"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudScfFunctions_Desc(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: scfFunctionCodeEmbed("main.py", TestAccDataSourceTencentCloudScfFunctionsDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_functions.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_scf_functions.foo", "functions.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_functions.foo", "functions.0.description", "test"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.handler"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.mem_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.timeout"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.runtime"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.namespace"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.modify_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.code_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.err_no"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.install_dependency"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.eip_fixed"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.eips.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.l5_enable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.trigger_info.#"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudScfFunctions_tag(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: scfFunctionCodeEmbed("main.py", TestAccDataSourceTencentCloudScfFunctionsTag),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_functions.foo"),
					resource.TestMatchResourceAttr("data.tencentcloud_scf_functions.foo", "functions.#", regexp.MustCompile(`^[1-9]\d*$`)),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.handler"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.mem_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.timeout"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.runtime"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.namespace"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.modify_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.code_size"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.err_no"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.install_dependency"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.eip_fixed"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.eips.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.l5_enable"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_functions.foo", "functions.0.trigger_info.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_functions.foo", "functions.0.tags.test", "test"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudScfFunctions = `
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function-basic"
  handler = "main.do_it"
  runtime = "Python3.6"

  zip_file = "%s"
}

data "tencentcloud_scf_functions" "foo" {
  name = tencentcloud_scf_function.foo.name
}
`

const TestAccDataSourceTencentCloudScfFunctionsNamespace = `
resource "tencentcloud_scf_namespace" "foo" {
  namespace = "ci-test-scf"
}

resource "tencentcloud_scf_function" "foo" {
  name      = "ci-test-function-ns"
  namespace = tencentcloud_scf_namespace.foo.id
  handler   = "main.do_it"
  runtime   = "Python3.6"

  zip_file = "%s"
}

data "tencentcloud_scf_functions" "foo" {
  namespace = tencentcloud_scf_function.foo.namespace
}
`

const TestAccDataSourceTencentCloudScfFunctionsDesc = `
resource "tencentcloud_scf_function" "foo" {
  name        = "ci-test-function-desc"
  handler     = "main.do_it"
  runtime     = "Python3.6"
  description = "test"

  zip_file = "%s"
}

data "tencentcloud_scf_functions" "foo" {
  description = tencentcloud_scf_function.foo.description
}
`

const TestAccDataSourceTencentCloudScfFunctionsTag = `
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function-tag"
  handler = "main.do_it"
  runtime = "Python3.6"

  zip_file = "%s"

  tags = {
    "test" = "test"
  }
}

data "tencentcloud_scf_functions" "foo" {
  tags = tencentcloud_scf_function.foo.tags
}
`
