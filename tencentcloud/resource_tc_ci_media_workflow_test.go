package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
					output {
						region = &lt;nil&gt;
						bucket = &lt;nil&gt;
						object = &lt;nil&gt;
						au_object = &lt;nil&gt;
						sprite_object = &lt;nil&gt;
					}
					watermark_template_id = &lt;nil&gt;
					delogo_param {
						switch = &lt;nil&gt;
						dx = &lt;nil&gt;
						dy = &lt;nil&gt;
						width = &lt;nil&gt;
						height = &lt;nil&gt;
					}
					s_d_rto_h_d_r {
						hdr_mode = &lt;nil&gt;
					}
					s_c_f {
						region = &lt;nil&gt;
						function_name = &lt;nil&gt;
						namespace = &lt;nil&gt;
					}
					hls_pack_info {
						video_stream_config {
							video_stream_name = &lt;nil&gt;
							band_width = &lt;nil&gt;
						}
					}
					transcode_template_id = &lt;nil&gt;
					smart_cover {
						format = &lt;nil&gt;
						width = &lt;nil&gt;
						height = &lt;nil&gt;
						count = &lt;nil&gt;
						delete_duplicates = &lt;nil&gt;
					}
					segment_config {
						format = &lt;nil&gt;
						duration = &lt;nil&gt;
					}
					digital_watermark {
						message = &lt;nil&gt;
						type = &lt;nil&gt;
						version = &lt;nil&gt;
					}
					stream_pack_config_info {
						pack_type = &lt;nil&gt;
						ignore_failed_stream = &lt;nil&gt;
						reserve_all_stream_node = &lt;nil&gt;
					}
					stream_pack_info {
						video_stream_config {
							video_stream_name = &lt;nil&gt;
							band_width = &lt;nil&gt;
						}
					}
				}
			}
		}

  }
      bucket_id = &lt;nil&gt;
}

`
