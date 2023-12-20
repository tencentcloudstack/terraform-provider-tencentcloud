package dayuv2

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDayuEip() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuEipRead,
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the resource.",
			},
			"bind_status": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tccommon.ValidateAllowedStringValue(DDOS_EIP_BIND_STATUS),
				},
				Optional:    true,
				Description: "The binding state of the instance, value range [BINDING, BIND, UNBINDING, UNBIND], default is [BINDING, BIND, UNBINDING, UNBIND].",
			},
			"offset": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "The page start offset, default is `0`.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     10,
				Description: "The number of pages, default is `10`.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of layer 4 rules. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"eip_list": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "",
							},
						},
						"instance_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "ID of the resource instance.",
						},
						"eip_bound_rsc_ins": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The ID of the resource instance for the binding.",
						},
						"eip_bound_rsc_eni": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the bound ENI.",
						},
						"eip_bound_rsc_vip": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Bind the resource intranet IP.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The created time of resource.",
						},
						"expired_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The expired time of resource.",
						},
						"modify_time": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The modify time of resource.",
						},
						"protection_status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The protection status of the asset instance.",
						},
						"eip_address_status": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Eip PUBLIC IP status.",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The region where the asset instance is located.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuEipRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dayu_l4_rules.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	resourceId := d.Get("resource_id").(string)
	var bindStatus []string
	if v, ok := d.GetOk("bind_status"); ok {
		tmpBindStatusList := v.([]interface{})
		for _, tmpBindStatus := range tmpBindStatusList {
			bindStatus = append(bindStatus, tmpBindStatus.(string))
		}
	} else {
		bindStatus = DDOS_EIP_BIND_STATUS
	}
	offset := d.Get("offset").(int)
	limit := d.Get("limit").(int)

	antiddosService := AntiddosService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	result, err := antiddosService.DescribeListBGPIPInstances(ctx, resourceId, bindStatus, offset, limit)
	if err != nil {
		return err
	}
	resultList := make([]map[string]interface{}, 0)
	for _, eipItem := range result {
		tmpEipInfo := make(map[string]interface{})
		eipList := make([]string, 0)
		for _, eip := range eipItem.InstanceDetail.EipList {
			eipList = append(eipList, *eip)
		}
		tmpEipInfo["eip_list"] = eipList
		tmpEipInfo["instance_id"] = *eipItem.InstanceDetail.InstanceId
		tmpEipInfo["region"] = *eipItem.Region.Region
		tmpEipInfo["eip_bound_rsc_ins"] = *eipItem.EipAddressInfo.EipBoundRscIns
		tmpEipInfo["eip_bound_rsc_eni"] = *eipItem.EipAddressInfo.EipBoundRscEni
		tmpEipInfo["eip_bound_rsc_vip"] = *eipItem.EipAddressInfo.EipBoundRscVip
		tmpEipInfo["eip_address_status"] = *eipItem.EipAddressStatus
		tmpEipInfo["protection_status"] = *eipItem.Status
		tmpEipInfo["created_time"] = *eipItem.CreatedTime
		tmpEipInfo["expired_time"] = *eipItem.ExpiredTime
		tmpEipInfo["modify_time"] = *eipItem.EipAddressInfo.ModifyTime

		resultList = append(resultList, tmpEipInfo)
	}
	ids := make([]string, 0, len(resultList))
	for _, listItem := range resultList {
		ids = append(ids, listItem["instance_id"].(string))
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	_ = d.Set("list", resultList)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), resultList)
	}
	return nil

}
