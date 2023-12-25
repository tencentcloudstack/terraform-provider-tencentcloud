package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfMsApiListDataSource_basic -v
func TestAccTencentCloudTsfMsApiListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfMsApiListDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_ms_api_list.ms_api_list"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_ms_api_list.ms_api_list", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_ms_api_list.ms_api_list", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_ms_api_list.ms_api_list", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_ms_api_list.ms_api_list", "result.0.content.0.method"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_ms_api_list.ms_api_list", "result.0.content.0.path"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_ms_api_list.ms_api_list", "result.0.content.0.status"),
				),
			},
		},
	})
}

const testAccTsfMsApiListDataSource = `

data "tencentcloud_tsf_ms_api_list" "ms_api_list" {
	microservice_id = "ms-yq3jo6jd"
	search_word = "echo"
}

`
