package teo

import (
	"context"
	"fmt"
	"log"
	"strings"

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
				Description: "Instance name, can be 1-200 characters, allowed characters are a-z, A-Z, 0-9, _, -.",
			},

			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance type, valid values: `HTTP` (HTTP dedicated type), `GENERAL` (general type).",
			},

			"origin_groups": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Origin group list and corresponding failover scheduling priority.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"priority": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Priority, format is 'priority_' + 'number', the highest priority is 'priority_1'. Valid values: priority_1 to priority_10.",
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
				MaxItems:    1,
				Optional:    true,
				Description: "Health check policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Health check policy type, valid values: `HTTP`, `HTTPS`, `TCP`, `UDP`, `ICMP Ping`, `NoCheck`.",
						},
						"port": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Check port. Required when Type=HTTP, HTTPS, TCP or UDP.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Check frequency, how often to initiate a health check task, in seconds. Valid values: 30, 60, 180, 300, 600.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Timeout for each health check, in seconds, default is 5s, must be less than Interval.",
						},
						"health_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Health threshold, the number of consecutive health checks that are 'healthy' before judging the origin as 'healthy', default 3, minimum 1.",
						},
						"critical_threshold": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Unhealthy threshold, the number of consecutive health checks that are 'unhealthy' before judging the origin as 'unhealthy', default 2.",
						},
						"path": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Detection path, only valid when Type=HTTP or HTTPS. Need to fill in the complete host/path, excluding the protocol part.",
						},
						"method": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request method, only valid when Type=HTTP or HTTPS. Valid values: `GET`, `HEAD`.",
						},
						"expected_codes": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Expected status codes for health determination, only valid when Type=HTTP or HTTPS.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"headers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Custom HTTP request headers for detection, only valid when Type=HTTP or HTTPS, up to 10.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"key": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Custom header Key.",
									},
									"value": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Custom header Value.",
									},
								},
							},
						},
						"follow_redirect": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Whether to enable 301/302 redirect following, only valid when Type=HTTP or HTTPS.",
						},
						"send_context": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Content sent by health check, only valid when Type=UDP. Only ASCII visible characters, max 500 characters.",
						},
						"recv_context": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Expected response from origin for health check, only valid when Type=UDP. Only ASCII visible characters, max 500 characters.",
						},
					},
				},
			},

			"steering_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Traffic scheduling policy between origin groups, valid values: `Pritory` (failover by priority order). Default is Pritory.",
			},

			"failover_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Request retry policy when accessing an origin fails, valid values: `OtherOriginGroup` (retry next priority origin group), `OtherRecordInOriginGroup` (retry other origins in the same group). Default is OtherRecordInOriginGroup.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancer instance ID.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Load balancer status, valid values: `Pending` (deploying), `Deleting` (deleting), `Running` (effective).",
			},

			"origin_group_health_status": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Origin group health status.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin group ID.",
						},
						"origin_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin group name.",
						},
						"origin_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin group type, valid values: `HTTP`, `GENERAL`.",
						},
						"priority": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Priority.",
						},
						"origin_health_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Health status of origins in the origin group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Origin.",
									},
									"healthy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Origin health status, valid values: `Healthy`, `Unhealthy`, `Undetected`.",
									},
								},
							},
						},
					},
				},
			},

			"l4_used_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of L4 proxy instances bound to this load balancer instance.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"l7_used_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of L7 domain names bound to this load balancer instance.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"references": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of instances that reference this load balancer.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference service type.",
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Reference instance name.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTeoLoadBalancerCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_load_balancer.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	var (
		zoneId     string
		instanceId string
	)
	var (
		request  = teo.NewCreateLoadBalancerRequest()
		response = teo.NewCreateLoadBalancerResponse()
	)

	if v, ok := d.GetOk("zone_id"); ok {
		zoneId = v.(string)
	}
	request.ZoneId = helper.String(zoneId)

	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("type"); ok {
		request.Type = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_groups"); ok {
		for _, item := range v.([]interface{}) {
			originGroupsMap := item.(map[string]interface{})
			originGroupInLB := teo.OriginGroupInLoadBalancer{}
			if v, ok := originGroupsMap["priority"]; ok {
				originGroupInLB.Priority = helper.String(v.(string))
			}
			if v, ok := originGroupsMap["origin_group_id"]; ok {
				originGroupInLB.OriginGroupId = helper.String(v.(string))
			}
			request.OriginGroups = append(request.OriginGroups, &originGroupInLB)
		}
	}

	if v, ok := d.GetOk("health_checker"); ok {
		for _, item := range v.([]interface{}) {
			healthCheckerMap := item.(map[string]interface{})
			healthChecker := teo.HealthChecker{}
			if v, ok := healthCheckerMap["type"]; ok {
				healthChecker.Type = helper.String(v.(string))
			}
			if v, ok := healthCheckerMap["port"]; ok {
				healthChecker.Port = helper.IntUint64(v.(int))
			}
			if v, ok := healthCheckerMap["interval"]; ok {
				healthChecker.Interval = helper.IntUint64(v.(int))
			}
			if v, ok := healthCheckerMap["timeout"]; ok {
				healthChecker.Timeout = helper.IntUint64(v.(int))
			}
			if v, ok := healthCheckerMap["health_threshold"]; ok {
				healthChecker.HealthThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := healthCheckerMap["critical_threshold"]; ok {
				healthChecker.CriticalThreshold = helper.IntUint64(v.(int))
			}
			if v, ok := healthCheckerMap["path"]; ok {
				healthChecker.Path = helper.String(v.(string))
			}
			if v, ok := healthCheckerMap["method"]; ok {
				healthChecker.Method = helper.String(v.(string))
			}
			if v, ok := healthCheckerMap["expected_codes"]; ok {
				expectedCodesList := v.([]interface{})
				for _, code := range expectedCodesList {
					healthChecker.ExpectedCodes = append(healthChecker.ExpectedCodes, helper.String(code.(string)))
				}
			}
			if v, ok := healthCheckerMap["headers"]; ok {
				for _, item := range v.([]interface{}) {
					headersMap := item.(map[string]interface{})
					customizedHeader := teo.CustomizedHeader{}
					if v, ok := headersMap["key"]; ok {
						customizedHeader.Key = helper.String(v.(string))
					}
					if v, ok := headersMap["value"]; ok {
						customizedHeader.Value = helper.String(v.(string))
					}
					healthChecker.Headers = append(healthChecker.Headers, &customizedHeader)
				}
			}
			if v, ok := healthCheckerMap["follow_redirect"]; ok {
				healthChecker.FollowRedirect = helper.String(v.(string))
			}
			if v, ok := healthCheckerMap["send_context"]; ok {
				healthChecker.SendContext = helper.String(v.(string))
			}
			if v, ok := healthCheckerMap["recv_context"]; ok {
				healthChecker.RecvContext = helper.String(v.(string))
			}
			request.HealthChecker = &healthChecker
		}
	}

	if v, ok := d.GetOk("steering_policy"); ok {
		request.SteeringPolicy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("failover_policy"); ok {
		request.FailoverPolicy = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().CreateLoadBalancerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create teo load balancer failed, reason:%+v", logId, err)
		return err
	}

	if response.Response == nil {
		return fmt.Errorf("create teo load balancer response is nil")
	}

	instanceId = *response.Response.InstanceId
	if instanceId == "" {
		return fmt.Errorf("create teo load balancer response InstanceId is empty")
	}

	d.SetId(strings.Join([]string{zoneId, instanceId}, tccommon.FILED_SP))

	return resourceTencentCloudTeoLoadBalancerRead(d, meta)
}

func resourceTencentCloudTeoLoadBalancerRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_load_balancer.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	service := TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	instanceId := idSplit[1]

	_ = d.Set("zone_id", zoneId)

	respData, err := service.DescribeTeoLoadBalancerById(ctx, zoneId, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `teo_load_balancer` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

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

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	originGroupsList := make([]map[string]interface{}, 0)
	if respData.OriginGroupHealthStatus != nil {
		for _, ogHealthStatus := range respData.OriginGroupHealthStatus {
			ogMap := map[string]interface{}{}
			if ogHealthStatus.OriginGroupID != nil {
				ogMap["origin_group_id"] = ogHealthStatus.OriginGroupID
			}
			if ogHealthStatus.Priority != nil {
				ogMap["priority"] = ogHealthStatus.Priority
			}
			originGroupsList = append(originGroupsList, ogMap)
		}
		_ = d.Set("origin_groups", originGroupsList)
	}

	if respData.HealthChecker != nil {
		healthCheckerList := make([]map[string]interface{}, 0, 1)
		healthCheckerMap := map[string]interface{}{}

		if respData.HealthChecker.Type != nil {
			healthCheckerMap["type"] = respData.HealthChecker.Type
		}
		if respData.HealthChecker.Port != nil {
			healthCheckerMap["port"] = int(*respData.HealthChecker.Port)
		}
		if respData.HealthChecker.Interval != nil {
			healthCheckerMap["interval"] = int(*respData.HealthChecker.Interval)
		}
		if respData.HealthChecker.Timeout != nil {
			healthCheckerMap["timeout"] = int(*respData.HealthChecker.Timeout)
		}
		if respData.HealthChecker.HealthThreshold != nil {
			healthCheckerMap["health_threshold"] = int(*respData.HealthChecker.HealthThreshold)
		}
		if respData.HealthChecker.CriticalThreshold != nil {
			healthCheckerMap["critical_threshold"] = int(*respData.HealthChecker.CriticalThreshold)
		}
		if respData.HealthChecker.Path != nil {
			healthCheckerMap["path"] = respData.HealthChecker.Path
		}
		if respData.HealthChecker.Method != nil {
			healthCheckerMap["method"] = respData.HealthChecker.Method
		}
		if respData.HealthChecker.ExpectedCodes != nil {
			expectedCodes := make([]string, 0, len(respData.HealthChecker.ExpectedCodes))
			for _, code := range respData.HealthChecker.ExpectedCodes {
				if code != nil {
					expectedCodes = append(expectedCodes, *code)
				}
			}
			healthCheckerMap["expected_codes"] = expectedCodes
		}
		if respData.HealthChecker.Headers != nil {
			headersList := make([]map[string]interface{}, 0, len(respData.HealthChecker.Headers))
			for _, header := range respData.HealthChecker.Headers {
				headerMap := map[string]interface{}{}
				if header.Key != nil {
					headerMap["key"] = header.Key
				}
				if header.Value != nil {
					headerMap["value"] = header.Value
				}
				headersList = append(headersList, headerMap)
			}
			healthCheckerMap["headers"] = headersList
		}
		if respData.HealthChecker.FollowRedirect != nil {
			healthCheckerMap["follow_redirect"] = respData.HealthChecker.FollowRedirect
		}
		if respData.HealthChecker.SendContext != nil {
			healthCheckerMap["send_context"] = respData.HealthChecker.SendContext
		}
		if respData.HealthChecker.RecvContext != nil {
			healthCheckerMap["recv_context"] = respData.HealthChecker.RecvContext
		}

		healthCheckerList = append(healthCheckerList, healthCheckerMap)
		_ = d.Set("health_checker", healthCheckerList)
	}

	originGroupHealthStatusList := make([]map[string]interface{}, 0, len(respData.OriginGroupHealthStatus))
	if respData.OriginGroupHealthStatus != nil {
		for _, ogHealthStatus := range respData.OriginGroupHealthStatus {
			ogHealthStatusMap := map[string]interface{}{}

			if ogHealthStatus.OriginGroupID != nil {
				ogHealthStatusMap["origin_group_id"] = ogHealthStatus.OriginGroupID
			}
			if ogHealthStatus.OriginGroupName != nil {
				ogHealthStatusMap["origin_group_name"] = ogHealthStatus.OriginGroupName
			}
			if ogHealthStatus.OriginType != nil {
				ogHealthStatusMap["origin_type"] = ogHealthStatus.OriginType
			}
			if ogHealthStatus.Priority != nil {
				ogHealthStatusMap["priority"] = ogHealthStatus.Priority
			}

			originHealthStatusList := make([]map[string]interface{}, 0, len(ogHealthStatus.OriginHealthStatus))
			if ogHealthStatus.OriginHealthStatus != nil {
				for _, ohs := range ogHealthStatus.OriginHealthStatus {
					ohsMap := map[string]interface{}{}
					if ohs.Origin != nil {
						ohsMap["origin"] = ohs.Origin
					}
					if ohs.Healthy != nil {
						ohsMap["healthy"] = ohs.Healthy
					}
					originHealthStatusList = append(originHealthStatusList, ohsMap)
				}
				ogHealthStatusMap["origin_health_status"] = originHealthStatusList
			}

			originGroupHealthStatusList = append(originGroupHealthStatusList, ogHealthStatusMap)
		}

		_ = d.Set("origin_group_health_status", originGroupHealthStatusList)
	}

	if respData.L4UsedList != nil {
		l4UsedList := make([]string, 0, len(respData.L4UsedList))
		for _, v := range respData.L4UsedList {
			if v != nil {
				l4UsedList = append(l4UsedList, *v)
			}
		}
		_ = d.Set("l4_used_list", l4UsedList)
	}

	if respData.L7UsedList != nil {
		l7UsedList := make([]string, 0, len(respData.L7UsedList))
		for _, v := range respData.L7UsedList {
			if v != nil {
				l7UsedList = append(l7UsedList, *v)
			}
		}
		_ = d.Set("l7_used_list", l7UsedList)
	}

	referencesList := make([]map[string]interface{}, 0, len(respData.References))
	if respData.References != nil {
		for _, ref := range respData.References {
			refMap := map[string]interface{}{}

			if ref.InstanceType != nil {
				refMap["instance_type"] = ref.InstanceType
			}
			if ref.InstanceId != nil {
				refMap["instance_id"] = ref.InstanceId
			}
			if ref.InstanceName != nil {
				refMap["instance_name"] = ref.InstanceName
			}

			referencesList = append(referencesList, refMap)
		}

		_ = d.Set("references", referencesList)
	}

	_ = zoneId
	return nil
}

func resourceTencentCloudTeoLoadBalancerUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_load_balancer.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	instanceId := idSplit[1]

	needChange := false
	mutableArgs := []string{"name", "origin_groups", "health_checker", "steering_policy", "failover_policy"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := teo.NewModifyLoadBalancerRequest()

		request.ZoneId = helper.String(zoneId)
		request.InstanceId = helper.String(instanceId)

		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}

		if v, ok := d.GetOk("origin_groups"); ok {
			for _, item := range v.([]interface{}) {
				originGroupsMap := item.(map[string]interface{})
				originGroupInLB := teo.OriginGroupInLoadBalancer{}
				if v, ok := originGroupsMap["priority"]; ok {
					originGroupInLB.Priority = helper.String(v.(string))
				}
				if v, ok := originGroupsMap["origin_group_id"]; ok {
					originGroupInLB.OriginGroupId = helper.String(v.(string))
				}
				request.OriginGroups = append(request.OriginGroups, &originGroupInLB)
			}
		}

		if v, ok := d.GetOk("health_checker"); ok {
			for _, item := range v.([]interface{}) {
				healthCheckerMap := item.(map[string]interface{})
				healthChecker := teo.HealthChecker{}
				if v, ok := healthCheckerMap["type"]; ok {
					healthChecker.Type = helper.String(v.(string))
				}
				if v, ok := healthCheckerMap["port"]; ok {
					healthChecker.Port = helper.IntUint64(v.(int))
				}
				if v, ok := healthCheckerMap["interval"]; ok {
					healthChecker.Interval = helper.IntUint64(v.(int))
				}
				if v, ok := healthCheckerMap["timeout"]; ok {
					healthChecker.Timeout = helper.IntUint64(v.(int))
				}
				if v, ok := healthCheckerMap["health_threshold"]; ok {
					healthChecker.HealthThreshold = helper.IntUint64(v.(int))
				}
				if v, ok := healthCheckerMap["critical_threshold"]; ok {
					healthChecker.CriticalThreshold = helper.IntUint64(v.(int))
				}
				if v, ok := healthCheckerMap["path"]; ok {
					healthChecker.Path = helper.String(v.(string))
				}
				if v, ok := healthCheckerMap["method"]; ok {
					healthChecker.Method = helper.String(v.(string))
				}
				if v, ok := healthCheckerMap["expected_codes"]; ok {
					expectedCodesList := v.([]interface{})
					for _, code := range expectedCodesList {
						healthChecker.ExpectedCodes = append(healthChecker.ExpectedCodes, helper.String(code.(string)))
					}
				}
				if v, ok := healthCheckerMap["headers"]; ok {
					for _, item := range v.([]interface{}) {
						headersMap := item.(map[string]interface{})
						customizedHeader := teo.CustomizedHeader{}
						if v, ok := headersMap["key"]; ok {
							customizedHeader.Key = helper.String(v.(string))
						}
						if v, ok := headersMap["value"]; ok {
							customizedHeader.Value = helper.String(v.(string))
						}
						healthChecker.Headers = append(healthChecker.Headers, &customizedHeader)
					}
				}
				if v, ok := healthCheckerMap["follow_redirect"]; ok {
					healthChecker.FollowRedirect = helper.String(v.(string))
				}
				if v, ok := healthCheckerMap["send_context"]; ok {
					healthChecker.SendContext = helper.String(v.(string))
				}
				if v, ok := healthCheckerMap["recv_context"]; ok {
					healthChecker.RecvContext = helper.String(v.(string))
				}
				request.HealthChecker = &healthChecker
			}
		}

		if v, ok := d.GetOk("steering_policy"); ok {
			request.SteeringPolicy = helper.String(v.(string))
		}

		if v, ok := d.GetOk("failover_policy"); ok {
			request.FailoverPolicy = helper.String(v.(string))
		}

		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().ModifyLoadBalancerWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s update teo load balancer failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudTeoLoadBalancerRead(d, meta)
}

func resourceTencentCloudTeoLoadBalancerDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_load_balancer.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	zoneId := idSplit[0]
	instanceId := idSplit[1]

	var (
		request  = teo.NewDeleteLoadBalancerRequest()
		response = teo.NewDeleteLoadBalancerResponse()
	)

	request.ZoneId = helper.String(zoneId)
	request.InstanceId = helper.String(instanceId)

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoClient().DeleteLoadBalancerWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete teo load balancer failed, reason:%+v", logId, err)
		return err
	}

	_ = response
	return nil
}
