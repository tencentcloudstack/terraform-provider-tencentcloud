package vpc

import (
	"log"

	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceTencentCloudNatGatewayFlowMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudNatGatewayFlowMonitorCreate,
		Read:   resourceTencentCloudNatGatewayFlowMonitorRead,
		Update: resourceTencentCloudNatGatewayFlowMonitorUpdate,
		Delete: resourceTencentCloudNatGatewayFlowMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				ForceNew:    true,
				Type:        schema.TypeString,
				Description: "ID of Gateway.",
			},
			"enable": {
				Required:    true,
				Type:        schema.TypeBool,
				Description: "Whether to enable flow monitor.",
			},
			"bandwidth": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "Bandwidth of flow monitor.",
			},
		},
	}
}

func resourceTencentCloudNatGatewayFlowMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_nat_gateway_flow_monitor.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var gatewayId string

	if v, ok := d.GetOk("gateway_id"); ok {
		gatewayId = v.(string)
	}

	d.SetId(gatewayId)

	return resourceTencentCloudNatGatewayFlowMonitorUpdate(d, meta)
}

func resourceTencentCloudNatGatewayFlowMonitorRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_nat_gateway_flow_monitor.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	gatewayId := d.Id()

	var (
		checkRequest = vpc.NewCheckGatewayFlowMonitorRequest()
		response     = vpc.NewCheckGatewayFlowMonitorResponse()
	)

	checkRequest.GatewayId = &gatewayId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().CheckGatewayFlowMonitor(checkRequest)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, checkRequest.GetAction(), checkRequest.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s enable nat flow monitor failed, reason:%+v", logId, err)
		return err
	}

	_ = d.Set("gateway_id", gatewayId)

	if response.Response.Enabled != nil {
		if *response.Response.Enabled {
			_ = d.Set("enable", true)
		} else {
			_ = d.Set("enable", false)
		}
	}

	if response.Response.Bandwidth != nil {
		_ = d.Set("bandwidth", *response.Response.Bandwidth)
	}
	return nil
}

func resourceTencentCloudNatGatewayFlowMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_nat_gateway_flow_monitor.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	var (
		enable         bool
		enableRequest  = vpc.NewEnableGatewayFlowMonitorRequest()
		disableRequest = vpc.NewDisableGatewayFlowMonitorRequest()
	)

	gatewayId := d.Id()

	if v, ok := d.GetOkExists("enable"); ok {
		enable = v.(bool)
	}

	if enable {
		enableRequest.GatewayId = &gatewayId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().EnableGatewayFlowMonitor(enableRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, enableRequest.GetAction(), enableRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s enable nat flow monitor failed, reason:%+v", logId, err)
			return err
		}
	} else {
		disableRequest.GatewayId = &gatewayId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseVpcClient().DisableGatewayFlowMonitor(disableRequest)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, disableRequest.GetAction(), disableRequest.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s disable nat flow monitor failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudNatGatewayFlowMonitorRead(d, meta)
}

func resourceTencentCloudNatGatewayFlowMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_nat_gateway_flow_monitor.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}
