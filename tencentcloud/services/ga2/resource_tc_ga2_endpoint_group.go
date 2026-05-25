package ga2

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudGa2EndpointGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGa2EndpointGroupCreate,
		Read:   resourceTencentCloudGa2EndpointGroupRead,
		Update: resourceTencentCloudGa2EndpointGroupUpdate,
		Delete: resourceTencentCloudGa2EndpointGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"global_accelerator_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Global accelerator instance ID.",
			},

			"listener_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Listener ID.",
			},

			"endpoint_group_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Endpoint group type. Valid values: `VIRTUAL`, `DEFAULT`.",
			},

			"endpoint_group_configuration": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Endpoint group configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Name. Maximum length is 60 bytes.",
						},
						"endpoint_group_region": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Region of the endpoint group.",
						},
						"endpoint_configurations": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Endpoint configurations under this group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"endpoint_type": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Endpoint type. Valid values: `Domain`, `PublicIp`.",
									},
									"endpoint_service": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Endpoint domain or IP.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Endpoint weight.",
									},
									"health_check_status": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Health check status. Valid values: `HEALTH`, `UNHEALTH`.",
									},
								},
							},
						},
						"check_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Health check protocol. Valid values: `TCP`, `HTTP`, `HTTPS`, `PING`, `CUSTOM`.",
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Description. Maximum length is 100 bytes.",
						},
						"check_port": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Health check port.",
						},
						"context_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Health check content type.",
						},
						"check_send_context": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Health check request payload.",
						},
						"check_recv_context": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Health check expected response.",
						},
						"enable_health_check": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether to enable health check.",
						},
						"connect_timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Response timeout in milliseconds.",
						},
						"health_check_interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Health check interval in seconds.",
						},
						"unhealthy_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Unhealthy threshold count.",
						},
						"healthy_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Healthy threshold count.",
						},
						"forward_protocol": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Forward protocol back to origin.",
						},
						"check_domain": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Health check domain.",
						},
						"check_path": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Health check URL path.",
						},
						"check_method": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Health check request method.",
						},
						"status_mask": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Status code masks for health check.",
						},
						"port_overrides": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Port overrides for the endpoint group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener_port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Listener port.",
									},
									"endpoint_port": {
										Type:        schema.TypeInt,
										Optional:    true,
										Computed:    true,
										Description: "Endpoint port.",
									},
								},
							},
						},
						"isp_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "ISP type.",
						},
						"cipher_policy_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "HTTPS cipher policy ID.",
						},
					},
				},
			},

			// computed
			"endpoint_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Endpoint group instance ID.",
			},
		},
	}
}

func resourceTencentCloudGa2EndpointGroupCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_endpoint_group.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId           = tccommon.GetLogId(tccommon.ContextNil)
		ctx             = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request         = ga2v20250115.NewCreateEndpointGroupRequest()
		response        = ga2v20250115.NewCreateEndpointGroupResponse()
		gaId            string
		listenerId      string
		endpointGroupId string
		taskId          string
	)

	if v, ok := d.GetOk("global_accelerator_id"); ok {
		gaId = v.(string)
		request.GlobalAcceleratorId = helper.String(gaId)
	}

	if v, ok := d.GetOk("listener_id"); ok {
		listenerId = v.(string)
		request.ListenerId = helper.String(listenerId)
	}

	if v, ok := d.GetOk("endpoint_group_type"); ok {
		request.EndpointGroupType = helper.String(v.(string))
	}

	if v, ok := d.GetOk("endpoint_group_configuration"); ok {
		request.EndpointGroupConfiguration = buildEndpointGroupConfiguration(v.([]interface{}))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().CreateEndpointGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create ga2 endpoint group failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create ga2 endpoint group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.EndpointGroupId == nil {
		return fmt.Errorf("EndpointGroupId is nil.")
	}
	endpointGroupId = *response.Response.EndpointGroupId

	if response.Response.TaskId == nil {
		return fmt.Errorf("TaskId is nil.")
	}
	taskId = *response.Response.TaskId

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId); err != nil {
		return err
	}

	d.SetId(strings.Join([]string{gaId, listenerId, endpointGroupId}, tccommon.FILED_SP))
	return resourceTencentCloudGa2EndpointGroupRead(d, meta)
}

