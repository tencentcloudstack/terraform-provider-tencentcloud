package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

const testAccPgInstanceNetworkAccessAttachmentObject = "tencentcloud_postgresql_instance_network_access_attachment.instance_network_access_attachment"

func TestAccTencentCloudPostgresqlInstanceNetworkAccessAttachmentResource_autoassign(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPgInstanceNetworkAccessAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlInstanceNetworkAccessAttachment_auto_assign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPgInstanceNetworkAccessAttachmentExists(testAccPgInstanceNetworkAccessAttachmentObject),
					resource.TestCheckResourceAttrSet(testAccPgInstanceNetworkAccessAttachmentObject, "db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccPgInstanceNetworkAccessAttachmentObject, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPgInstanceNetworkAccessAttachmentObject, "subnet_id"),
					resource.TestCheckResourceAttr(testAccPgInstanceNetworkAccessAttachmentObject, "is_assign_vip", "false"),
					resource.TestCheckResourceAttrSet(testAccPgInstanceNetworkAccessAttachmentObject, "vip"),
				),
			},
			{
				ResourceName:            testAccPgInstanceNetworkAccessAttachmentObject,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_assign_vip"},
			},
		},
	})
}

func TestAccTencentCloudPostgresqlInstanceNetworkAccessAttachmentResource_userassign(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPgInstanceNetworkAccessAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlInstanceNetworkAccessAttachment_user_assign,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPgInstanceNetworkAccessAttachmentExists(testAccPgInstanceNetworkAccessAttachmentObject),
					resource.TestCheckResourceAttrSet(testAccPgInstanceNetworkAccessAttachmentObject, "db_instance_id"),
					resource.TestCheckResourceAttrSet(testAccPgInstanceNetworkAccessAttachmentObject, "vpc_id"),
					resource.TestCheckResourceAttrSet(testAccPgInstanceNetworkAccessAttachmentObject, "subnet_id"),
					resource.TestCheckResourceAttr(testAccPgInstanceNetworkAccessAttachmentObject, "is_assign_vip", "true"),
					resource.TestCheckResourceAttrSet(testAccPgInstanceNetworkAccessAttachmentObject, "vip"),
				),
			},
			{
				ResourceName:            testAccPgInstanceNetworkAccessAttachmentObject,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"is_assign_vip"},
			},
		},
	})
}

func testAccCheckPgInstanceNetworkAccessAttachmentDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := PostgresqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_postgresql_instance_network_access_attachment" {
			continue
		}

		object, err := service.DescribePostgresqlDBInstanceNetInfoById(ctx, rs.Primary.ID)

		if err != nil {
			err, ok := err.(*sdkErrors.TencentCloudSDKError)
			if ok && err.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
				// it is ok
				return nil
			}
			return err
		}

		if object != nil {
			return fmt.Errorf("The tencentcloud_postgresql_instance_network_access_attachment instance [%s] still exist.", rs.Primary.ID)
		}

	}
	return nil
}

func testAccCheckPgInstanceNetworkAccessAttachmentExists(n string) resource.TestCheckFunc {
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

		object, err := service.DescribePostgresqlDBInstanceNetInfoById(ctx, rs.Primary.ID)

		if err != nil {
			err, ok := err.(*sdkErrors.TencentCloudSDKError)
			if ok && err.GetCode() == "ResourceNotFound.InstanceNotFoundError" {
				// it is ok
				return nil
			}
			return err
		}

		if object == nil {
			return fmt.Errorf("The tencentcloud_postgresql_instance_network_access_attachment instance [%s] not found.", rs.Primary.ID)
		}

		return nil
	}
}

const testAccPostgresqlInstanceNetworkAccess = defaultAzVariable + defaultSecurityGroupData + `
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

resource "tencentcloud_postgresql_instance" "test" {
	name 				= "tf_postsql_instance_network_access_attachment"
	availability_zone =  var.default_az
	charge_type 		= "POSTPAID_BY_HOUR"
	vpc_id  	  		= local.my_vpc_id
	subnet_id 		= local.my_subnet_id
	engine_version	= "13.3"
	root_password	    = "t1qaA2k1wgvfa3?ZZZ"
	security_groups   = [local.sg_id]
	charset			= "LATIN1"
	project_id 		= 0
	memory 			= 4
	storage 			= 20
  
	db_kernel_version = "v13.3_r1.1"
  
	tags = {
	  tf = "test"
	}
  }
`

const testAccPostgresqlInstanceNetworkAccessAttachment_auto_assign = testAccPostgresqlInstanceNetworkAccess + `

resource "tencentcloud_postgresql_instance_network_access_attachment" "instance_network_access_attachment" {
  db_instance_id = tencentcloud_postgresql_instance.test.id
  vpc_id = local.my_vpc_id
  subnet_id = local.my_subnet_id
  is_assign_vip = false
  tags = {
    "createdBy" = "terraform"
  }
}

`

const testAccPostgresqlInstanceNetworkAccessAttachment_user_assign = testAccPostgresqlInstanceNetworkAccess + `

resource "tencentcloud_postgresql_instance_network_access_attachment" "instance_network_access_attachment" {
  db_instance_id = tencentcloud_postgresql_instance.test.id
  vpc_id = local.my_vpc_id
  subnet_id = local.my_subnet_id
  is_assign_vip = true
  vip = "172.18.111.111"
  tags = {
    "createdBy" = "terraform"
  }
}

`
