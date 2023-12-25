package scf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixScfAsyncEventStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfAsyncEventStatusDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_scf_async_event_status.async_event_status")),
			},
		},
	})
}

const testAccScfAsyncEventStatusDataSource = `

data "tencentcloud_scf_async_event_status" "async_event_status" {
  invoke_request_id = "9de9405a-e33a-498d-bb59-e80b7bed1191"
}

`
