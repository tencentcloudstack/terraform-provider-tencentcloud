package clb

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/errors"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudClbListener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbListenerCreate,
		Read:   resourceTencentCloudClbListenerRead,
		Update: resourceTencentCloudClbListenerUpdate,
		Delete: resourceTencentCloudClbListenerDelete,
		Importer: &schema.ResourceImporter{
			State: helper.ImportWithDefaultValue(map[string]interface{}{
				"scheduler": CLB_LISTENER_SCHEDULER_WRR,
			}),
		},
		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the CLB instance.",
			},
			"listener_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the CLB listener, 1-80 characters. Supports letters, Chinese and other common international language characters, digits, hyphen '-' and underscore '_' (Unicode supplementary characters such as emoji are not allowed).",
			},
			"port": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 65535),
				Description:  "Port of the CLB listener. Port range: [1 - 65535].",
			},
			"protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_LISTENER_PROTOCOL),
				Description:  "Type of protocol within the listener. Valid values: `TCP`, `UDP`, `HTTP`, `HTTPS`, `TCP_SSL` and `QUIC`.",
			},
			"health_check_switch": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether health check is enabled.",
			},
			"health_check_time_out": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 60),
				Description:  "Response timeout of health check in seconds. Value range: 2-60, default 2. The response timeout must be less than the check interval.",
			},
			"health_check_interval_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Health check probe interval in seconds. Default 5. Value range: 2-300 for IPv4 CLB instances and 5-300 for IPv6 CLB instances. Note: some older IPv4 CLB instances have a range of 5-300.",
			},
			"health_check_health_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 10),
				Description:  "Health threshold. Default 3, meaning the backend is considered healthy after 3 consecutive successful probes. Value range: 2-10.",
			},
			"health_check_unhealth_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(2, 10),
				Description:  "Unhealthy threshold. Default 3, meaning the backend is considered unhealthy after 3 consecutive failed probes. Value range: 2-10.",
			},
			"health_check_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(HEALTH_CHECK_TYPE),
				Description:  "Health check protocol. Valid values: `TCP`, `HTTP`, `HTTPS`, `GRPC`, `PING`, `CUSTOM`. UDP listeners support `PING`/`CUSTOM`; TCP listeners support `TCP`/`HTTP`/`CUSTOM`; TCP_SSL/QUIC listeners support `TCP`/`HTTP`; HTTP rules support `HTTP`/`GRPC`; HTTPS rules support `HTTP`/`HTTPS`/`GRPC`. Defaults: `HTTP` for HTTP listeners, `TCP` for TCP/TCP_SSL/QUIC listeners, `PING` for UDP listeners; for HTTPS listeners the default matches the backend forwarding protocol.",
			},
			"health_check_port": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Health check port. Defaults to the backend service port; leave blank unless a specific port is required. Pass `-1` to restore the default. Only applicable to `TCP`/`UDP` listeners.",
			},
			"health_check_http_version": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(HTTP_VERSION),
				Description:  "HTTP version of the backend service. Required when `health_check_type` is `HTTP`. Valid values: `HTTP/1.0`, `HTTP/1.1`. Only applicable to `TCP` listeners.",
			},
			"health_check_http_code": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(1, 31),
				Description:  "Health check status code (only applicable to HTTP/HTTPS forwarding rules and the HTTP health check method of TCP listeners). Value range: 1-31, default 31. `1`=1xx healthy, `2`=2xx, `4`=3xx, `8`=4xx, `16`=5xx. To treat multiple return codes as healthy, add the corresponding values together.",
			},
			"health_check_http_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Health check path (only applicable to HTTP/HTTPS forwarding rules and the HTTP health check method of TCP listeners).",
			},
			"health_check_http_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Health check domain, carried in the HTTP Host header (only applicable to HTTP/HTTPS listeners and the HTTP health check method of TCP listeners; for TCP listeners using HTTP health check, this field is required).",
			},
			"health_check_http_method": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_HTTP_METHOD),
				Description:  "Health check method (only applicable to HTTP/HTTPS forwarding rules and the HTTP health check method of TCP listeners). Default `HEAD`. Valid values: `HEAD`, `GET`.",
			},
			"health_check_context_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CONTEX_TYPE),
				Description:  "Custom probe parameter. Required when `health_check_type` is `CUSTOM`, representing the input format of the health check. Valid values: `HEX`, `TEXT`. When `HEX`, the characters of `send_context`/`recv_context` can only be selected from `0123456789ABCDEF` and the length must be even. Only applicable to `TCP`/`UDP` listeners.",
			},
			"health_check_send_context": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom probe parameter. Required when `health_check_type` is `CUSTOM`, representing the request content sent by the health check. Only ASCII visible characters are allowed, max length 500. Only applicable to `TCP`/`UDP` listeners.",
			},
			"health_check_recv_context": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Custom probe parameter. Required when `health_check_type` is `CUSTOM`, representing the result returned by the health check. Only ASCII visible characters are allowed, max length 500. Only applicable to `TCP`/`UDP` listeners.",
			},
			"health_source_ip_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedIntValue([]int{0, 1}),
				Description:  "Health check source IP type. `0`: use the CLB VIP as the source IP, `1`: use a 100.64 IP range as the source IP.",
			},
			"certificate_ssl_mode": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"multi_cert_info"},
				ValidateFunc:  tccommon.ValidateAllowedStringValue(CERT_SSL_MODE),
				Description:   "Type of certificate. Valid values: `UNIDIRECTIONAL`, `MUTUAL`. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when it is available.",
			},
			"certificate_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"multi_cert_info"},
				Description:   "ID of the server certificate. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when it is available.",
			},
			"certificate_ca_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"multi_cert_info"},
				Description:   "ID of the client certificate. NOTES: Only supports listeners of `HTTPS` and `TCP_SSL` protocol and must be set when the ssl mode is `MUTUAL`.",
			},
			"multi_cert_info": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				ConflictsWith: []string{"certificate_ssl_mode", "certificate_id", "certificate_ca_id"},
				Description:   "Certificate information, supporting multiple server certificates with different algorithm types at the same time. Only applicable to `TCP_SSL` listeners and `HTTPS` listeners with SNI disabled. When creating a `TCP_SSL` listener or an `HTTPS` listener with SNI disabled, at least one of `certificate`/`multi_cert_info` must be specified, but they cannot be specified at the same time.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ssl_mode": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue(CERT_SSL_MODE),
							Description:  "Authentication type. Values: UNIDIRECTIONAL (one-way authentication), MUTUAL (two-way authentication).",
						},
						"cert_id_list": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "List of server certificate ID.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"session_expire_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: tccommon.ValidateIntegerInRange(30, 3600),
				Description:  "Session persistence time in seconds. Value range: 30-3600, default 0 (disabled). Only applicable to `TCP`/`UDP` listeners.",
			},
			"scheduler": {
				Type:         schema.TypeString,
				Default:      CLB_LISTENER_SCHEDULER_WRR,
				Optional:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(CLB_LISTENER_SCHEDULER),
				Description:  "Scheduling method. Valid values: `WRR` (weighted round-robin), `LEAST_CONN` (least connections). Default is `WRR`. Only applicable to `TCP`/`UDP`/`TCP_SSL`/`QUIC` listeners.",
			},
			"sni_switch": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Optional:    true,
				Description: "Indicates whether SNI is enabled. Only applicable to `HTTPS` listeners. `0`: disabled, `1`: enabled.",
			},
			"target_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CLB_TARGET_TYPE_NODE, CLB_TARGET_TYPE_TARGETGROUP, CLB_TARGET_TYPE_TARGETGROUP_V2}),
				Description:  "Backend target type. Valid values: `NODE`, `TARGETGROUP`, `TARGETGROUP-V2`. `NODE` means binding ordinary nodes, `TARGETGROUP` means binding a target group. Only applicable to `TCP`/`UDP` listeners; L7 listeners should be configured in the forwarding rule.",
			},
			"session_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{CLB_SESSION_TYPE_NORMAL, CLB_SESSION_TYPE_QUIC}),
				Description:  "Session persistence type. `NORMAL` (default): default session persistence type; `QUIC_CID`: session persistence by QUIC connection ID. Only applicable to `TCP`/`UDP` listeners; L7 listeners should be configured in the forwarding rule. If `QUIC_CID` is selected, `protocol` must be `UDP`, `scheduler` must be `WRR`, and only IPv4 is supported.",
			},
			"keepalive_enable": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Whether to enable persistent connection (long connection). Only applicable to `HTTP`/`HTTPS` listeners. Valid values: `0` (disable, default), `1` (enable). This feature is currently in beta.",
			},
			"end_port": {
				Type:        schema.TypeInt,
				ForceNew:    true,
				Computed:    true,
				Optional:    true,
				Description: "This parameter is used to specify the end port and is required when creating a port range listener. Only one member can be passed in when inputting the `Ports` parameter, which is used to specify the start port. If you want to try the port range feature, please [submit a ticket](https://console.cloud.tencent.com/workorder/category).",
			},
			"h2c_switch": {
				Type:        schema.TypeBool,
				ForceNew:    true,
				Computed:    true,
				Optional:    true,
				Description: "Whether to enable H2C for intranet `HTTP` listeners. `true`: enable, `false`: disable (default). When enabled, the listener only supports creating L7 rules with backend forwarding type `GRPC` or `GRPCS`; `GRPC` or `GRPCS` must be explicitly specified in the forwarding type when creating rules.",
			},
			"snat_enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Whether to enable SNAT (source IP replacement). `true`: enable, `false`: disable (default). Note: when SNAT is enabled, the client source IP is replaced and the pass-through client source IP option is disabled, and vice versa.",
			},
			"deregister_target_rst": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Reschedule function: the switch for unbinding backend services. When enabled, rescheduling is triggered when a backend service is unbound. Only supported by `TCP`/`UDP` listeners.",
			},
			"idle_connect_timeout": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Idle connection timeout. This parameter is only available for TCP/UDP listeners, in seconds. Default: 900s for TCP listeners, 300s for UDP listeners. Value range: 10-900 for shared and dedicated instances; 10-1980 for LCU-supported CLB instances. To set a value beyond the range, please submit a ticket for application.",
			},
			"max_conn": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Listener-level maximum concurrent connections. Currently only supported for performance capacity-type CLB instances with TCP/UDP/TCP_SSL/QUIC listeners. Pass -1 to indicate no limit at the listener level. Basic network instances do not support this parameter.",
			},
			"max_cps": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Listener-level maximum new connections per second. Currently only supported for performance capacity-type CLB instances with TCP/UDP/TCP_SSL/QUIC listeners. Pass -1 to indicate no limit at the listener level. Basic network instances do not support this parameter.",
			},
			"proxy_protocol": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable proxy protocol for TCP_SSL and QUIC listeners. Note: this field is not returned by the DescribeListeners API, so it will not be refreshed in state after creation.",
			},
			"data_compress_mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Data compression mode. Valid values: `transparent`, `compatibility`.",
			},
			"reschedule_target_zero_weight": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "The rescheduling function, with a weight of 0 as a switch, triggers rescheduling when the weight of the backend server is set to 0. Only supported by TCP/UDP listeners.",
			},
			"reschedule_unhealthy": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "Rescheduling function, health check exception switch. Enabling this switch triggers rescheduling when a backend server fails a health check. Supported only by TCP/UDP listeners.",
			},
			"reschedule_expand_target": {
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
				Description: "The rescheduling function, a switch for scaling backend services, triggers rescheduling when backend servers are added or removed. Only supported by TCP/UDP listeners.",
			},
			"reschedule_start_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Reschedule the trigger start time, with a value ranging from 0 to 3600 seconds. Only supported by TCP/UDP listeners.",
			},
			"reschedule_interval": {
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				Description: "Rescheduled trigger duration, ranging from 0 to 3600 seconds. Supported only by TCP/UDP listeners.",
			},
			//computed
			"listener_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of this CLB listener.",
			},
		},
	}
}

func resourceTencentCloudClbListenerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener.create")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	clbId := d.Get("clb_id").(string)
	listenerName := d.Get("listener_name").(string)
	request := clb.NewCreateListenerRequest()

	request.LoadBalancerId = helper.String(clbId)
	request.ListenerNames = []*string{&listenerName}

	port := int64(d.Get("port").(int))
	ports := []*int64{&port}
	request.Ports = ports
	protocol := d.Get("protocol").(string)
	request.Protocol = helper.String(protocol)

	healthSetFlag, healthCheck, healthErr := checkHealthCheckPara(ctx, d, protocol, HEALTH_APPLY_TYPE_LISTENER)
	if healthErr != nil {
		return healthErr
	}
	if healthSetFlag {
		request.HealthCheck = healthCheck
	}

	certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d, meta)

	if certErr != nil {
		return certErr
	}
	if certificateSetFlag {
		request.Certificate = certificateInput
	} else {
		if protocol == CLB_LISTENER_PROTOCOL_TCPSSL {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: certificated need to be set when protocol is TCPSSL")
		}
	}

	multiCertificateSetFlag, multiCertInput, certErr := checkMultiCertificateInputPara(ctx, d, meta)
	if certErr != nil {
		return certErr
	}

	if multiCertificateSetFlag {
		request.MultiCertInfo = multiCertInput
	} else {
		if protocol == CLB_LISTENER_PROTOCOL_TCPSSL {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: certificated need to be set when protocol is TCPSSL")
		}
	}

	scheduler := ""
	if v, ok := d.GetOk("scheduler"); ok {
		if v == CLB_LISTENER_SCHEDULER_IP_HASH {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: Scheduler 'IP_HASH' can only be set with rule of listener HTTP/HTTPS")
		}
		scheduler = v.(string)
		request.Scheduler = helper.String(scheduler)
	}
	if v, ok := d.GetOk("target_type"); ok {
		targetType := v.(string)
		request.TargetType = &targetType
	} else if protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP ||
		protocol == CLB_LISTENER_PROTOCOL_TCPSSL || protocol == CLB_LISTENER_PROTOCOL_QUIC {
		targetType := CLB_TARGET_TYPE_NODE
		request.TargetType = &targetType
	}

	if v, ok := d.GetOk("session_expire_time"); ok {
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP) {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS")
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: session_expire_time can only be set when scheduler is WRR ")
		}
		vv := int64(v.(int))
		request.SessionExpireTime = &vv
	}
	if v, ok := d.GetOkExists("sni_switch"); ok {
		if protocol != CLB_LISTENER_PROTOCOL_HTTPS {
			return fmt.Errorf("[CHECK][CLB listener][Create] check: sni_switch can only be set with protocol HTTPS ")
		} else {
			vv := v.(bool)
			vvv := int64(0)
			if vv {
				vvv = 1
			} else {
				if !certificateSetFlag && !multiCertificateSetFlag {
					return fmt.Errorf("[CHECK][CLB listener][Create] check: certificated need to be set when protocol is HTTPS")
				}
			}
			request.SniSwitch = &vvv
		}
	}

	if v, ok := d.GetOk("session_type"); ok {
		request.SessionType = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("keepalive_enable"); ok {
		request.KeepaliveEnable = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("end_port"); ok {
		request.EndPort = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("h2c_switch"); ok {
		request.H2cSwitch = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("snat_enable"); ok {
		request.SnatEnable = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("deregister_target_rst"); ok {
		request.DeregisterTargetRst = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("idle_connect_timeout"); ok {
		request.IdleConnectTimeout = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("reschedule_target_zero_weight"); ok {
		request.RescheduleTargetZeroWeight = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("reschedule_unhealthy"); ok {
		request.RescheduleUnhealthy = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("reschedule_expand_target"); ok {
		request.RescheduleExpandTarget = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("reschedule_start_time"); ok {
		request.RescheduleStartTime = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("reschedule_interval"); ok {
		request.RescheduleInterval = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("max_conn"); ok {
		request.MaxConn = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("max_cps"); ok {
		request.MaxCps = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("proxy_protocol"); ok {
		request.ProxyProtocol = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOkExists("data_compress_mode"); ok {
		request.DataCompressMode = helper.String(v.(string))
	}

	var response *clb.CreateListenerResponse
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().CreateListener(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

			if result == nil || result.Response == nil || result.Response.RequestId == nil {
				return resource.NonRetryableError(fmt.Errorf("Create CLB listener failed, Response si nil."))
			}

			requestId := *result.Response.RequestId
			retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
			if retryErr != nil {
				return resource.NonRetryableError(errors.WithStack(retryErr))
			}
		}

		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CLB listener failed, reason:%+v", logId, err)
		return err
	}
	if response.Response.ListenerIds == nil || len(response.Response.ListenerIds) < 1 {
		return fmt.Errorf("[CHECK][CLB listener][Create] check: Response error, listener id is null")
	}
	listenerId := *response.Response.ListenerIds[0]

	//this ID style changes since terraform 1.47.0
	d.SetId(clbId + tccommon.FILED_SP + listenerId)
	return resourceTencentCloudClbListenerRead(d, meta)
}

func resourceTencentCloudClbListenerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	resourceId := d.Id()
	var listenerId = resourceId
	items := strings.Split(resourceId, tccommon.FILED_SP)
	itemLength := len(items)
	clbId := d.Get("clb_id").(string)

	if itemLength == 1 && clbId == "" {
		return fmt.Errorf("the old style listenerId %s does not support import, please use clbId#listenerId style", resourceId)
	} else if itemLength == 2 && clbId == "" {
		listenerId = items[1]
		clbId = items[0]
	} else if itemLength == 2 && clbId != "" {
		listenerId = items[1]
	} else if itemLength != 1 && itemLength != 2 {
		return fmt.Errorf("broken ID %s", resourceId)
	}

	var instance *clb.Listener
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := clbService.DescribeListenerById(ctx, listenerId, clbId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		instance = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLB listener failed, reason:%s\n ", logId, err.Error())
		return err
	}

	if instance == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("clb_id", clbId)
	_ = d.Set("listener_id", instance.ListenerId)
	_ = d.Set("port", instance.Port)
	_ = d.Set("protocol", instance.Protocol)
	if instance.ListenerName != nil {
		_ = d.Set("listener_name", instance.ListenerName)
	}
	if instance.TargetType != nil {
		_ = d.Set("target_type", instance.TargetType)
	}
	if instance.SessionExpireTime != nil {
		_ = d.Set("session_expire_time", instance.SessionExpireTime)
	}
	if *instance.Protocol == CLB_LISTENER_PROTOCOL_TCP || *instance.Protocol == CLB_LISTENER_PROTOCOL_TCPSSL ||
		*instance.Protocol == CLB_LISTENER_PROTOCOL_UDP || *instance.Protocol == CLB_LISTENER_PROTOCOL_QUIC {
		_ = d.Set("scheduler", instance.Scheduler)
	}
	_ = d.Set("sni_switch", *instance.SniSwitch > 0)

	//health check
	if instance.HealthCheck != nil {
		healthCheckSwitch := false
		if *instance.HealthCheck.HealthSwitch == int64(1) {
			healthCheckSwitch = true
		}
		_ = d.Set("health_check_switch", healthCheckSwitch)
		if instance.HealthCheck.IntervalTime != nil {
			_ = d.Set("health_check_interval_time", instance.HealthCheck.IntervalTime)
		}
		if instance.HealthCheck.TimeOut != nil {
			_ = d.Set("health_check_time_out", instance.HealthCheck.TimeOut)
		}
		if instance.HealthCheck.HealthNum != nil {
			_ = d.Set("health_check_health_num", instance.HealthCheck.HealthNum)
		}
		if instance.HealthCheck.UnHealthNum != nil {
			_ = d.Set("health_check_unhealth_num", instance.HealthCheck.UnHealthNum)
		}
		if instance.HealthCheck.CheckPort != nil {
			_ = d.Set("health_check_port", instance.HealthCheck.CheckPort)
		}
		if instance.HealthCheck.CheckType != nil {
			_ = d.Set("health_check_type", instance.HealthCheck.CheckType)
		}
		if instance.HealthCheck.HttpCode != nil {
			_ = d.Set("health_check_http_code", instance.HealthCheck.HttpCode)
		}
		if instance.HealthCheck.HttpCheckPath != nil {
			_ = d.Set("health_check_http_path", instance.HealthCheck.HttpCheckPath)
		}
		if instance.HealthCheck.HttpCheckDomain != nil {
			_ = d.Set("health_check_http_domain", instance.HealthCheck.HttpCheckDomain)
		}
		if instance.HealthCheck.HttpCheckMethod != nil {
			_ = d.Set("health_check_http_method", instance.HealthCheck.HttpCheckMethod)
		}
		if instance.HealthCheck.HttpVersion != nil {
			_ = d.Set("health_check_http_version", instance.HealthCheck.HttpVersion)
		}
		if instance.HealthCheck.ContextType != nil {
			_ = d.Set("health_check_context_type", instance.HealthCheck.ContextType)
		}
		if instance.HealthCheck.SendContext != nil {
			_ = d.Set("health_check_send_context", instance.HealthCheck.SendContext)
		}
		if instance.HealthCheck.RecvContext != nil {
			_ = d.Set("health_check_recv_context", instance.HealthCheck.RecvContext)
		}
		if instance.HealthCheck.SourceIpType != nil {
			_ = d.Set("health_source_ip_type", instance.HealthCheck.SourceIpType)
		}
	}

	if instance.Certificate != nil {
		// check single cert or multi cert
		if instance.Certificate.ExtCertIds != nil && len(instance.Certificate.ExtCertIds) > 0 {
			multiCertInfo := make([]map[string]interface{}, 0, 1)
			multiCert := make(map[string]interface{}, 0)
			certIds := make([]string, 0)
			if instance.Certificate.SSLMode != nil {
				multiCert["ssl_mode"] = *instance.Certificate.SSLMode
			}

			if instance.Certificate.CertId != nil {
				certIds = append(certIds, *instance.Certificate.CertId)
			}

			for _, item := range instance.Certificate.ExtCertIds {
				certIds = append(certIds, *item)
			}

			multiCert["cert_id_list"] = certIds
			multiCertInfo = append(multiCertInfo, multiCert)
			_ = d.Set("multi_cert_info", multiCertInfo)
		} else {
			_ = d.Set("certificate_ssl_mode", instance.Certificate.SSLMode)
			_ = d.Set("certificate_id", instance.Certificate.CertId)
			if instance.Certificate.CertCaId != nil {
				_ = d.Set("certificate_ca_id", instance.Certificate.CertCaId)
			}
		}
	}

	if instance.SessionType != nil {
		_ = d.Set("session_type", instance.SessionType)
	}
	if instance.KeepaliveEnable != nil {
		_ = d.Set("keepalive_enable", instance.KeepaliveEnable)
	}

	if instance.EndPort != nil {
		_ = d.Set("end_port", instance.EndPort)
	}

	if instance.AttrFlags != nil && len(instance.AttrFlags) > 0 {
		if tccommon.IsContains(helper.PStrings(instance.AttrFlags), "H2cSwitch") {
			_ = d.Set("h2c_switch", true)
		} else {
			_ = d.Set("h2c_switch", false)
		}

		if tccommon.IsContains(helper.PStrings(instance.AttrFlags), "SnatEnable") {
			_ = d.Set("snat_enable", true)
		} else {
			_ = d.Set("snat_enable", false)
		}
	} else {
		_ = d.Set("h2c_switch", false)
		_ = d.Set("snat_enable", false)
	}

	if instance.DeregisterTargetRst != nil {
		_ = d.Set("deregister_target_rst", instance.DeregisterTargetRst)
	}

	if instance.IdleConnectTimeout != nil {
		_ = d.Set("idle_connect_timeout", instance.IdleConnectTimeout)
	}

	_ = d.Set("reschedule_target_zero_weight", false)
	_ = d.Set("reschedule_unhealthy", false)
	_ = d.Set("reschedule_expand_target", false)
	if instance.AttrFlags != nil {
		for _, item := range instance.AttrFlags {
			if item != nil && *item == "RescheduleTargetZeroWeight" {
				_ = d.Set("reschedule_target_zero_weight", true)
			}

			if item != nil && *item == "RescheduleUnhealthy" {
				_ = d.Set("reschedule_unhealthy", true)
			}

			if item != nil && *item == "RescheduleExpandTarget" {
				_ = d.Set("reschedule_expand_target", true)
			}
		}
	}

	if instance.RescheduleStartTime != nil {
		_ = d.Set("reschedule_start_time", instance.RescheduleStartTime)
	}

	if instance.RescheduleInterval != nil {
		_ = d.Set("reschedule_interval", instance.RescheduleInterval)
	}

	if instance.MaxConn != nil {
		_ = d.Set("max_conn", instance.MaxConn)
	}

	if instance.MaxCps != nil {
		_ = d.Set("max_cps", instance.MaxCps)
	}

	if instance.DataCompressMode != nil {
		_ = d.Set("data_compress_mode", instance.DataCompressMode)
	}

	return nil
}

func resourceTencentCloudClbListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener.update")()

	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	resourceId := d.Id()
	items := strings.Split(resourceId, tccommon.FILED_SP)
	itemLength := len(items)
	listenerId := items[itemLength-1]
	clbId := d.Get("clb_id").(string)
	changed := false
	scheduler := ""
	listenerName := ""
	sessionExpireTime := 0
	protocol := d.Get("protocol").(string)

	request := clb.NewModifyListenerRequest()
	request.ListenerId = helper.String(listenerId)
	request.LoadBalancerId = helper.String(clbId)

	if d.HasChange("listener_name") {
		changed = true
		listenerName = d.Get("listener_name").(string)
		request.ListenerName = helper.String(listenerName)
	}

	if d.HasChange("scheduler") {
		changed = true
		scheduler = d.Get("scheduler").(string)
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP ||
			protocol == CLB_LISTENER_PROTOCOL_TCPSSL || protocol == CLB_LISTENER_PROTOCOL_QUIC) {
			return fmt.Errorf("[CHECK][CLB listener %s][Update] check: Scheduler can only be set with listener protocol TCP/UDP/TCP_SSL/QUIC or rule of listener HTTP/HTTPS", listenerId)
		}
		if scheduler == CLB_LISTENER_SCHEDULER_IP_HASH {
			return fmt.Errorf("[CHECK][CLB listener %s][Update] check: Scheduler 'IP_HASH' can only be set with rule of listener HTTP/HTTPS", listenerId)
		}
		request.Scheduler = helper.String(scheduler)
	}

	if d.HasChange("session_expire_time") {
		changed = true
		sessionExpireTime = d.Get("session_expire_time").(int)
		if !(protocol == CLB_LISTENER_PROTOCOL_TCP || protocol == CLB_LISTENER_PROTOCOL_UDP) {
			return fmt.Errorf("[CHECK][CLB listener %s][Update] check: session_expire_time can only be set with protocol TCP/UDP or rule of listener HTTP/HTTPS", listenerId)
		}
		if scheduler != CLB_LISTENER_SCHEDULER_WRR && scheduler != "" {
			return fmt.Errorf("[CHECK][CLB listener %s][Update] check: session_expire_time can only be set when scheduler is WRR", listenerId)
		}
		sessionExpireTime64 := int64(sessionExpireTime)
		request.SessionExpireTime = &sessionExpireTime64
	}

	healthSetFlag, healthCheck, healthErr := checkHealthCheckParaForUpdate(ctx, d, protocol, HEALTH_APPLY_TYPE_LISTENER)
	if healthErr != nil {
		return healthErr
	}
	if healthSetFlag {
		changed = true
		request.HealthCheck = healthCheck
	}

	certificateSetFlag, certificateInput, certErr := checkCertificateInputPara(ctx, d, meta)
	if certErr != nil {
		return certErr
	}
	if certificateSetFlag {
		changed = true
		request.Certificate = certificateInput
	}

	multiCertificateSetFlag, multiCertInput, certErr := checkMultiCertificateInputPara(ctx, d, meta)
	if certErr != nil {
		return certErr
	}

	if multiCertificateSetFlag {
		changed = true
		request.MultiCertInfo = multiCertInput
	}

	if d.HasChange("target_type") {
		changed = true
		targetType := d.Get("target_type").(string)
		request.TargetType = helper.String(targetType)
	}

	if d.HasChange("session_type") {
		changed = true
		sessionType := d.Get("session_type").(string)
		request.SessionType = helper.String(sessionType)
	}

	if d.HasChange("keepalive_enable") {
		changed = true
		keepaliveEnable := d.Get("keepalive_enable").(int)
		request.KeepaliveEnable = helper.IntInt64(keepaliveEnable)
	}

	if d.HasChange("snat_enable") {
		changed = true
		if v, ok := d.GetOkExists("snat_enable"); ok {
			request.SnatEnable = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("deregister_target_rst") {
		changed = true
		if v, ok := d.GetOkExists("deregister_target_rst"); ok {
			request.DeregisterTargetRst = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("idle_connect_timeout") {
		changed = true
		if v, ok := d.GetOkExists("idle_connect_timeout"); ok {
			request.IdleConnectTimeout = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("reschedule_target_zero_weight") {
		changed = true
		if v, ok := d.GetOkExists("reschedule_target_zero_weight"); ok {
			request.RescheduleTargetZeroWeight = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("reschedule_unhealthy") {
		changed = true
		if v, ok := d.GetOkExists("reschedule_unhealthy"); ok {
			request.RescheduleUnhealthy = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("reschedule_expand_target") {
		changed = true
		if v, ok := d.GetOkExists("reschedule_expand_target"); ok {
			request.RescheduleExpandTarget = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("reschedule_start_time") {
		changed = true
		if v, ok := d.GetOkExists("reschedule_start_time"); ok {
			request.RescheduleStartTime = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("reschedule_interval") {
		changed = true
		if v, ok := d.GetOkExists("reschedule_interval"); ok {
			request.RescheduleInterval = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("max_conn") {
		changed = true
		if v, ok := d.GetOkExists("max_conn"); ok {
			request.MaxConn = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("max_cps") {
		changed = true
		if v, ok := d.GetOkExists("max_cps"); ok {
			request.MaxCps = helper.IntInt64(v.(int))
		}
	}

	if d.HasChange("proxy_protocol") {
		changed = true
		if v, ok := d.GetOkExists("proxy_protocol"); ok {
			request.ProxyProtocol = helper.Bool(v.(bool))
		}
	}

	if d.HasChange("data_compress_mode") {
		changed = true
		if v, ok := d.GetOkExists("data_compress_mode"); ok {
			request.DataCompressMode = helper.String(v.(string))
		}
	}

	if changed {
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			response, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient().ModifyListener(request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
					logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

				if response == nil || response.Response == nil || response.Response.RequestId == nil {
					return resource.NonRetryableError(fmt.Errorf("Modify CLB listener failed, Response si nil."))
				}

				requestId := *response.Response.RequestId
				retryErr := waitForTaskFinish(requestId, meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseClbClient())
				if retryErr != nil {
					return resource.NonRetryableError(errors.WithStack(retryErr))
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update CLB listener failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudClbListenerRead(d, meta)
}

func resourceTencentCloudClbListenerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_clb_listener.delete")()
	clbActionMu.Lock()
	defer clbActionMu.Unlock()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	resourceId := d.Id()
	items := strings.Split(resourceId, tccommon.FILED_SP)
	itemLength := len(items)
	listenerId := items[itemLength-1]
	clbId := d.Get("clb_id").(string)

	clbService := ClbService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		e := clbService.DeleteListenerById(ctx, clbId, listenerId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CLB listener failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
