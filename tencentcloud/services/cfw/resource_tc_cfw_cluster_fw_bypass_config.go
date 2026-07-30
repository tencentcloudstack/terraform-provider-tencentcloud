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
		},
	}
}

func resourceTencentCloudCfwClusterFwBypassConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_fw_bypass_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	fwType := d.Get("fw_type").(string)
	ccnId := d.Get("ccn_id").(string)

	var id string
	if fwType == "NAT_FW" {
		natInsId := d.Get("nat_ins_id").(string)
		if natInsId == "" {
			return fmt.Errorf("`nat_ins_id` is required when fw_type is `NAT_FW`")
		}

		res := strings.Join([]string{natInsId, ccnId}, tccommon.COMMA_SP)
		id = strings.Join([]string{fwType, res}, tccommon.FILED_SP)
	} else if fwType == "VPC_FW" {
		id = strings.Join([]string{fwType, ccnId}, tccommon.FILED_SP)
	} else {
		return fmt.Errorf("invalid fw_type: %s", fwType)
	}

	d.SetId(id)
	return resourceTencentCloudCfwClusterFwBypassConfigUpdate(d, meta)
}

func resourceTencentCloudCfwClusterFwBypassConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_fw_bypass_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cfwv20190904.NewDescribeClusterNatCcnFwSwitchListRequest()
		id      = d.Id()
	)

	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 {
		return fmt.Errorf("invalid id format: %s", id)
	}

	fwType := parts[0]
	tmp := parts[1]
	var natInsId string
	var ccnId string
	if fwType == "NAT_FW" {
		res := strings.Split(tmp, tccommon.COMMA_SP)
		if len(res) != 2 {
			return fmt.Errorf("invalid id format: %s", id)
		}

		natInsId = res[0]
		ccnId = res[1]
	} else if fwType == "VPC_FW" {
		ccnId = tmp
	} else {
		return fmt.Errorf("invalid fw_type: %s", fwType)
	}

	request.Limit = helper.IntInt64(100)
	request.Offset = helper.IntInt64(0)
	request.NatType = helper.String("nat_ccn")
	if natInsId != "" {
		request.Filters = []*cfwv20190904.CommonFilter{
			{
				Name:         helper.String("InsObj"),
				OperatorType: helper.IntInt64(1),
				Values:       []*string{helper.String(natInsId)},
			},
			{
				Name:         helper.String("AttachId"),
				OperatorType: helper.IntInt64(1),
				Values:       []*string{helper.String(ccnId)},
			},
		}
	} else {
		request.Filters = []*cfwv20190904.CommonFilter{
			{
				Name:         helper.String("AttachId"),
				OperatorType: helper.IntInt64(1),
				Values:       []*string{helper.String(ccnId)},
			},
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
		log.Printf("[CRITAL]%s read cfw_cluster_fw_bypass_config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response == nil || response.Response == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cfw_cluster_fw_bypass_config` id=%s not found, response is nil", logId, d.Id())
		d.SetId("")
		return nil
	}

	// find matching instance
	var matchedItem *cfwv20190904.NatFwSwitchDetailS
	for _, item := range response.Response.Data {
		if item == nil {
			continue
		}

		if fwType == "NAT_FW" && natInsId != "" {
			if item.InsObj != nil && *item.InsObj == natInsId &&
				item.AttachId != nil && *item.AttachId == ccnId {
				matchedItem = item
				break
			}
		} else {
			if item.AttachId != nil && *item.AttachId == ccnId {
				matchedItem = item
				break
			}
		}
	}

	if matchedItem == nil {
		log.Printf("[WARN]%s resource `tencentcloud_cfw_cluster_fw_bypass_config` id=%s not found in response data", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("fw_type", fwType)
	_ = d.Set("ccn_id", ccnId)
	if natInsId != "" {
		_ = d.Set("nat_ins_id", natInsId)
	}

	// Set enable based on Bypass field: 1 means bypass enabled (true), 0 means bypass disabled (false)
	if matchedItem.Bypass != nil {
		if *matchedItem.Bypass == 1 {
			_ = d.Set("enable", true)
		} else {
			_ = d.Set("enable", false)
		}
	}

	return nil
}

func resourceTencentCloudCfwClusterFwBypassConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_fw_bypass_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = cfwv20190904.NewModifyClusterFwBypassRequest()
		id      = d.Id()
	)

	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 {
		return fmt.Errorf("invalid id format: %s", id)
	}

	fwType := parts[0]
	tmp := parts[1]
	var natInsId string
	var ccnId string
	var enable bool
	if fwType == "NAT_FW" {
		res := strings.Split(tmp, tccommon.COMMA_SP)
		if len(res) != 2 {
			return fmt.Errorf("invalid id format: %s", id)
		}

		natInsId = res[0]
		ccnId = res[1]
	} else if fwType == "VPC_FW" {
		ccnId = tmp
	} else {
		return fmt.Errorf("invalid fw_type: %s", fwType)
	}

	request.FwType = helper.String(fwType)
	request.CcnId = helper.String(ccnId)

	if v, ok := d.GetOkExists("enable"); ok {
		enable = v.(bool)
		request.Enable = helper.Bool(enable)
	}

	if natInsId != "" {
		request.NatInsId = helper.String(natInsId)
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
		log.Printf("[CRITAL]%s update cfw_cluster_fw_bypass_config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// wait
	waitReq := cfwv20190904.NewDescribeClusterNatCcnFwSwitchListRequest()
	waitReq.Limit = helper.IntInt64(100)
	waitReq.Offset = helper.IntInt64(0)
	waitReq.NatType = helper.String("nat_ccn")
	if natInsId != "" {
		waitReq.Filters = []*cfwv20190904.CommonFilter{
			{
				Name:         helper.String("InsObj"),
				OperatorType: helper.IntInt64(1),
				Values:       []*string{helper.String(natInsId)},
			},
			{
				Name:         helper.String("AttachId"),
				OperatorType: helper.IntInt64(1),
				Values:       []*string{helper.String(ccnId)},
			},
		}
	} else {
		waitReq.Filters = []*cfwv20190904.CommonFilter{
			{
				Name:         helper.String("AttachId"),
				OperatorType: helper.IntInt64(1),
				Values:       []*string{helper.String(ccnId)},
			},
		}
	}

	reqErr = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCfwV20190904Client().DescribeClusterNatCcnFwSwitchListWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil || result.Response.Data == nil {
			return resource.NonRetryableError(fmt.Errorf("describe cfw_cluster_fw_bypass_config failed, response is nil"))
		}

		var matchedItem *cfwv20190904.NatFwSwitchDetailS
		for _, item := range result.Response.Data {
			if item == nil {
				continue
			}

			if fwType == "NAT_FW" && natInsId != "" {
				if item.InsObj != nil && *item.InsObj == natInsId &&
					item.AttachId != nil && *item.AttachId == ccnId {
					matchedItem = item
					break
				}
			} else {
				if item.AttachId != nil && *item.AttachId == ccnId {
					matchedItem = item
					break
				}
			}
		}

		if matchedItem == nil {
			return resource.NonRetryableError(fmt.Errorf("describe cfw_cluster_fw_bypass_config failed, matchedItem is nil"))
		}

		if matchedItem.Bypass != nil {
			if enable && *matchedItem.Bypass != 1 {
				return resource.RetryableError(fmt.Errorf("waiting for cfw_cluster_fw_bypass_config to be enabled"))
			} else if !enable && *matchedItem.Bypass != 0 {
				return resource.RetryableError(fmt.Errorf("waiting for cfw_cluster_fw_bypass_config to be disabled"))
			}
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s read cfw_cluster_fw_bypass_config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudCfwClusterFwBypassConfigRead(d, meta)
}

func resourceTencentCloudCfwClusterFwBypassConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_fw_bypass_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
