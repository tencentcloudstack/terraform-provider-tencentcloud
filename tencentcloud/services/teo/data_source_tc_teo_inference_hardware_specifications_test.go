package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaTeoInferenceHardwareSpecificationsDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoInferenceHardwareSpecificationsDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoInferenceHardwareSpecificationsDS{}

func newMockMetaTeoInferenceHardwareSpecificationsDS() *mockMetaTeoInferenceHardwareSpecificationsDS {
	return &mockMetaTeoInferenceHardwareSpecificationsDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStrTeoInfHwSpecDS(s string) *string {
	return &s
}

func ptrFloat64TeoInfHwSpecDS(f float64) *float64 {
	return &f
}

func ptrInt64TeoInfHwSpecDS(i int64) *int64 {
	return &i
}

// TestTeoInferenceHardwareSpecificationsDataSource_Read_Basic tests that all fields are correctly
// flattened when the API returns a fully populated hardware specification.
func TestTeoInferenceHardwareSpecificationsDataSource_Read_Basic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoInferenceHardwareSpecificationsByFilter", func(_ context.Context, _ map[string]interface{}) ([]*teov20220901.InferenceHardwareSpecification, error) {
		return []*teov20220901.InferenceHardwareSpecification{
			{
				Spec:           ptrStrTeoInfHwSpecDS("spec-001"),
				HardwareSpecId: ptrStrTeoInfHwSpecDS("hw-spec-001"),
				Name:           ptrStrTeoInfHwSpecDS("GPU Standard"),
				GPUNum:         ptrFloat64TeoInfHwSpecDS(1.0),
				CPUNum:         ptrFloat64TeoInfHwSpecDS(8.0),
				MemSize:        ptrInt64TeoInfHwSpecDS(16384),
				GPUMemSize:     ptrInt64TeoInfHwSpecDS(16384),
				DiskSize:       ptrInt64TeoInfHwSpecDS(51200),
				AllowedGPUNums: []*float64{ptrFloat64TeoInfHwSpecDS(1.0), ptrFloat64TeoInfHwSpecDS(2.0), ptrFloat64TeoInfHwSpecDS(4.0)},
			},
		}, nil
	})

	meta := newMockMetaTeoInferenceHardwareSpecificationsDS()
	res := teo.DataSourceTencentCloudTeoInferenceHardwareSpecifications()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	hardwareSpecifications := d.Get("hardware_specifications").([]interface{})
	assert.Equal(t, 1, len(hardwareSpecifications))
	spec := hardwareSpecifications[0].(map[string]interface{})
	assert.Equal(t, "spec-001", spec["spec"])
	assert.Equal(t, "hw-spec-001", spec["hardware_spec_id"])
	assert.Equal(t, "GPU Standard", spec["name"])
	assert.Equal(t, 1.0, spec["gpu_num"])
	assert.Equal(t, 8.0, spec["cpu_num"])
	assert.Equal(t, 16384, spec["mem_size"])
	assert.Equal(t, 16384, spec["gpu_mem_size"])
	assert.Equal(t, 51200, spec["disk_size"])

	allowedGpuNums := spec["allowed_gpu_nums"].([]interface{})
	assert.Equal(t, 3, len(allowedGpuNums))
	assert.Equal(t, 1.0, allowedGpuNums[0])
	assert.Equal(t, 2.0, allowedGpuNums[1])
	assert.Equal(t, 4.0, allowedGpuNums[2])
}

// TestTeoInferenceHardwareSpecificationsDataSource_Read_NilFields tests that nil fields are safely
// omitted and no error is raised.
func TestTeoInferenceHardwareSpecificationsDataSource_Read_NilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&teo.TeoService{}, "DescribeTeoInferenceHardwareSpecificationsByFilter", func(_ context.Context, _ map[string]interface{}) ([]*teov20220901.InferenceHardwareSpecification, error) {
		return []*teov20220901.InferenceHardwareSpecification{
			{
				Spec:           nil,
				HardwareSpecId: ptrStrTeoInfHwSpecDS("hw-spec-002"),
				Name:           ptrStrTeoInfHwSpecDS("GPU Lite"),
				GPUNum:         nil,
				CPUNum:         nil,
				MemSize:        nil,
				GPUMemSize:     nil,
				DiskSize:       nil,
				AllowedGPUNums: nil,
			},
		}, nil
	})

	meta := newMockMetaTeoInferenceHardwareSpecificationsDS()
	res := teo.DataSourceTencentCloudTeoInferenceHardwareSpecifications()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	hardwareSpecifications := d.Get("hardware_specifications").([]interface{})
	assert.Equal(t, 1, len(hardwareSpecifications))
	spec := hardwareSpecifications[0].(map[string]interface{})
	assert.Equal(t, "hw-spec-002", spec["hardware_spec_id"])
	assert.Equal(t, "GPU Lite", spec["name"])
	assert.Equal(t, "", spec["spec"])
	assert.Equal(t, 0.0, spec["gpu_num"])
	assert.Nil(t, spec["allowed_gpu_nums"])
}

// TestTeoInferenceHardwareSpecificationsDataSource_Schema tests the schema definition.
func TestTeoInferenceHardwareSpecificationsDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoInferenceHardwareSpecifications()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	zoneIdSchema, ok := res.Schema["zone_id"]
	assert.True(t, ok)
	assert.NotNil(t, zoneIdSchema)
	assert.Equal(t, schema.TypeString, zoneIdSchema.Type)
	assert.True(t, zoneIdSchema.Required)

	hardwareSpecifications := res.Schema["hardware_specifications"]
	assert.NotNil(t, hardwareSpecifications)
	assert.Equal(t, schema.TypeList, hardwareSpecifications.Type)
	assert.True(t, hardwareSpecifications.Computed)
	elem := hardwareSpecifications.Elem.(*schema.Resource)

	hardwareSpecIdSchema, ok := elem.Schema["hardware_spec_id"]
	assert.True(t, ok)
	assert.NotNil(t, hardwareSpecIdSchema)
	assert.Equal(t, schema.TypeString, hardwareSpecIdSchema.Type)

	gpuNumSchema, ok := elem.Schema["gpu_num"]
	assert.True(t, ok)
	assert.NotNil(t, gpuNumSchema)
	assert.Equal(t, schema.TypeFloat, gpuNumSchema.Type)

	memSizeSchema, ok := elem.Schema["mem_size"]
	assert.True(t, ok)
	assert.NotNil(t, memSizeSchema)
	assert.Equal(t, schema.TypeInt, memSizeSchema.Type)

	allowedGpuNumsSchema, ok := elem.Schema["allowed_gpu_nums"]
	assert.True(t, ok)
	assert.NotNil(t, allowedGpuNumsSchema)
	assert.Equal(t, schema.TypeList, allowedGpuNumsSchema.Type)
}
