package dlc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDlcEngineNodeSpecifications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDlcEngineNodeSpecificationsRead,
		Schema: map[string]*schema.Schema{
			"data_engine_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Engine Name.",
			},

			"driver_spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Driver available specifications.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specification name.",
						},
						"cu": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Current specification of CU number.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Current CPU specifications.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The current memory size, in GB.",
						},
					},
				},
			},

			"executor_spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Available executor specifications.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specification name.",
						},
						"cu": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Current specification of CU number.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Current CPU specifications.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The current memory size, in GB.",
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

func dataSourceTencentCloudDlcEngineNodeSpecificationsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dlc_engine_node_specifications.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = DlcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("data_engine_name"); ok {
		paramMap["DataEngineName"] = helper.String(v.(string))
	}

	var respData *dlcv20210125.DescribeEngineNodeSpecResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDlcEngineNodeSpecificationsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	driverSpecList := make([]map[string]interface{}, 0, len(respData.DriverSpec))
	if respData.DriverSpec != nil {
		for _, driverSpec := range respData.DriverSpec {
			driverSpecMap := map[string]interface{}{}
			if driverSpec.Name != nil {
				driverSpecMap["name"] = driverSpec.Name
			}

			if driverSpec.Cu != nil {
				driverSpecMap["cu"] = driverSpec.Cu
			}

			if driverSpec.Cpu != nil {
				driverSpecMap["cpu"] = driverSpec.Cpu
			}

			if driverSpec.Memory != nil {
				driverSpecMap["memory"] = driverSpec.Memory
			}

			driverSpecList = append(driverSpecList, driverSpecMap)
		}

		_ = d.Set("driver_spec", driverSpecList)
	}

	executorSpecList := make([]map[string]interface{}, 0, len(respData.ExecutorSpec))
	if respData.ExecutorSpec != nil {
		for _, executorSpec := range respData.ExecutorSpec {
			executorSpecMap := map[string]interface{}{}
			if executorSpec.Name != nil {
				executorSpecMap["name"] = executorSpec.Name
			}

			if executorSpec.Cu != nil {
				executorSpecMap["cu"] = executorSpec.Cu
			}

			if executorSpec.Cpu != nil {
				executorSpecMap["cpu"] = executorSpec.Cpu
			}

			if executorSpec.Memory != nil {
				executorSpecMap["memory"] = executorSpec.Memory
			}

			executorSpecList = append(executorSpecList, executorSpecMap)
		}

		_ = d.Set("executor_spec", executorSpecList)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
