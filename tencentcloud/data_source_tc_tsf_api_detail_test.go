package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfApiDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfApiDetailDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_api_detail.api_detail")),
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
