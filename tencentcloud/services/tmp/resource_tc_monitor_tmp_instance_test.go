package tmp_test

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	monitor "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20180724"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tmp"
)

// go test -i; go test -test.run TestAccTencentCloudMonitorInstance_basic -v
func TestAccTencentCloudMonitorInstance_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { tcacctest.AccPreCheck(t) },
		Providers:    tcacctest.AccProviders,
		CheckDestroy: testAccCheckMonInstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("tencentcloud_monitor_tmp_instance.example"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.example", "instance_name", "tf-tmp-instance"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_instance.example", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_instance.example", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.example", "data_retention_time", "30"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_instance.example", "zone"),
				),
			},
			{
				Config: testInstance_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists("tencentcloud_monitor_tmp_instance.example"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.example", "instance_name", "tf-tmp-instance-update"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_instance.example", "vpc_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_instance.example", "subnet_id"),
					resource.TestCheckResourceAttr("tencentcloud_monitor_tmp_instance.example", "data_retention_time", "90"),
					resource.TestCheckResourceAttrSet("tencentcloud_monitor_tmp_instance.example", "zone"),
				),
			},
			{
				ResourceName:      "tencentcloud_monitor_tmp_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckMonInstanceDestroy(s *terraform.State) error {
	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tencentcloud_monitor_tmp_instance" {
			continue
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		instance, err := service.DescribeMonitorTmpInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance != nil {
			status := strconv.FormatInt(*instance.InstanceStatus, 10)
			if strings.Contains("5,6,8,9", status) {
				return nil
			}
			return fmt.Errorf("instance %s still exists: %v", rs.Primary.ID, *instance.InstanceStatus)
		}
	}

	return nil
}

func testAccCheckInstanceExists(r string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		logId := tccommon.GetLogId(tccommon.ContextNil)
		ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

		rs, ok := s.RootModule().Resources[r]
		if !ok {
			return fmt.Errorf("resource %s is not found", r)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id is not set")
		}

		service := svcmonitor.NewMonitorService(tcacctest.AccProvider.Meta().(tccommon.ProviderMeta).GetAPIV3Conn())
		instance, err := service.DescribeMonitorTmpInstance(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if instance == nil || *instance.InstanceStatus != 2 {
			return fmt.Errorf("instance %s is not found", rs.Primary.ID)
		}

		return nil
	}
}

const testInstance_basic = `
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_monitor_tmp_instance" "example" {
  instance_name       = "tf-tmp-instance"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  data_retention_time = 30
  zone                = var.availability_zone
  tags = {
    "createdBy" = "terraform"
  }
}
`

const testInstance_update = `
variable "availability_zone" {
  default = "ap-guangzhou-4"
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "tf_monitor_vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = var.availability_zone
  name              = "tf_monitor_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_monitor_tmp_instance" "example" {
  instance_name       = "tf-tmp-instance-update"
  vpc_id              = tencentcloud_vpc.vpc.id
  subnet_id           = tencentcloud_subnet.subnet.id
  data_retention_time = 90
  zone                = var.availability_zone
  tags = {
    "createdBy" = "terraformUpdate"
  }
}
`

// mockMetaTmpInstance implements tccommon.ProviderMeta
type mockMetaTmpInstance struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTmpInstance) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTmpInstance{}

func newMockMetaTmpInstance() *mockMetaTmpInstance {
	return &mockMetaTmpInstance{client: &connectivity.TencentCloudClient{Region: "ap-guangzhou"}}
}

func ptrStringTmpInstance(s string) *string {
	return &s
}

