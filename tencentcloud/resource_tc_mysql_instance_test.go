package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

const TestAccTencentCloudMysqlMasterInstance_availability_zone = "ap-guangzhou-3"

func TestAccTencentCloudMysqlMasterInstance_fullslave(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlMasterInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlMasterInstance_fullslave(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "instance_name", "testAccMysql"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_deploy_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_sync_mode", "2"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "availability_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "first_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "second_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
				),
			},
		},
	})
}

func TestAccTencentCloudMysqlMasterInstance_basic_and_update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlMasterInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlMasterInstance_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "instance_name", "testAccMysql"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "mem_size", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "volume_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "intranet_port", "3360"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "engine_version", "5.7"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "internet_service", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_deploy_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_sync_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "project_id", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "availability_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "first_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),

					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "intranet_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "task_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "gtid"),
				),
			},
			// add tag
			{
				Config: testAccMysqlMasterInstance_multiTags("master"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "tags.role", "master"),
				),
			},
			// update tag
			{
				Config: testAccMysqlMasterInstance_multiTags("master-version2"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "tags.role", "master-version2"),
				),
			},
			// remove tag
			{
				Config: testAccMysqlMasterInstance_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckNoResourceAttr("tencentcloud_mysql_instance.mysql_master", "tags.role"),
				),
			},

			// open internet service
			{
				Config: testAccMysqlMasterInstance_internet_service(true),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "internet_service", "1"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "internet_host"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "internet_port"),
				),
			},

			//close internet  service
			{
				Config: testAccMysqlMasterInstance_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "internet_service", "0")),
			},

			//modify  parameters
			{
				Config: testAccMysqlMasterInstance_parameters(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "parameters.max_connections", "1000")),
			},
			//remove parameters and  restore
			{
				Config: testAccMysqlMasterInstance_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master")),
			},
			// update instance_name
			{
				Config: testAccMysqlMasterInstance_update("testAccMysql-version1", "3360"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "instance_name", "testAccMysql-version1"),
				),
			},
			// update intranet_port
			{
				Config: testAccMysqlMasterInstance_update("testAccMysql-version1", "3361"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "intranet_port", "3361"),
				),
			},
		},
	})
}

func TestAccTencentCloudMysqlPrepaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlMasterInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlMasterInstance_prepaid(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.prepaid", "instance_name", "testAccMysqlPrepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.prepaid", "slave_deploy_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.prepaid", "slave_sync_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.prepaid", "availability_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.prepaid", "first_slave_zone", TestAccTencentCloudMysqlMasterInstance_availability_zone),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.prepaid", "auto_renew_flag", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.prepaid", "charge_type", "PREPAID"),
				),
			},
			// update auto_renew_flag
			{
				Config: testAccMysqlMasterInstance_prepaidupdate(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.prepaid", "auto_renew_flag", "1"),
				),
			},
		},
	})
}

func testAccCheckMysqlMasterInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	mysqlService := MysqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mysql_instance" {
			continue
		}
		instance, err := mysqlService.DescribeRunningDBInstanceById(ctx, rs.Primary.ID)
		if instance != nil {
			return fmt.Errorf("mysql instance still exist")
		}
		if err != nil {
			sdkErr, ok := err.(*errors.TencentCloudSDKError)
			if ok && sdkErr.Code == MysqlInstanceIdNotFound {
				continue
			}
			return err
		}
	}
	return nil
}

func testAccCheckMysqlMasterInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("mysql instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("mysql instance id is not set")
		}

		mysqlService := MysqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
		instance, err := mysqlService.DescribeDBInstanceById(ctx, rs.Primary.ID)
		if instance == nil {
			return fmt.Errorf("mysql instance %s is not found", rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccMysqlMasterInstance_basic() string {
	return `
resource "tencentcloud_mysql_instance" "mysql_master" {
  pay_type          = 1
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
}`
}

func testAccMysqlMasterInstance_prepaid() string {
	return `
resource "tencentcloud_mysql_instance" "prepaid" {
  charge_type 		= "PREPAID"
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysqlPrepaid"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  force_delete = false
}`
}

func testAccMysqlMasterInstance_prepaidupdate() string {
	return `
resource "tencentcloud_mysql_instance" "prepaid" {
  charge_type 		= "PREPAID"
  mem_size          = 1000
  volume_size       = 50
  auto_renew_flag	= 1
  instance_name     = "testAccMysqlPrepaid"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  force_delete		= false
}`
}

func testAccMysqlMasterInstance_fullslave() string {
	return `
resource "tencentcloud_mysql_instance" "mysql_master" {
  pay_type          = 1
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  slave_deploy_mode = 0
  first_slave_zone  = "ap-guangzhou-3"
  second_slave_zone = "ap-guangzhou-3"
  slave_sync_mode   = 2
}`
}

func testAccMysqlMasterInstance_internet_service(open bool) string {
	tag := "0"
	if open {
		tag = "1"
	}
	return `
resource "tencentcloud_mysql_instance" "mysql_master" {
  pay_type          = 1
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  internet_service  = ` + tag + `
}`

}

func testAccMysqlMasterInstance_parameters() string {
	return `
resource "tencentcloud_mysql_instance" "mysql_master" {
  pay_type          = 1
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  
  parameters = {
    max_connections = "1000"
  }
}`
}

func testAccMysqlMasterInstance_multiTags(value string) string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_instance" "mysql_master" {
  pay_type          = 1
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "testAccMysql"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = 3360
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
  tags = {
    test = "test-tf"
    role = "%s"
  }
}
	`, value)
}

func testAccMysqlMasterInstance_update(instance_name, instranet_port string) string {
	tpl := `
resource "tencentcloud_mysql_instance" "mysql_master" {
  pay_type          = 1
  mem_size          = 1000
  volume_size       = 50
  instance_name     = "%s"
  engine_version    = "5.7"
  root_password     = "test1234"
  intranet_port     = %s
  availability_zone = "ap-guangzhou-3"
  first_slave_zone  = "ap-guangzhou-3"
}`
	return fmt.Sprintf(tpl, instance_name, instranet_port)
}
