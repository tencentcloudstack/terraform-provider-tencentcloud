/*
Use this data source to query instances type.

Example Usage

```hcl
data "tencentcloud_instance_types" "foo" {
  availability_zone = "ap-guangzhou-2"
  cpu_core_count    = 2
  memory_size       = 4
}

data tencentcloud_instance_types "t1c1g" {
  cpu_core_count    = 1
  memory_size       = 1
  exclude_sold_out=true
  filter {
    name   = "instance-charge-type"
    values = ["POSTPAID_BY_HOUR"]
  }
  filter {
    name   = "zone"
    values = ["ap-shanghai-2"]
  }
}
```
*/
package tencentcloud

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"cpu_core_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of CPU cores of the instance.",
			},
			"gpu_core_count": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The number of GPU cores of the instance.",
			},
			"memory_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Instance memory capacity, unit in GB.",
			},
			"availability_zone": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"filter"},
				Description:   "The available zone that the CVM instance locates at. This field is conflict with `filter`.",
			},
			"filter": {
				Type:          schema.TypeSet,
				Optional:      true,
				MaxItems:      10,
				ConflictsWith: []string{"availability_zone"},
				Description:   "One or more name/value pairs to filter. This field is conflict with `availability_zone`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The filter name. Valid values: `zone`, `instance-family` and `instance-charge-type`.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The filter values.",
						},
					},
				},
			},
			"exclude_sold_out": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicate to filter instances types that is sold out or not, default is false.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed values.
			"instance_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of cvm instance. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"availability_zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The available zone that the CVM instance locates at.",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the instance.",
						},
						"cpu_core_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of CPU cores of the instance.",
						},
						"gpu_core_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of GPU cores of the instance.",
						},
						"memory_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance memory capacity, unit in GB.",
						},
						"family": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type series of the instance.",
						},
						"instance_charge_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Charge type of the instance.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Sell status of the instance.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_instance_types.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	isExcludeSoldOut := d.Get("exclude_sold_out").(bool)
	cpu, cpuOk := d.GetOk("cpu_core_count")
	gpu, gpuOk := d.GetOk("gpu_core_count")
	memory, memoryOk := d.GetOk("memory_size")
	var instanceSellTypes []*cvm.InstanceTypeQuotaItem
	var errRet error
	var err error
	typeList := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

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
	err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		instanceSellTypes, errRet = cvmService.DescribeInstancesSellTypeByFilter(ctx, filterMap)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
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

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), typeList); err != nil {
			return err
		}
	}
	return nil
}
