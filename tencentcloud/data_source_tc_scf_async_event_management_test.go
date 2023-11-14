package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
  function_name = "test_function"
  namespace = "test_namespace"
  qualifier = "$LATEST"
  invoke_type = &lt;nil&gt;
  status = &lt;nil&gt;
  start_time_interval {
		start = &lt;nil&gt;
		end = &lt;nil&gt;

  }
  end_time_interval {
		start = "2020-02-02 04:03:03"
		end = "2020-02-02 05:03:03"

  }
  order = "ASC"
  orderby = "StartTime"
  offset = 0
  limit = 20
  invoke_request_id = "xxxxxxxx-xxxx-xxxx-xxxxxxxxxxxx"
  }

`
