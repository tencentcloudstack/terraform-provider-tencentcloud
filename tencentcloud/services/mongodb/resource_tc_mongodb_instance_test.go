package mongodb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svcmongodb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"context"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	mongodb_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func init() {
	// go test -v ./tencentcloud -sweep=ap-guangzhou -sweep-run=tencentcloud_mongodb_instance
	resource.AddTestSweepers("tencentcloud_mongodb_instance", &resource.Sweeper{
		Name: "tencentcloud_mongodb_instance",
		F: func(r string) error {
			logId := tccommon.GetLogId(tccommon.ContextNil)
			ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
			cli, _ := tcacctest.SharedClientForRegion(r)
			client := cli.(tccommon.ProviderMeta).GetAPIV3Conn()
			service := svcmongodb.NewMongodbService(client)

			instances, err := service.DescribeInstancesByFilter(ctx, "", -1)
			if err != nil {
				return err
			}

			var isolated []string

			for i := range instances {
				ins := instances[i]
				id := *ins.InstanceId
				name := *ins.InstanceName

				if strings.HasPrefix(name, tcacctest.KeepResource) || strings.HasPrefix(name, tcacctest.DefaultResource) {
					continue
				}

				created, err := time.Parse("2006-01-02 15:04:05", *ins.CreateTime)
				if err != nil {
					created = time.Time{}
				}
				if tcacctest.IsResourcePersist(name, &created) {
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

func TestAccTencentCloudMongodbInstanceResource_PostPaid(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMongodbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstance,
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
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "charge_type", svcmongodb.MONGODB_CHARGE_TYPE_POSTPAID),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_instance.mongodb", "prepaid_period"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "security_groups.0", "sg-if748odn"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_instance.mongodb",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"security_groups", "password", "auto_renew_flag"},
			},
			{
				Config: testAccMongodbInstance_updateConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "instance_name", "tf-mongodb-update"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "memory", "8"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "volume", "512"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.abc", "abc"),
				),
			},
			{
				Config: testAccMongodbInstance_updateNode,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "instance_name", "tf-mongodb-update"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "memory", "8"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "volume", "512"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "node_num", "5"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "tags.abc", "abc"),
				),
			},
			{
				Config: testAccMongodbInstance_updateSecurityGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "security_groups.0", "sg-05f7wnhn"),
				),
			},
			{
				Config: testAccMongodbInstance_updateMaintenance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "maintenance_start", "05:00"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb", "maintenance_end", "06:00"),
				),
			},
		},
	})
}

func TestAccTencentCloudMongodbInstanceResource_MultiZone(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMongodbInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstance_multiZone,
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

func TestAccTencentCloudMongodbInstanceResource_Prepaid(t *testing.T) {
	// Avoid to set Parallel to make sure EnvVar secure
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheck(t) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMongodbInstancePrepaid,
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
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "charge_type", svcmongodb.MONGODB_CHARGE_TYPE_PREPAID),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "prepaid_period", "1"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "auto_renew_flag", "1"),
				),
			},
			{
				Config: testAccMongodbInstancePrepaid_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "instance_name", "tf-mongodb-test-prepaid-update"),
					resource.TestCheckNoResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "tags.test"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "tags.prepaid", "prepaid"),
				),
			},
			{
				Config: testAccMongodbInstancePrepaid_updateMaintenance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "maintenance_start", "05:00"),
					resource.TestCheckResourceAttr("tencentcloud_mongodb_instance.mongodb_prepaid", "maintenance_end", "06:00"),
				),
			},
			{
				ResourceName:            "tencentcloud_mongodb_instance.mongodb_prepaid",
				ImportState:             true,
				ImportStateVerifyIgnore: []string{"security_groups", "password", "auto_renew_flag", "prepaid_period"},
			},
		},
	})
}

func testAccCheckMongodbInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	mongodbService := svcmongodb.NewMongodbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("mongodb instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("mongodb instance id is not set")
		}
		mongodbService := svcmongodb.NewMongodbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
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

