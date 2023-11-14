package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudTeoRule_engineResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoRule_engine,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_teo_rule_engine.rule_engine", "id")),
			},
			{
				ResourceName:      "tencentcloud_teo_rule_engine.rule_engine",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoRule_engine = `

resource "tencentcloud_teo_rule_engine" "rule_engine" {
  zone_id = ""
  rule_name = ""
  status = ""
  rules {
		actions {
			normal_action {
				action = ""
				parameters {
					name = ""
					values = 
				}
			}
			rewrite_action {
				action = ""
				parameters {
					action = ""
					name = ""
					values = 
				}
			}
			code_action {
				action = ""
				parameters {
					status_code = 
					name = ""
					values = 
				}
			}
		}
		conditions {
			conditions {
				operator = ""
				target = ""
				values = 
				ignore_case = 
				name = ""
				ignore_name_case = 
			}
		}
		sub_rules {
			rules {
				conditions {
					conditions {
						operator = ""
						target = ""
						values = 
						ignore_case = 
						name = ""
						ignore_name_case = 
					}
				}
				actions {
					normal_action {
						action = ""
						parameters {
							name = ""
							values = 
						}
					}
					rewrite_action {
						action = ""
						parameters {
							action = ""
							name = ""
							values = 
						}
					}
					code_action {
						action = ""
						parameters {
							status_code = 
							name = ""
							values = 
						}
					}
				}
			}
			tags = 
		}

  }
}

`
