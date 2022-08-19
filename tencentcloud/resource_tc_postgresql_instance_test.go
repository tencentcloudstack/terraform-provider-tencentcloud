package tencentcloud

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

var testPostgresqlInstanceResourceName = "tencentcloud_postgresql_instance"
var testPostgresqlInstanceResourceKey = testPostgresqlInstanceResourceName + ".test"

func init() {
	resource.AddTestSweepers(testPostgresqlInstanceResourceName, &resource.Sweeper{
		Name: testPostgresqlInstanceResourceName,
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			postgresqlService := PostgresqlService{client: client}
			vpcService := VpcService{client: client}

			instances, err := postgresqlService.DescribePostgresqlInstances(ctx, nil)
			if err != nil {
				return err
			}

			var vpcs []string

			for _, v := range instances {
				id := *v.DBInstanceId
				name := *v.DBInstanceName
				vpcId := *v.VpcId

				now := time.Now()
				createTime := stringTotime(*v.CreateTime)
				interval := now.Sub(createTime).Minutes()

				if strings.HasPrefix(name, keepResource) || strings.HasPrefix(name, defaultResource) {
					continue
				}

				// less than 30 minute, not delete
				if needProtect == 1 && int64(interval) < 30 {
					continue
				}
				// isolate
				err := postgresqlService.IsolatePostgresqlInstance(ctx, id)
				if err != nil {
					continue
				}
				// describe status
				err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
					instance, has, err := postgresqlService.DescribePostgresqlInstanceById(ctx, id)
					if err != nil {
						return retryError(err)
					}
					if !has {
						return resource.NonRetryableError(fmt.Errorf("instance %s removed", id))
					}
					if *instance.DBInstanceStatus != "isolated" {
						return resource.RetryableError(fmt.Errorf("waiting for instance isolated, now is %s", *instance.DBInstanceStatus))
					}
					return nil
				})
				if err != nil {
					continue
				}
				// delete
				err = postgresqlService.DeletePostgresqlInstance(ctx, id)
				if err != nil {
					continue
				}
				vpcs = append(vpcs, vpcId)
			}
			for _, v := range vpcs {
				_ = vpcService.DeleteVpc(ctx, v)
			}

			return nil
		},
	})
}

func TestAccTencentCloudPostgresqlInstanceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPostgresqlInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlInstanceExists(testPostgresqlInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "name", "tf_postsql_instance"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "storage", "100"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "public_access_switch", "false"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "root_password", "t1qaA2k1wgvfa3?ZZZ"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "private_access_port"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "db_major_vesion"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "db_major_version"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "backup_plan.#", "1"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "backup_plan.0.min_backup_start_time", "00:10:11"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "backup_plan.0.max_backup_start_time", "01:10:11"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "backup_plan.0.backup_period.#", "2"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "backup_plan.0.base_backup_retention_period", "7"),
					//resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "tags.tf", "test"),
				),
			},
			{
				ResourceName:            testPostgresqlInstanceResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "spec_code", "public_access_switch", "charset", "backup_plan"},
			},
			{
				Config: testAccPostgresqlInstanceOpenPublic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlInstanceExists(testPostgresqlInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "public_access_switch", "true"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "private_access_ip"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "private_access_port"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "public_access_host"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "public_access_port"),
					//resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "tags.tf", "teest"),
				),
			},
			{
				Config: testAccPostgresqlInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlInstanceExists(testPostgresqlInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "name", "tf_postsql_instance_update"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "charge_type", "POSTPAID_BY_HOUR"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "vpc_id"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "subnet_id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "memory", "4"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "storage", "250"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "create_time"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "project_id", "0"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "public_access_switch", "false"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "root_password", "t1qaA2k1wgvfa3?ZZZ"),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "availability_zone"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "backup_plan.0.min_backup_start_time", "01:10:11"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "backup_plan.0.max_backup_start_time", "02:10:11"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "backup_plan.0.backup_period.#", "3"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "backup_plan.0.base_backup_retention_period", "5"),
					//resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "tags.tf", "teest"),
				),
			},
		},
	})
}

func TestAccTencentCloudPostgresqlMAZInstanceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckPostgresqlInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPostgresqlMAZInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlInstanceExists(testPostgresqlInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "id"),
					// SDK 1.0 cannot provide set test expected "db_node_set.*.role" , "Primary"
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "db_node_set.#", "2"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "availability_zone", "ap-guangzhou-6"),
				),
			},
			{
				ResourceName:            testPostgresqlInstanceResourceKey,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"root_password", "spec_code", "public_access_switch", "charset"},
			},

			{
				Config: testAccPostgresqlMAZInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPostgresqlInstanceExists(testPostgresqlInstanceResourceKey),
					resource.TestCheckResourceAttrSet(testPostgresqlInstanceResourceKey, "id"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "db_node_set.#", "2"),
					resource.TestCheckResourceAttr(testPostgresqlInstanceResourceKey, "availability_zone", "ap-guangzhou-6"),
				),
			},
		},
	})
}

func testAccCheckPostgresqlInstanceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != testPostgresqlInstanceResourceName {
			continue
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := PostgresqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribePostgresqlInstanceById(ctx, rs.Primary.ID)

		if !has {
			return nil
		} else {
			if err != nil {
				return err
			}
			return fmt.Errorf("delete postgresql instance %s fail", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckPostgresqlInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s is not found", n)
		}
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		service := PostgresqlService{client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn}

		_, has, err := service.DescribePostgresqlInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			_, has, err = service.DescribePostgresqlInstanceById(ctx, rs.Primary.ID)
		}
		if err != nil {
			return err
		}
		if has {
			return nil
		} else {
			return fmt.Errorf("postgresql instance %s is not found", rs.Primary.ID)
		}
	}
}

const testAccPostgresqlInstanceBasic = `
data "tencentcloud_availability_zones_by_product" "zone" {
  product = "postgres"
}
`

const testAccPostgresqlInstance string = testAccPostgresqlInstanceBasic + defaultVpcSubnets + `
resource "tencentcloud_postgresql_instance" "test" {
  name 				= "tf_postsql_instance"
  availability_zone = data.tencentcloud_availability_zones_by_product.zone.zones[5].name
  charge_type 		= "POSTPAID_BY_HOUR"
  vpc_id  	  		= local.vpc_id
  subnet_id 		= local.subnet_id
  engine_version	= "10.4"
  root_password	    = "t1qaA2k1wgvfa3?ZZZ"
  charset			= "LATIN1"
  project_id 		= 0
  memory 			= 4
  storage 			= 100

  backup_plan {
	min_backup_start_time = "00:10:11"
	max_backup_start_time = "01:10:11"
	base_backup_retention_period = 7
	backup_period = ["tuesday", "wednesday"]
  }

  tags = {
	tf = "test"
  }
}
`

const testAccPostgresqlInstanceOpenPublic string = testAccPostgresqlInstanceBasic + defaultVpcSubnets + `
resource "tencentcloud_postgresql_instance" "test" {
  name = "tf_postsql_instance_update"
  availability_zone = data.tencentcloud_availability_zones_by_product.zone.zones[5].name
  charge_type	    = "POSTPAID_BY_HOUR"
  vpc_id  	  		= local.vpc_id
  subnet_id 		= local.subnet_id
  engine_version	= "10.4"
  root_password	    = "t1qaA2k1wgvfa3?ZZZ"
  charset 			= "LATIN1"
  project_id 		= 0
  public_access_switch = true
  memory 			= 4
  storage 			= 250
  backup_plan {
	min_backup_start_time 		 = "01:10:11"
	max_backup_start_time		 = "02:10:11"
	base_backup_retention_period = 5
	backup_period 			     = ["monday", "thursday", "sunday"]
  }

  tags = {
	tf = "teest"
  }
}
`

const testAccPostgresqlInstanceUpdate string = testAccPostgresqlInstanceBasic + defaultVpcSubnets + `
resource "tencentcloud_postgresql_instance" "test" {
  name = "tf_postsql_instance_update"
  availability_zone = data.tencentcloud_availability_zones_by_product.zone.zones[5].name
  charge_type	    = "POSTPAID_BY_HOUR"
  vpc_id  	  		= local.vpc_id
  subnet_id 		= local.subnet_id
  engine_version	= "10.4"
  root_password	    = "t1qaA2k1wgvfa3?ZZZ"
  charset 			= "LATIN1"
  project_id 		= 0
  public_access_switch = false
  memory 			= 4
  storage 			= 250
  backup_plan {
	min_backup_start_time 		 = "01:10:11"
	max_backup_start_time		 = "02:10:11"
	base_backup_retention_period = 5
	backup_period 			     = ["monday", "thursday", "sunday"]
  }

  tags = {
	tf = "teest"
  }
}
`

const testAccPostgresqlMAZInstance string = `
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/24"
  name       = "test-pg-vpc"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.0.0/24"
  name              = "pg-sub1"
  vpc_id            = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_postgresql_instance" "test" {
  name = "tf_postsql_maz_instance"
  availability_zone = "ap-guangzhou-6"
  charge_type = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  engine_version		= "10.4"
  root_password                 = "t1qaA2k1wgvfa3?ZZZ"
  charset = "LATIN1"
  memory = 4
  storage = 100
  db_node_set {
    role = "Primary"
    zone = "ap-guangzhou-6"
  }
  db_node_set {
    zone = "ap-guangzhou-7"
  }
}
`

const testAccPostgresqlMAZInstanceUpdate string = `
resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/24"
  name       = "test-pg-vpc"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.0.0/24"
  name              = "pg-sub1"
  vpc_id            = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_postgresql_instance" "test" {
  name = "tf_postsql_maz_instance"
  availability_zone = "ap-guangzhou-6"
  charge_type = "POSTPAID_BY_HOUR"
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
  engine_version		= "10.4"
  root_password                 = "t1qaA2k1wgvfa3?ZZZ"
  charset = "LATIN1"
  memory = 4
  storage = 250
  db_node_set {
    role = "Primary"
    zone = "ap-guangzhou-6"
  }
  db_node_set {
    zone = "ap-guangzhou-6"
  }
}
`