const testAccMongodbInstance = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-test"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  tags = {
    test = "test"
  }
}
`

const testAccMongodbInstance_updateConfig = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-update"
  memory         = local.memory * 2
  volume         = local.volume * 2
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234update"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id
  tags = {
    abc = "abc"
  }
}
`

const testAccMongodbInstance_updateNode = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-update"
  memory         = local.memory * 2
  volume         = local.volume * 2
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234update"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  node_num = 5
  add_node_list {
    role = "SECONDARY"
    zone = "ap-guangzhou-3"
  }
  add_node_list {
    role = "SECONDARY"
    zone = "ap-guangzhou-3"
  }
  tags = {
    abc = "abc"
  }
}
`

const testAccMongodbInstance_updateSecurityGroup = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-update"
  memory         = local.memory * 2
  volume         = local.volume * 2
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = ["sg-05f7wnhn"]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234update"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  node_num = 5
  add_node_list {
    role = "SECONDARY"
    zone = "ap-guangzhou-3"
  }
  add_node_list {
    role = "SECONDARY"
    zone = "ap-guangzhou-3"
  }
  tags = {
    abc = "abc"
  }
}
`

const testAccMongodbInstance_updateMaintenance = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb" {
  instance_name  = "tf-mongodb-update"
  memory         = local.memory * 2
  volume         = local.volume * 2
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = ["sg-05f7wnhn"]
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234update"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  node_num = 5
  add_node_list {
    role = "SECONDARY"
    zone = "ap-guangzhou-3"
  }
  add_node_list {
    role = "SECONDARY"
    zone = "ap-guangzhou-3"
  }
  tags = {
    abc = "abc"
  }
  maintenance_start = "05:00"
  maintenance_end = "06:00"
}
`

const testAccMongodbInstancePrepaid = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-prepaid-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-prepaid-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb_prepaid" {
  instance_name   = "tf-mongodb-test-prepaid"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234"
  charge_type     = "PREPAID"
  prepaid_period  = 1
  auto_renew_flag = 1
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  tags = {
    test = "test-prepaid"
  }
}
`

const testAccMongodbInstancePrepaid_update = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-prepaid-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-prepaid-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb_prepaid" {
  instance_name   = "tf-mongodb-test-prepaid-update"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234update"
  charge_type     = "PREPAID"
  prepaid_period  = 1
  auto_renew_flag = 1
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  tags = {
    prepaid = "prepaid"
  }
}
`

const testAccMongodbInstancePrepaid_updateMaintenance = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-instance-prepaid-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-instance-prepaid-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb_prepaid" {
  instance_name   = "tf-mongodb-test-prepaid-update"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  security_groups = [local.security_group_id]
  available_zone  = "ap-guangzhou-3"
  project_id      = 0
  password        = "test1234update"
  charge_type     = "PREPAID"
  prepaid_period  = 1
  auto_renew_flag = 1
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id

  tags = {
    prepaid = "prepaid"
  }
  maintenance_start = "05:00"
  maintenance_end = "06:00"
}
`

const testAccMongodbInstance_multiZone = tcacctest.DefaultMongoDBSpec + `
resource "tencentcloud_vpc" "vpc" {
	name       = "mongodb-multi-zone-vpc"
	cidr_block = "10.0.0.0/16"
  }
  
resource "tencentcloud_subnet" "subnet" {
	vpc_id            = tencentcloud_vpc.vpc.id
	name              = "mongodb-multi-zone-subnet"
	cidr_block        = "10.0.0.0/16"
	availability_zone = "ap-guangzhou-3"
}

resource "tencentcloud_mongodb_instance" "mongodb_mutil_zone" {
  instance_name   = "mongodb-mutil-zone-test"
  memory         = local.memory
  volume         = local.volume
  engine_version = local.engine_version
  machine_type   = local.machine_type
  available_zone = "ap-guangzhou-3"
  project_id     = 0
  password       = "test1234"
  vpc_id         = tencentcloud_vpc.vpc.id
  subnet_id      = tencentcloud_subnet.subnet.id
  node_num = 5
  availability_zone_list = ["ap-guangzhou-3", "ap-guangzhou-3", "ap-guangzhou-4", "ap-guangzhou-4", "ap-guangzhou-6"]
  hidden_zone = "ap-guangzhou-6"
  tags = {
    test = "test"
  }
}
`

