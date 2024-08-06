package tpulsar_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudTdmqNamespceRoleAttachmentResource_basic -v
func TestAccTencentCloudTdmqNamespceRoleAttachmentResource_basic(t *testing.T) {
	terraformId := "tencentcloud_tdmq_namespace_role_attachment.example"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTdmqNamespaceRoleAttachment,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(terraformId, "environ_id"),
					resource.TestCheckResourceAttrSet(terraformId, "role_name"),
					resource.TestCheckResourceAttrSet(terraformId, "permissions.#"),
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
				),
			},
			{
				Config: testAccTdmqNamespaceRoleAttachmentUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(terraformId, "environ_id"),
					resource.TestCheckResourceAttrSet(terraformId, "role_name"),
					resource.TestCheckResourceAttrSet(terraformId, "permissions.#"),
					resource.TestCheckResourceAttrSet(terraformId, "cluster_id"),
				),
			},
		},
	})
}

const testAccTdmqNamespaceRoleAttachment = `
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags         = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_namespace" "example" {
  environ_name = "tf_example"
  msg_ttl      = 300
  cluster_id   = tencentcloud_tdmq_instance.example.id
  retention_policy {
    time_in_minutes = 60
    size_in_mb      = 10
  }
  remark = "remark."
}

resource "tencentcloud_tdmq_role" "example" {
  role_name  = "tf_example"
  cluster_id = tencentcloud_tdmq_instance.example.id
  remark     = "remark."
}

resource "tencentcloud_tdmq_namespace_role_attachment" "example" {
  environ_id  = tencentcloud_tdmq_namespace.example.environ_name
  role_name   = tencentcloud_tdmq_role.example.role_name
  permissions = ["produce", "consume"]
  cluster_id  = tencentcloud_tdmq_instance.example.id
}
`

const testAccTdmqNamespaceRoleAttachmentUpdate = `
resource "tencentcloud_tdmq_instance" "example" {
  cluster_name = "tf_example"
  remark       = "remark."
  tags         = {
    "createdBy" = "terraform"
  }
}

resource "tencentcloud_tdmq_namespace" "example" {
  environ_name = "tf_example"
  msg_ttl      = 300
  cluster_id   = tencentcloud_tdmq_instance.example.id
  retention_policy {
    time_in_minutes = 60
    size_in_mb      = 10
  }
  remark = "remark."
}

resource "tencentcloud_tdmq_role" "example" {
  role_name  = "tf_example"
  cluster_id = tencentcloud_tdmq_instance.example.id
  remark     = "remark."
}

resource "tencentcloud_tdmq_namespace_role_attachment" "example" {
  environ_id  = tencentcloud_tdmq_namespace.example.environ_name
  role_name   = tencentcloud_tdmq_role.example.role_name
  permissions = ["produce"]
  cluster_id  = tencentcloud_tdmq_instance.example.id
}
`
