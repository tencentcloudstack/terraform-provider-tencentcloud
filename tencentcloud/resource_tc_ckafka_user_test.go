package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_kafka
	resource.AddTestSweepers("tencentcloud_kafka", &resource.Sweeper{
		Name: "tencentcloud_kafka",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			sharedClient, err := sharedClientForRegion(r)
			if err != nil {
				return fmt.Errorf("getting tencentcloud client error: %s", err.Error())
			}
			client := sharedClient.(*TencentCloudClient)

			ckafkaService := CkafkaService{client: client.apiV3Conn}
			params := make(map[string]interface{})
			params["instance_id"] = defaultKafkaInstanceId
			userInfos, err := ckafkaService.DescribeUserByFilter(ctx, params)
			if err != nil {
				return nil
			}
			for _, userInfo := range userInfos {
				userName := *userInfo.Name
				now := time.Now()
				createTime := stringTotime(*userInfo.CreateTime)
				interval := now.Sub(createTime).Minutes()
				// less than 30 minute, not delete
				if needProtect == 1 && int64(interval) < 30 {
					continue
				}

				if strings.HasPrefix(userName, keepResource) || strings.HasPrefix(userName, defaultResource) {
					continue
				}
				userIdStr := fmt.Sprintf("%v#%v", defaultKafkaInstanceId, userName)
				err := ckafkaService.DeleteUser(ctx, userIdStr)
				if err != nil {
					return nil
				}
			}
			return nil
		},
	})
}

func TestAccTencentCloudCkafkaUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCkafkaUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCkafkaUser,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCkafkaUserExists("tencentcloud_ckafka_user.foo"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_user.foo", "account_name", "tf-test"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_user.foo", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_user.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_user.foo", "update_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_user.foo", "password"),
				),
			},
			{
				Config: testAccCkafkaUser_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCkafkaUserExists("tencentcloud_ckafka_user.foo"),
					resource.TestCheckResourceAttr("tencentcloud_ckafka_user.foo", "account_name", "tf-test"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_user.foo", "instance_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_user.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_user.foo", "update_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_ckafka_user.foo", "password"),
				),
			},
			{
				ResourceName:            "tencentcloud_ckafka_user.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password"},
			},
		},
	})
}

func testAccCheckCkafkaUserExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		ckafkaService := CkafkaService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("ckafka user %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("ckafka user id is not set")
		}

		_, has, err := ckafkaService.DescribeUserByUserId(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("ckafka user doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckCkafkaUserDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ckafkaService := CkafkaService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_ckafka_user" {
			continue
		}

		_, has, err := ckafkaService.DescribeUserByUserId(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("ckafka user still exists: %s", rs.Primary.ID)
	}
	return nil
}

const testAccCkafkaUser = defaultKafkaVariable + `
resource "tencentcloud_ckafka_user" "foo" {
  instance_id  = var.instance_id
  account_name = "tf-test"
  password     = "test1234"
}
`

const testAccCkafkaUser_update = defaultKafkaVariable + `
resource "tencentcloud_ckafka_user" "foo" {
  instance_id  = var.instance_id
  account_name = "tf-test"
  password     = "test1234update"
}
`
