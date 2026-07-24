package ga2

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdkErrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// ResourceTencentCloudGa2Listener manages a Tencent Cloud GA2 listener.
//
// All write APIs (CreateListener / ModifyListener / DeleteListener) are asynchronous:
// each returns a TaskId that must be polled via DescribeTaskResult until
// Status == "SUCCESS" before this resource considers the operation complete.
func ResourceTencentCloudGa2Listener() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGa2ListenerCreate,
		Read:   resourceTencentCloudGa2ListenerRead,
		Update: resourceTencentCloudGa2ListenerUpdate,
		Delete: resourceTencentCloudGa2ListenerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_accelerator_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Global accelerator instance ID this listener belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Listener name. Maximum length is 60 bytes.",
			},
			"port_ranges": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Listening port range. Cannot be modified after creation; modifying it forces a new resource.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Inclusive start port.",
						},
						"to_port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Inclusive end port.",
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Listener description. Maximum length is 100 bytes.",
			},
			"listener_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Listener routing type. Valid values: `Standard` (smart routing). Default: `Standard`. Cannot be modified after creation; modifying it forces a new resource.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Listener protocol. Valid values: `TCP`, `UDP`, `HTTP`, `HTTPS`. Default: `TCP`. Cannot be modified after creation.",
			},
			"idle_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Connection idle timeout in seconds. Valid range and default value depend on the listener protocol: " +
					"`1-60` for HTTP/HTTPS listeners (default `15`), `10-900` for TCP listeners (default `900`), " +
					"`10-20` for UDP listeners (default `20`).",
			},
			"get_real_ip_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Method used to retrieve the real client IP for layer-4 listeners. Valid values: `TOA`, `ProxyProtocol`, " +
					"`ProxyProtocolV2`, `Close`. Only takes effect when the layer-4 real-IP feature is enabled. " +
					"Only TCP listeners support modifying this field after creation.",
			},
			"client_affinity": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Whether to enable session stickiness. Valid values: `Open`, `Close`. " +
					"Only TCP/UDP listeners support modifying this field.",
			},
			"client_affinity_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Session-stickiness duration in seconds. Valid range: [60, 3600]. " +
					"NOTE: this field is silently ignored on Create (the SDK CreateListener API has no equivalent slot) " +
					"and forwarded only on Update via ModifyListener.",
			},
			"request_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: "Request timeout in seconds. Valid range: [1, 180]. Default: `60`. " +
					"Only applicable to HTTP/HTTPS listeners.",
			},
			"x_forwarded_for_real_ip": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
				Description: "Whether to enable layer-7 real-IP forwarding (X-Forwarded-For). " +
					"Only HTTP/HTTPS listeners support modifying this field.",
			},
			"certification_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "SSL authentication mode. Valid values: `UNIDIRECTIONAL` (one-way authentication, server certificate only), " +
					"`MUTUAL` (mutual/two-way authentication, requires both server and client certificates). " +
					"Required when the listener protocol is `HTTPS`. Only HTTP/HTTPS listeners support modifying this field.",
			},
			"cipher_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "TLS cipher suite policy. Valid values: `tls_policy_1.0-2`, `tls_policy_1.1-2`, `tls_policy_1.2`, " +
					"`tls_policy_1.2_strict`, `tls_policy_1.2_strict-1.3`. Only HTTPS listeners support configuring/modifying this field.",
			},
			"server_certificates": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Server certificate ID list. Required when the listener protocol is `HTTPS`. " +
					"Only HTTPS listeners support modifying this field. Treated as an unordered set; HCL element order has no semantic meaning.",
			},
			"client_ca_certificates": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Client CA certificate ID list. Required when the listener protocol is `HTTPS` and `certification_type` " +
					"is `MUTUAL`. Only HTTPS listeners support modifying this field. Treated as an unordered set; HCL element order has no semantic meaning.",
			},
			"http_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "HTTP version negotiated for this listener. Valid values: `HTTP/1.1`, `HTTP/2`. Only applicable to HTTPS listeners.",
			},

			// Computed
			"listener_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Listener instance ID.",
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Listener creation time.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Listener operational status.",
			},
			"endpoint_group_counts": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of endpoint groups attached to this listener.",
			},
		},
	}
}

