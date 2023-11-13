package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
  function_name = ""
  function_request_id = ""
  namespace = ""
  start_time = ""
  end_time = ""
  }

`
