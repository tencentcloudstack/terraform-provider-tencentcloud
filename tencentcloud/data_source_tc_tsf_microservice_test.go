package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfMicroserviceDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfMicroserviceDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_microservice.microservice")),
			},
		},
	})
}

const testAccTsfMicroserviceDataSource = `

data "tencentcloud_tsf_microservice" "microservice" {
  namespace_id = "ns-123456"
  search_word = ""
  order_by = ""
  order_type = 0
  status = 
  microservice_id_list = 
  microservice_name_list = 
  }

`