func resourceTencentCloudGa2EndpointGroupRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_endpoint_group.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	gaId, listenerId, endpointGroupId, err := parseGa2EndpointGroupId(d.Id())
	if err != nil {
		return err
	}

	respData, err := service.DescribeGa2EndpointGroupById(ctx, gaId, listenerId, endpointGroupId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_ga2_endpoint_group` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.GlobalAcceleratorId != nil {
		_ = d.Set("global_accelerator_id", respData.GlobalAcceleratorId)
	}

	if respData.ListenerId != nil {
		_ = d.Set("listener_id", respData.ListenerId)
	}

	if respData.EndpointGroupId != nil {
		_ = d.Set("endpoint_group_id", respData.EndpointGroupId)
	}

	if respData.EndpointGroupType != nil {
		_ = d.Set("endpoint_group_type", respData.EndpointGroupType)
	}

	configurationList := []map[string]interface{}{flattenEndpointGroupConfigurationFromSet(respData)}
	_ = d.Set("endpoint_group_configuration", configurationList)

	return nil
}

func resourceTencentCloudGa2EndpointGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_endpoint_group.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	gaId, listenerId, endpointGroupId, err := parseGa2EndpointGroupId(d.Id())
	if err != nil {
		return err
	}

	if !d.HasChange("endpoint_group_configuration") {
		return resourceTencentCloudGa2EndpointGroupRead(d, meta)
	}

	request := ga2v20250115.NewModifyEndpointGroupRequest()
	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)
	request.EndpointGroupId = helper.String(endpointGroupId)

	if v, ok := d.GetOk("endpoint_group_configuration"); ok {
		applyConfigurationToModifyRequest(request, v.([]interface{}))
	}

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().ModifyEndpointGroupWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Modify ga2 endpoint group failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update ga2 endpoint group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId); err != nil {
		return err
	}

	return resourceTencentCloudGa2EndpointGroupRead(d, meta)
}

func resourceTencentCloudGa2EndpointGroupDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_ga2_endpoint_group.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = ga2v20250115.NewDeleteEndpointGroupsRequest()
	)

	gaId, listenerId, endpointGroupId, err := parseGa2EndpointGroupId(d.Id())
	if err != nil {
		return err
	}

	request.GlobalAcceleratorId = helper.String(gaId)
	request.ListenerId = helper.String(listenerId)
	request.EndpointGroupIds = []*string{helper.String(endpointGroupId)}

	var taskId string
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseGa2V20250115Client().DeleteEndpointGroupsWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.TaskId == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete ga2 endpoint group failed, Response is nil."))
		}

		taskId = *result.Response.TaskId
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete ga2 endpoint group failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	service := Ga2Service{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if err := service.WaitForGa2TaskFinish(ctx, taskId); err != nil {
		return err
	}

	return nil
}

// parseGa2EndpointGroupId splits the composite resource ID into its three components.
func parseGa2EndpointGroupId(id string) (gaId, listenerId, endpointGroupId string, err error) {
	parts := strings.Split(id, tccommon.FILED_SP)
	if len(parts) != 3 {
		err = fmt.Errorf("invalid resource id %q, expected format <global_accelerator_id>%s<listener_id>%s<endpoint_group_id>", id, tccommon.FILED_SP, tccommon.FILED_SP)
		return
	}
	gaId, listenerId, endpointGroupId = parts[0], parts[1], parts[2]
	if gaId == "" || listenerId == "" || endpointGroupId == "" {
		err = fmt.Errorf("invalid resource id %q, components must all be non-empty", id)
		return
	}
	return
}

// buildEndpointGroupConfiguration converts the schema list into the SDK struct.
func buildEndpointGroupConfiguration(rawList []interface{}) *ga2v20250115.EndpointGroupConfiguration {
	if len(rawList) == 0 || rawList[0] == nil {
		return nil
	}
	m := rawList[0].(map[string]interface{})
	cfg := &ga2v20250115.EndpointGroupConfiguration{}

	if v, ok := m["name"].(string); ok && v != "" {
		cfg.Name = helper.String(v)
	}
	if v, ok := m["endpoint_group_region"].(string); ok && v != "" {
		cfg.EndpointGroupRegion = helper.String(v)
	}
	if v, ok := m["endpoint_configurations"]; ok {
		cfg.EndpointConfigurations = buildEndpointConfigurations(v.([]interface{}))
	}
	if v, ok := m["check_type"].(string); ok && v != "" {
		cfg.CheckType = helper.String(v)
	}
	if v, ok := m["description"].(string); ok && v != "" {
		cfg.Description = helper.String(v)
	}
	if v, ok := m["check_port"].(string); ok && v != "" {
		cfg.CheckPort = helper.String(v)
	}
	if v, ok := m["context_type"].(string); ok && v != "" {
		cfg.ContextType = helper.String(v)
	}
	if v, ok := m["check_send_context"].(string); ok && v != "" {
		cfg.CheckSendContext = helper.String(v)
	}
	if v, ok := m["check_recv_context"].(string); ok && v != "" {
		cfg.CheckRecvContext = helper.String(v)
	}
	if v, ok := m["enable_health_check"].(bool); ok {
		cfg.EnableHealthCheck = helper.Bool(v)
	}
	if v, ok := m["connect_timeout"].(int); ok && v > 0 {
		cfg.ConnectTimeout = helper.IntUint64(v)
	}
	if v, ok := m["health_check_interval"].(int); ok && v > 0 {
		cfg.HealthCheckInterval = helper.IntUint64(v)
	}
	if v, ok := m["unhealthy_threshold"].(int); ok && v > 0 {
		cfg.UnhealthyThreshold = helper.IntUint64(v)
	}
	if v, ok := m["healthy_threshold"].(int); ok && v > 0 {
		cfg.HealthyThreshold = helper.IntUint64(v)
	}
	if v, ok := m["forward_protocol"].(string); ok && v != "" {
		cfg.ForwardProtocol = helper.String(v)
	}
	if v, ok := m["check_domain"].(string); ok && v != "" {
		cfg.CheckDomain = helper.String(v)
	}
	if v, ok := m["check_path"].(string); ok && v != "" {
		cfg.CheckPath = helper.String(v)
	}
	if v, ok := m["check_method"].(string); ok && v != "" {
		cfg.CheckMethod = helper.String(v)
	}
	if v, ok := m["status_mask"]; ok {
		cfg.StatusMask = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := m["port_overrides"]; ok {
		cfg.PortOverrides = buildPortOverrides(v.([]interface{}))
	}
	if v, ok := m["isp_type"].(string); ok && v != "" {
		cfg.IspType = helper.String(v)
	}
	if v, ok := m["cipher_policy_id"].(string); ok && v != "" {
		cfg.CipherPolicyId = helper.String(v)
	}
	return cfg
}

// buildEndpointConfigurations converts the schema list into the SDK struct slice.
func buildEndpointConfigurations(rawList []interface{}) []*ga2v20250115.EndpointConfigurations {
	result := make([]*ga2v20250115.EndpointConfigurations, 0, len(rawList))
	for _, item := range rawList {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		ec := &ga2v20250115.EndpointConfigurations{}

		if v, ok := m["endpoint_type"].(string); ok && v != "" {
			ec.EndpointType = helper.String(v)
		}
		if v, ok := m["endpoint_service"].(string); ok && v != "" {
			ec.EndpointService = helper.String(v)
		}
		if v, ok := m["weight"].(int); ok {
			ec.Weight = helper.IntUint64(v)
		}

		result = append(result, ec)
	}
	return result
}

// buildPortOverrides converts the schema list into the SDK struct slice.
func buildPortOverrides(rawList []interface{}) []*ga2v20250115.PortOverride {
	result := make([]*ga2v20250115.PortOverride, 0, len(rawList))
	for _, item := range rawList {
		if item == nil {
			continue
		}
		m := item.(map[string]interface{})
		po := &ga2v20250115.PortOverride{}

		if v, ok := m["listener_port"].(int); ok && v > 0 {
			po.ListenerPort = helper.IntUint64(v)
		}
		if v, ok := m["endpoint_port"].(int); ok && v > 0 {
			po.EndpointPort = helper.IntUint64(v)
		}

		result = append(result, po)
	}
	return result
}

// applyConfigurationToModifyRequest copies the nested configuration block onto the flat ModifyEndpointGroup request.
func applyConfigurationToModifyRequest(request *ga2v20250115.ModifyEndpointGroupRequest, rawList []interface{}) {
	if len(rawList) == 0 || rawList[0] == nil {
		return
	}
	m := rawList[0].(map[string]interface{})

	if v, ok := m["name"].(string); ok && v != "" {
		request.Name = helper.String(v)
	}
	if v, ok := m["description"].(string); ok && v != "" {
		request.Description = helper.String(v)
	}
	if v, ok := m["endpoint_configurations"]; ok {
		request.EndpointConfigurations = buildEndpointConfigurations(v.([]interface{}))
	}
	if v, ok := m["enable_health_check"].(bool); ok {
		request.EnableHealthCheck = helper.Bool(v)
	}
	if v, ok := m["connect_timeout"].(int); ok && v > 0 {
		request.ConnectTimeout = helper.IntUint64(v)
	}
	if v, ok := m["health_check_interval"].(int); ok && v > 0 {
		request.HealthCheckInterval = helper.IntUint64(v)
	}
	if v, ok := m["unhealthy_threshold"].(int); ok && v > 0 {
		request.UnhealthyThreshold = helper.IntUint64(v)
	}
	if v, ok := m["healthy_threshold"].(int); ok && v > 0 {
		request.HealthyThreshold = helper.IntUint64(v)
	}
	if v, ok := m["check_type"].(string); ok && v != "" {
		request.CheckType = helper.String(v)
	}
	if v, ok := m["check_port"].(string); ok && v != "" {
		// ModifyEndpointGroup uses uint64 CheckPort; parse the string from schema.
		if port, err := strconv.ParseUint(v, 10, 64); err == nil {
			request.CheckPort = &port
		}
	}
	if v, ok := m["context_type"].(string); ok && v != "" {
		request.ContextType = helper.String(v)
	}
	if v, ok := m["check_send_context"].(string); ok && v != "" {
		request.CheckSendContext = helper.String(v)
	}
	if v, ok := m["check_recv_context"].(string); ok && v != "" {
		request.CheckRecvContext = helper.String(v)
	}
	if v, ok := m["check_domain"].(string); ok && v != "" {
		request.CheckDomain = helper.String(v)
	}
	if v, ok := m["check_path"].(string); ok && v != "" {
		request.CheckPath = helper.String(v)
	}
	if v, ok := m["check_method"].(string); ok && v != "" {
		request.CheckMethod = helper.String(v)
	}
	if v, ok := m["status_mask"]; ok {
		request.StatusMask = helper.InterfacesStringsPoint(v.([]interface{}))
	}
	if v, ok := m["forward_protocol"].(string); ok && v != "" {
		request.ForwardProtocol = helper.String(v)
	}
	if v, ok := m["port_overrides"]; ok {
		request.PortOverrides = buildPortOverrides(v.([]interface{}))
	}
	if v, ok := m["cipher_policy_id"].(string); ok && v != "" {
		request.CipherPolicyId = helper.String(v)
	}
}

// flattenEndpointGroupConfigurationFromSet maps EndpointGroupConfigurationSet (Describe response) into a single nested block payload.
func flattenEndpointGroupConfigurationFromSet(set *ga2v20250115.EndpointGroupConfigurationSet) map[string]interface{} {
	m := map[string]interface{}{}
	if set == nil {
		return m
	}

	if set.Name != nil {
		m["name"] = *set.Name
	}
	if set.EndpointGroupRegion != nil {
		m["endpoint_group_region"] = *set.EndpointGroupRegion
	}
	if set.Description != nil {
		m["description"] = *set.Description
	}
	if set.CheckType != nil {
		m["check_type"] = *set.CheckType
	}
	if set.CheckPort != nil {
		// Response carries CheckPort as uint64; schema stores it as string for symmetry with Create input.
		m["check_port"] = strconv.FormatUint(*set.CheckPort, 10)
	}
	if set.ContextType != nil {
		m["context_type"] = *set.ContextType
	}
	if set.CheckSendContext != nil {
		m["check_send_context"] = *set.CheckSendContext
	}
	if set.CheckRecvContext != nil {
		m["check_recv_context"] = *set.CheckRecvContext
	}
	if set.EnableHealthCheck != nil {
		m["enable_health_check"] = *set.EnableHealthCheck
	}
	if set.ConnectTimeout != nil {
		m["connect_timeout"] = int(*set.ConnectTimeout)
	}
	if set.HealthCheckInterval != nil {
		m["health_check_interval"] = int(*set.HealthCheckInterval)
	}
	if set.UnhealthyThreshold != nil {
		m["unhealthy_threshold"] = int(*set.UnhealthyThreshold)
	}
	if set.HealthyThreshold != nil {
		m["healthy_threshold"] = int(*set.HealthyThreshold)
	}
	if set.ForwardProtocol != nil {
		m["forward_protocol"] = *set.ForwardProtocol
	}
	if set.CheckDomain != nil {
		m["check_domain"] = *set.CheckDomain
	}
	if set.CheckPath != nil {
		m["check_path"] = *set.CheckPath
	}
	if set.CheckMethod != nil {
		m["check_method"] = *set.CheckMethod
	}
	if len(set.StatusMask) > 0 {
		m["status_mask"] = helper.PStrings(set.StatusMask)
	}
	if len(set.EndpointConfigurations) > 0 {
		m["endpoint_configurations"] = flattenEndpointConfigurations(set.EndpointConfigurations)
	}
	if len(set.PortOverrides) > 0 {
		m["port_overrides"] = flattenPortOverrides(set.PortOverrides)
	}
	if set.IspType != nil {
		m["isp_type"] = *set.IspType
	}
	if set.CipherPolicyId != nil {
		m["cipher_policy_id"] = *set.CipherPolicyId
	}
	return m
}

// flattenEndpointConfigurations maps SDK struct slice back into terraform list.
func flattenEndpointConfigurations(items []*ga2v20250115.EndpointConfigurations) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		m := map[string]interface{}{}
		if item.EndpointType != nil {
			m["endpoint_type"] = *item.EndpointType
		}
		if item.EndpointService != nil {
			m["endpoint_service"] = *item.EndpointService
		}
		if item.Weight != nil {
			m["weight"] = int(*item.Weight)
		}
		if item.HealthCheckStatus != nil {
			m["health_check_status"] = *item.HealthCheckStatus
		}
		result = append(result, m)
	}
	return result
}

// flattenPortOverrides maps SDK struct slice back into terraform list.
func flattenPortOverrides(items []*ga2v20250115.PortOverride) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		m := map[string]interface{}{}
		if item.ListenerPort != nil {
			m["listener_port"] = int(*item.ListenerPort)
		}
		if item.EndpointPort != nil {
			m["endpoint_port"] = int(*item.EndpointPort)
		}
		result = append(result, m)
	}
	return result
}
