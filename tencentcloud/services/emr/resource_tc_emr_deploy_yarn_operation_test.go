package emr_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudEmrDeployYarnOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: emrDeployYarn,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_emr_deploy_yarn_operation.emr_yarn", "id"),
				),
			},
		},
	})
}

const emrDeployYarn = `
resource "tencentcloud_emr_deploy_yarn_operation" "emr_yarn" {
  instance_id = "emr-rzrochgp"
}
`
