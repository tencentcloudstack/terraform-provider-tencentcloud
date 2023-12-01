package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixScfAsyncEventStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfAsyncEventStatusDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_scf_async_event_status.async_event_status")),
			},
		},
	})
}

const testAccScfAsyncEventStatusDataSource = `

data "tencentcloud_scf_async_event_status" "async_event_status" {
  invoke_request_id = "9de9405a-e33a-498d-bb59-e80b7bed1191"
}

`
