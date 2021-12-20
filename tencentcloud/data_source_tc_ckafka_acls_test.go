package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDataSourceCkafkaAcls(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCkafkaAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudDataSourceCkafkaAcl,
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCkafkaAclExists("tencentcloud_ckafka_acl.foo"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_acls.foo", "acl_list.0.operation_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_acls.foo", "acl_list.0.permission_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_acls.foo", "acl_list.0.resource_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_acls.foo", "acl_list.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_acls.foo", "acl_list.0.host"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ckafka_acls.foo", "acl_list.0.principal"),
				),
			},
		},
	})
}

const testAccTencentCloudDataSourceCkafkaAcl = testAccCkafkaAcl + `
data "tencentcloud_ckafka_acls" "foo" {
	instance_id   = tencentcloud_ckafka_acl.foo.instance_id
    resource_type = tencentcloud_ckafka_acl.foo.resource_type
	resource_name = tencentcloud_ckafka_acl.foo.resource_name
}
`