// ---- gomonkey-based unit tests for the cpu parameter ----
// Run with: go test ./tencentcloud/services/mongodb/ -run "TestMongodbInstanceCpu" -v -count=1 -gcflags="all=-l"
//
// These tests mock the cloud API methods (ModifyDBInstanceSpec, DescribeDBInstanceDeal,
// DescribeDBInstances) on the mongodb SDK client returned by UseMongodbClient, so that
// the service-layer mapping (params["cpu"] -> request.Cpu) and the resource Read/Update
// business logic are exercised end-to-end.

type mockMetaForMongodbInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForMongodbInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForMongodbInstance{}

func newMockMetaForMongodbInstance() *mockMetaForMongodbInstance {
	return &mockMetaForMongodbInstance{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

// buildMongodbInstanceDetailForCpu builds an InstanceDetail with all mandatory
// fields populated (so the Read CheckNil validation passes) and CpuNum set,
// suitable for the cpu Read backfill tests.
func buildMongodbInstanceDetailForCpu(cpuNum uint64) *mongodb_sdk.InstanceDetail {
	return &mongodb_sdk.InstanceDetail{
		InstanceId:   helper.String("cmgo-cpu-test"),
		InstanceName: helper.String("tf-mongodb-cpu-test"),
		PayMode:      helper.Uint64(0),
		ProjectId:    helper.Uint64(0),
		Zone:         helper.String("ap-guangzhou-3"),
		VpcId:        helper.String("vpc-test"),
		SubnetId:     helper.String("subnet-test"),
		Status:       helper.Int64(2),
		Vip:          helper.String("10.0.0.1"),
		Vport:        helper.Uint64(27017),
		CreateTime:   helper.String("2024-01-01 00:00:00"),
		MongoVersion: helper.String("MONGO_40_WT"),
		Memory:       helper.Uint64(2048),
		Volume:       helper.Uint64(25600),
		CpuNum:       helper.Uint64(cpuNum),
		MachineType:  helper.String("HIO10G"),
		SecondaryNum: helper.Uint64(2),
	}
}

// buildMongodbDescribeDBInstancesResponse wraps an InstanceDetail into a
// DescribeDBInstancesResponse, matching the shape consumed by DescribeInstanceById.
func buildMongodbDescribeDBInstancesResponse(instance *mongodb_sdk.InstanceDetail) *mongodb_sdk.DescribeDBInstancesResponse {
	resp := mongodb_sdk.NewDescribeDBInstancesResponse()
	resp.Response = &mongodb_sdk.DescribeDBInstancesResponseParams{
		InstanceDetails: []*mongodb_sdk.InstanceDetail{instance},
	}
	return resp
}

// mockCommonMongodbReadAPIMethods mocks the auxiliary cloud API methods invoked
// during resourceTencentCloudMongodbInstance_read (security groups, node property)
// on the mongodb SDK client. Tags are fetched via the tag service which is mocked
// separately.
func mockCommonMongodbReadAPIMethods(patches *gomonkey.Patches, mongodbClient *mongodb_sdk.Client) {
	patches.ApplyMethodFunc(mongodbClient, "DescribeSecurityGroup", func(request *mongodb_sdk.DescribeSecurityGroupRequest) (*mongodb_sdk.DescribeSecurityGroupResponse, error) {
		resp := mongodb_sdk.NewDescribeSecurityGroupResponse()
		resp.Response = &mongodb_sdk.DescribeSecurityGroupResponseParams{}
		return resp, nil
	})
	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstanceNodeProperty", func(request *mongodb_sdk.DescribeDBInstanceNodePropertyRequest) (*mongodb_sdk.DescribeDBInstanceNodePropertyResponse, error) {
		resp := mongodb_sdk.NewDescribeDBInstanceNodePropertyResponse()
		resp.Response = &mongodb_sdk.DescribeDBInstanceNodePropertyResponseParams{}
		return resp, nil
	})
	var tagService svctag.TagService
	patches.ApplyMethodFunc(&tagService, "DescribeResourceTags", func(ctx context.Context, serviceType, resourceType, region, resourceId string) (map[string]string, error) {
		return map[string]string{}, nil
	})
}

