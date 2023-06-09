package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudScfRequestStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfRequestStatusDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_request_status.request_status")),
			},
		},
	})
}

const testAccScfRequestStatusDataSource = `

data "tencentcloud_scf_request_status" "request_status" {
  function_name       = "keep-1676351130"
  function_request_id = "9de9405a-e33a-498d-bb59-e80b7bed1191"
  namespace           = "default"
}

`