func ptrInt64TmpInstance(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/tmp/ -run "TestUnitMonitorTmpInstance" -v -count=1 -gcflags="all=-l"

// TestUnitMonitorTmpInstance_CreateWithLongTermStorageRetentionTime tests Create with long_term_storage_retention_time specified
func TestUnitMonitorTmpInstance_CreateWithLongTermStorageRetentionTime(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTmpInstance()
	monitorClient := &monitor.Client{}
	patches.ApplyMethodReturn(meta.client, "UseMonitorClient", monitorClient)

	// Patch CreatePrometheusMultiTenantInstancePostPayMode
	patches.ApplyMethodFunc(monitorClient, "CreatePrometheusMultiTenantInstancePostPayMode", func(request *monitor.CreatePrometheusMultiTenantInstancePostPayModeRequest) (*monitor.CreatePrometheusMultiTenantInstancePostPayModeResponse, error) {
		assert.NotNil(t, request.InstanceAttributes)
		assert.Equal(t, 1, len(request.InstanceAttributes))
		assert.Equal(t, "LongTermStorageRetentionTime", *request.InstanceAttributes[0].Key)
		assert.Equal(t, "90", *request.InstanceAttributes[0].Value)
		assert.Equal(t, "test-instance", *request.InstanceName)
		assert.Equal(t, "vpc-123", *request.VpcId)
		assert.Equal(t, "subnet-123", *request.SubnetId)
		assert.Equal(t, int64(30), *request.DataRetentionTime)
		assert.Equal(t, "ap-guangzhou-4", *request.Zone)

		resp := &monitor.CreatePrometheusMultiTenantInstancePostPayModeResponse{}
		resp.Response = &monitor.CreatePrometheusMultiTenantInstancePostPayModeResponseParams{
			InstanceId: ptrStringTmpInstance("prom-test-001"),
			RequestId:  ptrStringTmpInstance("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribePrometheusInstances for status polling and Read
	patches.ApplyMethodFunc(monitorClient, "DescribePrometheusInstances", func(request *monitor.DescribePrometheusInstancesRequest) (*monitor.DescribePrometheusInstancesResponse, error) {
		resp := &monitor.DescribePrometheusInstancesResponse{}
		resp.Response = &monitor.DescribePrometheusInstancesResponseParams{
			InstanceSet: []*monitor.PrometheusInstancesItem{
				{
					InstanceId:        ptrStringTmpInstance("prom-test-001"),
					InstanceName:      ptrStringTmpInstance("test-instance"),
					InstanceStatus:    ptrInt64TmpInstance(2),
					VpcId:             ptrStringTmpInstance("vpc-123"),
					SubnetId:          ptrStringTmpInstance("subnet-123"),
					DataRetentionTime: ptrInt64TmpInstance(30),
					Zone:              ptrStringTmpInstance("ap-guangzhou-4"),
					IPv4Address:       ptrStringTmpInstance("10.0.0.1"),
					InstanceAttributes: []*monitor.PrometheusRuleKV{
						{
							Key:   ptrStringTmpInstance("LongTermStorageRetentionTime"),
							Value: ptrStringTmpInstance("90"),
						},
					},
				},
			},
			TotalCount: ptrInt64TmpInstance(1),
			RequestId:  ptrStringTmpInstance("fake-request-id"),
		}
		return resp, nil
	})

	// Patch TagService.DescribeResourceTags
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	res := tmp.ResourceTencentCloudMonitorTmpInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":                    "test-instance",
		"vpc_id":                           "vpc-123",
		"subnet_id":                        "subnet-123",
		"data_retention_time":              30,
		"zone":                             "ap-guangzhou-4",
		"long_term_storage_retention_time": 90,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "prom-test-001", d.Id())
	assert.Equal(t, 90, d.Get("long_term_storage_retention_time"))
}

// TestUnitMonitorTmpInstance_ReadWithLongTermStorageRetentionTime tests Read populating long_term_storage_retention_time
func TestUnitMonitorTmpInstance_ReadWithLongTermStorageRetentionTime(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTmpInstance()
	monitorClient := &monitor.Client{}
	patches.ApplyMethodReturn(meta.client, "UseMonitorClient", monitorClient)

	// Patch DescribePrometheusInstances
	patches.ApplyMethodFunc(monitorClient, "DescribePrometheusInstances", func(request *monitor.DescribePrometheusInstancesRequest) (*monitor.DescribePrometheusInstancesResponse, error) {
		resp := &monitor.DescribePrometheusInstancesResponse{}
		resp.Response = &monitor.DescribePrometheusInstancesResponseParams{
			InstanceSet: []*monitor.PrometheusInstancesItem{
				{
					InstanceId:        ptrStringTmpInstance("prom-test-001"),
					InstanceName:      ptrStringTmpInstance("test-instance"),
					InstanceStatus:    ptrInt64TmpInstance(2),
					VpcId:             ptrStringTmpInstance("vpc-123"),
					SubnetId:          ptrStringTmpInstance("subnet-123"),
					DataRetentionTime: ptrInt64TmpInstance(30),
					Zone:              ptrStringTmpInstance("ap-guangzhou-4"),
					IPv4Address:       ptrStringTmpInstance("10.0.0.1"),
					RemoteWrite:       ptrStringTmpInstance("http://remote-write.example.com"),
					ApiRootPath:       ptrStringTmpInstance("http://api-root.example.com"),
					ProxyAddress:      ptrStringTmpInstance("http://proxy.example.com"),
					InstanceAttributes: []*monitor.PrometheusRuleKV{
						{
							Key:   ptrStringTmpInstance("LongTermStorageRetentionTime"),
							Value: ptrStringTmpInstance("120"),
						},
						{
							Key:   ptrStringTmpInstance("CreatedFrom"),
							Value: ptrStringTmpInstance("0"),
						},
					},
				},
			},
			TotalCount: ptrInt64TmpInstance(1),
			RequestId:  ptrStringTmpInstance("fake-request-id"),
		}
		return resp, nil
	})

	// Patch TagService.DescribeResourceTags
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{"env": "test"}, nil
	})

	res := tmp.ResourceTencentCloudMonitorTmpInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":       "test-instance",
		"vpc_id":              "vpc-123",
		"subnet_id":           "subnet-123",
		"data_retention_time": 30,
		"zone":                "ap-guangzhou-4",
	})
	d.SetId("prom-test-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 120, d.Get("long_term_storage_retention_time"))
	assert.Equal(t, "test-instance", d.Get("instance_name"))
}

// TestUnitMonitorTmpInstance_ReadWithNilInstanceAttributes tests Read when InstanceAttributes is nil
func TestUnitMonitorTmpInstance_ReadWithNilInstanceAttributes(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	meta := newMockMetaTmpInstance()
	monitorClient := &monitor.Client{}
	patches.ApplyMethodReturn(meta.client, "UseMonitorClient", monitorClient)

	// Patch DescribePrometheusInstances - InstanceAttributes is nil
	patches.ApplyMethodFunc(monitorClient, "DescribePrometheusInstances", func(request *monitor.DescribePrometheusInstancesRequest) (*monitor.DescribePrometheusInstancesResponse, error) {
		resp := &monitor.DescribePrometheusInstancesResponse{}
		resp.Response = &monitor.DescribePrometheusInstancesResponseParams{
			InstanceSet: []*monitor.PrometheusInstancesItem{
				{
					InstanceId:         ptrStringTmpInstance("prom-test-001"),
					InstanceName:       ptrStringTmpInstance("test-instance"),
					InstanceStatus:     ptrInt64TmpInstance(2),
					VpcId:              ptrStringTmpInstance("vpc-123"),
					SubnetId:           ptrStringTmpInstance("subnet-123"),
					DataRetentionTime:  ptrInt64TmpInstance(30),
					Zone:               ptrStringTmpInstance("ap-guangzhou-4"),
					InstanceAttributes: nil,
				},
			},
			TotalCount: ptrInt64TmpInstance(1),
			RequestId:  ptrStringTmpInstance("fake-request-id"),
		}
		return resp, nil
	})

	// Patch TagService.DescribeResourceTags
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	res := tmp.ResourceTencentCloudMonitorTmpInstance()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_name":       "test-instance",
		"vpc_id":              "vpc-123",
		"subnet_id":           "subnet-123",
		"data_retention_time": 30,
		"zone":                "ap-guangzhou-4",
	})
	d.SetId("prom-test-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 0, d.Get("long_term_storage_retention_time"))
}

