package cvm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudPlacementGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudPlacementGroupsRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the placement group to be queried.",
			},

			"placement_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "ID of the placement group to be queried.",
			},

			"placement_group_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of placement group. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the placement group.",
						},
						"current_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of hosts in the placement group.",
						},
						"cvm_quota_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum number of hosts in the placement group.",
						},
						"instance_ids": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Host IDs in the placement group.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the placement group.",
						},
						"placement_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the placement group.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the placement group.",
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

func dataSourceTencentCloudPlacementGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_placement_groups.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(nil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := CvmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("placement_group_id"); ok {
		paramMap["DisasterRecoverGroupIds"] = []*string{helper.String(v.(string))}
	}

	if v, ok := d.GetOk("name"); ok {
		paramMap["Name"] = helper.String(v.(string))
	}

	var respData *cvm.DescribeDisasterRecoverGroupsResponseParams
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribePlacementGroupsByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	if err := dataSourceTencentCloudPlacementGroupsReadPostHandleResponse0(ctx, paramMap, respData); err != nil {
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), dataSourceTencentCloudPlacementGroupsReadOutputContent(ctx)); e != nil {
			return e
		}
	}

	return nil
}
