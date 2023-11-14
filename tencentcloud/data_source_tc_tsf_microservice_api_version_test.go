package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfMicroserviceApiVersionDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfMicroserviceApiVersionDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_microservice_api_version.microservice_api_version")),
			},
		},
	})
}

const testAccTsfMicroserviceApiVersionDataSource = `

data "tencentcloud_tsf_microservice_api_version" "microservice_api_version" {
  microservice_id = "ms-yq3jo6jd"
  path = ""
  method = "get"
  }

`
