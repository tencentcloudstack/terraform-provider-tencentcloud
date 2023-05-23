package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTsfRepositoryDataSource_basic -v
func TestAccTencentCloudTsfRepositoryDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_TSF) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfRepositoryDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_tsf_repository.repository"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_repository.repository", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_repository.repository", "result.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_repository.repository", "result.0.total_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_repository.repository", "result.0.content.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_repository.repository", "result.0.content.0.create_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_repository.repository", "result.0.content.0.is_used"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_repository.repository", "result.0.content.0.repository_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_tsf_repository.repository", "result.0.content.0.repository_type"),
				),
			},
		},
	})
}

const testAccTsfRepositoryDataSource = `

data "tencentcloud_tsf_repository" "repository" {
	search_word = "test"
	repository_type = "default"
}

`
