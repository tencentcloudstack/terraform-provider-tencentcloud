package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudScfAsyncEventManagementDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfAsyncEventManagementDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_async_event_management.async_event_management")),
			},
		},
	})
}

const testAccScfAsyncEventManagementDataSource = `

data "tencentcloud_scf_async_event_management" "async_event_management" {
  function_name = "keep-1676351130"
  namespace     = "default"
  qualifier     = "$LATEST"
  order   = "ASC"
  orderby = "StartTime"
}

`
