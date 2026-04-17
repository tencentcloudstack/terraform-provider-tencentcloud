package tke

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tke "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tke/v20180525"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func ResourceTencentCloudKubernetesClusterSchedulerPolicyConfig() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKubernetesClusterSchedulerPolicyConfigCreate,
		Read:   resourceTencentCloudKubernetesClusterSchedulerPolicyConfigRead,
		Update: resourceTencentCloudKubernetesClusterSchedulerPolicyConfigUpdate,
		Delete: resourceTencentCloudKubernetesClusterSchedulerPolicyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Cluster ID.",
			},

			"scheduler_policy_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Scheduler policy configuration list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheduler_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Scheduler name.",
						},
						"plugin_configs": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Scheduler plugin configuration list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "Plugin name.",
									},
									"args": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										DiffSuppressFunc: func(k, oldVal, newVal string, d *schema.ResourceData) bool {
											return suppressJSONWhitespaceDiff(oldVal, newVal)
										},
										Description: "Plugin args in raw JSON format. Terraform will automatically base64-encode it before calling the API and decode it on read.",
									},
								},
							},
						},
						"plugin_set": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Plugin set configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeSet,
										Optional:    true,
										Computed:    true,
										Description: "List of plugins to enable.",
										Set: func(v interface{}) int {
											m := v.(map[string]interface{})
											return schema.HashString(m["name"].(string))
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Plugin name.",
												},
												"weight": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Plugin weight.",
												},
											},
										},
									},
									"disabled": {
										Type:        schema.TypeSet,
										Optional:    true,
										Computed:    true,
										Description: "List of plugins to disable.",
										Set: func(v interface{}) int {
											m := v.(map[string]interface{})
											return schema.HashString(m["name"].(string))
										},
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"name": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "Plugin name.",
												},
												"weight": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Plugin weight.",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"extenders": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "Extender scheduler configuration list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"filter_verb": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Filter stage interface.",
						},
						"prioritize_verb": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Prioritize stage interface.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Weight for prioritize stage.",
						},
						"preempt_verb": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Preempt stage interface.",
						},
						"node_cache_capable": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Whether node cache capability is enabled.",
						},
						"extender_client_config": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
							Description: "Extender client configuration.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service": {
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										MaxItems:    1,
										Description: "Service reference configuration.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"namespace": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Service namespace.",
												},
												"name": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Service name.",
												},
												"port": {
													Type:        schema.TypeInt,
													Optional:    true,
													Computed:    true,
													Description: "Service port.",
												},
												"path": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Service path.",
												},
												"scheme": {
													Type:        schema.TypeString,
													Optional:    true,
													Computed:    true,
													Description: "Service protocol scheme (e.g. http, https).",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},

			"client_connection": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Description: "Client connection configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"qps": {
							Type:        schema.TypeFloat,
							Optional:    true,
							Computed:    true,
							Description: "Maximum queries per second.",
						},
						"burst": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Burst request limit.",
						},
					},
				},
			},

			"high_performance": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "High performance mode switch.",
			},

			// Computed
			"policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Raw scheduler policy JSON string.",
			},
		},
	}
}

func resourceTencentCloudKubernetesClusterSchedulerPolicyConfigCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_scheduler_policy_config.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	clusterId := d.Get("cluster_id").(string)
	d.SetId(clusterId)

	return resourceTencentCloudKubernetesClusterSchedulerPolicyConfigUpdate(d, meta)
}

