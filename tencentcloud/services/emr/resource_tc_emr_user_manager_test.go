package emr_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixEmrUserManagerResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEmrUserManager,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_emr_user_manager.user_manager", "id"),
					resource.TestCheckResourceAttr("tencentcloud_emr_user_manager.user_manager", "user_name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_emr_user_manager.user_manager", "user_group", "group1"),
				),
			},
			{
				Config: testAccEmrUserManagerUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_emr_user_manager.user_manager", "id"),
					resource.TestCheckResourceAttr("tencentcloud_emr_user_manager.user_manager", "user_name", "tf-test"),
					resource.TestCheckResourceAttr("tencentcloud_emr_user_manager.user_manager", "user_group", "group1"),
				),
			},
			{
				ResourceName:            "tencentcloud_emr_user_manager.user_manager",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"password"},
				ImportStateVerify:       true,
			},
		},
	})
}

const testAccEmrUserManager = `

data "tencentcloud_emr" "my_emr" {
  display_strategy = "clusterList"
}

resource "tencentcloud_emr_user_manager" "user_manager" {
  instance_id = data.tencentcloud_emr.my_emr.clusters.0.cluster_id
  user_name   = "tf-test"
  user_group  = "group1"
  password    = "tf@123456"
}


`

const testAccEmrUserManagerUpdate = `

data "tencentcloud_emr" "my_emr" {
  display_strategy = "clusterList"
}

resource "tencentcloud_emr_user_manager" "user_manager" {
  instance_id = data.tencentcloud_emr.my_emr.clusters.0.cluster_id
  user_name   = "tf-test"
  user_group  = "group1"
  password    = "tf@12345678"
}


`
