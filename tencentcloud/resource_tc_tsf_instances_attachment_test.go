package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTsfInstancesAttachmentResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTsfInstancesAttachment,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_tsf_instances_attachment.instances_attachment", "id")),
			},
			{
				ResourceName:      "tencentcloud_tsf_instances_attachment.instances_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTsfInstancesAttachment = `

resource "tencentcloud_tsf_add_cluster_instances_operation" "add_cluster_instances_operation" {
	cluster_id = "cluster-123456"
	instance_id_list = [""]
	os_name = "Ubuntu 20.04"
	image_id = "img-123456"
	password = "MyP@ssw0rd"
	key_id = "key-123456"
	sg_id = "sg-123456"
	instance_import_mode = "R"
	os_customize_type = "my_customize"
	feature_id_list =
	instance_advanced_settings {
		  mount_target = "/mnt/data"
		  docker_graph_path = "/var/lib/docker"
	}
	security_group_ids = [""]
}

`
