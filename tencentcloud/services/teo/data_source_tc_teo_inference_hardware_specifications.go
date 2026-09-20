package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoInferenceHardwareSpecifications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoInferenceHardwareSpecificationsRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone ID.",
			},

			"hardware_specifications": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Inference hardware specification list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"spec": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specification identifier. Deprecated, refer to `hardware_spec_id`.",
						},
						"hardware_spec_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specification unique identifier ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specification name.",
						},
						"gpu_num": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Number of GPU cards allocated by default for the specification.",
						},
						"cpu_num": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Number of CPU cores allocated by default for the specification.",
						},
						"mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory size allocated by default for the specification. Unit: MB.",
						},
						"gpu_mem_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "GPU memory size allocated by default for the specification. Unit: MB.",
						},
						"disk_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Disk size allocated by default for the specification. Unit: MB.",
						},
						"allowed_gpu_nums": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of GPU card counts currently supported by the specification.",
							Elem: &schema.Schema{
								Type: schema.TypeFloat,
							},
						},
					},
				},
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoInferenceHardwareSpecificationsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_inference_hardware_specifications.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
	}

	var respData []*teov20220901.InferenceHardwareSpecification
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoInferenceHardwareSpecificationsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		if result == nil || len(result) < 1 {
			log.Printf("[DATASOURCE] read empty, skip SetId, teo_inference_hardware_specifications zone_id=%v", paramMap["ZoneId"])
			return resource.NonRetryableError(fmt.Errorf("teo_inference_hardware_specifications response is empty, zone_id=%v", paramMap["ZoneId"]))
		}
		respData = result
		return nil
	})
	if reqErr != nil {
		return reqErr
	}

	ids := make([]string, 0, len(respData))
	hardwareSpecificationsList := make([]map[string]interface{}, 0, len(respData))
	for _, hardwareSpecification := range respData {
		hardwareSpecificationMap := map[string]interface{}{}

		if hardwareSpecification.Spec != nil {
			hardwareSpecificationMap["spec"] = hardwareSpecification.Spec
		}

		if hardwareSpecification.HardwareSpecId != nil {
			hardwareSpecificationMap["hardware_spec_id"] = hardwareSpecification.HardwareSpecId
			ids = append(ids, *hardwareSpecification.HardwareSpecId)
		}

		if hardwareSpecification.Name != nil {
			hardwareSpecificationMap["name"] = hardwareSpecification.Name
		}

		if hardwareSpecification.GPUNum != nil {
			hardwareSpecificationMap["gpu_num"] = *hardwareSpecification.GPUNum
		}

		if hardwareSpecification.CPUNum != nil {
			hardwareSpecificationMap["cpu_num"] = *hardwareSpecification.CPUNum
		}

		if hardwareSpecification.MemSize != nil {
			hardwareSpecificationMap["mem_size"] = int(*hardwareSpecification.MemSize)
		}

		if hardwareSpecification.GPUMemSize != nil {
			hardwareSpecificationMap["gpu_mem_size"] = int(*hardwareSpecification.GPUMemSize)
		}

		if hardwareSpecification.DiskSize != nil {
			hardwareSpecificationMap["disk_size"] = int(*hardwareSpecification.DiskSize)
		}

		if hardwareSpecification.AllowedGPUNums != nil {
			allowedGpuNumsList := make([]interface{}, 0, len(hardwareSpecification.AllowedGPUNums))
			for _, allowedGpuNum := range hardwareSpecification.AllowedGPUNums {
				if allowedGpuNum != nil {
					allowedGpuNumsList = append(allowedGpuNumsList, *allowedGpuNum)
				}
			}
			hardwareSpecificationMap["allowed_gpu_nums"] = allowedGpuNumsList
		}

		hardwareSpecificationsList = append(hardwareSpecificationsList, hardwareSpecificationMap)
	}

	_ = d.Set("hardware_specifications", hardwareSpecificationsList)

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), hardwareSpecificationsList); e != nil {
			return e
		}
	}

	return nil
}
