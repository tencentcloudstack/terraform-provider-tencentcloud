package fl

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
)

func ResourceTencentCloudVpcFlowLogConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudVpcFlowLogConfigCreate,
		Read:   resourceTencentCloudVpcFlowLogConfigRead,
		Update: resourceTencentCloudVpcFlowLogConfigUpdate,
		Delete: resourceTencentCloudVpcFlowLogConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"flow_log_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Flow log ID.",
			},

			"enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "If enable snapshot policy.",
			},
		},
	}
}

func resourceTencentCloudVpcFlowLogConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_flow_log_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	flowLogId := d.Get("flow_log_id").(string)

	d.SetId(flowLogId)

	return resourceTencentCloudVpcFlowLogConfigUpdate(d, meta)
}

func resourceTencentCloudVpcFlowLogConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_flow_log_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := VpcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	flowLogId := d.Id()

	request := vpc.NewDescribeFlowLogsRequest()
	request.FlowLogId = &flowLogId

	flowLogs, err := service.DescribeFlowLogs(ctx, request)
	if err != nil {
		return err
	}

	if len(flowLogs) < 1 {
		d.SetId("")
		log.Printf("[WARN]%s resource `VpcFlowLogConfig` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	_ = d.Set("flow_log_id", flowLogId)

	flowLogConfig := flowLogs[0]

	if flowLogConfig.Enable != nil {
		_ = d.Set("enable", flowLogConfig.Enable)
	}

	return nil
}

func resourceTencentCloudVpcFlowLogConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_flow_log_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		enable         bool
		enableRequest  = vpc.NewEnableFlowLogsRequest()
		disableRequest = vpc.NewDisableFlowLogsRequest()
	)

	flowLogId := d.Id()

	if v, ok := d.GetOkExists("enable"); ok {
		enable = v.(bool)
	}

	if enable {
		enableRequest.FlowLogIds = []*string{&flowLogId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().EnableFlowLogs(enableRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableRequest.GetAction(), enableRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update vpc flowLogConfig failed, reason:%+v", logId, err)
			return err
		}
	} else {
		disableRequest.FlowLogIds = []*string{&flowLogId}
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DisableFlowLogs(disableRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableRequest.GetAction(), disableRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update vpc flowLogConfig failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudVpcFlowLogConfigRead(d, meta)
}

func resourceTencentCloudVpcFlowLogConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_vpc_flow_log_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
