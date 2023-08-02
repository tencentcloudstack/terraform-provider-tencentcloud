package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_mongodb_instance
	resource.AddTestSweepers("tencentcloud_mongodb_instance", &resource.Sweeper{
		Name: "tencentcloud_mongodb_instance",
		F: func(r string) error {
			logId := getLogId(contextNil)
			ctx := context.WithValue(context.TODO(), logIdKey, logId)
			cli, _ := sharedClientForRegion(r)
			client := cli.(*TencentCloudClient).apiV3Conn
			service := MongodbService{client}

			instances, err := service.DescribeInstancesByFilter(ctx, "", -1)
			if err != nil {
				return err
			}

			var isolated []string

			for i := range instances {
				ins := instances[i]
				id := *ins.InstanceId
				name := *ins.InstanceName

				if strings.HasPrefix(name, keepResource) || strings.HasPrefix(name, defaultResource) {
					continue
				}

				created, err := time.Parse("2006-01-02 15:04:05", *ins.CreateTime)
				if err != nil {
					created = time.Time{}
				}
				if isResourcePersist(name, &created) {
					continue
				}
				log.Printf("%s (%s) will Isolated", id, name)
				err = service.IsolateInstance(ctx, id)
				if err != nil {
					continue
				}
				isolated = append(isolated, id)
			}

			log.Printf("Offline isolated instance %v", isolated)
			for _, id := range isolated {
				err = service.OfflineIsolatedDBInstance(ctx, id, true)
				if err != nil {
					continue
				}
			}

			return nil
		},
	})
}

func TestAccTencentCloudMongodbInstanceResourcePostPaid(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_PREPAY) },
				Config:    testAccMongodbInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_instance.mongodb"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "instance_name", "tf-mongodb-test"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "memory"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "volume"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "engine_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "machine_type"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "available_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.test", "test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "charge_type", MONGODB_CHARGE_TYPE_POSTPAID),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_instance.mongodb", "prepaid_period"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_instance.mongodb",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_groups", "password", "auto_renew_flag"},
			},
			{
				SkipFunc: func() (bool, error) {
					log.Printf("[WARN] MongoDB Update Need DealID query available, skip checking.")
					return true, nil
				},
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_PREPAY) },
				Config:    testAccMongodbInstance_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "instance_name", "tf-mongodb-update"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "memory", "8"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "volume", "200"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.abc", "abc"),
				),
			},
		},
	})
}

func TestAccTencentCloudMongodbInstanceResource_multiZone(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckMongodbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_PREPAY) },
				Config:    testAccMongodbInstance_multiZone,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_instance.mongodb_mutil_zone"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_mutil_zone", "node_num", "5"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_mutil_zone", "availability_zone_list.#", "5"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_mutil_zone", "hidden_zone", "ap-guangzhou-6"),
				),
			},
		},
	})
}

func TestAccTencentCloudMongodbInstanceResourcePrepaid(t *testing.T) {
	// Avoid to set Parallel to make sure EnvVar secure
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_PREPAY) },
				Config:    testAccMongodbInstancePrepaid,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckMongodbInstanceExists("tencentcloud_mongodb_instance.mongodb_prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "instance_name", "tf-mongodb-test-prepaid"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "memory"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "volume"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "engine_version"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "machine_type"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "available_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "project_id", "0"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "status"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "vip"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "vport"),
					resource.TestCheckResourceAttrSet("tencentcloud_mongodb_instance.mongodb_prepaid", "create_time"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "tags.test", "test-prepaid"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "charge_type", MONGODB_CHARGE_TYPE_PREPAID),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "prepaid_period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "auto_renew_flag", "1"),
				),
			},
			{
				PreConfig: func() { testAccStepPreConfigSetTempAKSK(t, ACCOUNT_TYPE_PREPAY) },
				Config:    testAccMongodbInstancePrepaid_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "instance_name", "tf-mongodb-test-prepaid-update"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "tags.prepaid", "prepaid"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_instance.mongodb_prepaid",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_groups", "password", "auto_renew_flag", "prepaid_period"},
			},
		},
	})
}

func testAccCheckMongodbInstanceDestroy(s *terraform.State) error {
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	mongodbService := MongodbService{
		client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_mongodb_instance" {
			continue
		}

		_, has, err := mongodbService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("mongodb instance still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckMongodbInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := getLogId(contextNil)
		ctx := context.WithValue(context.TODO(), logIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("mongodb instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("mongodb instance id is not set")
		}
		mongodbService := MongodbService{
			client: testAccProvider.Meta().(*TencentCloudClient).apiV3Conn,
		}
		_, has, err := mongodbService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("mongodb instance doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const testAccMongodbInstance = DefaultMongoDBSpec + `
resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-test"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"
  vpc_id         = var.vpc_id
  subnet_id      = var.subnet_id

  tags = {
    test = "test"
  }
}
`

const testAccMongodbInstance_update = DefaultMongoDBSpec + `
resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-update"
  memory         = local.memory * 2
  volume         = local.volume * 2
  engine_version = local.engine_version
  machine_type   = local.machine_type
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234update"
  vpc_id         = var.vpc_id
  subnet_id      = var.subnet_id

  tags = {
    abc = "abc"
  }
}
`

const testAccMongodbInstancePrepaid = DefaultMongoDBSpec + `
resource "tencentcloud_mongodb_instance" "mongodb_prepaid" {
  instance_name   = "tf-mongodb-test-prepaid"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234"
  charge_type     = "PREPAID"
  prepaid_period  = 1
  auto_renew_flag = 1
  vpc_id         = var.vpc_id
  subnet_id      = var.subnet_id

  tags = {
    test = "test-prepaid"
  }
}
`

const testAccMongodbInstancePrepaid_update = DefaultMongoDBSpec + `
resource "tencentcloud_mongodb_instance" "mongodb_prepaid" {
  instance_name   = "tf-mongodb-test-prepaid-update"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234update"
  charge_type     = "PREPAID"
  prepaid_period  = 1
  auto_renew_flag = 1
  vpc_id         = var.vpc_id
  subnet_id      = var.subnet_id

  tags = {
    prepaid = "prepaid"
  }
}
`

const testAccMongodbInstance_multiZone = DefaultMongoDBSpec + `
resource "tencentcloud_mongodb_instance" "mongodb_mutil_zone" {
  instance_name   = "mongodb-mutil-zone-test"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"
  vpc_id         = var.vpc_id
  subnet_id      = var.subnet_id
  node_num = 5
  availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-4", "ap-guangzhou-6"]
  hidden_zone = "ap-guangzhou-6"
  tags = {
    test = "test"
  }
}
`
