package waf_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudWafDomainPostActionConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccWafDomainPostAction,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_domain_post_action_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_domain_post_action_config.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_domain_post_action_config.example", "post_cls_action"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_domain_post_action_config.example", "post_ckafka_action"),
				),
			},
			{
				Config: testAccWafDomainPostActionUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_waf_domain_post_action_config.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_domain_post_action_config.example", "domain"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_domain_post_action_config.example", "post_cls_action"),
					resource.TestCheckResourceAttrSet("tencentcloud_waf_domain_post_action_config.example", "post_ckafka_action"),
				),
			},
			{
				ResourceName:      "tencentcloud_waf_domain_post_action_config.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccWafDomainPostAction = `
resource "tencentcloud_waf_domain_post_action_config" "example" {
  domain             = "example.com"
  post_cls_action    = 1
  post_ckafka_action = 0
}
`

const testAccWafDomainPostActionUpdate = `
resource "tencentcloud_waf_domain_post_action_config" "example" {
  domain             = "example.com"
  post_cls_action    = 0
  post_ckafka_action = 1
}
`
