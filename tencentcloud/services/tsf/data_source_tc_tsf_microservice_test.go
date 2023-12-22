package tsf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfMicroserviceDataSource_basic -v
func TestAccTencentCloudTsfMicroserviceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_TSF) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfMicroserviceDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_microservice.microservice"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_microservice.microservice", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_microservice.microservice", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_microservice.microservice", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_microservice.microservice", "result.0.content.0.microservice_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_microservice.microservice", "result.0.content.0.microservice_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_microservice.microservice", "result.0.content.0.namespace_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_microservice.microservice", "result.0.content.0.run_instance_count"),
				),
			},
		},
	})
}

const testAccTsfMicroserviceDataSourceVar = `
variable "namespace_id" {
	default = "` + tcacctest.DefaultNamespaceId + `"
}
`

const testAccTsfMicroserviceDataSource = testAccTsfMicroserviceDataSourceVar + `

data "tencentcloud_tsf_microservice" "microservice" {
	namespace_id = var.namespace_id
	# status =
	microservice_id_list = ["ms-yq3jo6jd"]
	microservice_name_list = ["provider-demo"]
}

`