// TestUnitMonitorTmpInstance_UpdateLongTermStorageRetentionTime tests Update modifying long_term_storage_retention_time
func TestUnitMonitorTmpInstance_UpdateLongTermStorageRetentionTime(t *testing.T) {
	res := tmp.ResourceTencentCloudMonitorTmpInstance()
	d := res.TestResourceData()
	d.SetId("prom-test-001")

	// Simulate old state
	_ = d.Set("instance_name", "test-instance")
	_ = d.Set("vpc_id", "vpc-123")
	_ = d.Set("subnet_id", "subnet-123")
	_ = d.Set("data_retention_time", 30)
	_ = d.Set("zone", "ap-guangzhou-4")
	_ = d.Set("long_term_storage_retention_time", 90)

	// Mark the resource as not new to enable HasChange detection
	d.MarkNewResource()

	meta := newMockMetaTmpInstance()

	// Verify the schema has the field
	longTermSchema := res.Schema["long_term_storage_retention_time"]
	assert.NotNil(t, longTermSchema)
	assert.True(t, longTermSchema.Optional)
	assert.True(t, longTermSchema.Computed)
	assert.Equal(t, schema.TypeInt, longTermSchema.Type)

	// Verify the Update function exists and works
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	monitorClient := &monitor.Client{}
	patches.ApplyMethodReturn(meta.client, "UseMonitorClient", monitorClient)

	// Patch ModifyPrometheusInstanceAttributes
	patches.ApplyMethodFunc(monitorClient, "ModifyPrometheusInstanceAttributes", func(request *monitor.ModifyPrometheusInstanceAttributesRequest) (*monitor.ModifyPrometheusInstanceAttributesResponse, error) {
		resp := &monitor.ModifyPrometheusInstanceAttributesResponse{}
		resp.Response = &monitor.ModifyPrometheusInstanceAttributesResponseParams{
			RequestId: ptrStringTmpInstance("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribePrometheusInstances for Read after Update
	patches.ApplyMethodFunc(monitorClient, "DescribePrometheusInstances", func(request *monitor.DescribePrometheusInstancesRequest) (*monitor.DescribePrometheusInstancesResponse, error) {
		resp := &monitor.DescribePrometheusInstancesResponse{}
		resp.Response = &monitor.DescribePrometheusInstancesResponseParams{
			InstanceSet: []*monitor.PrometheusInstancesItem{
				{
					InstanceId:        ptrStringTmpInstance("prom-test-001"),
					InstanceName:      ptrStringTmpInstance("test-instance"),
					InstanceStatus:    ptrInt64TmpInstance(2),
					VpcId:             ptrStringTmpInstance("vpc-123"),
					SubnetId:          ptrStringTmpInstance("subnet-123"),
					DataRetentionTime: ptrInt64TmpInstance(30),
					Zone:              ptrStringTmpInstance("ap-guangzhou-4"),
					InstanceAttributes: []*monitor.PrometheusRuleKV{
						{
							Key:   ptrStringTmpInstance("LongTermStorageRetentionTime"),
							Value: ptrStringTmpInstance("90"),
						},
					},
				},
			},
			TotalCount: ptrInt64TmpInstance(1),
			RequestId:  ptrStringTmpInstance("fake-request-id"),
		}
		return resp, nil
	})

	// Patch TagService.DescribeResourceTags
	patches.ApplyMethodFunc(&svctag.TagService{}, "DescribeResourceTags", func(_ context.Context, _, _, _, _ string) (map[string]string, error) {
		return map[string]string{}, nil
	})

	// Update with no changes should succeed
	err := res.Update(d, meta)
	assert.NoError(t, err)
}