// mockMongodbUpdateOtherBranches mocks the cloud API methods that the
// resourceTencentCloudMongodbInstanceUpdate branches (other than the cpu spec
// modification branch) may invoke when TestResourceDataRaw reports all fields as
// changed. This isolates the cpu-triggered spec modification branch so the test
// focuses on the new cpu parameter behavior.
func mockMongodbUpdateOtherBranches(patches *gomonkey.Patches, mongodbClient *mongodb_sdk.Client) {
	// ModifyInstanceAz: return the "already distribution" error so the branch
	// skips the deal-polling path (needModifyFlag=false, returns nil).
	patches.ApplyMethodFunc(mongodbClient, "ModifyInstanceAz", func(request *mongodb_sdk.ModifyInstanceAzRequest) (*mongodb_sdk.ModifyInstanceAzResponse, error) {
		return nil, sdkErrors.NewTencentCloudSDKError("InvalidParameter", "The target az already distribution.", "")
	})
	// RenameInstance (instance_name change).
	patches.ApplyMethodFunc(mongodbClient, "RenameInstance", func(request *mongodb_sdk.RenameInstanceRequest) (*mongodb_sdk.RenameInstanceResponse, error) {
		resp := mongodb_sdk.NewRenameInstanceResponse()
		resp.Response = &mongodb_sdk.RenameInstanceResponseParams{}
		return resp, nil
	})
	// AssignProject (project_id change).
	patches.ApplyMethodFunc(mongodbClient, "AssignProject", func(request *mongodb_sdk.AssignProjectRequest) (*mongodb_sdk.AssignProjectResponse, error) {
		resp := mongodb_sdk.NewAssignProjectResponse()
		resp.Response = &mongodb_sdk.AssignProjectResponseParams{}
		return resp, nil
	})
	// ResetDBInstancePassword (password change).
	patches.ApplyMethodFunc(mongodbClient, "ResetDBInstancePassword", func(request *mongodb_sdk.ResetDBInstancePasswordRequest) (*mongodb_sdk.ResetDBInstancePasswordResponse, error) {
		resp := mongodb_sdk.NewResetDBInstancePasswordResponse()
		resp.Response = &mongodb_sdk.ResetDBInstancePasswordResponseParams{}
		return resp, nil
	})
	// ModifyDBInstanceNetworkAddress (vpc_id/subnet_id change).
	patches.ApplyMethodFunc(mongodbClient, "ModifyDBInstanceNetworkAddress", func(request *mongodb_sdk.ModifyDBInstanceNetworkAddressRequest) (*mongodb_sdk.ModifyDBInstanceNetworkAddressResponse, error) {
		resp := mongodb_sdk.NewModifyDBInstanceNetworkAddressResponse()
		resp.Response = &mongodb_sdk.ModifyDBInstanceNetworkAddressResponseParams{}
		return resp, nil
	})
	// ModifyDBInstanceSecurityGroup (security_groups change).
	patches.ApplyMethodFunc(mongodbClient, "ModifyDBInstanceSecurityGroup", func(request *mongodb_sdk.ModifyDBInstanceSecurityGroupRequest) (*mongodb_sdk.ModifyDBInstanceSecurityGroupResponse, error) {
		resp := mongodb_sdk.NewModifyDBInstanceSecurityGroupResponse()
		resp.Response = &mongodb_sdk.ModifyDBInstanceSecurityGroupResponseParams{}
		return resp, nil
	})
	// SetInstanceMaintenance (maintenance_start/end change).
	patches.ApplyMethodFunc(mongodbClient, "SetInstanceMaintenance", func(request *mongodb_sdk.SetInstanceMaintenanceRequest) (*mongodb_sdk.SetInstanceMaintenanceResponse, error) {
		resp := mongodb_sdk.NewSetInstanceMaintenanceResponse()
		resp.Response = &mongodb_sdk.SetInstanceMaintenanceResponseParams{}
		return resp, nil
	})
	// UpgradeDbInstanceVersion (engine_version change): return a flow id, then
	// DescribeAsyncRequestInfo returns success status.
	patches.ApplyMethodFunc(mongodbClient, "UpgradeDbInstanceVersionWithContext", func(ctx context.Context, request *mongodb_sdk.UpgradeDbInstanceVersionRequest) (*mongodb_sdk.UpgradeDbInstanceVersionResponse, error) {
		resp := mongodb_sdk.NewUpgradeDbInstanceVersionResponse()
		resp.Response = &mongodb_sdk.UpgradeDbInstanceVersionResponseParams{
			FlowId: helper.Uint64(1),
		}
		return resp, nil
	})
	patches.ApplyMethodFunc(mongodbClient, "DescribeAsyncRequestInfo", func(request *mongodb_sdk.DescribeAsyncRequestInfoRequest) (*mongodb_sdk.DescribeAsyncRequestInfoResponse, error) {
		resp := mongodb_sdk.NewDescribeAsyncRequestInfoResponse()
		resp.Response = &mongodb_sdk.DescribeAsyncRequestInfoResponseParams{
			Status: helper.String("success"),
		}
		return resp, nil
	})
	// ModifyTags (tags change).
	var tagService svctag.TagService
	patches.ApplyMethodFunc(&tagService, "ModifyTags", func(ctx context.Context, resourceName string, replaceTags map[string]string, deleteKeys []string) error {
		return nil
	})
}

