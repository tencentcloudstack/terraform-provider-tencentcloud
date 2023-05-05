package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfApplicationDataSource_basic -v
func TestAccTencentCloudTsfApplicationDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_application.application"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.0.content.0.application_desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.0.content.0.application_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.0.content.0.application_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.0.content.0.application_resource_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.0.content.0.application_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.0.content.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.0.content.0.ignore_create_image_repository"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.0.content.0.microservice_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.0.content.0.prog_lang"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_application.application", "result.0.content.0.update_time"),
				),
			},
		},
	})
}

const testAccTsfApplicationDataSourceVar = `
variable "application_id" {
	default = "` + defaultTsfApplicationId + `"
}
`

const testAccTsfApplicationDataSource = testAccTsfApplicationDataSourceVar + `

data "tencentcloud_tsf_application" "application" {
	application_type = "V"
	microservice_type = "N"
	# application_resource_type_list = [""]
	application_id_list = [var.application_id]
}

`
