package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudScfLogs_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: scfFunctionCodeEmbed("main.py", TestAccDataSourceTencentCloudScfLogs),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_logs.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_logs.foo", "function_name", "ci-test-function"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_logs.foo", "logs.#"),
				),
			},
		},
	})
}

func TestAccDataSourceTencentCloudScfLogs_allArgWithoutReqId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: scfFunctionCodeEmbed("main.py", TestAccDataSourceTencentCloudScfLogsAllArgWithoutReqId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_logs.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_logs.foo", "function_name", "ci-test-function"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_logs.foo", "offset", "0"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_logs.foo", "limit", "100"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_logs.foo", "order", "desc"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_logs.foo", "order_by", "duration"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_logs.foo", "ret_code", "UserCodeException"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_logs.foo", "start_time", "2017-05-16 20:00:00"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_logs.foo", "end_time", "2017-05-17 20:00:00"),
					resource.TestCheckResourceAttr("data.tencentcloud_scf_logs.foo", "offset", "0"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_scf_logs.foo", "logs.#"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudScfLogs = `
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  zip_file = "%s"
}

data "tencentcloud_scf_logs" "foo" {
  function_name = tencentcloud_scf_function.foo.name
}
`

const TestAccDataSourceTencentCloudScfLogsAllArgWithoutReqId = `
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  zip_file = "%s"
}

data "tencentcloud_scf_logs" "foo" {
  function_name = tencentcloud_scf_function.foo.name
  offset        = 0
  limit         = 100
  order         = "desc"
  order_by      = "duration"
  ret_code      = "UserCodeException"
  start_time    = "2017-05-16 20:00:00"
  end_time      = "2017-05-17 20:00:00"
}
`
