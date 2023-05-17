package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testDataTCRRepositoriesNameAll = "data.tencentcloud_tcr_repositories.id_test"

func TestAccTencentCloudTcrRepositoriesData(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config:    testAccTencentCloudDataTCRRepositoriesBasic,
				PreConfig: func() { testAccStepSetRegion(t, "ap-shanghai") },
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(testDataTCRRepositoriesNameAll, "repository_list.0.name"),
					resource.TestCheckResourceAttrSet(testDataTCRRepositoriesNameAll, "repository_list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataTCRRepositoriesNameAll, "repository_list.0.url"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTCRRepositoriesBasic = defaultTCRInstanceData + `
data "tencentcloud_tcr_repositories" "id_test" {
  instance_id = local.tcr_id
  namespace_name = var.tcr_namespace
}
`
