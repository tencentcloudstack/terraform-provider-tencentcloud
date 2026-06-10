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
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_fw_bypass_config.read")()
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
		log.Printf("[WARN]%s resource `cfw_cluster_fw_bypass_config` id=%s not found, response is nil", logId, d.Id())
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
		log.Printf("[WARN]%s resource `cfw_cluster_fw_bypass_config` id=%s not found in response data", logId, d.Id())
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
		_ = d.Set("enable", *matchedItem.Bypass == 1)
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
		log.Printf("[CRITAL]%s update cfw_cluster_fw_bypass_config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudCfwClusterFwBypassConfigRead(d, meta)
}

func resourceTencentCloudCfwClusterFwBypassConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_cfw_cluster_fw_bypass_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
