package ga2

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

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
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"global_accelerator_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Global accelerator instance ID.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Listener name. Maximum length is 60 bytes.",
			},
			"port_ranges": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Port ranges.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Start port.",
						},
						"to_port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "End port.",
						},
					},
				},
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description. Maximum length is 100 bytes.",
			},
			"listener_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Listener type, default is intelligent routing.",
			},
			"protocol": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Protocol, default is TCP.",
			},
			"idle_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Connection idle timeout.",
			},
			"get_real_ip_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Layer 4 method to get real IP. Valid values: `TOA`, `ProxyProtocol`.",
			},
			"client_affinity": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable session persistence.",
			},
			"client_affinity_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Session persistence time.",
			},
			"request_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Request timeout.",
			},
			"x_forwarded_for_real_ip": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable layer 7 method to get real IP.",
			},
			"certification_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Certificate type. `UNIDIRECTIONAL`: one-way. `MUTUAL`: two-way.",
			},
			"cipher_policy_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Cipher policy ID.",
			},
			"server_certificates": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Server certificates.",
			},
			"client_ca_certificates": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Client CA certificates.",
			},
			// computed
			"listener_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Listener ID.",
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
		portRangesList := v.([]interface{})
		if len(portRangesList) > 0 && portRangesList[0] != nil {
			m := portRangesList[0].(map[string]interface{})
			portRanges := &ga2v20250115.PortRanges{}
			if fromPort, ok := m["from_port"].(int); ok {
				portRanges.FromPort = helper.IntUint64(fromPort)
			}
			if toPort, ok := m["to_port"].(int); ok {
				portRanges.ToPort = helper.IntUint64(toPort)
			}
			request.PortRanges = portRanges
		}
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

	if v, ok := d.GetOk("client_affinity"); ok {
		request.ClientAffinity = helper.String(v.(string))
	}

	if v, ok := d.GetOk("request_timeout"); ok {
		request.RequestTimeout = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("x_forwarded_for_real_ip"); ok {
		request.XForwardedForRealIp = helper.Bool(v.(bool))
	}

	if v, ok := d.GetOk("certification_type"); ok {
		request.CertificationType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cipher_policy_id"); ok {
		request.CipherPolicyId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("server_certificates"); ok {
		request.ServerCertificates = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	if v, ok := d.GetOk("client_ca_certificates"); ok {
		request.ClientCaCertificates = helper.InterfacesStringsPoint(v.([]interface{}))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
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
		portRangesMap := map[string]interface{}{}
		if respData.PortRanges.FromPort != nil {
			portRangesMap["from_port"] = int(*respData.PortRanges.FromPort)
		}
		if respData.PortRanges.ToPort != nil {
			portRangesMap["to_port"] = int(*respData.PortRanges.ToPort)
		}
		_ = d.Set("port_ranges", []interface{}{portRangesMap})
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

	if respData.ServerCertificates != nil {
		_ = d.Set("server_certificates", helper.PStrings(respData.ServerCertificates))
	}

	if respData.ClientCaCertificates != nil {
		_ = d.Set("client_ca_certificates", helper.PStrings(respData.ClientCaCertificates))
	}

	if respData.CipherPolicyId != nil {
		_ = d.Set("cipher_policy_id", respData.CipherPolicyId)
	}

	if respData.RequestTimeout != nil {
		_ = d.Set("request_timeout", int(*respData.RequestTimeout))
	}

	if respData.ListenerType != nil {
		_ = d.Set("listener_type", respData.ListenerType)
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

	needChange := false
	request := ga2v20250115.NewModifyListenerRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)

	if d.HasChange("name") {
		needChange = true
		request.Name = helper.String(d.Get("name").(string))
	}

	if d.HasChange("description") {
		needChange = true
		request.Description = helper.String(d.Get("description").(string))
	}

	if d.HasChange("idle_timeout") {
		needChange = true
		request.IdleTimeout = helper.IntUint64(d.Get("idle_timeout").(int))
	}

	if d.HasChange("client_affinity") {
		needChange = true
		request.ClientAffinity = helper.String(d.Get("client_affinity").(string))
	}

	if d.HasChange("client_affinity_time") {
		needChange = true
		request.ClientAffinityTime = helper.IntUint64(d.Get("client_affinity_time").(int))
	}

	if d.HasChange("request_timeout") {
		needChange = true
		request.RequestTimeout = helper.IntUint64(d.Get("request_timeout").(int))
	}

	if d.HasChange("x_forwarded_for_real_ip") {
		needChange = true
		request.XForwardedForRealIp = helper.Bool(d.Get("x_forwarded_for_real_ip").(bool))
	}

	if d.HasChange("certification_type") {
		needChange = true
		request.CertificationType = helper.String(d.Get("certification_type").(string))
	}

	if d.HasChange("cipher_policy_id") {
		needChange = true
		request.CipherPolicyId = helper.String(d.Get("cipher_policy_id").(string))
	}

	if d.HasChange("server_certificates") {
		needChange = true
		request.ServerCertificates = helper.InterfacesStringsPoint(d.Get("server_certificates").([]interface{}))
	}

	if d.HasChange("client_ca_certificates") {
		needChange = true
		request.ClientCaCertificates = helper.InterfacesStringsPoint(d.Get("client_ca_certificates").([]interface{}))
	}

	if d.HasChange("get_real_ip_type") {
		needChange = true
		request.GetRealIpType = helper.String(d.Get("get_real_ip_type").(string))
	}

	if !needChange {
		return resourceTencentCloudGa2ListenerRead(d, meta)
	}

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyListenerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
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
			return tccommon.RetryError(e)
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
