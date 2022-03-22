package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

func TestAccTencentCloudMysqlReadonlyInstance(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlReadonlyInstance(CommonPresetMysql),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "instance_name", "mysql-readonly-test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "mem_size", "2000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "volume_size", "200"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_port", "3360"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "task_status"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.test", "test-tf"),
				),
			},
			// add tag
			{
				Config: testAccMysqlReadonlyInstance_multiTags(CommonPresetMysql, "read"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.role", "read"),
				),
			},
			// update tag
			{
				Config: testAccMysqlReadonlyInstance_multiTags(CommonPresetMysql, "readonly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.role", "readonly"),
				),
			},
			// remove tag
			{
				Config: testAccMysqlReadonlyInstance(CommonPresetMysql),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckNoResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.role"),
				),
			},
			// update instance_name
			{
				Config: testAccMysqlReadonlyInstance_update(CommonPresetMysql, "mysql-readonly-update", "3360"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "instance_name", "mysql-readonly-update"),
				),
			},
			// // update intranet_port
			// {
			// 	Config: testAccMysqlReadonlyInstance_update(CommonPresetMysql, "mysql-readonly-update", "3361"),
			// 	Check: resource.ComposeAggregateTestCheckFunc(
			// 		testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
			// 		resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_port", "3361"),
			// 	),
			// },
		},
	})
}

func testAccCheckMysqlReadonlyInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mysqlService := MysqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mysql_readonly_instance" {
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

func testAccCheckMysqlInstanceExists(n string) resource.TestCheckFunc {
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

func testAccMysqlReadonlyInstance(mysqlTestCase string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
  master_instance_id = local.mysql_id
  mem_size           = 2000
  volume_size        = 200
  instance_name      = "mysql-readonly-test"
  intranet_port      = 3360
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
  }
}
	`, mysqlTestCase)
}

func testAccMysqlReadonlyInstance_multiTags(mysqlTestCase, value string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
  master_instance_id = local.mysql_id
  mem_size           = 2000
  cpu                = 1
  volume_size        = 200
  instance_name      = "mysql-readonly-test"
  intranet_port      = 3360
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
    role = "%s"
  }
}
	`, mysqlTestCase, value)
}

func testAccMysqlReadonlyInstance_update(mysqlTestCase, instance_name, instranet_port string) string {
	return fmt.Sprintf(`
%s
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
  master_instance_id = local.mysql_id
  mem_size           = 2000
  cpu                = 1
  volume_size        = 200
  instance_name      = "%s"
  intranet_port      = %s 
  master_region = var.region
  zone = var.availability_zone
  tags = {
    test = "test-tf"
  }
}
	`, mysqlTestCase, instance_name, instranet_port)
}
