package cfw

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cfwv20190904 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cfw/v20190904"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudCfwClusterFwBypassConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCfwClusterFwBypassConfigCreate,
		Read:   resourceTencentCloudCfwClusterFwBypassConfigRead,
		Update: resourceTencentCloudCfwClusterFwBypassConfigUpdate,
		Delete: resourceTencentCloudCfwClusterFwBypassConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"fw_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Firewall type. `VPC_FW` - VPC firewall, `NAT_FW` - NAT firewall.",
			},
			"ccn_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "CCN instance ID.",
			},
			"enable": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "Bypass switch. `true` - enable Bypass (traffic bypasses firewall), `false` - disable Bypass (traffic goes through firewall).",
			},
			"nat_ins_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "NAT firewall instance ID. Required when fw_type is `NAT_FW`.",
			},
			"nat_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "NAT firewall type filter. `nat` - VPC internal protection type, `nat_ccn` - CCN cluster mode type. If not specified, both types are queried.",
			},
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter condition list. Supports filtering by Common (general search), InsObj (instance ID), ObjName (instance name), etc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter key.",
						},
						"operator_type": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Operator type. 1: equal, 2: greater than, 3: less than, 4: greater than or equal, 5: less than or equal, 6: not equal, 8: not in, 9: fuzzy match.",
						},
						"values": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Filter values.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Total number of records matching the conditions.",
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "NAT firewall switch detail list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ins_obj": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "NAT instance ID.",
						},
						"obj_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"fw_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Firewall type.",
						},
						"asset_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Asset type. `nat` - VPC internal protection, `nat_ccn` - CCN cluster mode.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"switch_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Switch access mode. 1: automatic, 2: manual.",
						},
						"routing_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Traffic routing method. 0: multi-route table, 1: policy routing.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Switch status. 0: not enabled, 1: enabled, 2: enabling, 3: disabling.",
						},
						"bypass": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bypass status. 0: disabled, 1: enabled.",
						},
						"attach_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Associated ID. For nat_ccn asset type, this is the CCN instance ID.",
						},
						"attach_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Associated ID instance name. For nat_ccn asset type, this is the CCN name.",
						},
					},
				},
			},
			"region_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Region list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"text": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Display name.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Actual value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudCfwClusterFwBypassConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_fw_bypass.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	fwType := d.Get("fw_type").(string)
	ccnId := d.Get("ccn_id").(string)

	var id string
	if v, ok := d.GetOk("nat_ins_id"); ok && fwType == "NAT_FW" {
		natInsId := v.(string)
		id = strings.Join([]string{fwType, ccnId, natInsId}, tccommon.FILED_SP)
	} else {
		id = strings.Join([]string{fwType, ccnId}, tccommon.FILED_SP)
	}

	d.SetId(id)
	return resourceTencentCloudCfwClusterFwBypassConfigUpdate(d, meta)
}

