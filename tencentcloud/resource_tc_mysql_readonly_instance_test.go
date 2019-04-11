package tencentcloud

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
)

const masterMysqlId = "cdb-d4on2aox"

func TestAccTencentCloudMysqlReadonlyInstance(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlReadonlyInstance(MysqlInstanceCommonTestCase),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "instance_name", "mysql-readonly-test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "pay_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "mem_size", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "volume_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "vpc_id", "vpc-fzdzrsir"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "subnet_id", "subnet-he8ldxx6"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_port", "3360"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "task_status"),
				),
			},
		},
	})
}

func TestAccTencentCloudMysqlReadonlyInstance_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlReadonlyInstance_basic(masterMysqlId),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "instance_name", "mysql-readonly-test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "pay_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "mem_size", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "volume_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "vpc_id", "vpc-fzdzrsir"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "subnet_id", "subnet-he8ldxx6"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_port", "3360"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "task_status"),
				),
			},
		},
	})
}

func TestAccTencentCloudMysqlReadonlyInstance_securityGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlReadonlyInstance_securityGroup(masterMysqlId),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "instance_name", "mysql-readonly-test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "pay_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "mem_size", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "volume_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "vpc_id", "vpc-fzdzrsir"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "subnet_id", "subnet-he8ldxx6"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_port", "3360"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "security_groups."+fmt.Sprintf("%d", hashcode.String("sg-eeks0kgz")), "sg-eeks0kgz"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "task_status"),
				),
			},
		},
	})
}

func TestAccTencentCloudMysqlReadonlyInstance_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlReadonlyInstance_tags(masterMysqlId),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "instance_name", "mysql-readonly-test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "pay_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "mem_size", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "volume_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "vpc_id", "vpc-fzdzrsir"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "subnet_id", "subnet-he8ldxx6"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_port", "3360"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "task_status"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.test", "test-tf"),
				),
			},
			// add tag
			{
				Config: testAccMysqlReadonlyInstance_multiTags(masterMysqlId, "read"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.role", "read"),
				),
			},
			// update tag
			{
				Config: testAccMysqlReadonlyInstance_multiTags(masterMysqlId, "readonly"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.role", "readonly"),
				),
			},
			// remove tag
			{
				Config: testAccMysqlReadonlyInstance_tags(masterMysqlId),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckNoResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "tags.role"),
				),
			},
		},
	})
}

func TestAccTencentCloudMysqlReadonlyInstance_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlReadonlyInstance_basic(masterMysqlId),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "instance_name", "mysql-readonly-test"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "pay_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "mem_size", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "volume_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "vpc_id", "vpc-fzdzrsir"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "subnet_id", "subnet-he8ldxx6"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_port", "3360"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_readonly_instance.mysql_readonly", "task_status"),
				),
			},
			// update mem_size and volume_size
			{
				Config: testAccMysqlReadonlyInstance_update(masterMysqlId, "2000", "100", "mysql-readonly-test", "vpc-fzdzrsir",
					"subnet-he8ldxx6", "3360"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "mem_size", "2000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "volume_size", "100"),
				),
			},
			// update instance_name
			{
				Config: testAccMysqlReadonlyInstance_update(masterMysqlId, "2000", "100", "mysql-readonly-update", "vpc-fzdzrsir",
					"subnet-he8ldxx6", "3360"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "instance_name", "mysql-readonly-update"),
				),
			},
			// update intranet_port
			{
				Config: testAccMysqlReadonlyInstance_update(masterMysqlId, "2000", "100", "mysql-readonly-update", "vpc-fzdzrsir",
					"subnet-he8ldxx6", "3361"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlInstanceExists("tencentcloud_mysql_readonly_instance.mysql_readonly"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_readonly_instance.mysql_readonly", "intranet_port", "3361"),
				),
			},
		},
	})
}

func testAccCheckMysqlReadonlyInstanceDestroy(s *terraform.State) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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
		logId := GetLogId(nil)
		ctx := context.WithValue(context.TODO(), "logId", logId)

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
    master_instance_id = "${tencentcloud_mysql_instance.default.id}"
    pay_type = 1
    mem_size = 1000
    volume_size = 50
    instance_name = "mysql-readonly-test"
    vpc_id = "vpc-fzdzrsir"
    subnet_id = "subnet-he8ldxx6"
	intranet_port = 3360
}
	`, mysqlTestCase)
}

func testAccMysqlReadonlyInstance_basic(mysqlId string) string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
    master_instance_id = "%s"
    pay_type = 1
    mem_size = 1000
    volume_size = 50
    instance_name = "mysql-readonly-test"
    vpc_id = "vpc-fzdzrsir"
    subnet_id = "subnet-he8ldxx6"
	intranet_port = 3360
}
	`, mysqlId)
}

func testAccMysqlReadonlyInstance_prePaid(mysqlId string) string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
    master_instance_id = "%s"
	pay_type = 0
	period = 1
	auto_renew_flag = 0
    mem_size = 1000
    volume_size = 50
    instance_name = "mysql-readonly-test"
    vpc_id = "vpc-fzdzrsir"
    subnet_id = "subnet-he8ldxx6"
	intranet_port = 3360
}
	`, mysqlId)
}

func testAccMysqlReadonlyInstance_securityGroup(mysqlId string) string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
    master_instance_id = "%s"
	pay_type = 1
    mem_size = 1000
    volume_size = 50
    instance_name = "mysql-readonly-test"
    vpc_id = "vpc-fzdzrsir"
    subnet_id = "subnet-he8ldxx6"
    intranet_port = 3360
	security_groups = ["sg-eeks0kgz"]
}
	`, mysqlId)
}

func testAccMysqlReadonlyInstance_tags(mysqlId string) string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
    master_instance_id = "%s"
	pay_type = 1
    mem_size = 1000
    volume_size = 50
    instance_name = "mysql-readonly-test"
    vpc_id = "vpc-fzdzrsir"
    subnet_id = "subnet-he8ldxx6"
	intranet_port = 3360
	tags = {
		test = "test-tf"
	}
}
	`, mysqlId)
}

func testAccMysqlReadonlyInstance_multiTags(mysqlId, value string) string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
    master_instance_id = "%s"
	pay_type = 1
    mem_size = 1000
    volume_size = 50
    instance_name = "mysql-readonly-test"
    vpc_id = "vpc-fzdzrsir"
    subnet_id = "subnet-he8ldxx6"
	intranet_port = 3360
	tags = {
		test = "test-tf"
		role = "%s"
	}
}
	`, mysqlId, value)
}

func testAccMysqlReadonlyInstance_update(mysqlId, mem_size, volume_size, instance_name, vpc_id, subnet_id, instranet_port string) string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_readonly_instance" "mysql_readonly" {
    master_instance_id = "%s"
	pay_type = 1
    mem_size = %s
    volume_size = %s
    instance_name = "%s"
    vpc_id = "%s"
    subnet_id = "%s"
    intranet_port = %s
}
	`, mysqlId, mem_size, volume_size, instance_name, vpc_id, subnet_id, instranet_port)
}
