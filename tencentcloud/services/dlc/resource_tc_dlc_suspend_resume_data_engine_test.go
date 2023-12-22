package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcSuspendResumeDataEngineResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcSuspendResumeDataEngineSuspend,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_suspend_resume_data_engine.suspend_resume_data_engine", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_suspend_resume_data_engine.suspend_resume_data_engine", "data_engine_name", "iac-test"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_suspend_resume_data_engine.suspend_resume_data_engine", "operate", "suspend")),
			},
			{
				Config: testAccDlcSuspendResumeDataEngineResume,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_suspend_resume_data_engine.suspend_resume_data_engine", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_suspend_resume_data_engine.suspend_resume_data_engine", "data_engine_name", "iac-test"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_suspend_resume_data_engine.suspend_resume_data_engine", "operate", "resume")),
			},
		},
	})
}

const testAccDlcSuspendResumeDataEngineSuspend = `

resource "tencentcloud_dlc_suspend_resume_data_engine" "suspend_resume_data_engine" {
  data_engine_name = "iac-test"
  operate = "suspend"
}

`
const testAccDlcSuspendResumeDataEngineResume = `

resource "tencentcloud_dlc_suspend_resume_data_engine" "suspend_resume_data_engine" {
  data_engine_name = "iac-test"
  operate = "resume"
}
`
