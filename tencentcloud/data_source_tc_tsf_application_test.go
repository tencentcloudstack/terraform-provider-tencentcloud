package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApplicationDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApplicationDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_application.application")),
			},
		},
	})
}

const testAccTsfApplicationDataSource = `

data "tencentcloud_tsf_application" "application" {
  search_word = "example"
  order_by = "name"
  order_type = 1
  application_type = "V"
  microservice_type = "M"
  application_resource_type_list = 
  application_id_list = 
  }

`