func resourceTencentCloudGa2ListenerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_listener.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = ga2v20250115.NewCreateListenerRequest()
		response   = ga2v20250115.NewCreateListenerResponse()
		gaId       string
		listenerId string
		taskId     string
	)

	if v, ok := d.GetOk("global_accelerator_id"); ok {
		gaId = v.(string)
		request.GlobalAcceleratorId = helper.String(gaId)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("port_ranges"); ok {
		request.PortRanges = buildGa2ListenerPortRanges(v.([]interface{}))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("listener_type"); ok {
		request.ListenerType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("protocol"); ok {
		request.Protocol = helper.String(v.(string))
	}

	if v, ok := d.GetOk("idle_timeout"); ok {
		request.IdleTimeout = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("get_real_ip_type"); ok {
		request.GetRealIpType = helper.String(v.(string))
	}

	layer7 := isGa2ListenerLayer7(d)

	// Layer-4-only field. The cloud API rejects this on HTTP/HTTPS listeners with
	// `InvalidParameter.ApplicationLayerListenerCannotCarryParameters`.
	if !layer7 {
		if v, ok := d.GetOk("client_affinity"); ok {
			request.ClientAffinity = helper.String(v.(string))
		}
	}

	// Layer-7-only fields (HTTP & HTTPS). The cloud API rejects these on TCP/UDP listeners with
	// `InvalidParameter.TransportLayerListenerCannotCarryParameters`.
	if layer7 {
		if v, ok := d.GetOk("request_timeout"); ok {
			request.RequestTimeout = helper.IntUint64(v.(int))
		}

		//nolint:staticcheck
		if v, ok := d.GetOkExists("x_forwarded_for_real_ip"); ok {
			request.XForwardedForRealIp = helper.Bool(v.(bool))
		}
	}

	// HTTPS-only fields. The cloud API rejects these on HTTP listeners with
	// `InvalidParameter.HttpListenerCannotCarryParameters` and on TCP/UDP listeners with
	// `InvalidParameter.TransportLayerListenerCannotCarryParameters`.
	if isGa2ListenerHttps(d) {
		if v, ok := d.GetOk("certification_type"); ok {
			request.CertificationType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("cipher_policy_id"); ok {
			request.CipherPolicyId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("server_certificates"); ok {
			request.ServerCertificates = expandGa2ListenerStringSet(v.(*schema.Set))
		}

		if v, ok := d.GetOk("client_ca_certificates"); ok {
			request.ClientCaCertificates = expandGa2ListenerStringSet(v.(*schema.Set))
		}

		if v, ok := d.GetOk("http_version"); ok {
			request.HttpVersion = helper.String(v.(string))
		}
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e, "UnsupportedOperation.InstanceNotRunning", "UnsupportedOperation.InstanceStateNotAllowedOperate")
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ga2 listener failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create ga2 listener failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.ListenerId == nil {
		return fmt.Errorf("ListenerId is nil.")
	}
	listenerId = *response.Response.ListenerId

	if response.Response.TaskId == nil {
		return fmt.Errorf("TaskId is nil.")
	}
	taskId = *response.Response.TaskId

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{gaId, listenerId}, tccommon.FILED_SP))
	return resourceTencentCloudGa2ListenerRead(d, meta)
}

func resourceTencentCloudGa2ListenerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_listener.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	gaId, listenerId, err := parseGa2ListenerId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeGa2ListenerById(ctx, gaId, listenerId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_ga2_listener` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.GlobalAcceleratorId != nil {
		_ = d.Set("global_accelerator_id", respData.GlobalAcceleratorId)
	}

	if respData.ListenerId != nil {
		_ = d.Set("listener_id", respData.ListenerId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Description != nil {
		_ = d.Set("description", respData.Description)
	}

	if respData.Protocol != nil {
		_ = d.Set("protocol", respData.Protocol)
	}

	if respData.PortRanges != nil {
		_ = d.Set("port_ranges", flattenGa2ListenerPortRanges(respData.PortRanges))
	}

	if respData.XForwardedForRealIp != nil {
		_ = d.Set("x_forwarded_for_real_ip", respData.XForwardedForRealIp)
	}

	if respData.ClientAffinity != nil {
		_ = d.Set("client_affinity", respData.ClientAffinity)
	}

	if respData.ClientAffinityTime != nil {
		_ = d.Set("client_affinity_time", int(*respData.ClientAffinityTime))
	}

	if respData.CertificationType != nil {
		_ = d.Set("certification_type", respData.CertificationType)
	}

	if len(respData.ServerCertificates) > 0 {
		_ = d.Set("server_certificates", helper.PStrings(respData.ServerCertificates))
	}

	if len(respData.ClientCaCertificates) > 0 {
		_ = d.Set("client_ca_certificates", helper.PStrings(respData.ClientCaCertificates))
	}

	if respData.CipherPolicyId != nil {
		_ = d.Set("cipher_policy_id", respData.CipherPolicyId)
	}

	if respData.HttpVersion != nil {
		_ = d.Set("http_version", respData.HttpVersion)
	}

	if respData.RequestTimeout != nil {
		_ = d.Set("request_timeout", int(*respData.RequestTimeout))
	}

	if respData.CreateTime != nil {
		_ = d.Set("create_time", respData.CreateTime)
	}

	if respData.ListenerType != nil {
		_ = d.Set("listener_type", respData.ListenerType)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.EndpointGroupCounts != nil {
		_ = d.Set("endpoint_group_counts", int(*respData.EndpointGroupCounts))
	}

	if respData.GetRealIpType != nil {
		_ = d.Set("get_real_ip_type", respData.GetRealIpType)
	}

	if respData.IdleTimeout != nil {
		_ = d.Set("idle_timeout", int(*respData.IdleTimeout))
	}

	return nil
}

func resourceTencentCloudGa2ListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_listener.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	gaId, listenerId, err := parseGa2ListenerId(d.Id())
	if err != nil {
		return err
	}

	// Modify-supported fields per the SDK ModifyListenerRequest definition.
	// Layer-specific fields are tracked separately so we can both:
	//   (1) skip ModifyListener entirely when only an inapplicable-layer field changed, and
	//   (2) avoid carrying them in the request body, which the API rejects with one of:
	//       - `InvalidParameter.TransportLayerListenerCannotCarryParameters` (L7 fields on L4)
	//       - `InvalidParameter.ApplicationLayerListenerCannotCarryParameters` (L4 fields on L7)
	//       - `InvalidParameter.HttpListenerCannotCarryParameters` (HTTPS-only fields on HTTP)
	commonModifyFields := []string{
		"name", "description", "idle_timeout",
	}
	layer4ModifyFields := []string{
		"client_affinity", "client_affinity_time", "get_real_ip_type",
	}
	layer7ModifyFields := []string{
		"request_timeout", "x_forwarded_for_real_ip",
	}
	httpsModifyFields := []string{
		"certification_type", "cipher_policy_id", "server_certificates", "client_ca_certificates",
	}

	layer7 := isGa2ListenerLayer7(d)
	https := isGa2ListenerHttps(d)

	hasChange := func(fields []string) bool {
		for _, f := range fields {
			if d.HasChange(f) {
				return true
			}
		}
		return false
	}

	needModify := hasChange(commonModifyFields)
	if !needModify {
		if layer7 {
			needModify = hasChange(layer7ModifyFields)
			if !needModify && https {
				needModify = hasChange(httpsModifyFields)
			}
		} else {
			needModify = hasChange(layer4ModifyFields)
		}
	}

	if !needModify {
		return resourceTencentCloudGa2ListenerRead(d, meta)
	}

	request := ga2v20250115.NewModifyListenerRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("idle_timeout"); ok {
		request.IdleTimeout = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("get_real_ip_type"); ok {
		request.GetRealIpType = helper.String(v.(string))
	}

	// Layer-4-only fields. The cloud API rejects these on HTTP/HTTPS listeners with
	// `InvalidParameter.ApplicationLayerListenerCannotCarryParameters`.
	if !layer7 {
		if v, ok := d.GetOk("client_affinity"); ok {
			request.ClientAffinity = helper.String(v.(string))
		}

		if v, ok := d.GetOk("client_affinity_time"); ok {
			request.ClientAffinityTime = helper.IntUint64(v.(int))
		}
	}

	// Layer-7-only fields (HTTP & HTTPS). The cloud API rejects these on TCP/UDP listeners with
	// `InvalidParameter.TransportLayerListenerCannotCarryParameters`.
	if layer7 {
		if v, ok := d.GetOk("request_timeout"); ok {
			request.RequestTimeout = helper.IntUint64(v.(int))
		}

		//nolint:staticcheck
		if v, ok := d.GetOkExists("x_forwarded_for_real_ip"); ok {
			request.XForwardedForRealIp = helper.Bool(v.(bool))
		}
	}

	// HTTPS-only fields. The cloud API rejects these on HTTP listeners with
	// `InvalidParameter.HttpListenerCannotCarryParameters` and on TCP/UDP listeners with
	// `InvalidParameter.TransportLayerListenerCannotCarryParameters`.
	if https {
		if v, ok := d.GetOk("certification_type"); ok {
			request.CertificationType = helper.String(v.(string))
		}

		if v, ok := d.GetOk("cipher_policy_id"); ok {
			request.CipherPolicyId = helper.String(v.(string))
		}

		if v, ok := d.GetOk("server_certificates"); ok {
			request.ServerCertificates = expandGa2ListenerStringSet(v.(*schema.Set))
		}

		if v, ok := d.GetOk("client_ca_certificates"); ok {
			request.ClientCaCertificates = expandGa2ListenerStringSet(v.(*schema.Set))
		}
	}

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e, "UnsupportedOperation.InstanceNotRunning", "UnsupportedOperation.InstanceStateNotAllowedOperate")
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify ga2 listener failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update ga2 listener failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return err
	}

	return resourceTencentCloudGa2ListenerRead(d, meta)
}

func resourceTencentCloudGa2ListenerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_listener.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = ga2v20250115.NewDeleteListenerRequest()
	)

	gaId, listenerId, err := parseGa2ListenerId(d.Id())
	if err != nil {
		return err
	}

	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteListenerWithContext(ctx, request)
		if e != nil {
			if sdkerr, ok := e.(*sdkErrors.TencentCloudSDKError); ok {
				if sdkerr.Code == "ResourceNotFound" {
					return nil
				}
			}

			return tccommon.RetryError(e, "UnsupportedOperation.InstanceNotRunning", "UnsupportedOperation.InstanceStateNotAllowedOperate")
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ga2 listener failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete ga2 listener failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId, d.Timeout(schema.TimeoutDelete)); err != nil {
		return err
	}

	return nil
}

// parseGa2ListenerId splits the composite resource ID into its two components.
func parseGa2ListenerId(id string) (gaId, listenerId string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid resource id %q, expected format <global_accelerator_id>%s<listener_id>", id, tccommon.FILED_SP)
		return
	}
	gaId, listenerId = parts[0], parts[1]
	if gaId == "" || listenerId == "" {
		err = fmt.Errorf("invalid resource id %q, components must all be non-empty", id)
		return
	}
	return
}

// buildGa2ListenerPortRanges converts the schema list (MaxItems=1) into the SDK PortRanges struct.
func buildGa2ListenerPortRanges(rawList []interface{}) *ga2v20250115.PortRanges {
	if len(rawList) == 0 || rawList[0] == nil {
		return nil
	}
	m := rawList[0].(map[string]interface{})
	pr := &ga2v20250115.PortRanges{}

	if v, ok := m["from_port"].(int); ok {
		pr.FromPort = helper.IntUint64(v)
	}
	if v, ok := m["to_port"].(int); ok {
		pr.ToPort = helper.IntUint64(v)
	}
	return pr
}

// flattenGa2ListenerPortRanges maps SDK PortRanges back into the nested schema block payload.
func flattenGa2ListenerPortRanges(pr *ga2v20250115.PortRanges) []map[string]interface{} {
	if pr == nil {
		return nil
	}
	m := map[string]interface{}{}
	if pr.FromPort != nil {
		m["from_port"] = int(*pr.FromPort)
	}
	if pr.ToPort != nil {
		m["to_port"] = int(*pr.ToPort)
	}
	return []map[string]interface{}{m}
}

// isGa2ListenerLayer7 returns true when the listener protocol is HTTP / HTTPS.
// `protocol` is ForceNew, so its value is stable for the lifetime of the resource;
// however the field is Optional+Computed, which means a freshly-imported resource
// may briefly carry an empty string before Read populates it. We treat empty as L4
// (the default protocol per the API docs is TCP).
func isGa2ListenerLayer7(d *schema.ResourceData) bool {
	switch strings.ToUpper(d.Get("protocol").(string)) {
	case "HTTP", "HTTPS":
		return true
	default:
		return false
	}
}

// isGa2ListenerHttps returns true when the listener protocol is HTTPS.
// Certificate-related fields (`certification_type`, `cipher_policy_id`, `server_certificates`,
// `client_ca_certificates`) are accepted only by HTTPS listeners; the API rejects them on
// HTTP listeners with `InvalidParameter.HttpListenerCannotCarryParameters` and on TCP/UDP
// listeners with `InvalidParameter.TransportLayerListenerCannotCarryParameters`.
func isGa2ListenerHttps(d *schema.ResourceData) bool {
	return strings.EqualFold(d.Get("protocol").(string), "HTTPS")
}

// expandGa2ListenerStringSet converts a TypeSet of strings into a []*string suitable for the SDK.
func expandGa2ListenerStringSet(set *schema.Set) []*string {
	if set == nil {
		return nil
	}
	raw := set.List()
	if len(raw) == 0 {
		return nil
	}
	result := make([]*string, 0, len(raw))
	for _, item := range raw {
		s, ok := item.(string)
		if !ok || s == "" {
			continue
		}
		v := s
		result = append(result, &v)
	}
	return result
}
