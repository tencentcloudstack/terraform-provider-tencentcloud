package tcr_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataTCRRepositoriesNameAll = "data.tencentcloud_tcr_repositories.id_test"

func TestAccTencentCloudTcrRepositoriesData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckTCRRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTCRRepositoriesBasic,
				PreConfig: func() {
					tcacctest.AccStepSetRegion(t, "ap-shanghai")
					tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_COMMON)
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataTCRRepositoriesNameAll, "repository_list.0.name"),
					resource.TestCheckResourceAttrSet(testDataTCRRepositoriesNameAll, "repository_list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataTCRRepositoriesNameAll, "repository_list.0.url"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTCRRepositoriesBasic = tcacctest.DefaultTCRInstanceData + `
data "tencentcloud_tcr_repositories" "id_test" {
  instance_id = local.tcr_id
  namespace_name = var.tcr_namespace
}
`
