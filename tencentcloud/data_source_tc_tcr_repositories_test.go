package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

var testDataTCRRepositoriesNameAll = "data.tencentcloud_tcr_repositories.id_test"

func TestAccTencentCloudDataTCRRepositories(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTCRRepositoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataTCRRepositoriesBasic,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckTCRRepositoryExists("tencentcloud_tcr_repository.mytcr_repository"),
					resource.TestCheckResourceAttr(testDataTCRRepositoriesNameAll, "repository_list.0.name", "test"),
					resource.TestCheckResourceAttr(testDataTCRRepositoriesNameAll, "repository_list.0.brief_desc", "2222"),
					resource.TestCheckResourceAttr(testDataTCRRepositoriesNameAll, "repository_list.0.description", "211111111111111111111111111111111111"),
					resource.TestCheckResourceAttrSet(testDataTCRRepositoriesNameAll, "repository_list.0.create_time"),
					resource.TestCheckResourceAttrSet(testDataTCRRepositoriesNameAll, "repository_list.0.url"),
				),
			},
		},
	})
}

const testAccTencentCloudDataTCRRepositoriesBasic = `
resource "tencentcloud_tcr_instance" "mytcr_instance" {
  name        = "testacctcrinstance"
  instance_type = "standard"
  delete_bucket = true

  tags ={
	test = "test"
  }
}

resource "tencentcloud_tcr_namespace" "mytcr_namespace" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  name        = "test"
  is_public   = false
}

resource "tencentcloud_tcr_repository" "mytcr_repository" {
  instance_id = tencentcloud_tcr_instance.mytcr_instance.id
  namespace_name        = tencentcloud_tcr_namespace.mytcr_namespace.name
  name = "test"
  brief_desc = "2222"
  description = "211111111111111111111111111111111111"
}

data "tencentcloud_tcr_repositories" "id_test" {
  instance_id = tencentcloud_tcr_repository.mytcr_repository.instance_id
  namespace_name = tencentcloud_tcr_namespace.mytcr_namespace.name
}
`
