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

const TestAccTencentCloudMysqlMasterInstance_availability_zone = "ap-guangzhou-4"
const TestAccTencentCloudMysqlMasterInstance_vpc_id = "vpc-fzdzrsir"
const TestAccTencentCloudMysqlMasterInstance_subnet_id = "subnet-he8ldxx6"
const TestAccTencentCloudMysqlMasterInstance_security_groups = "sg-eeks0kgz"
const TestAccTencentCloudMysqlMasterInstance_project_id = "1144023"

func TestAccTencentCloudMysqlMasterInstance_basic(t *testing.T) {
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
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "pay_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "mem_size", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "volume_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "vpc_id", TestAccTencentCloudMysqlMasterInstance_vpc_id),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "subnet_id", TestAccTencentCloudMysqlMasterInstance_subnet_id),
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
		},
	})
}

func TestAccTencentCloudMysqlMasterInstance_securityGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMysqlMasterInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMysqlMasterInstance_securityGroup(),
				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "security_groups."+
						fmt.Sprintf("%d", hashcode.String(TestAccTencentCloudMysqlMasterInstance_security_groups)),
						TestAccTencentCloudMysqlMasterInstance_security_groups),
				),
			},
		},
	})
}

func TestAccTencentCloudMysqlMasterInstance_fullslave(t *testing.T) {
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

func TestAccTencentCloudMysqlMasterInstance_update(t *testing.T) {
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
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "pay_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "mem_size", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "volume_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "vpc_id", "vpc-fzdzrsir"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "subnet_id", "subnet-he8ldxx6"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "intranet_port", "3360"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "engine_version", "5.7"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "internet_service", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_deploy_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_sync_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "availability_zone"),

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
			//modify  project_id
			{
				Config: testAccMysqlMasterInstance_project_id(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "project_id",
						TestAccTencentCloudMysqlMasterInstance_project_id)),
			},
			//modify project_id to default
			{
				Config: testAccMysqlMasterInstance_basic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "project_id",
						"0")),
			},
		},
	})
}

func TestAccTencentCloudMysqlMasterInstance_updateMemAndSize(t *testing.T) {
	InitLogConfig(true)
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
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "pay_type", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "mem_size", "1000"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "volume_size", "50"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "vpc_id", "vpc-fzdzrsir"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "subnet_id", "subnet-he8ldxx6"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "intranet_port", "3360"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "engine_version", "5.7"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "internet_service", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_deploy_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "slave_sync_mode", "0"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "availability_zone"),

					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "intranet_ip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "task_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mysql_instance.mysql_master", "gtid"),
				),
			},
			// update mem_size and volume_size
			{
				Config: testAccMysqlMasterInstance_update("2000", "100", "testAccMysql", "3360"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "mem_size", "2000"),
					//resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "volume_size", "100"),
				),
			},
			// update instance_name
			{
				Config: testAccMysqlMasterInstance_update("2000", "100", "testAccMysql-version1", "3360"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "instance_name", "testAccMysql-version1"),
				),
			},
			// update intranet_port
			{
				Config: testAccMysqlMasterInstance_update("2000", "100", "testAccMysql-version1", "3361"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckMysqlMasterInstanceExists("tencentcloud_mysql_instance.mysql_master"),
					resource.TestCheckResourceAttr("tencentcloud_mysql_instance.mysql_master", "intranet_port", "3361"),
				),
			},
		}})
}

