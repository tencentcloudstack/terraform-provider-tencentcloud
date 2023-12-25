package scf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudScfAsyncEventManagementDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfAsyncEventManagementDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_scf_async_event_management.async_event_management")),
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
