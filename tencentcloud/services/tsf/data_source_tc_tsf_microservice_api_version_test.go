package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfMicroserviceApiVersionDataSource_basic -v
func TestAccTencentCloudTsfMicroserviceApiVersionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfMicroserviceApiVersionDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_microservice_api_version.microservice_api_version"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_microservice_api_version.microservice_api_version", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_microservice_api_version.microservice_api_version", "result.0.application_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_microservice_api_version.microservice_api_version", "result.0.application_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_microservice_api_version.microservice_api_version", "result.0.pkg_version"),
				),
			},
		},
	})
}

const testAccTsfMicroserviceApiVersionDataSource = `

data "tencentcloud_tsf_microservice_api_version" "microservice_api_version" {
  microservice_id = "ms-yq3jo6jd"
  path = "/printRequest"
  method = "get"
}

`
