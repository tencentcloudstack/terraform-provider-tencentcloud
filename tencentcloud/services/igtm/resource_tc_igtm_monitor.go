package igtm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIgtmMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIgtmMonitorCreate,
		Read:   resourceTencentCloudIgtmMonitorRead,
		Update: resourceTencentCloudIgtmMonitorUpdate,
		Delete: resourceTencentCloudIgtmMonitorDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"monitor_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Monitor name.",
			},

			"check_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Detection protocol, optional values `PING`, `TCP`, `HTTP`, `HTTPS`.",
			},

			"check_interval": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Check interval (seconds), optional values 15 60 120 300.",
			},

			"timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Timeout time, unit seconds, optional values 2 3 5 10.",
			},

			"fail_times": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Retry count, optional values 0, 1, 2.",
			},

			"fail_rate": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Failure rate, values 20 30 40 50 60 70 80 100, default value 50.",
			},

			"detector_style": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Monitoring node type, optional values AUTO INTERNAL OVERSEAS IPV6 ALL.",
			},

			"detector_group_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Detector group ID list separated by commas.",
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
			},

			"ping_num": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "PING packet count, required when CheckProtocol=ping, optional values 20 50 100.",
			},

			"tcp_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Check port, optional values between 1-65535.",
			},

			"host": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host setting, default is business domain name.",
			},

			"path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "URL path, default is \"/\".",
			},

			"return_code_threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Return error code threshold, optional values 400 and 500, default value 500.",
			},

			"enable_redirect": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Follow 3XX redirect, DISABLED for disabled, ENABLED for enabled, default disabled.",
			},

			"enable_sni": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Enable SNI, DISABLED for disabled, ENABLED for enabled, default disabled.",
			},

			"packet_loss_rate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Packet loss rate alarm threshold, required when CheckProtocol=ping, values 10 30 50 80 90 100.",
			},

			"continue_period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Continuous period count, optional values 1-5.",
			},

			// computed
			"monitor_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Monitor ID.",
			},
		},
	}
}

func resourceTencentCloudIgtmMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_monitor.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = igtmv20231024.NewCreateMonitorRequest()
		response  = igtmv20231024.NewCreateMonitorResponse()
		monitorId string
	)

	if v, ok := d.GetOk("monitor_name"); ok {
		request.MonitorName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("check_protocol"); ok {
		request.CheckProtocol = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("check_interval"); ok {
		request.CheckInterval = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("timeout"); ok {
		request.Timeout = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("fail_times"); ok {
		request.FailTimes = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("fail_rate"); ok {
		request.FailRate = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("detector_style"); ok {
		request.DetectorStyle = helper.String(v.(string))
	}

	if v, ok := d.GetOk("detector_group_ids"); ok {
		detectorGroupIdsSet := v.(*schema.Set).List()
		for i := range detectorGroupIdsSet {
			detectorGroupIds := detectorGroupIdsSet[i].(int)
			request.DetectorGroupIds = append(request.DetectorGroupIds, helper.IntUint64(detectorGroupIds))
		}
	}

	if v, ok := d.GetOkExists("ping_num"); ok {
		request.PingNum = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("tcp_port"); ok {
		request.TcpPort = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("host"); ok {
		request.Host = helper.String(v.(string))
	}

	if v, ok := d.GetOk("path"); ok {
		request.Path = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("return_code_threshold"); ok {
		request.ReturnCodeThreshold = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("enable_redirect"); ok {
		request.EnableRedirect = helper.String(v.(string))
	}

	if v, ok := d.GetOk("enable_sni"); ok {
		request.EnableSni = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("packet_loss_rate"); ok {
		request.PacketLossRate = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("continue_period"); ok {
		request.ContinuePeriod = helper.IntUint64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().CreateMonitorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create igtm monitor failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create igtm monitor failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.MonitorId == nil {
		return fmt.Errorf("MonitorId is nil.")
	}

	monitorId = helper.UInt64ToStr(*response.Response.MonitorId)
	d.SetId(monitorId)
	return resourceTencentCloudIgtmMonitorRead(d, meta)
}

func resourceTencentCloudIgtmMonitorRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_monitor.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service   = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		monitorId = d.Id()
	)

	respData, err := service.DescribeIgtmMonitorById(ctx, monitorId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_igtm_monitor` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.MonitorName != nil {
		_ = d.Set("monitor_name", respData.MonitorName)
	}

	if respData.CheckProtocol != nil {
		_ = d.Set("check_protocol", respData.CheckProtocol)
	}

	if respData.CheckInterval != nil {
		_ = d.Set("check_interval", respData.CheckInterval)
	}

	if respData.Timeout != nil {
		_ = d.Set("timeout", respData.Timeout)
	}

	if respData.FailTimes != nil {
		_ = d.Set("fail_times", respData.FailTimes)
	}

	if respData.FailRate != nil {
		_ = d.Set("fail_rate", respData.FailRate)
	}

	if respData.DetectorStyle != nil {
		_ = d.Set("detector_style", respData.DetectorStyle)
	}

	if respData.DetectorGroupIds != nil {
		_ = d.Set("detector_group_ids", respData.DetectorGroupIds)
	}

	if respData.PingNum != nil {
		_ = d.Set("ping_num", respData.PingNum)
	}

	if respData.TcpPort != nil {
		_ = d.Set("tcp_port", respData.TcpPort)
	}

	if respData.Host != nil {
		_ = d.Set("host", respData.Host)
	}

	if respData.Path != nil {
		_ = d.Set("path", respData.Path)
	}

	if respData.ReturnCodeThreshold != nil {
		_ = d.Set("return_code_threshold", respData.ReturnCodeThreshold)
	}

	if respData.EnableRedirect != nil {
		_ = d.Set("enable_redirect", respData.EnableRedirect)
	}

	if respData.EnableSni != nil {
		_ = d.Set("enable_sni", respData.EnableSni)
	}

	if respData.PacketLossRate != nil {
		_ = d.Set("packet_loss_rate", respData.PacketLossRate)
	}

	if respData.ContinuePeriod != nil {
		_ = d.Set("continue_period", respData.ContinuePeriod)
	}

	if respData.MonitorId != nil {
		_ = d.Set("monitor_id", respData.MonitorId)
	}

	return nil
}

func resourceTencentCloudIgtmMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_monitor.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		monitorId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"monitor_name", "check_protocol", "check_interval", "timeout", "fail_times", "fail_rate", "detector_style", "detector_group_ids", "ping_num", "tcp_port", "host", "path", "return_code_threshold", "enable_redirect", "enable_sni", "packet_loss_rate", "continue_period"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := igtmv20231024.NewModifyMonitorRequest()
		if v, ok := d.GetOk("monitor_name"); ok {
			request.MonitorName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("check_protocol"); ok {
			request.CheckProtocol = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("check_interval"); ok {
			request.CheckInterval = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("timeout"); ok {
			request.Timeout = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("fail_times"); ok {
			request.FailTimes = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("fail_rate"); ok {
			request.FailRate = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("detector_style"); ok {
			request.DetectorStyle = helper.String(v.(string))
		}

		if v, ok := d.GetOk("detector_group_ids"); ok {
			detectorGroupIdsSet := v.(*schema.Set).List()
			for i := range detectorGroupIdsSet {
				detectorGroupIds := detectorGroupIdsSet[i].(int)
				request.DetectorGroupIds = append(request.DetectorGroupIds, helper.IntUint64(detectorGroupIds))
			}
		}

		if v, ok := d.GetOkExists("ping_num"); ok {
			request.PingNum = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("tcp_port"); ok {
			request.TcpPort = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("host"); ok {
			request.Host = helper.String(v.(string))
		}

		if v, ok := d.GetOk("path"); ok {
			request.Path = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("return_code_threshold"); ok {
			request.ReturnCodeThreshold = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("enable_redirect"); ok {
			request.EnableRedirect = helper.String(v.(string))
		}

		if v, ok := d.GetOk("enable_sni"); ok {
			request.EnableSni = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("packet_loss_rate"); ok {
			request.PacketLossRate = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOkExists("continue_period"); ok {
			request.ContinuePeriod = helper.IntUint64(v.(int))
		}

		request.MonitorId = helper.StrToUint64Point(monitorId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().ModifyMonitorWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Msg == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify igtm monitor failed, Response is nil."))
			}

			if *result.Response.Msg == "success" {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("Modify igtm monitor failed, Msg is %s.", *result.Response.Msg))
			}
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update igtm monitor failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudIgtmMonitorRead(d, meta)
}

func resourceTencentCloudIgtmMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_monitor.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId     = tccommon.GetLogId(tccommon.ContextNil)
		ctx       = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request   = igtmv20231024.NewDeleteMonitorRequest()
		monitorId = d.Id()
	)

	request.MonitorId = helper.StrToUint64Point(monitorId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().DeleteMonitorWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Msg == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete igtm monitor failed, Response is nil."))
		}

		if *result.Response.Msg == "success" {
			return nil
		} else {
			return resource.NonRetryableError(fmt.Errorf("Delete igtm monitor failed, Msg is %s.", *result.Response.Msg))
		}
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete igtm monitor failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