// TestMongodbInstanceCpu_Schema verifies the cpu field schema definition.
func TestMongodbInstanceCpu_Schema(t *testing.T) {
	res := svcmongodb.ResourceTencentCloudMongodbInstance()

	s, ok := res.Schema["cpu"]
	assert.True(t, ok, "cpu field should exist in schema")
	assert.Equal(t, schema.TypeInt, s.Type)
	assert.True(t, s.Optional)
	assert.True(t, s.Computed)
	assert.False(t, s.ForceNew)
}

// TestMongodbInstanceCpu_Update_TriggersUpgrade verifies that changing only the cpu
// field triggers ModifyDBInstanceSpec with request.Cpu set, and that the deal is
// polled to success via DescribeDBInstanceDeal.
func TestMongodbInstanceCpu_Update_TriggersUpgrade(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForMongodbInstance().client, "UseMongodbClient", mongodbClient)

	var capturedRequest *mongodb_sdk.ModifyDBInstanceSpecRequest
	patches.ApplyMethodFunc(mongodbClient, "ModifyDBInstanceSpec", func(request *mongodb_sdk.ModifyDBInstanceSpecRequest) (*mongodb_sdk.ModifyDBInstanceSpecResponse, error) {
		capturedRequest = request
		resp := mongodb_sdk.NewModifyDBInstanceSpecResponse()
		resp.Response = &mongodb_sdk.ModifyDBInstanceSpecResponseParams{
			DealId: helper.String("test-deal-id"),
		}
		return resp, nil
	})
	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstanceDeal", func(request *mongodb_sdk.DescribeDBInstanceDealRequest) (*mongodb_sdk.DescribeDBInstanceDealResponse, error) {
		resp := mongodb_sdk.NewDescribeDBInstanceDealResponse()
		status := int64(4)
		resp.Response = &mongodb_sdk.DescribeDBInstanceDealResponseParams{
			Status: &status,
		}
		return resp, nil
	})
	// Mock the Read that runs at the end of Update.
	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstances", func(request *mongodb_sdk.DescribeDBInstancesRequest) (*mongodb_sdk.DescribeDBInstancesResponse, error) {
		return buildMongodbDescribeDBInstancesResponse(buildMongodbInstanceDetailForCpu(4)), nil
	})
	mockCommonMongodbReadAPIMethods(patches, mongodbClient)
	// Mock the other Update branches that trigger because TestResourceDataRaw
	// reports all fields as changed against an empty prior state.
	mockMongodbUpdateOtherBranches(patches, mongodbClient)

	meta := newMockMetaForMongodbInstance()
	res := svcmongodb.ResourceTencentCloudMongodbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":  "tf-mongodb-cpu-test",
		"memory":         2,
		"volume":         25,
		"engine_version": "MONGO_40_WT",
		"machine_type":   "HIO10G",
		"available_zone": "ap-guangzhou-3",
		"cpu":            2,
	})
	d.SetId("cmgo-cpu-test")

	// Simulate a cpu-only change from 2 to 4.
	d.Set("cpu", 4)

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify ModifyDBInstanceSpec was called and request.Cpu was set to 4.
	assert.NotNil(t, capturedRequest, "ModifyDBInstanceSpec should be called")
	assert.NotNil(t, capturedRequest.Cpu, "request.Cpu should be set")
	assert.Equal(t, int64(4), *capturedRequest.Cpu)
}

