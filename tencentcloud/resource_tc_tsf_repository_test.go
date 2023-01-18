package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudTsfRepositoryResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfRepository,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_repository.repository", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_repository.repository",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfRepository = `

resource "tencentcloud_tsf_repository" "repository" {
    repository_name = ""
  repository_type = ""
  bucket_name = ""
  bucket_region = ""
  directory = ""
  repository_desc = ""
    }

`
