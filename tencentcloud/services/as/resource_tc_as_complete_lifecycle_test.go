package as_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixAsCompleteLifecycleResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAsCompleteLifecycle,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_as_complete_lifecycle.complete_lifecycle", "id")),
			},
		},
	})
}

const testAccAsCompleteLifecycle = `

resource "tencentcloud_as_complete_lifecycle" "complete_lifecycle" {
  lifecycle_hook_id = "ash-xxxxxxxx"
  lifecycle_action_result = "CONTINUE"
  instance_id = "ins-xxxxxxxx"
}
`