// TestMongodbInstanceCpu_Read_Backfill verifies that Read populates the cpu field
// from DescribeDBInstances (DescribeInstanceById) when CpuNum is not nil.
func TestMongodbInstanceCpu_Read_Backfill(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForMongodbInstance().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstances", func(request *mongodb_sdk.DescribeDBInstancesRequest) (*mongodb_sdk.DescribeDBInstancesResponse, error) {
		return buildMongodbDescribeDBInstancesResponse(buildMongodbInstanceDetailForCpu(8)), nil
	})
	mockCommonMongodbReadAPIMethods(patches, mongodbClient)

	meta := newMockMetaForMongodbInstance()
	res := svcmongodb.ResourceTencentCloudMongodbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":  "tf-mongodb-cpu-test",
		"memory":         2,
		"volume":         25,
		"engine_version": "MONGO_40_WT",
		"machine_type":   "HIO10G",
		"available_zone": "ap-guangzhou-3",
	})
	d.SetId("cmgo-cpu-test")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 8, d.Get("cpu"))
}

// TestMongodbInstanceCpu_Read_NilCpuNum verifies that Read does not fail and does
// not set cpu when CpuNum is nil.
func TestMongodbInstanceCpu_Read_NilCpuNum(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForMongodbInstance().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstances", func(request *mongodb_sdk.DescribeDBInstancesRequest) (*mongodb_sdk.DescribeDBInstancesResponse, error) {
		instance := buildMongodbInstanceDetailForCpu(0)
		instance.CpuNum = nil
		return buildMongodbDescribeDBInstancesResponse(instance), nil
	})
	mockCommonMongodbReadAPIMethods(patches, mongodbClient)

	meta := newMockMetaForMongodbInstance()
	res := svcmongodb.ResourceTencentCloudMongodbInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":  "tf-mongodb-cpu-test",
		"memory":         2,
		"volume":         25,
		"engine_version": "MONGO_40_WT",
		"machine_type":   "HIO10G",
		"available_zone": "ap-guangzhou-3",
	})
	d.SetId("cmgo-cpu-test")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// cpu should remain at its zero value (not set from API)
	_, ok := d.GetOk("cpu")
	assert.False(t, ok, "cpu should not be set when CpuNum is nil")
}
