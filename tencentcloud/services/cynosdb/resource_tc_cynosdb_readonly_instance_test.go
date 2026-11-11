package cynosdb_test

import (
	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svccynosdb "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cynosdb"

	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

type mockMetaCynosdbReadonlyInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaCynosdbReadonlyInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaCynosdbReadonlyInstance{}

func newMockMetaCynosdbReadonlyInstance() *mockMetaCynosdbReadonlyInstance {
	return &mockMetaCynosdbReadonlyInstance{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func TestAccTencentCloudCynosdbReadonlyInstanceResource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckCynosdbReadonlyInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbReadonlyInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCynosdbReadonlyInstanceExists("tencentcloud_cynosdb_readonly_instance.foo"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_name", "tf-cynosdb-readonly-instance"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "force_delete", "true"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_cpu_core", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_memory_size", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_maintain_duration", "3600"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_maintain_start_time", "10800"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_maintain_weekdays.#", "7"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_readonly_instance.foo", "instance_memory_size"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_readonly_instance.foo", "instance_status"),
					resource.TestCheckResourceAttrSet("tencentcloud_cynosdb_readonly_instance.foo", "instance_storage_size"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "vpc_id", "vpc-m0d2dbnn"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "subnet_id", "subnet-j10lsueq"),
				),
			},
			{
				ResourceName:            "tencentcloud_cynosdb_readonly_instance.foo",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force_delete"},
			},
			{
				Config: testAccCynosdbReadonlyInstance_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_maintain_duration", "7200"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_maintain_start_time", "21600"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_maintain_weekdays.#", "6"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_cpu_core", "2"),
					resource.TestCheckResourceAttr("tencentcloud_cynosdb_readonly_instance.foo", "instance_memory_size", "4"),
				),
			},
		},
	})
}

func testAccCheckCynosdbReadonlyInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	cynosdbService := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_cynosdb_readonly_instance" {
			continue
		}

		_, _, has, err := cynosdbService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return nil
		}
		return fmt.Errorf("cynosdb readonly instance still exists: %s", rs.Primary.ID)
	}
	return nil
}

func testAccCheckCynosdbReadonlyInstanceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("cynosdb readonly instance %s is not found", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("cynosdb readonly instance id is not set")
		}
		cynosdbService := svccynosdb.NewCynosdbService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		_, _, has, err := cynosdbService.DescribeInstanceById(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}
		if !has {
			return fmt.Errorf("cynosdb readonly instance doesn't exist: %s", rs.Primary.ID)
		}
		return nil
	}
}

const readonlyInstanceVar = `
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "my_vpc" {
  default = "vpc-m0d2dbnn"
}

variable "my_subnet" {
  default = "subnet-j10lsueq"
}

variable "readonly_subnet" {
  default = "subnet-j10lsueq"
}
`

const testAccCynosdbReadonlyInstance = readonlyInstanceVar + `
resource "tencentcloud_cynosdb_cluster" "foo" {
  available_zone               = var.availability_zone
  vpc_id                       = var.my_vpc
  subnet_id                    = var.my_subnet
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf-cynosdb"
  password                     = "cynos@123"
  instance_maintain_duration   = 3600
  instance_maintain_start_time = 10800
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  instance_cpu_core    = 1
  instance_memory_size = 2
  param_items {
    name = "character_set_server"
    current_value = "utf8"
  }

  tags = {
    test = "test"
  }

  force_delete = true

  rw_group_sg = [
    "` + tcacctest.DefaultSecurityGroup + `",
  ]
}

resource "tencentcloud_cynosdb_readonly_instance" "foo" {
  cluster_id           = tencentcloud_cynosdb_cluster.foo.id
  instance_name        = "tf-cynosdb-readonly-instance"
  force_delete         = true
  instance_cpu_core    = 1
  instance_memory_size = 2
  vpc_id               = var.my_vpc
  subnet_id            = var.readonly_subnet

  instance_maintain_duration   = 3600
  instance_maintain_start_time = 10800
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]
}
`

const testAccCynosdbReadonlyInstance_update = readonlyInstanceVar + `
resource "tencentcloud_cynosdb_cluster" "foo" {
  available_zone               = var.availability_zone
  vpc_id                       = var.my_vpc
  subnet_id                    = var.my_subnet
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf-cynosdb"
  password                     = "cynos@123"
  instance_maintain_duration   = 3600
  instance_maintain_start_time = 10800
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  instance_cpu_core    = 1
  instance_memory_size = 2
  param_items {
    name = "character_set_server"
    current_value = "utf8"
  }

  tags = {
    test = "test"
  }

  force_delete = true

  rw_group_sg = [
    "` + tcacctest.DefaultSecurityGroup + `",
  ]
}

resource "tencentcloud_cynosdb_readonly_instance" "foo" {
  cluster_id           = tencentcloud_cynosdb_cluster.foo.id
  instance_name        = "tf-cynosdb-readonly-instance"
  force_delete         = true
  instance_cpu_core    = 2
  instance_memory_size = 4

  instance_maintain_duration   = 7200
  instance_maintain_start_time = 21600
  instance_maintain_weekdays   = [
    "Fri",
    "Mon",
    "Sat",
    "Thu",
    "Wed",
    "Tue",
  ]
}
`

