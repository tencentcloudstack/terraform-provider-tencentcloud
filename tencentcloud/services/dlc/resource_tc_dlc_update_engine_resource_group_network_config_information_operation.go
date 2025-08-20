package dlc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dlcv20210125 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dlc/v20210125"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudDlcUpdateEngineResourceGroupNetworkConfigInformationOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudDlcUpdateEngineResourceGroupNetworkConfigInformationOperationCreate,
		Read:   resourceTencentCloudDlcUpdateEngineResourceGroupNetworkConfigInformationOperationRead,
		Delete: resourceTencentCloudDlcUpdateEngineResourceGroupNetworkConfigInformationOperationDelete,
		Schema: map[string]*schema.Schema{
			"engine_resource_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Engine resource group ID.",
			},

			"network_config_names": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Description: "A collection of network configuration names bound to the resource group.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTencentCloudDlcUpdateEngineResourceGroupNetworkConfigInformationOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_update_engine_resource_group_network_config_information_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                 = tccommon.GetLogId(tccommon.ContextNil)
		ctx                   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request               = dlcv20210125.NewUpdateEngineResourceGroupNetworkConfigInfoRequest()
		engineResourceGroupId string
	)

	if v, ok := d.GetOk("engine_resource_group_id"); ok {
		request.EngineResourceGroupId = helper.String(v.(string))
		engineResourceGroupId = v.(string)
	}

	if v, ok := d.GetOk("network_config_names"); ok {
		networkConfigNamesSet := v.(*schema.Set).List()
		for i := range networkConfigNamesSet {
			networkConfigNames := networkConfigNamesSet[i].(string)
			request.NetworkConfigNames = append(request.NetworkConfigNames, helper.String(networkConfigNames))
		}
	}

	request.IsEffectiveNow = helper.IntInt64(0)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().UpdateEngineResourceGroupNetworkConfigInfoWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create dlc update engine resource group network config information operation failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	d.SetId(engineResourceGroupId)

	// wait
	waitReq := dlcv20210125.NewDescribeStandardEngineResourceGroupsRequest()
	waitReq.Filters = []*dlcv20210125.Filter{
		{
			Name:   helper.String("engine-resource-group-id"),
			Values: helper.Strings([]string{engineResourceGroupId}),
		},
	}
	reqErr = resource.Retry(tccommon.ReadRetryTimeout*7, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseDlcClient().DescribeStandardEngineResourceGroupsWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe standard engine resource groups failed, Response is nil."))
		}

		if result.Response.UserEngineResourceGroupInfos == nil || len(result.Response.UserEngineResourceGroupInfos) == 0 {
			return resource.NonRetryableError(fmt.Errorf("UserEngineResourceGroupInfos is nil."))
		}

		if result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState == nil {
			return resource.NonRetryableError(fmt.Errorf("ResourceGroupState is nil."))
		}

		resourceGroupState := result.Response.UserEngineResourceGroupInfos[0].ResourceGroupState
		if *resourceGroupState == 2 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("ResourceGroupState is not running, current state is %d...", *resourceGroupState))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s describe standard engine resource groups failed,, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudDlcUpdateEngineResourceGroupNetworkConfigInformationOperationRead(d, meta)
}

func resourceTencentCloudDlcUpdateEngineResourceGroupNetworkConfigInformationOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_update_engine_resource_group_network_config_information_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudDlcUpdateEngineResourceGroupNetworkConfigInformationOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_dlc_update_engine_resource_group_network_config_information_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
