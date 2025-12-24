package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpcv20170312 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudVpcRoutePolicyEntries() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcRoutePolicyEntriesCreate,
		Read:   resourceTencentCloudVpcRoutePolicyEntriesRead,
		Update: resourceTencentCloudVpcRoutePolicyEntriesUpdate,
		Delete: resourceTencentCloudVpcRoutePolicyEntriesDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"route_policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the instance ID of the route reception policy.",
			},

			"route_policy_entry_set": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Route reception policy entry list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_policy_entry_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the unique ID of the IPv4 routing strategy entry.\nNote: This field may return null, indicating that no valid value was found.",
						},
						"cidr_block": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Destination ip range.\nNote: This field may return null, indicating that no valid value was found.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Describes the routing strategy rule.\nNote: This field may return null, indicating that no valid value was found.",
						},
						"route_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Routing Type\n\nSpecifies the USER-customized data type.\nNETD: specifies the route for network detection.\nCCN: CCN route.\nNote: This field may return null, indicating that no valid value was found.",
						},
						"gateway_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Next hop type. types currently supported:.\nCVM: cloud virtual machine with public network gateway type.\nVPN: vpn gateway.\nDIRECTCONNECT: direct connect gateway.\nPEERCONNECTION: peering connection.\nHAVIP: high availability virtual ip.\nNAT: specifies the nat gateway. \nEIP: specifies the public ip address of the cloud virtual machine.\nLOCAL_GATEWAY: specifies the local gateway.\nPVGW: pvgw gateway.\nNote: This field may return null, indicating that no valid value was found.",
						},
						"gateway_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Gateway unique ID.\nNote: This field may return null, indicating that no valid value was found.",
						},
						"priority": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Priority. a smaller value indicates a higher priority.\nNote: This field may return null, indicating that no valid value was found.",
						},
						"action": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Action.\nDROP: drop.\nDISABLE: receive and disable.\nACCEPT: receive and enable.\nNote: This field may return null, indicating that no valid value was found.",
						},
						"created_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.\n\nNote: This field may return null, indicating that no valid value was found.",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Specifies the region.\nNote: This field may return null, indicating that no valid value was found.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudVpcRoutePolicyEntriesCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy_entries.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request       = vpcv20170312.NewCreateRoutePolicyEntriesRequest()
		routePolicyId string
	)

	if v, ok := d.GetOk("route_policy_id"); ok {
		request.RoutePolicyId = helper.String(v.(string))
		routePolicyId = v.(string)
	}

	if v, ok := d.GetOk("route_policy_entry_set"); ok {
		for _, item := range v.(*schema.Set).List() {
			routePolicyEntrySetMap := item.(map[string]interface{})
			routePolicyEntry := vpcv20170312.RoutePolicyEntry{}
			if v, ok := routePolicyEntrySetMap["cidr_block"].(string); ok && v != "" {
				routePolicyEntry.CidrBlock = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["description"].(string); ok && v != "" {
				routePolicyEntry.Description = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["route_type"].(string); ok && v != "" {
				routePolicyEntry.RouteType = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["gateway_type"].(string); ok && v != "" {
				routePolicyEntry.GatewayType = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["gateway_id"].(string); ok && v != "" {
				routePolicyEntry.GatewayId = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["priority"].(int); ok {
				routePolicyEntry.Priority = helper.IntUint64(v)
			}

			if v, ok := routePolicyEntrySetMap["action"].(string); ok && v != "" {
				routePolicyEntry.Action = helper.String(v)
			}

			request.RoutePolicyEntrySet = append(request.RoutePolicyEntrySet, &routePolicyEntry)
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CreateRoutePolicyEntriesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create vpc route policy entries failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(routePolicyId)
	return resourceTencentCloudVpcRoutePolicyEntriesRead(d, meta)
}

func resourceTencentCloudVpcRoutePolicyEntriesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy_entries.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service       = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		routePolicyId = d.Id()
	)

	respData, err := service.DescribeVpcRoutePolicyEntriesById(ctx, routePolicyId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vpc_route_policy_entries` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("route_policy_id", routePolicyId)

	if len(respData) > 0 {
		routePolicyEntrySet := make([]map[string]interface{}, 0, len(respData))
		for _, item := range respData {
			entry := make(map[string]interface{})
			if item.RoutePolicyEntryId != nil {
				entry["route_policy_entry_id"] = *item.RoutePolicyEntryId
			}

			if item.CidrBlock != nil {
				entry["cidr_block"] = *item.CidrBlock
			}

			if item.Description != nil {
				entry["description"] = *item.Description
			}

			if item.RouteType != nil {
				entry["route_type"] = *item.RouteType
			}

			if item.GatewayType != nil {
				entry["gateway_type"] = *item.GatewayType
			}

			if item.GatewayId != nil {
				entry["gateway_id"] = *item.GatewayId
			}

			if item.Priority != nil {
				entry["priority"] = *item.Priority
			}

			if item.Action != nil {
				entry["action"] = *item.Action
			}

			if item.CreatedTime != nil {
				entry["created_time"] = *item.CreatedTime
			}

			if item.Region != nil {
				entry["region"] = *item.Region
			}

			routePolicyEntrySet = append(routePolicyEntrySet, entry)
		}

		_ = d.Set("route_policy_entry_set", routePolicyEntrySet)
	}

	return nil
}

func resourceTencentCloudVpcRoutePolicyEntriesUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy_entries.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service       = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request       = vpcv20170312.NewResetRoutePolicyEntriesRequest()
		routePolicyId = d.Id()
	)

	// temp get RoutePolicyDescription and RoutePolicyName
	respData, err := service.DescribeVpcRoutePolicyById(ctx, routePolicyId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_vpc_route_policy_entries` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return fmt.Errorf("resource `tencentcloud_vpc_route_policy_entries` [%s] not found, please check if it has been deleted.", d.Id())
	}

	if respData.RoutePolicyDescription != nil {
		request.RoutePolicyDescription = respData.RoutePolicyDescription
	}

	if respData.RoutePolicyName != nil {
		request.RoutePolicyName = respData.RoutePolicyName
	}

	if v, ok := d.GetOk("route_policy_entry_set"); ok {
		for _, item := range v.(*schema.Set).List() {
			routePolicyEntrySetMap := item.(map[string]interface{})
			routePolicyEntry := vpcv20170312.RoutePolicyEntry{}
			if v, ok := routePolicyEntrySetMap["route_policy_entry_id"].(string); ok && v != "" {
				routePolicyEntry.RoutePolicyEntryId = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["cidr_block"].(string); ok && v != "" {
				routePolicyEntry.CidrBlock = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["description"].(string); ok && v != "" {
				routePolicyEntry.Description = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["route_type"].(string); ok && v != "" {
				routePolicyEntry.RouteType = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["gateway_type"].(string); ok && v != "" {
				routePolicyEntry.GatewayType = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["gateway_id"].(string); ok && v != "" {
				routePolicyEntry.GatewayId = helper.String(v)
			}

			if v, ok := routePolicyEntrySetMap["priority"].(int); ok {
				routePolicyEntry.Priority = helper.IntUint64(v)
			}

			if v, ok := routePolicyEntrySetMap["action"].(string); ok && v != "" {
				routePolicyEntry.Action = helper.String(v)
			}

			request.RoutePolicyEntrySet = append(request.RoutePolicyEntrySet, &routePolicyEntry)
		}
	}

	request.RoutePolicyId = &routePolicyId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().ResetRoutePolicyEntriesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update vpc route policy entries failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudVpcRoutePolicyEntriesRead(d, meta)
}

func resourceTencentCloudVpcRoutePolicyEntriesDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_route_policy_entries.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		ctx           = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service       = VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		request       = vpcv20170312.NewDeleteRoutePolicyEntriesRequest()
		routePolicyId = d.Id()
	)

	// get all entryIds first
	respData, err := service.DescribeVpcRoutePolicyEntriesById(ctx, routePolicyId)
	if err != nil {
		return err
	}

	if respData == nil || len(respData) == 0 {
		return nil
	}

	for _, item := range respData {
		routePolicyEntry := vpcv20170312.RoutePolicyEntry{}
		if item.RoutePolicyEntryId != nil {
			routePolicyEntry.RoutePolicyEntryId = item.RoutePolicyEntryId
		}

		request.RoutePolicyEntrySet = append(request.RoutePolicyEntrySet, &routePolicyEntry)
	}

	request.RoutePolicyId = &routePolicyId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DeleteRoutePolicyEntriesWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete vpc route policy entries failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
