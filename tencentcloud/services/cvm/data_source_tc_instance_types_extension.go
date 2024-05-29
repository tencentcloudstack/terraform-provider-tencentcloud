package cvm

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func dataSourceTencentCloudInstanceTypesReadPostFillRequest0(ctx context.Context, req map[string]interface{}) error {
	d := tccommon.ResourceDataFromContext(ctx)
	var zone string
	var zone_in = 0
	if v, ok := d.GetOk("availability_zone"); ok {
		zone = v.(string)
		zone_in = 1
	}
	filters := d.Get("filter").(*schema.Set).List()
	filterMap := make(map[string][]string, len(filters)+zone_in)
	for _, v := range filters {
		item := v.(map[string]interface{})
		name := item["name"].(string)
		values := item["values"].([]interface{})
		filterValues := make([]string, 0, len(values))
		for _, value := range values {
			filterValues = append(filterValues, value.(string))
		}
		filterMap[name] = filterValues
	}
	if zone != "" {
		filterMap["zone"] = []string{zone}
	}

	return nil
}

func dataSourceTencentCloudInstanceTypesReadPostHandleResponse0(ctx context.Context, req map[string]interface{}, resp *cvm.DescribeZoneInstanceConfigInfosResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)
	logId := tccommon.GetLogId(tccommon.ContextNil)
	isExcludeSoldOut := d.Get("exclude_sold_out").(bool)
	cpu, cpuOk := d.GetOk("cpu_core_count")
	gpu, gpuOk := d.GetOk("gpu_core_count")
	memory, memoryOk := d.GetOk("memory_size")
	var err error
	typeList := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

	instanceSellTypes := resp.InstanceTypeQuotaSet
	for _, instanceType := range instanceSellTypes {
		flag := true
		if cpuOk && int64(cpu.(int)) != *instanceType.Cpu {
			flag = false
		}
		if gpuOk && int64(gpu.(int)) != *instanceType.Gpu {
			flag = false
		}
		if memoryOk && int64(memory.(int)) != *instanceType.Memory {
			flag = false
		}
		if isExcludeSoldOut && CVM_SOLD_OUT_STATUS == *instanceType.Status {
			flag = false
		}

		if flag {
			mapping := map[string]interface{}{
				"availability_zone":    instanceType.Zone,
				"cpu_core_count":       instanceType.Cpu,
				"gpu_core_count":       instanceType.Gpu,
				"memory_size":          instanceType.Memory,
				"family":               instanceType.InstanceFamily,
				"instance_type":        instanceType.InstanceType,
				"instance_charge_type": instanceType.InstanceChargeType,
				"status":               instanceType.Status,
			}
			typeList = append(typeList, mapping)
			ids = append(ids, *instanceType.InstanceType)
		}
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("instance_types", typeList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set instance type list fail, reason:%s\n ", logId, err.Error())
		return err
	}
	context.WithValue(ctx, "typeList", typeList)
	return nil
}

func dataSourceTencentCloudInstanceTypesReadOutputContent(ctx context.Context) interface{} {
	typeList := ctx.Value("typeList").([]interface{})
	return typeList
}
