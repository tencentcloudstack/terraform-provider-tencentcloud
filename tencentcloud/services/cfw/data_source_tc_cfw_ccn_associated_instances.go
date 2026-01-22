package cfw

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCfwCcnAssociatedInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCfwCcnAssociatedInstancesRead,
		Schema: map[string]*schema.Schema{
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CCN ID.",
			},

			"ccn_associated_instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Information of instances associated with CCN.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"ins_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type.",
						},
						"cidr_lst": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "List of network segments for the instance.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"instance_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region where the instance belongs.",
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

func dataSourceTencentCloudCfwCcnAssociatedInstancesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cfw_ccn_associated_instances.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = CfwService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		ccnId   string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("ccn_id"); ok {
		paramMap["CcnId"] = helper.String(v.(string))
		ccnId = v.(string)
	}

	var respData []*cfwv20190904.CcnAssociatedInstance
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCfwCcnAssociatedInstancesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	ccnAssociatedInstancesList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, ccnAssociatedInstances := range respData {
			ccnAssociatedInstancesMap := map[string]interface{}{}
			if ccnAssociatedInstances.InstanceId != nil {
				ccnAssociatedInstancesMap["instance_id"] = ccnAssociatedInstances.InstanceId
			}

			if ccnAssociatedInstances.InstanceName != nil {
				ccnAssociatedInstancesMap["instance_name"] = ccnAssociatedInstances.InstanceName
			}

			if ccnAssociatedInstances.InsType != nil {
				ccnAssociatedInstancesMap["ins_type"] = ccnAssociatedInstances.InsType
			}

			if ccnAssociatedInstances.CidrLst != nil {
				ccnAssociatedInstancesMap["cidr_lst"] = ccnAssociatedInstances.CidrLst
			}

			if ccnAssociatedInstances.InstanceRegion != nil {
				ccnAssociatedInstancesMap["instance_region"] = ccnAssociatedInstances.InstanceRegion
			}

			ccnAssociatedInstancesList = append(ccnAssociatedInstancesList, ccnAssociatedInstancesMap)
		}

		_ = d.Set("ccn_associated_instances", ccnAssociatedInstancesList)
	}

	d.SetId(ccnId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
