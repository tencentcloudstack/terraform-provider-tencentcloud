package cvm

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func TestCpuTopologySchemaDefinition(t *testing.T) {
	res := ResourceTencentCloudInstance()
	s := res.Schema

	cpuTopologySchema, ok := s["cpu_topology"]
	assert.True(t, ok, "cpu_topology field should exist in schema")
	assert.Equal(t, schema.TypeList, cpuTopologySchema.Type)
	assert.True(t, cpuTopologySchema.Optional)
	assert.True(t, cpuTopologySchema.ForceNew)
	assert.Equal(t, 1, cpuTopologySchema.MaxItems)

	elemResource, ok := cpuTopologySchema.Elem.(*schema.Resource)
	assert.True(t, ok, "cpu_topology Elem should be *schema.Resource")

	coreCountSchema, ok := elemResource.Schema["core_count"]
	assert.True(t, ok, "core_count field should exist")
	assert.Equal(t, schema.TypeInt, coreCountSchema.Type)
	assert.True(t, coreCountSchema.Optional)
	assert.True(t, coreCountSchema.ForceNew)

	threadPerCoreSchema, ok := elemResource.Schema["thread_per_core"]
	assert.True(t, ok, "thread_per_core field should exist")
	assert.Equal(t, schema.TypeInt, threadPerCoreSchema.Type)
	assert.True(t, threadPerCoreSchema.Optional)
	assert.True(t, threadPerCoreSchema.ForceNew)
}

func TestCpuTopologyCreateRequestPopulation(t *testing.T) {
	request := cvm.NewRunInstancesRequest()

	cpuTopologyList := []interface{}{
		map[string]interface{}{
			"core_count":      4,
			"thread_per_core": 1,
		},
	}

	if len(cpuTopologyList) > 0 {
		cpuTopologyMap := cpuTopologyList[0].(map[string]interface{})
		cpuTopology := &cvm.CpuTopology{}
		if coreCount, ok := cpuTopologyMap["core_count"]; ok && coreCount.(int) > 0 {
			cpuTopology.CoreCount = helper.IntInt64(coreCount.(int))
		}
		if threadPerCore, ok := cpuTopologyMap["thread_per_core"]; ok && threadPerCore.(int) > 0 {
			cpuTopology.ThreadPerCore = helper.IntInt64(threadPerCore.(int))
		}
		request.CpuTopology = cpuTopology
	}

	assert.NotNil(t, request.CpuTopology)
	assert.Equal(t, int64(4), *request.CpuTopology.CoreCount)
	assert.Equal(t, int64(1), *request.CpuTopology.ThreadPerCore)
}

func TestCpuTopologyCreateRequestWithoutTopology(t *testing.T) {
	request := cvm.NewRunInstancesRequest()

	cpuTopologyList := []interface{}{}

	if len(cpuTopologyList) > 0 {
		cpuTopologyMap := cpuTopologyList[0].(map[string]interface{})
		cpuTopology := &cvm.CpuTopology{}
		if coreCount, ok := cpuTopologyMap["core_count"]; ok && coreCount.(int) > 0 {
			cpuTopology.CoreCount = helper.IntInt64(coreCount.(int))
		}
		if threadPerCore, ok := cpuTopologyMap["thread_per_core"]; ok && threadPerCore.(int) > 0 {
			cpuTopology.ThreadPerCore = helper.IntInt64(threadPerCore.(int))
		}
		request.CpuTopology = cpuTopology
	}

	assert.Nil(t, request.CpuTopology)
}

func TestCpuTopologyReadStatePopulation(t *testing.T) {
	instance := &cvm.Instance{
		CpuTopology: &cvm.CpuTopology{
			CoreCount:     helper.Int64(4),
			ThreadPerCore: helper.Int64(2),
		},
	}

	var cpuTopologyResult []interface{}
	if instance.CpuTopology != nil {
		cpuTopologyMap := map[string]interface{}{}
		if instance.CpuTopology.CoreCount != nil {
			cpuTopologyMap["core_count"] = *instance.CpuTopology.CoreCount
		}
		if instance.CpuTopology.ThreadPerCore != nil {
			cpuTopologyMap["thread_per_core"] = *instance.CpuTopology.ThreadPerCore
		}
		cpuTopologyResult = []interface{}{cpuTopologyMap}
	}

	assert.Equal(t, 1, len(cpuTopologyResult))
	topologyMap := cpuTopologyResult[0].(map[string]interface{})
	assert.Equal(t, int64(4), topologyMap["core_count"])
	assert.Equal(t, int64(2), topologyMap["thread_per_core"])
}

func TestCpuTopologyReadStateNil(t *testing.T) {
	instance := &cvm.Instance{
		CpuTopology: nil,
	}

	var cpuTopologyResult []interface{}
	if instance.CpuTopology != nil {
		cpuTopologyMap := map[string]interface{}{}
		if instance.CpuTopology.CoreCount != nil {
			cpuTopologyMap["core_count"] = *instance.CpuTopology.CoreCount
		}
		if instance.CpuTopology.ThreadPerCore != nil {
			cpuTopologyMap["thread_per_core"] = *instance.CpuTopology.ThreadPerCore
		}
		cpuTopologyResult = []interface{}{cpuTopologyMap}
	}

	assert.Equal(t, 0, len(cpuTopologyResult))
}

func TestCpuTopologyReadStatePartialNil(t *testing.T) {
	instance := &cvm.Instance{
		CpuTopology: &cvm.CpuTopology{
			CoreCount:     helper.Int64(8),
			ThreadPerCore: nil,
		},
	}

	var cpuTopologyResult []interface{}
	if instance.CpuTopology != nil {
		cpuTopologyMap := map[string]interface{}{}
		if instance.CpuTopology.CoreCount != nil {
			cpuTopologyMap["core_count"] = *instance.CpuTopology.CoreCount
		}
		if instance.CpuTopology.ThreadPerCore != nil {
			cpuTopologyMap["thread_per_core"] = *instance.CpuTopology.ThreadPerCore
		}
		cpuTopologyResult = []interface{}{cpuTopologyMap}
	}

	assert.Equal(t, 1, len(cpuTopologyResult))
	topologyMap := cpuTopologyResult[0].(map[string]interface{})
	assert.Equal(t, int64(8), topologyMap["core_count"])
	_, hasThreadPerCore := topologyMap["thread_per_core"]
	assert.False(t, hasThreadPerCore)
}
