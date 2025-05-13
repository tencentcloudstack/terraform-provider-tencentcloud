package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafInstanceAttackLogPostResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafInstanceAttackLogPost,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_instance_attack_log_post.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_instance_attack_log_post.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_instance_attack_log_post.example", "attack_log_post"),
				),
			},
			{
				Config: testAccWafInstanceAttackLogPostUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_instance_attack_log_post.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_instance_attack_log_post.example", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_instance_attack_log_post.example", "attack_log_post"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_instance_attack_log_post.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafInstanceAttackLogPost = `
resource "tencentcloud_waf_instance_attack_log_post" "example" {
  instance_id     = "waf_2kxtlbky11b4wcrb"
  attack_log_post = 1
}
`

const testAccWafInstanceAttackLogPostUpdate = `
resource "tencentcloud_waf_instance_attack_log_post" "example" {
  instance_id     = "waf_2kxtlbky11b4wcrb"
  attack_log_post = 0
}
`