func resourceTencentCloudCfwClusterFwBypassConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_fw_bypass.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cfwv20190904.NewDescribeClusterNatCcnFwSwitchListRequest()
		id      = d.Id()
	)

	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) < 2 {
		return fmt.Errorf("invalid id format: %s", id)
	}
	fwType := parts[0]
	ccnId := parts[1]
	var natInsId string
	if len(parts) >= 3 {
		natInsId = parts[2]
	}

	request.Limit = helper.IntInt64(100)
	request.Offset = helper.IntInt64(0)

	if v, ok := d.GetOk("nat_type"); ok {
		request.NatType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		for _, item := range v.([]interface{}) {
			filterMap := item.(map[string]interface{})
			commonFilter := cfwv20190904.CommonFilter{}
			if name, ok := filterMap["name"].(string); ok && name != "" {
				commonFilter.Name = helper.String(name)
			}
			if opType, ok := filterMap["operator_type"].(int); ok {
				commonFilter.OperatorType = helper.IntInt64(opType)
			}
			if vals, ok := filterMap["values"]; ok {
				for _, v := range vals.([]interface{}) {
					commonFilter.Values = append(commonFilter.Values, helper.String(v.(string)))
				}
			}
			request.Filters = append(request.Filters, &commonFilter)
		}
	}

	var response *cfwv20190904.DescribeClusterNatCcnFwSwitchListResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().DescribeClusterNatCcnFwSwitchListWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read cfw_cluster_fw_bypass failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil {
		log.Printf("[WARN]%s resource `cfw_cluster_fw_bypass` id=%s not found, response is nil", logId, d.Id())
		d.SetId("")
		return nil
	}

	// find matching instance
	var matchedItem *cfwv20190904.NatFwSwitchDetailS
	for _, item := range response.Response.Data {
		if item == nil {
			continue
		}
		// match by ccn_id (AttachId for nat_ccn type) or InsObj
		if fwType == "NAT_FW" && natInsId != "" {
			if item.InsObj != nil && *item.InsObj == natInsId &&
				item.AttachId != nil && *item.AttachId == ccnId {
				matchedItem = item
				break
			}
		} else {
			// VPC_FW: match by AttachId == ccnId
			if item.AttachId != nil && *item.AttachId == ccnId {
				matchedItem = item
				break
			}
		}
	}

	if matchedItem == nil {
		log.Printf("[WARN]%s resource `cfw_cluster_fw_bypass` id=%s not found in response data", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("fw_type", fwType)
	_ = d.Set("ccn_id", ccnId)
	if natInsId != "" {
		_ = d.Set("nat_ins_id", natInsId)
	}

	if response.Response.Total != nil {
		_ = d.Set("total", response.Response.Total)
	}

	dataList := make([]map[string]interface{}, 0, len(response.Response.Data))
	for _, item := range response.Response.Data {
		if item == nil {
			continue
		}
		itemMap := map[string]interface{}{}
		if item.InsObj != nil {
			itemMap["ins_obj"] = item.InsObj
		}
		if item.ObjName != nil {
			itemMap["obj_name"] = item.ObjName
		}
		if item.FwType != nil {
			itemMap["fw_type"] = item.FwType
		}
		if item.AssetType != nil {
			itemMap["asset_type"] = item.AssetType
		}
		if item.Region != nil {
			itemMap["region"] = item.Region
		}
		if item.SwitchMode != nil {
			itemMap["switch_mode"] = item.SwitchMode
		}
		if item.RoutingMode != nil {
			itemMap["routing_mode"] = item.RoutingMode
		}
		if item.Status != nil {
			itemMap["status"] = item.Status
		}
		if item.Bypass != nil {
			itemMap["bypass"] = item.Bypass
		}
		if item.AttachId != nil {
			itemMap["attach_id"] = item.AttachId
		}
		if item.AttachName != nil {
			itemMap["attach_name"] = item.AttachName
		}
		dataList = append(dataList, itemMap)
	}
	_ = d.Set("data", dataList)

	if response.Response.RegionList != nil {
		regionList := make([]map[string]interface{}, 0, len(response.Response.RegionList))
		for _, region := range response.Response.RegionList {
			if region == nil {
				continue
			}
			regionMap := map[string]interface{}{}
			if region.Text != nil {
				regionMap["text"] = region.Text
			}
			if region.Value != nil {
				regionMap["value"] = region.Value
			}
			regionList = append(regionList, regionMap)
		}
		_ = d.Set("region_list", regionList)
	}

	return nil
}

func resourceTencentCloudCfwClusterFwBypassConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_fw_bypass.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cfwv20190904.NewModifyClusterFwBypassRequest()
		id      = d.Id()
	)

	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) < 2 {
		return fmt.Errorf("invalid id format: %s", id)
	}
	fwType := parts[0]
	ccnId := parts[1]

	request.FwType = helper.String(fwType)
	request.CcnId = helper.String(ccnId)

	if v, ok := d.GetOkExists("enable"); ok {
		request.Enable = helper.Bool(v.(bool))
	}

	if fwType == "NAT_FW" {
		if len(parts) >= 3 {
			request.NatInsId = helper.String(parts[2])
		} else if v, ok := d.GetOk("nat_ins_id"); ok {
			request.NatInsId = helper.String(v.(string))
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().ModifyClusterFwBypassWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update cfw_cluster_fw_bypass failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudCfwClusterFwBypassConfigRead(d, meta)
}

func resourceTencentCloudCfwClusterFwBypassConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_fw_bypass.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	d.SetId("")
	return nil
}
