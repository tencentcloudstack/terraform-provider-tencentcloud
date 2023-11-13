package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfDescribeMsApiListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfDescribeMsApiListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_describe_ms_api_list.describe_ms_api_list")),
			},
		},
	})
}

const testAccTsfDescribeMsApiListDataSource = `

data "tencentcloud_tsf_describe_ms_api_list" "describe_ms_api_list" {
  microservice_id = "ms-xxxxxxxx"
  search_word = ""
  }

`