// go test ./tencentcloud/services/cynosdb/ -run "TestUnitCynosdbReadonlyInstance_UpdateInstanceName" -v -count=1 -gcflags="all=-l"
func TestUnitCynosdbReadonlyInstance_UpdateInstanceName(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaCynosdbReadonlyInstance()
	cynosdbClient := &cynosdb.Client{}
	patches.ApplyMethodReturn(meta.client, "UseCynosdbClient", cynosdbClient)

	modifyCalled := false
	patches.ApplyMethodFunc(cynosdbClient, "ModifyInstanceName", func(request *cynosdb.ModifyInstanceNameRequest) (*cynosdb.ModifyInstanceNameResponse, error) {
		assert.Equal(t, "cynosdbmysql-ins-abcdefgh", *request.InstanceId)
		assert.Equal(t, "tf-cynosdb-readonly-instance-new", *request.InstanceName)
		modifyCalled = true
		resp := &cynosdb.ModifyInstanceNameResponse{}
		resp.Response = &cynosdb.ModifyInstanceNameResponseParams{
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	// mock DescribeInstances used by DescribeInstanceById in Read
	patches.ApplyMethodFunc(cynosdbClient, "DescribeInstances", func(request *cynosdb.DescribeInstancesRequest) (*cynosdb.DescribeInstancesResponse, error) {
		resp := &cynosdb.DescribeInstancesResponse{}
		resp.Response = &cynosdb.DescribeInstancesResponseParams{
			TotalCount: helper.IntInt64(1),
			InstanceSet: []*cynosdb.CynosdbInstance{
				{
					ClusterId:    helper.String("cynosdbmysql-12345678"),
					InstanceId:   helper.String("cynosdbmysql-ins-abcdefgh"),
					InstanceName: helper.String("tf-cynosdb-readonly-instance-new"),
					Status:       helper.String("running"),
					Cpu:          helper.IntInt64(2),
					Memory:       helper.IntInt64(4),
					VpcId:        helper.String("vpc-m0d2dbnn"),
					SubnetId:     helper.String("subnet-j10lsueq"),
				},
			},
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	// mock DescribeInstanceDetail used by DescribeInstanceById in Read
	patches.ApplyMethodFunc(cynosdbClient, "DescribeInstanceDetail", func(request *cynosdb.DescribeInstanceDetailRequest) (*cynosdb.DescribeInstanceDetailResponse, error) {
		resp := &cynosdb.DescribeInstanceDetailResponse{}
		resp.Response = &cynosdb.DescribeInstanceDetailResponseParams{
			Detail: &cynosdb.CynosdbInstanceDetail{
				ClusterId:    helper.String("cynosdbmysql-12345678"),
				InstanceId:   helper.String("cynosdbmysql-ins-abcdefgh"),
				InstanceName: helper.String("tf-cynosdb-readonly-instance-new"),
				Status:       helper.String("running"),
				Cpu:          helper.IntInt64(2),
				Memory:       helper.IntInt64(4),
				Storage:      helper.IntInt64(50),
				VpcId:        helper.String("vpc-m0d2dbnn"),
				SubnetId:     helper.String("subnet-j10lsueq"),
			},
			RequestId: helper.String("fake-request-id"),
		}
		return resp, nil
	})

	// mock DescribeMaintainPeriod used by Read
	patches.ApplyMethodFunc(cynosdbClient, "DescribeMaintainPeriod", func(request *cynosdb.DescribeMaintainPeriodRequest) (*cynosdb.DescribeMaintainPeriodResponse, error) {
		resp := &cynosdb.DescribeMaintainPeriodResponse{}
		resp.Response = &cynosdb.DescribeMaintainPeriodResponseParams{
			MaintainWeekDays:  []*string{helper.String("Mon"), helper.String("Tue")},
			MaintainStartTime: helper.IntInt64(10800),
			MaintainDuration:  helper.IntInt64(3600),
			RequestId:         helper.String("fake-request-id"),
		}
		return resp, nil
	})

	res := svccynosdb.ResourceTencentCloudCynosdbReadonlyInstance()

	// Build a prior state where instance_name is the old value, and a new config
	// where only instance_name changes, so that d.HasChange("instance_name")
	// returns true while other fields remain unchanged.
	state := &terraform.InstanceState{
		ID: "cynosdbmysql-ins-abcdefgh",
		Attributes: map[string]string{
			"id":                           "cynosdbmysql-ins-abcdefgh",
			"cluster_id":                   "cynosdbmysql-12345678",
			"instance_name":                "tf-cynosdb-readonly-instance-old",
			"instance_cpu_core":            "2",
			"instance_memory_size":         "4",
			"instance_maintain_duration":   "3600",
			"instance_maintain_start_time": "10800",
		},
	}

	rawConfig := terraform.NewResourceConfigRaw(map[string]interface{}{
		"cluster_id":                   "cynosdbmysql-12345678",
		"instance_name":                "tf-cynosdb-readonly-instance-new",
		"instance_cpu_core":            2,
		"instance_memory_size":         4,
		"instance_maintain_duration":   3600,
		"instance_maintain_start_time": 10800,
	})

	diff, err := res.Diff(nil, state, rawConfig, meta)
	assert.NoError(t, err)
	assert.NotNil(t, diff)

	d, err := schema.InternalMap(res.Schema).Data(state, diff)
	assert.NoError(t, err)

	err = res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, "tf-cynosdb-readonly-instance-new", d.Get("instance_name"))
}