func testAccCheckMysqlMasterInstanceDestroy(s *terraform.State) error {
	logId := GetLogId(nil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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

func testAccMysqlMasterInstance_basic() string {
	return `
	resource "tencentcloud_mysql_instance" "mysql_master" {
	pay_type = 1
	mem_size = 1000
	volume_size = 50
	instance_name = "testAccMysql"
	vpc_id = "` + TestAccTencentCloudMysqlMasterInstance_vpc_id + `"
	subnet_id = "` + TestAccTencentCloudMysqlMasterInstance_subnet_id + `"
	engine_version = "5.7"
	root_password = "test1234"
	intranet_port =3360
	availability_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	first_slave_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
}`
}

func testAccMysqlMasterInstance_project_id() string {
	return `
	resource "tencentcloud_mysql_instance" "mysql_master" {
	pay_type = 1
	mem_size = 1000
	volume_size = 50
	instance_name = "testAccMysql"
	vpc_id = "` + TestAccTencentCloudMysqlMasterInstance_vpc_id + `"
	subnet_id = "` + TestAccTencentCloudMysqlMasterInstance_subnet_id + `"
	engine_version = "5.7"
	root_password = "test1234"
	intranet_port =3360
	availability_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	first_slave_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	project_id = ` + TestAccTencentCloudMysqlMasterInstance_project_id + `
}`

}

func testAccMysqlMasterInstance_fullslave() string {
	return `
	resource "tencentcloud_mysql_instance" "mysql_master" {
	pay_type = 1
	mem_size = 1000
	volume_size = 50
	instance_name = "testAccMysql"
	vpc_id = "` + TestAccTencentCloudMysqlMasterInstance_vpc_id + `"
	subnet_id = "` + TestAccTencentCloudMysqlMasterInstance_subnet_id + `"
	engine_version = "5.7"
	root_password = "test1234"
	intranet_port =3360
	availability_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	slave_deploy_mode = 0
	first_slave_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	second_slave_zone  = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	slave_sync_mode     = 2
}`
}

func testAccMysqlMasterInstance_internet_service(open bool) string {
	tag := "0"
	if open {
		tag = "1"
	}
	return `
	resource "tencentcloud_mysql_instance" "mysql_master" {
	pay_type = 1
	mem_size = 1000
	volume_size = 50
	instance_name = "testAccMysql"
	vpc_id = "` + TestAccTencentCloudMysqlMasterInstance_vpc_id + `"
	subnet_id = "` + TestAccTencentCloudMysqlMasterInstance_subnet_id + `"
	engine_version = "5.7"
	root_password = "test1234"
	intranet_port =3360
	availability_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	first_slave_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	internet_service = ` + tag + `
}`

}

func testAccMysqlMasterInstance_prePaid() string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_instance" "mysql_master" {
	pay_type = 0
	period = 1
	auto_renew_flag = 0
    instance_name = "testAccMysql"
	vpc_id = "` + TestAccTencentCloudMysqlMasterInstance_vpc_id + `"
	subnet_id = "` + TestAccTencentCloudMysqlMasterInstance_subnet_id + `"
	engine_version = "5.7"
	root_password = "test1234"
	intranet_port =3360
	availability_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	first_slave_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
}
	`)
}

func testAccMysqlMasterInstance_securityGroup() string {
	return fmt.Sprintf(`resource "tencentcloud_mysql_instance" "mysql_master" {
	pay_type = 1
	mem_size = 1000
	volume_size = 50
	instance_name = "testAccMysql"
	vpc_id = "` + TestAccTencentCloudMysqlMasterInstance_vpc_id + `"
	subnet_id =  "` + TestAccTencentCloudMysqlMasterInstance_subnet_id + `"
	engine_version = "5.7"
	root_password = "test1234"
	intranet_port =3360
	availability_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	first_slave_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	security_groups = ["sg-eeks0kgz"]
}
	`)
}

func testAccMysqlMasterInstance_parameters() string {
	return `
	resource "tencentcloud_mysql_instance" "mysql_master" {
	pay_type = 1
	mem_size = 1000
	volume_size = 50
	instance_name = "testAccMysql"
	vpc_id = "` + TestAccTencentCloudMysqlMasterInstance_vpc_id + `"
	subnet_id = "` + TestAccTencentCloudMysqlMasterInstance_subnet_id + `"
	engine_version = "5.7"
	root_password = "test1234"
	intranet_port =3360
	availability_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	first_slave_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	parameters = {
		max_connections = "1000"
	}
}`
}

func testAccMysqlMasterInstance_multiTags(value string) string {
	return fmt.Sprintf(`
resource "tencentcloud_mysql_instance" "mysql_master" {
	pay_type = 1
	mem_size = 1000
	volume_size = 50
	instance_name = "testAccMysql"
	vpc_id ="`+TestAccTencentCloudMysqlMasterInstance_vpc_id+`"
	subnet_id =  "`+TestAccTencentCloudMysqlMasterInstance_subnet_id+`"
	engine_version = "5.7"
	root_password = "test1234"
	intranet_port =3360
	availability_zone = "`+TestAccTencentCloudMysqlMasterInstance_availability_zone+`"
	first_slave_zone = "`+TestAccTencentCloudMysqlMasterInstance_availability_zone+`"
	tags = {
		test = "test-tf"
		role = "%s"
	}
}
	`, value)
}

func testAccMysqlMasterInstance_update(mem_size, volume_size, instance_name, instranet_port string) string {

	tpl := `
	resource "tencentcloud_mysql_instance" "mysql_master" {
	pay_type = 1
	mem_size = %s
	volume_size = %s
	instance_name =  "%s"
	vpc_id = "` + TestAccTencentCloudMysqlMasterInstance_vpc_id + `"
	subnet_id = "` + TestAccTencentCloudMysqlMasterInstance_subnet_id + `"
	engine_version = "5.7"
	root_password = "test1234"
	intranet_port = %s
	availability_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
	first_slave_zone = "` + TestAccTencentCloudMysqlMasterInstance_availability_zone + `"
}`
	return fmt.Sprintf(tpl, mem_size, volume_size, instance_name, instranet_port)
}
