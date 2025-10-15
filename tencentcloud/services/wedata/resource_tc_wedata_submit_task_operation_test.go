package wedata_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWedataSubmitTaskOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccWedataSubmitTaskOperation,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_wedata_submit_task_operation.wedata_submit_task_operation", "id")),
		}},
	})
}

const testAccWedataSubmitTaskOperation = `

resource "tencentcloud_wedata_submit_task_operation" "wedata_submit_task_operation" {
	project_id = "2905622749543821312"
	task_id = "20251015164958429"
	version_remark = "v1"
}
`
