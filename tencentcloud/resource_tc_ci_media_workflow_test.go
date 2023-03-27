package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudCiMediaWorkflowResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaWorkflow,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_workflow.media_workflow", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_workflow.media_workflow",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaWorkflow = `

resource "tencentcloud_ci_media_workflow" "media_workflow" {
  name = &lt;nil&gt;
  workflow_id = &lt;nil&gt;
  state = &lt;nil&gt;
  topology {
		dependencies {
			key = &lt;nil&gt;
			value = &lt;nil&gt;
		}
		nodes {
			key = &lt;nil&gt;
			node {
				type = &lt;nil&gt;
				input {
					queue_id = &lt;nil&gt;
					object_prefix = &lt;nil&gt;
					notify_config {
						u_r_l = &lt;nil&gt;
						event = &lt;nil&gt;
						type = &lt;nil&gt;
						result_format = &lt;nil&gt;
					}
					ext_filter {
						state = &lt;nil&gt;
						audio = &lt;nil&gt;
						custom = &lt;nil&gt;
						custom_exts = &lt;nil&gt;
						all_file = &lt;nil&gt;
					}
				}
				operation {
					template_id = &lt;nil&gt;
				}
			}
		}

  }
  create_time = &lt;nil&gt;
  update_time = &lt;nil&gt;
  bucket_id = &lt;nil&gt;
}

`
