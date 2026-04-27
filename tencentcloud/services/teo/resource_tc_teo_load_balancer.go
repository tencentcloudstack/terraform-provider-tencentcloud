package teo

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoLoadBalancer() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTeoLoadBalancerCreate,
		Read:   resourceTencentCloudTeoLoadBalancerRead,
		Update: resourceTencentCloudTeoLoadBalancerUpdate,
		Delete: resourceTencentCloudTeoLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Site ID.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Load balancer instance name, 1-200 characters, allowed characters: `a-z`, `A-Z`, `0-9`, `_`, `-`.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance type. Valid values: `HTTP` (HTTP-specific, supports HTTP-specific and general origin groups, only referenced by site acceleration services); `GENERAL` (general, only supports general origin groups, can be referenced by site acceleration and Layer-4 proxy).",
			},

			"origin_groups": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Source origin group list with failover priority.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Priority, format: `priority_` + number, highest priority is `priority_1`. Valid values: `priority_1`, `priority_2`..., `priority_10`.",
						},
						"origin_group_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Origin group ID.",
						},
					},
				},
			},

			"health_checker": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Health check policy. If not set, health check is disabled by default.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Health check type. Valid values: `HTTP`, `HTTPS`, `TCP`, `UDP`, `ICMP Ping`, `NoCheck`. `NoCheck` means health check is disabled.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Check port. Required when `type` is `HTTP`, `HTTPS`, `TCP`, or `UDP`.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Check interval in seconds. Valid values: `30`, `60`, `180`, `300`, `600`.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Timeout in seconds for each health check. Must be less than `interval`. Default: 5.",
						},
						"health_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of consecutive healthy results to mark the origin as healthy. Default: 3, minimum: 1.",
						},
						"critical_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Number of consecutive unhealthy results to mark the origin as unhealthy. Default: 2.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Probe path, valid when `type` is `HTTP` or `HTTPS`. Full host/path without protocol, e.g., `www.example.com/test`.",
						},
						"method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request method, valid when `type` is `HTTP` or `HTTPS`. Valid values: `GET`, `HEAD`.",
						},
						"expected_codes": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "Expected response status codes, valid when `type` is `HTTP` or `HTTPS`. e.g., `[\"200\", \"301\"]`.",
						},
						"follow_redirect": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to follow 301/302 redirects, valid when `type` is `HTTP` or `HTTPS`. Valid values: `true`, `false`.",
						},
						"headers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Custom HTTP request headers (up to 10), valid when `type` is `HTTP` or `HTTPS`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Custom header key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Custom header value.",
									},
								},
							},
						},
						"send_context": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Content to send during health check, valid when `type` is `UDP`. Only ASCII visible characters allowed, max length 500.",
						},
						"recv_context": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Expected response content from origin, valid when `type` is `UDP`. Only ASCII visible characters allowed, max length 500.",
						},
					},
				},
			},

			"steering_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Traffic steering policy between origin groups. Valid value: `Pritory` (failover by priority). Default: `Pritory`.",
			},

			"failover_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Retry policy on request failure. Valid values: `OtherOriginGroup` (retry next priority origin group); `OtherRecordInOriginGroup` (retry another origin in the same group). Default: `OtherRecordInOriginGroup`.",
			},

			// computed
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancer instance ID.",
			},
		},
	}
}

func resourceTencentCloudTeoLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_load_balancer.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = teo.NewCreateLoadBalancerRequest()
		response   = teo.NewCreateLoadBalancerResponse()
		zoneId     string
		instanceId string
	)

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
		zoneId = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_groups"); ok {
		for _, item := range v.([]interface{}) {
			ogMap := item.(map[string]interface{})
			og := &teo.OriginGroupInLoadBalancer{}
			if val, ok := ogMap["priority"].(string); ok && val != "" {
				og.Priority = helper.String(val)
			}
			if val, ok := ogMap["origin_group_id"].(string); ok && val != "" {
				og.OriginGroupId = helper.String(val)
			}
			request.OriginGroups = append(request.OriginGroups, og)
		}
	}

	if v, ok := d.GetOk("health_checker"); ok {
		request.HealthChecker = buildTeoHealthChecker(v.([]interface{}))
	}

	if v, ok := d.GetOk("steering_policy"); ok {
		request.SteeringPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("failover_policy"); ok {
		request.FailoverPolicy = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreateLoadBalancerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo load balancer failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create teo load balancer failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.InstanceId == nil {
		return fmt.Errorf("InstanceId is nil.")
	}

	instanceId = *response.Response.InstanceId
	d.SetId(strings.Join([]string{zoneId, instanceId}, tccommon.FILED_SP))

	// Wait for load balancer to become Running
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if _, err := (&resource.StateChangeConf{
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{"Pending"},
		Target:     []string{"Running"},
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Refresh: func() (interface{}, string, error) {
			lb, e := service.DescribeTeoLoadBalancerById(ctx, zoneId, instanceId)
			if e != nil {
				return nil, "", e
			}
			if lb == nil {
				return nil, "Pending", nil
			}
			status := "Pending"
			if lb.Status != nil {
				status = *lb.Status
			}
			return lb, status, nil
		},
	}).WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for teo load balancer (%s) to become Running: %s", d.Id(), err)
	}

	return resourceTencentCloudTeoLoadBalancerRead(d, meta)
}

func resourceTencentCloudTeoLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_load_balancer.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	zoneId := idSplit[0]
	instanceId := idSplit[1]

	respData, err := service.DescribeTeoLoadBalancerById(ctx, zoneId, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_load_balancer` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	_ = d.Set("zone_id", zoneId)

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	if respData.Name != nil {
		_ = d.Set("name", respData.Name)
	}

	if respData.Type != nil {
		_ = d.Set("type", respData.Type)
	}

	if respData.SteeringPolicy != nil {
		_ = d.Set("steering_policy", respData.SteeringPolicy)
	}

	if respData.FailoverPolicy != nil {
		_ = d.Set("failover_policy", respData.FailoverPolicy)
	}

	if respData.OriginGroupHealthStatus != nil {
		ogList := make([]map[string]interface{}, 0, len(respData.OriginGroupHealthStatus))
		for _, og := range respData.OriginGroupHealthStatus {
			if og == nil {
				continue
			}
			ogMap := map[string]interface{}{}
			if og.Priority != nil {
				ogMap["priority"] = og.Priority
			}
			if og.OriginGroupID != nil {
				ogMap["origin_group_id"] = og.OriginGroupID
			}
			ogList = append(ogList, ogMap)
		}
		_ = d.Set("origin_groups", ogList)
	}

	if respData.HealthChecker != nil {
		_ = d.Set("health_checker", flattenTeoHealthChecker(respData.HealthChecker))
	}

	return nil
}

func resourceTencentCloudTeoLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_load_balancer.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewModifyLoadBalancerRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	request.ZoneId = helper.String(idSplit[0])
	request.InstanceId = helper.String(idSplit[1])

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_groups"); ok {
		for _, item := range v.([]interface{}) {
			ogMap := item.(map[string]interface{})
			og := &teo.OriginGroupInLoadBalancer{}
			if val, ok := ogMap["priority"].(string); ok && val != "" {
				og.Priority = helper.String(val)
			}
			if val, ok := ogMap["origin_group_id"].(string); ok && val != "" {
				og.OriginGroupId = helper.String(val)
			}
			request.OriginGroups = append(request.OriginGroups, og)
		}
	}

	if v, ok := d.GetOk("health_checker"); ok {
		request.HealthChecker = buildTeoHealthChecker(v.([]interface{}))
	}

	if v, ok := d.GetOk("steering_policy"); ok {
		request.SteeringPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("failover_policy"); ok {
		request.FailoverPolicy = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyLoadBalancerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update teo load balancer failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Wait for load balancer to become Running after update
	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	if _, err := (&resource.StateChangeConf{
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
		Pending:    []string{"Pending"},
		Target:     []string{"Running"},
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Refresh: func() (interface{}, string, error) {
			lb, e := service.DescribeTeoLoadBalancerById(ctx, *request.ZoneId, *request.InstanceId)
			if e != nil {
				return nil, "", e
			}
			if lb == nil {
				return nil, "Pending", nil
			}
			status := "Pending"
			if lb.Status != nil {
				status = *lb.Status
			}
			return lb, status, nil
		},
	}).WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for teo load balancer (%s) to become Running: %s", d.Id(), err)
	}

	return resourceTencentCloudTeoLoadBalancerRead(d, meta)
}

func resourceTencentCloudTeoLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_load_balancer.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teo.NewDeleteLoadBalancerRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken, %s", d.Id())
	}

	request.ZoneId = helper.String(idSplit[0])
	request.InstanceId = helper.String(idSplit[1])

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DeleteLoadBalancerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete teo load balancer failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

// buildTeoHealthChecker converts the health_checker schema list to *teo.HealthChecker.
func buildTeoHealthChecker(v []interface{}) *teo.HealthChecker {
	if len(v) == 0 {
		return nil
	}
	m := v[0].(map[string]interface{})
	hc := &teo.HealthChecker{}

	if val, ok := m["type"].(string); ok && val != "" {
		hc.Type = helper.String(val)
	}
	if val, ok := m["port"].(int); ok && val != 0 {
		hc.Port = helper.IntUint64(val)
	}
	if val, ok := m["interval"].(int); ok && val != 0 {
		hc.Interval = helper.IntUint64(val)
	}
	if val, ok := m["timeout"].(int); ok && val != 0 {
		hc.Timeout = helper.IntUint64(val)
	}
	if val, ok := m["health_threshold"].(int); ok && val != 0 {
		hc.HealthThreshold = helper.IntUint64(val)
	}
	if val, ok := m["critical_threshold"].(int); ok && val != 0 {
		hc.CriticalThreshold = helper.IntUint64(val)
	}
	if val, ok := m["path"].(string); ok && val != "" {
		hc.Path = helper.String(val)
	}
	if val, ok := m["method"].(string); ok && val != "" {
		hc.Method = helper.String(val)
	}
	if val, ok := m["expected_codes"].([]interface{}); ok && len(val) > 0 {
		for _, code := range val {
			hc.ExpectedCodes = append(hc.ExpectedCodes, helper.String(code.(string)))
		}
	}
	if val, ok := m["follow_redirect"].(string); ok && val != "" {
		hc.FollowRedirect = helper.String(val)
	}
	if val, ok := m["headers"].([]interface{}); ok && len(val) > 0 {
		for _, h := range val {
			hMap := h.(map[string]interface{})
			header := &teo.CustomizedHeader{}
			if k, ok := hMap["key"].(string); ok && k != "" {
				header.Key = helper.String(k)
			}
			if v, ok := hMap["value"].(string); ok && v != "" {
				header.Value = helper.String(v)
			}
			hc.Headers = append(hc.Headers, header)
		}
	}
	if val, ok := m["send_context"].(string); ok && val != "" {
		hc.SendContext = helper.String(val)
	}
	if val, ok := m["recv_context"].(string); ok && val != "" {
		hc.RecvContext = helper.String(val)
	}
	return hc
}

// flattenTeoHealthChecker converts *teo.HealthChecker to the schema list format.
func flattenTeoHealthChecker(hc *teo.HealthChecker) []map[string]interface{} {
	if hc == nil {
		return nil
	}
	m := map[string]interface{}{}
	if hc.Type != nil {
		m["type"] = hc.Type
	}
	if hc.Port != nil {
		m["port"] = int(*hc.Port)
	}
	if hc.Interval != nil {
		m["interval"] = int(*hc.Interval)
	}
	if hc.Timeout != nil {
		m["timeout"] = int(*hc.Timeout)
	}
	if hc.HealthThreshold != nil {
		m["health_threshold"] = int(*hc.HealthThreshold)
	}
	if hc.CriticalThreshold != nil {
		m["critical_threshold"] = int(*hc.CriticalThreshold)
	}
	if hc.Path != nil {
		m["path"] = hc.Path
	}
	if hc.Method != nil {
		m["method"] = hc.Method
	}
	if len(hc.ExpectedCodes) > 0 {
		codes := make([]string, 0, len(hc.ExpectedCodes))
		for _, c := range hc.ExpectedCodes {
			if c != nil {
				codes = append(codes, *c)
			}
		}
		m["expected_codes"] = codes
	}
	if hc.FollowRedirect != nil {
		m["follow_redirect"] = hc.FollowRedirect
	}
	if len(hc.Headers) > 0 {
		headerList := make([]map[string]interface{}, 0, len(hc.Headers))
		for _, h := range hc.Headers {
			if h == nil {
				continue
			}
			hMap := map[string]interface{}{}
			if h.Key != nil {
				hMap["key"] = h.Key
			}
			if h.Value != nil {
				hMap["value"] = h.Value
			}
			headerList = append(headerList, hMap)
		}
		m["headers"] = headerList
	}
	if hc.SendContext != nil {
		m["send_context"] = hc.SendContext
	}
	if hc.RecvContext != nil {
		m["recv_context"] = hc.RecvContext
	}
	return []map[string]interface{}{m}
}
