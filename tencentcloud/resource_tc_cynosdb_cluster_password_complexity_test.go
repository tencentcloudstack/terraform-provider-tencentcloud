package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbClusterPasswordComplexityResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbClusterPasswordComplexity,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity", "id")),
			},
			{
				ResourceName:      "tencentcloud_cynosdb_cluster_password_complexity.cluster_password_complexity",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCynosdbClusterPasswordComplexity = `

resource "tencentcloud_cynosdb_cluster_password_complexity" "cluster_password_complexity" {
  cluster_id = "cynosdbpg-bzxxrmtq"
  validate_password_length = 8
  validate_password_mixed_case_count = 1
  validate_password_special_char_count = 1
  validate_password_number_count = 1
  validate_password_policy = "MEDIUM"
  validate_password_dictionary = 
}

`
