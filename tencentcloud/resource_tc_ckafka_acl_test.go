package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTencentCloudCkafkaAcl(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCkafkaAclDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaAcl,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCkafkaAclExists("tencentcloud_ckafka_acl.foo"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl.foo", "resource_type", "TOPIC"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl.foo", "operation_type", "WRITE"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl.foo", "permission_type", "ALLOW"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_acl.foo", "host", "10.10.10.0"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_acl.foo", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_acl.foo", "resource_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_acl.foo", "principal"),
				),
			},
			{
				ResourceName:      "tencentcloud_ckafka_acl.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckCkafkaAclExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		ckafkaService := CkafkaService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("ckafka acl %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ckafka acl id is not set")
		}

		_, has, err := ckafkaService.DescribeAclByAclId(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("ckafka acl doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckCkafkaAclDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkaService := CkafkaService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ckafka_acl" {
			continue
		}

		_, has, err := ckafkaService.DescribeAclByAclId(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("ckafka acl still exists: %s", rs.Primary.ID)
	}
	return nil
}

const testAccCkafkaAcl = testAccCkafkaUser + `
resource "tencentcloud_ckafka_acl" foo {
  instance_id     = "ckafka-f9ife4zz"
  resource_type   = "TOPIC"
  resource_name   = "topic-tf-test"
  operation_type  = "WRITE"
  permission_type = "ALLOW"
  host            = "10.10.10.0"
  principal       = tencentcloud_ckafka_user.foo.account_name
}
`
