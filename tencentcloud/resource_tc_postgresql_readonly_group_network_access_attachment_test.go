package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

const testAccReadonlyGroupNetworkAccessAttachmentObject = "tencentcloud_postgresql_readonly_group_network_access_attachment.readonly_group_network_access_attachment"

func TestAccTencentCloudPostgresqlReadonlyGroupNetworkAccessAttachmentResource_auto_assign(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckReadonlyGroupNetworkAccessAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlReadonlyGroupNetworkAccessAttachment_auto_assign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckReadonlyGroupNetworkAccessAttachmentExists(testAccReadonlyGroupNetworkAccessAttachmentObject),
					resource.TestCheckResourceAttrSet(testAccReadonlyGroupNetworkAccessAttachmentObject, "db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccReadonlyGroupNetworkAccessAttachmentObject, "readonly_group_id"),
					resource.TestCheckResourceAttrSet(testAccReadonlyGroupNetworkAccessAttachmentObject, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccReadonlyGroupNetworkAccessAttachmentObject, "subnet_id"),
					resource.TestCheckResourceAttr(testAccReadonlyGroupNetworkAccessAttachmentObject, "is_assign_vip", "false"),
					resource.TestCheckResourceAttrSet(testAccReadonlyGroupNetworkAccessAttachmentObject, "vip"),
				),
			},
			{
				ResourceName:            testAccReadonlyGroupNetworkAccessAttachmentObject,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_assign_vip"},
			},
		},
	})
}

func TestAccTencentCloudPostgresqlReadonlyGroupNetworkAccessAttachmentResource_user_assign(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckReadonlyGroupNetworkAccessAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlReadonlyGroupNetworkAccessAttachment_user_assign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckReadonlyGroupNetworkAccessAttachmentExists(testAccReadonlyGroupNetworkAccessAttachmentObject),
					resource.TestCheckResourceAttrSet(testAccReadonlyGroupNetworkAccessAttachmentObject, "db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccReadonlyGroupNetworkAccessAttachmentObject, "readonly_group_id"),
					resource.TestCheckResourceAttrSet(testAccReadonlyGroupNetworkAccessAttachmentObject, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccReadonlyGroupNetworkAccessAttachmentObject, "subnet_id"),
					resource.TestCheckResourceAttr(testAccReadonlyGroupNetworkAccessAttachmentObject, "is_assign_vip", "true"),
					resource.TestCheckResourceAttrSet(testAccReadonlyGroupNetworkAccessAttachmentObject, "vip"),
				),
			},
			{
				ResourceName:            testAccReadonlyGroupNetworkAccessAttachmentObject,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_assign_vip"},
			},
		},
	})
}

func testAccCheckReadonlyGroupNetworkAccessAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := PostgresqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_postgresql_readonly_group_network_access_attachment" {
			continue
		}

		object, err := service.DescribePostgresqlReadonlyGroupNetInfoById(ctx, rs.Primary.ID)

		if err != nil {
			err, ok := err.(*sdkErrors.TencentCloudSDKError)
			if ok && err.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
				// it is ok
				return nil
			}
			return err
		}

		if object != nil {
			return fmt.Errorf("The tencentcloud_postgresql_readonly_group_network_access_attachment instance [%s] still exist.", rs.Primary.ID)
		}

	}
	return nil
}

func testAccCheckReadonlyGroupNetworkAccessAttachmentExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)
		service := PostgresqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("vpc attachment %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("vpc attachment id is not set")
		}

		object, err := service.DescribePostgresqlReadonlyGroupNetInfoById(ctx, rs.Primary.ID)

		if err != nil {
			err, ok := err.(*sdkErrors.TencentCloudSDKError)
			if ok && err.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
				// it is ok
				return nil
			}
			return err
		}

		if object == nil {
			return fmt.Errorf("The tencentcloud_postgresql_readonly_group_network_access_attachment instance [%s] not found.", rs.Primary.ID)
		}

		return nil
	}
}

const testAccPostgresqlInstanceForROGroupNetworkAccess = defaultAzVariable + defaultSecurityGroupData + CommonPresetPGSQL + `
resource "tencentcloud_vpc" "vpc" {
	cidr_block = "172.18.111.0/24"
	name       = "test-pg-network-vpc"
  }
  
  resource "tencentcloud_subnet" "subnet" {
	availability_zone = var.default_az
	cidr_block        = "172.18.111.0/24"
	name              = "test-pg-network-sub1"
	vpc_id            = tencentcloud_vpc.vpc.id
  }

  locals {
	my_vpc_id = tencentcloud_subnet.subnet.vpc_id
	my_subnet_id = tencentcloud_subnet.subnet.id
  }

  resource "tencentcloud_postgresql_readonly_group" "group" {
	master_db_instance_id = local.pgsql_id
	name = "tf_test_postgresql_readonly_group"
	project_id = 0
	vpc_id  	  		= local.my_vpc_id
	subnet_id 		= local.my_subnet_id
	replay_lag_eliminate = 1
	replay_latency_eliminate =  1
	max_replay_lag = 100
	max_replay_latency = 512
	min_delay_eliminate_reserve = 1
  }
`

const testAccPostgresqlReadonlyGroupNetworkAccessAttachment_auto_assign = testAccPostgresqlInstanceForROGroupNetworkAccess + `

resource "tencentcloud_postgresql_readonly_group_network_access_attachment" "readonly_group_network_access_attachment" {
  db_instance_id = local.pgsql_id
  readonly_group_id = tencentcloud_postgresql_readonly_group.group.id
  vpc_id = local.my_vpc_id
  subnet_id = local.my_subnet_id
  is_assign_vip = false
  tags = {
    "createdBy" = "terraform"
  }
}

`

const testAccPostgresqlReadonlyGroupNetworkAccessAttachment_user_assign = testAccPostgresqlInstanceForROGroupNetworkAccess + `

resource "tencentcloud_postgresql_readonly_group_network_access_attachment" "readonly_group_network_access_attachment" {
  db_instance_id = local.pgsql_id
  readonly_group_id = tencentcloud_postgresql_readonly_group.group.id
  vpc_id = local.my_vpc_id
  subnet_id = local.my_subnet_id
  is_assign_vip = true
  vip = "172.18.111.111"
  tags = {
    "createdBy" = "terraform"
  }
}

`