func resourceTencentCloudKubernetesClusterSchedulerPolicyConfigRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_scheduler_policy_config.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TkeService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	clusterId := d.Id()
	_ = d.Set("cluster_id", clusterId)

	respData, err := service.DescribeKubernetesClusterSchedulerPolicy(ctx, clusterId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_kubernetes_cluster_scheduler_policy_config` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.Policy != nil {
		_ = d.Set("policy", respData.Policy)
	}

	if respData.HighPerformance != nil {
		_ = d.Set("high_performance", respData.HighPerformance)
	}

	if respData.ClientConnection != nil {
		cc := map[string]interface{}{}
		if respData.ClientConnection.QPS != nil {
			cc["qps"] = *respData.ClientConnection.QPS
		}

		if respData.ClientConnection.Burst != nil {
			cc["burst"] = int(*respData.ClientConnection.Burst)
		}

		_ = d.Set("client_connection", []map[string]interface{}{cc})
	}

	if respData.SchedulerPolicyConfig != nil {
		_ = d.Set("scheduler_policy_config", flattenSchedulerPolicyConfigList(respData.SchedulerPolicyConfig))
	}

	if respData.Extenders != nil {
		_ = d.Set("extenders", flattenExtendersList(respData.Extenders))
	}

	return nil
}

func resourceTencentCloudKubernetesClusterSchedulerPolicyConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_scheduler_policy_config.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = tke.NewModifyClusterSchedulerPolicyRequest()
	)

	clusterId := d.Id()
	request.ClusterId = &clusterId

	if v, ok := d.GetOk("scheduler_policy_config"); ok {
		request.SchedulerPolicyConfig = buildSchedulerPolicyConfigList(v.([]interface{}))
	}

	if v, ok := d.GetOk("extenders"); ok {
		request.Extenders = buildExtendersList(v.([]interface{}))
	}

	if v, ok := d.GetOk("client_connection"); ok {
		ccList := v.([]interface{})
		if len(ccList) > 0 {
			ccMap := ccList[0].(map[string]interface{})
			cc := &tke.ClientConnection{}
			if qps, ok := ccMap["qps"].(float64); ok {
				cc.QPS = &qps
			}

			if burst, ok := ccMap["burst"].(int); ok {
				cc.Burst = helper.IntUint64(burst)
			}

			request.ClientConnection = cc
		}
	}

	if v, ok := d.GetOkExists("high_performance"); ok {
		request.HighPerformance = helper.Bool(v.(bool))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().ModifyClusterSchedulerPolicyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s update kubernetes cluster scheduler policy config failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	// Poll DescribeTasks until LifeState == "done"
	if err := waitForClusterSchedulerPolicyTaskDone(ctx, meta, clusterId, logId); err != nil {
		return err
	}

	return resourceTencentCloudKubernetesClusterSchedulerPolicyConfigRead(d, meta)
}

func resourceTencentCloudKubernetesClusterSchedulerPolicyConfigDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_kubernetes_cluster_scheduler_policy_config.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return nil
}

func waitForClusterSchedulerPolicyTaskDone(ctx context.Context, meta interface{}, clusterId, logId string) error {
	request := tke.NewDescribeTasksRequest()
	request.Latest = helper.Bool(true)
	request.Filter = []*tke.Filter{
		{
			Name:   helper.String("TaskType"),
			Values: []*string{helper.String("scheduler_policy")},
		},
		{
			Name:   helper.String("ClusterId"),
			Values: []*string{&clusterId},
		},
	}

	return resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTkeV20180525Client().DescribeTasksWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || result.Response == nil {
			return resource.RetryableError(fmt.Errorf("waiting for cluster scheduler policy task, response is nil"))
		}

		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())

		if len(result.Response.Tasks) == 0 {
			return resource.RetryableError(fmt.Errorf("waiting for cluster scheduler policy task, no tasks found"))
		}

		task := result.Response.Tasks[0]
		if task.LifeState == nil {
			return resource.RetryableError(fmt.Errorf("waiting for cluster scheduler policy task, LifeState is nil"))
		}

		lifeState := *task.LifeState
		if lifeState == "done" {
			return nil
		}

		if lifeState == "abort" || lifeState == "aborted" {
			lastErr := ""
			if task.LastError != nil {
				lastErr = *task.LastError
			}

			return resource.NonRetryableError(fmt.Errorf("cluster scheduler policy task failed with state: %s, error: %s", lifeState, lastErr))
		}

		return resource.RetryableError(fmt.Errorf("waiting for cluster scheduler policy task, current LifeState: %s", lifeState))
	})
}

func buildSchedulerPolicyConfigList(rawList []interface{}) []*tke.SchedulerPolicyConfig {
	configs := make([]*tke.SchedulerPolicyConfig, 0, len(rawList))
	for _, item := range rawList {
		m := item.(map[string]interface{})
		cfg := &tke.SchedulerPolicyConfig{}

		if v, ok := m["scheduler_name"].(string); ok && v != "" {
			cfg.SchedulerName = helper.String(v)
		}

		if v, ok := m["plugin_configs"]; ok {
			pcList := v.([]interface{})
			for _, pcItem := range pcList {
				pcMap := pcItem.(map[string]interface{})
				pc := &tke.SchedulerPluginConfigs{}
				if name, ok := pcMap["name"].(string); ok && name != "" {
					pc.Name = helper.String(name)
				}

				if args, ok := pcMap["args"].(string); ok && args != "" {
					pc.Args = helper.String(base64.StdEncoding.EncodeToString([]byte(args)))
				}

				cfg.PluginConfigs = append(cfg.PluginConfigs, pc)
			}
		}

		if v, ok := m["plugin_set"]; ok {
			psList := v.([]interface{})
			if len(psList) > 0 {
				psMap := psList[0].(map[string]interface{})
				ps := &tke.PluginSet{}

				if enabled, ok := psMap["enabled"]; ok {
					for _, e := range enabled.(*schema.Set).List() {
						eMap := e.(map[string]interface{})
						p := &tke.SchedulerPolicyPriority{}
						if name, ok := eMap["name"].(string); ok && name != "" {
							p.Name = helper.String(name)
						}

						if weight, ok := eMap["weight"].(int); ok {
							p.Weight = helper.IntInt64(weight)
						}

						ps.Enabled = append(ps.Enabled, p)
					}
				}

				if disabled, ok := psMap["disabled"]; ok {
					for _, d := range disabled.(*schema.Set).List() {
						dMap := d.(map[string]interface{})
						p := &tke.SchedulerPolicyPriority{}
						if name, ok := dMap["name"].(string); ok && name != "" {
							p.Name = helper.String(name)
						}

						if weight, ok := dMap["weight"].(int); ok {
							p.Weight = helper.IntInt64(weight)
						}

						ps.Disabled = append(ps.Disabled, p)
					}
				}

				cfg.PluginSet = ps
			}
		}

		configs = append(configs, cfg)
	}

	return configs
}

func buildExtendersList(rawList []interface{}) []*tke.Extenders {
	extenders := make([]*tke.Extenders, 0, len(rawList))
	for _, item := range rawList {
		m := item.(map[string]interface{})
		ext := &tke.Extenders{}

		if v, ok := m["filter_verb"].(string); ok && v != "" {
			ext.FilterVerb = helper.String(v)
		}

		if v, ok := m["prioritize_verb"].(string); ok && v != "" {
			ext.PrioritizeVerb = helper.String(v)
		}

		if v, ok := m["weight"].(int); ok {
			ext.Weight = helper.IntInt64(v)
		}

		if v, ok := m["preempt_verb"].(string); ok && v != "" {
			ext.PreemptVerb = helper.String(v)
		}

		if v, ok := m["node_cache_capable"].(bool); ok {
			ext.NodeCacheCapable = helper.Bool(v)
		}

		if v, ok := m["extender_client_config"]; ok {
			eccList := v.([]interface{})
			if len(eccList) > 0 {
				eccMap := eccList[0].(map[string]interface{})
				ecc := &tke.ExtenderClientConfig{}

				if svc, ok := eccMap["service"]; ok {
					svcList := svc.([]interface{})
					if len(svcList) > 0 {
						svcMap := svcList[0].(map[string]interface{})
						sr := &tke.ServiceReference{}
						if ns, ok := svcMap["namespace"].(string); ok && ns != "" {
							sr.Namespace = helper.String(ns)
						}

						if name, ok := svcMap["name"].(string); ok && name != "" {
							sr.Name = helper.String(name)
						}

						if port, ok := svcMap["port"].(int); ok {
							sr.Port = helper.IntInt64(port)
						}

						if path, ok := svcMap["path"].(string); ok && path != "" {
							sr.Path = helper.String(path)
						}

						if scheme, ok := svcMap["scheme"].(string); ok && scheme != "" {
							sr.Scheme = helper.String(scheme)
						}

						ecc.Service = sr
					}
				}

				ext.ExtenderClientConfig = ecc
			}
		}

		extenders = append(extenders, ext)
	}

	return extenders
}

func flattenSchedulerPolicyConfigList(items []*tke.SchedulerPolicyConfig) []map[string]interface{} {
	list := make([]map[string]interface{}, 0, len(items))
	for _, cfg := range items {
		m := map[string]interface{}{}

		if cfg.SchedulerName != nil {
			m["scheduler_name"] = cfg.SchedulerName
		}

		if cfg.PluginConfigs != nil {
			pcList := make([]map[string]interface{}, 0, len(cfg.PluginConfigs))
			for _, pc := range cfg.PluginConfigs {
				pcMap := map[string]interface{}{}
				if pc.Name != nil {
					pcMap["name"] = pc.Name
				}

				if pc.Args != nil {
					if decoded, err := base64.StdEncoding.DecodeString(*pc.Args); err == nil {
						pcMap["args"] = string(decoded)
					} else {
						pcMap["args"] = *pc.Args
					}
				}

				pcList = append(pcList, pcMap)
			}

			m["plugin_configs"] = pcList
		}

		if cfg.PluginSet != nil {
			psMap := map[string]interface{}{}
			if cfg.PluginSet.Enabled != nil {
				enabled := make([]map[string]interface{}, 0, len(cfg.PluginSet.Enabled))
				for _, e := range cfg.PluginSet.Enabled {
					eMap := map[string]interface{}{}
					if e.Name != nil {
						eMap["name"] = e.Name
					}

					if e.Weight != nil {
						eMap["weight"] = int(*e.Weight)
					}

					enabled = append(enabled, eMap)
				}

				psMap["enabled"] = enabled
			}

			if cfg.PluginSet.Disabled != nil {
				disabled := make([]map[string]interface{}, 0, len(cfg.PluginSet.Disabled))
				for _, d := range cfg.PluginSet.Disabled {
					dMap := map[string]interface{}{}
					if d.Name != nil {
						dMap["name"] = d.Name
					}

					if d.Weight != nil {
						dMap["weight"] = int(*d.Weight)
					}

					disabled = append(disabled, dMap)
				}

				psMap["disabled"] = disabled
			}

			m["plugin_set"] = []map[string]interface{}{psMap}
		}

		list = append(list, m)
	}

	return list
}

func flattenExtendersList(items []*tke.Extenders) []map[string]interface{} {
	list := make([]map[string]interface{}, 0, len(items))
	for _, ext := range items {
		m := map[string]interface{}{}

		if ext.FilterVerb != nil {
			m["filter_verb"] = ext.FilterVerb
		}

		if ext.PrioritizeVerb != nil {
			m["prioritize_verb"] = ext.PrioritizeVerb
		}

		if ext.Weight != nil {
			m["weight"] = int(*ext.Weight)
		}

		if ext.PreemptVerb != nil {
			m["preempt_verb"] = ext.PreemptVerb
		}

		if ext.NodeCacheCapable != nil {
			m["node_cache_capable"] = ext.NodeCacheCapable
		}

		if ext.ExtenderClientConfig != nil {
			eccMap := map[string]interface{}{}
			if ext.ExtenderClientConfig.Service != nil {
				svc := ext.ExtenderClientConfig.Service
				svcMap := map[string]interface{}{}
				if svc.Namespace != nil {
					svcMap["namespace"] = svc.Namespace
				}

				if svc.Name != nil {
					svcMap["name"] = svc.Name
				}

				if svc.Port != nil {
					svcMap["port"] = int(*svc.Port)
				}

				if svc.Path != nil {
					svcMap["path"] = svc.Path
				}

				if svc.Scheme != nil {
					svcMap["scheme"] = svc.Scheme
				}

				eccMap["service"] = []map[string]interface{}{svcMap}
			}

			m["extender_client_config"] = []map[string]interface{}{eccMap}
		}

		list = append(list, m)
	}

	return list
}

// suppressJSONWhitespaceDiff returns true if oldVal and newVal represent
// the same JSON content (ignoring whitespace and key ordering differences).
func suppressJSONWhitespaceDiff(oldVal, newVal string) bool {
	if oldVal == newVal {
		return true
	}
	normalize := func(s string) string {
		var v interface{}
		if err := json.Unmarshal([]byte(s), &v); err != nil {
			return s
		}
		b, err := json.Marshal(v)
		if err != nil {
			return s
		}
		return string(b)
	}
	return normalize(oldVal) == normalize(newVal)
}
