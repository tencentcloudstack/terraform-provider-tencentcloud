package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfMicroserviceDataSource_basic -v
func TestAccTencentCloudTsfMicroserviceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfMicroserviceDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_microservice.microservice"),
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
	default = "` + defaultNamespaceId + `"
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
