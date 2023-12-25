package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApiDetailDataSource_basic -v
func TestAccTencentCloudTsfApiDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApiDetailDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_api_detail.api_detail"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_detail.api_detail", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_detail.api_detail", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_detail.api_detail", "result.0.can_run"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_detail.api_detail", "result.0.definitions.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_detail.api_detail", "result.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_detail.api_detail", "result.0.request.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_detail.api_detail", "result.0.request.0.required"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_detail.api_detail", "result.0.response.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_detail.api_detail", "result.0.response.0.description"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_detail.api_detail", "result.0.response.0.name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_detail.api_detail", "result.0.response.0.type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_api_detail.api_detail", "result.0.status"),
				),
			},
		},
	})
}

const testAccTsfApiDetailDataSource = `

data "tencentcloud_tsf_api_detail" "api_detail" {
  microservice_id = "ms-yq3jo6jd"
  path = "/printRequest"
  method = "GET"
  pkg_version = "20210625192923"
  application_id = "application-a24x29xv"
}

`
