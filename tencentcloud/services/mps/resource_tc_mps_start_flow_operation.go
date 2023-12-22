package mps

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudMpsStartFlowOperation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsStartFlowOperationCreate,
		Read:   resourceTencentCloudMpsStartFlowOperationRead,
		Delete: resourceTencentCloudMpsStartFlowOperationDelete,
		Schema: map[string]*schema.Schema{
			"flow_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "Flow Id.",
			},
			"start": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeBool,
				Description: "`true`: start mps stream link flow; `false`: stop.",
			},
		},
	}
}

func resourceTencentCloudMpsStartFlowOperationCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_start_flow_operation.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		startRequest = mps.NewStartStreamLinkFlowRequest()
		stopRequest  = mps.NewStopStreamLinkFlowRequest()
		flowId       string
		start        bool
	)
	if v, ok := d.GetOk("flow_id"); ok {
		startRequest.FlowId = helper.String(v.(string))
		stopRequest.FlowId = helper.String(v.(string))
		flowId = v.(string)
	}

	if v, ok := d.GetOkExists("start"); ok {
		start = v.(bool)
	}

	if start {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().StartStreamLinkFlow(startRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, startRequest.GetAction(), startRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate mps startFlowOperation failed, reason:%+v", logId, err)
			return err
		}
	} else {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseMpsClient().StopStreamLinkFlow(stopRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, stopRequest.GetAction(), stopRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s operate mps stopFlowOperation failed, reason:%+v", logId, err)
			return err
		}
	}

	d.SetId(flowId)

	return resourceTencentCloudMpsStartFlowOperationRead(d, meta)
}

func resourceTencentCloudMpsStartFlowOperationRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_start_flow_operation.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func resourceTencentCloudMpsStartFlowOperationDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_mps_start_flow_operation.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
