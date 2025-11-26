package igtm

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudIgtmStrategy() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIgtmStrategyCreate,
		Read:   resourceTencentCloudIgtmStrategyRead,
		Update: resourceTencentCloudIgtmStrategyUpdate,
		Delete: resourceTencentCloudIgtmStrategyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance ID.",
			},

			"strategy_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Strategy name, cannot be duplicated.",
			},

			"source": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Resolution lines.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dns_line_id": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Resolution request source line ID.",
						},
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Resolution request source line name.",
						},
					},
				},
			},

			"main_address_pool_set": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Main address pool set, up to four levels allowed.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_pools": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Address pool IDs and weights in the set, array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pool_id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Address pool ID.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Weight.",
									},
								},
							},
						},
						"main_address_pool_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Address pool set ID.",
						},
						"min_survive_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Switch threshold, cannot exceed the total number of addresses in the main set.",
						},
						"traffic_strategy": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Switch strategy: ALL resolves all addresses; WEIGHT: load balancing. When ALL, the weight value of resolved addresses is 1; when WEIGHT, weight is address pool weight * address weight.",
						},
					},
				},
			},

			"fallback_address_pool_set": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Fallback address pool set, only one level allowed and address pool count must be 1.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_pools": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "Address pool IDs and weights in the set, array.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"pool_id": {
										Type:        schema.TypeInt,
										Required:    true,
										Description: "Address pool ID.",
									},
									"weight": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Weight.",
									},
								},
							},
						},
						"main_address_pool_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Address pool set ID.",
						},
						"min_survive_num": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "Switch threshold, cannot exceed the total number of addresses in the main set.",
						},
						"traffic_strategy": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Switch strategy: ALL resolves all addresses; WEIGHT: load balancing. When ALL, the weight value of resolved addresses is 1; when WEIGHT, weight is address pool weight * address weight.",
						},
					},
				},
			},

			"keep_domain_records": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Whether to enable policy forced retention of default lines disabled, enabled, default is disabled and only one policy can be enabled.",
			},

			"switch_pool_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Policy scheduling mode: AUTO default switching; STOP only pause without switching.",
			},

			// computed
			"strategy_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Strategy ID.",
			},
		},
	}
}

func resourceTencentCloudIgtmStrategyCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_strategy.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = igtmv20231024.NewCreateStrategyRequest()
		response   = igtmv20231024.NewCreateStrategyResponse()
		instanceId string
		strategyId string
	)

	if v, ok := d.GetOk("instance_id"); ok {
		request.InstanceId = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("strategy_name"); ok {
		request.StrategyName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("source"); ok {
		for _, item := range v.([]interface{}) {
			sourceMap := item.(map[string]interface{})
			source := igtmv20231024.Source{}
			if v, ok := sourceMap["dns_line_id"].(int); ok {
				source.DnsLineId = helper.IntUint64(v)
			}

			if v, ok := sourceMap["name"].(string); ok && v != "" {
				source.Name = helper.String(v)
			}

			request.Source = append(request.Source, &source)
		}
	}

	if v, ok := d.GetOk("main_address_pool_set"); ok {
		for _, item := range v.([]interface{}) {
			mainAddressPoolSetMap := item.(map[string]interface{})
			mainAddressPool := igtmv20231024.MainAddressPool{}
			if v, ok := mainAddressPoolSetMap["address_pools"]; ok {
				for _, item := range v.([]interface{}) {
					addressPoolsMap := item.(map[string]interface{})
					mainPoolWeight := igtmv20231024.MainPoolWeight{}
					if v, ok := addressPoolsMap["pool_id"].(int); ok {
						mainPoolWeight.PoolId = helper.IntUint64(v)
					}

					if v, ok := addressPoolsMap["weight"].(int); ok && v != 0 {
						mainPoolWeight.Weight = helper.IntUint64(v)
					}

					mainAddressPool.AddressPools = append(mainAddressPool.AddressPools, &mainPoolWeight)
				}
			}

			if v, ok := mainAddressPoolSetMap["main_address_pool_id"].(int); ok && v != 0 {
				mainAddressPool.MainAddressPoolId = helper.IntUint64(v)
			}

			if v, ok := mainAddressPoolSetMap["min_survive_num"].(int); ok && v != 0 {
				mainAddressPool.MinSurviveNum = helper.IntUint64(v)
			}

			if v, ok := mainAddressPoolSetMap["traffic_strategy"].(string); ok && v != "" {
				mainAddressPool.TrafficStrategy = helper.String(v)
			}

			request.MainAddressPoolSet = append(request.MainAddressPoolSet, &mainAddressPool)
		}
	}

	if v, ok := d.GetOk("fallback_address_pool_set"); ok {
		for _, item := range v.([]interface{}) {
			fallbackAddressPoolSetMap := item.(map[string]interface{})
			mainAddressPool := igtmv20231024.MainAddressPool{}
			if v, ok := fallbackAddressPoolSetMap["address_pools"]; ok {
				for _, item := range v.([]interface{}) {
					addressPoolsMap := item.(map[string]interface{})
					mainPoolWeight := igtmv20231024.MainPoolWeight{}
					if v, ok := addressPoolsMap["pool_id"].(int); ok {
						mainPoolWeight.PoolId = helper.IntUint64(v)
					}

					if v, ok := addressPoolsMap["weight"].(int); ok && v != 0 {
						mainPoolWeight.Weight = helper.IntUint64(v)
					}

					mainAddressPool.AddressPools = append(mainAddressPool.AddressPools, &mainPoolWeight)
				}
			}

			if v, ok := fallbackAddressPoolSetMap["main_address_pool_id"].(int); ok && v != 0 {
				mainAddressPool.MainAddressPoolId = helper.IntUint64(v)
			}

			if v, ok := fallbackAddressPoolSetMap["min_survive_num"].(int); ok && v != 0 {
				mainAddressPool.MinSurviveNum = helper.IntUint64(v)
			}

			if v, ok := fallbackAddressPoolSetMap["traffic_strategy"].(string); ok && v != "" {
				mainAddressPool.TrafficStrategy = helper.String(v)
			}

			request.FallbackAddressPoolSet = append(request.FallbackAddressPoolSet, &mainAddressPool)
		}
	}

	if v, ok := d.GetOk("keep_domain_records"); ok {
		request.KeepDomainRecords = helper.String(v.(string))
	}

	if v, ok := d.GetOk("switch_pool_type"); ok {
		request.SwitchPoolType = helper.String(v.(string))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().CreateStrategyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create igtm strategy failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create igtm strategy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.StrategyId == nil {
		return fmt.Errorf("StrategyId is nil.")
	}

	strategyId = helper.Int64ToStr(*response.Response.StrategyId)
	d.SetId(strings.Join([]string{instanceId, strategyId}, tccommon.FILED_SP))
	return resourceTencentCloudIgtmStrategyRead(d, meta)
}

func resourceTencentCloudIgtmStrategyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_strategy.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	strategyId := idSplit[1]

	respData, err := service.DescribeIgtmStrategyById(ctx, instanceId, strategyId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_igtm_strategy` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	if respData.Name != nil {
		_ = d.Set("strategy_name", respData.Name)
	}

	if respData.Source != nil && len(respData.Source) > 0 {
		sourceList := make([]map[string]interface{}, 0, len(respData.Source))
		for _, source := range respData.Source {
			sourceMap := map[string]interface{}{}
			if source.DnsLineId != nil {
				sourceMap["dns_line_id"] = source.DnsLineId
			}

			if source.Name != nil {
				sourceMap["name"] = source.Name
			}

			sourceList = append(sourceList, sourceMap)
		}

		_ = d.Set("source", sourceList)
	}

	if respData.MainAddressPoolSet != nil && len(respData.MainAddressPoolSet) > 0 {
		mainAddressPoolSetList := make([]map[string]interface{}, 0, len(respData.MainAddressPoolSet))
		for _, mainAddressPoolSet := range respData.MainAddressPoolSet {
			mainAddressPoolSetMap := map[string]interface{}{}
			addressPoolsList := make([]map[string]interface{}, 0, len(mainAddressPoolSet.AddressPools))
			if mainAddressPoolSet.AddressPools != nil {
				for _, addressPools := range mainAddressPoolSet.AddressPools {
					addressPoolsMap := map[string]interface{}{}
					if addressPools.PoolId != nil {
						addressPoolsMap["pool_id"] = addressPools.PoolId
					}

					if addressPools.Weight != nil {
						addressPoolsMap["weight"] = addressPools.Weight
					}

					addressPoolsList = append(addressPoolsList, addressPoolsMap)
				}

				mainAddressPoolSetMap["address_pools"] = addressPoolsList
			}

			if mainAddressPoolSet.MainAddressPoolId != nil {
				mainAddressPoolSetMap["main_address_pool_id"] = mainAddressPoolSet.MainAddressPoolId
			}

			if mainAddressPoolSet.MinSurviveNum != nil {
				mainAddressPoolSetMap["min_survive_num"] = mainAddressPoolSet.MinSurviveNum
			}

			if mainAddressPoolSet.TrafficStrategy != nil {
				mainAddressPoolSetMap["traffic_strategy"] = mainAddressPoolSet.TrafficStrategy
			}

			mainAddressPoolSetList = append(mainAddressPoolSetList, mainAddressPoolSetMap)
		}

		_ = d.Set("main_address_pool_set", mainAddressPoolSetList)
	}

	if respData.FallbackAddressPoolSet != nil && len(respData.FallbackAddressPoolSet) > 0 {
		fallbackAddressPoolSetList := make([]map[string]interface{}, 0, len(respData.FallbackAddressPoolSet))
		for _, fallbackAddressPoolSet := range respData.FallbackAddressPoolSet {
			fallbackAddressPoolSetMap := map[string]interface{}{}
			addressPoolsList := make([]map[string]interface{}, 0, len(fallbackAddressPoolSet.AddressPools))
			if fallbackAddressPoolSet.AddressPools != nil {
				for _, addressPools := range fallbackAddressPoolSet.AddressPools {
					addressPoolsMap := map[string]interface{}{}
					if addressPools.PoolId != nil {
						addressPoolsMap["pool_id"] = addressPools.PoolId
					}

					if addressPools.Weight != nil {
						addressPoolsMap["weight"] = addressPools.Weight
					}

					addressPoolsList = append(addressPoolsList, addressPoolsMap)
				}

				fallbackAddressPoolSetMap["address_pools"] = addressPoolsList
			}

			if fallbackAddressPoolSet.MainAddressPoolId != nil {
				fallbackAddressPoolSetMap["main_address_pool_id"] = fallbackAddressPoolSet.MainAddressPoolId
			}

			if fallbackAddressPoolSet.MinSurviveNum != nil {
				fallbackAddressPoolSetMap["min_survive_num"] = fallbackAddressPoolSet.MinSurviveNum
			}

			if fallbackAddressPoolSet.TrafficStrategy != nil {
				fallbackAddressPoolSetMap["traffic_strategy"] = fallbackAddressPoolSet.TrafficStrategy
			}

			fallbackAddressPoolSetList = append(fallbackAddressPoolSetList, fallbackAddressPoolSetMap)
		}

		_ = d.Set("fallback_address_pool_set", fallbackAddressPoolSetList)
	}

	if respData.KeepDomainRecords != nil {
		_ = d.Set("keep_domain_records", respData.KeepDomainRecords)
	}

	if respData.SwitchPoolType != nil {
		_ = d.Set("switch_pool_type", respData.SwitchPoolType)
	}

	if respData.StrategyId != nil {
		_ = d.Set("strategy_id", respData.StrategyId)
	}

	return nil
}

func resourceTencentCloudIgtmStrategyUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_strategy.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId = tccommon.GetLogId(tccommon.ContextNil)
		ctx   = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	strategyId := idSplit[1]

	needChange := false
	mutableArgs := []string{"source", "main_address_pool_set", "fallback_address_pool_set", "strategy_name", "is_enabled", "keep_domain_records", "switch_pool_type"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := igtmv20231024.NewModifyStrategyRequest()
		if v, ok := d.GetOk("source"); ok {
			for _, item := range v.([]interface{}) {
				sourceMap := item.(map[string]interface{})
				source := igtmv20231024.Source{}
				if v, ok := sourceMap["dns_line_id"].(int); ok {
					source.DnsLineId = helper.IntUint64(v)
				}

				if v, ok := sourceMap["name"].(string); ok && v != "" {
					source.Name = helper.String(v)
				}
				request.Source = append(request.Source, &source)
			}
		}

		if v, ok := d.GetOk("main_address_pool_set"); ok {
			for _, item := range v.([]interface{}) {
				mainAddressPoolSetMap := item.(map[string]interface{})
				mainAddressPool := igtmv20231024.MainAddressPool{}
				if v, ok := mainAddressPoolSetMap["address_pools"]; ok {
					for _, item := range v.([]interface{}) {
						addressPoolsMap := item.(map[string]interface{})
						mainPoolWeight := igtmv20231024.MainPoolWeight{}
						if v, ok := addressPoolsMap["pool_id"].(int); ok {
							mainPoolWeight.PoolId = helper.IntUint64(v)
						}

						if v, ok := addressPoolsMap["weight"].(int); ok && v != 0 {
							mainPoolWeight.Weight = helper.IntUint64(v)
						}

						mainAddressPool.AddressPools = append(mainAddressPool.AddressPools, &mainPoolWeight)
					}
				}

				if v, ok := mainAddressPoolSetMap["main_address_pool_id"].(int); ok && v != 0 {
					mainAddressPool.MainAddressPoolId = helper.IntUint64(v)
				}

				if v, ok := mainAddressPoolSetMap["min_survive_num"].(int); ok && v != 0 {
					mainAddressPool.MinSurviveNum = helper.IntUint64(v)
				}

				if v, ok := mainAddressPoolSetMap["traffic_strategy"].(string); ok && v != "" {
					mainAddressPool.TrafficStrategy = helper.String(v)
				}

				request.MainAddressPoolSet = append(request.MainAddressPoolSet, &mainAddressPool)
			}
		}

		if v, ok := d.GetOk("fallback_address_pool_set"); ok {
			for _, item := range v.([]interface{}) {
				fallbackAddressPoolSetMap := item.(map[string]interface{})
				mainAddressPool := igtmv20231024.MainAddressPool{}
				if v, ok := fallbackAddressPoolSetMap["address_pools"]; ok {
					for _, item := range v.([]interface{}) {
						addressPoolsMap := item.(map[string]interface{})
						mainPoolWeight := igtmv20231024.MainPoolWeight{}
						if v, ok := addressPoolsMap["pool_id"].(int); ok {
							mainPoolWeight.PoolId = helper.IntUint64(v)
						}

						if v, ok := addressPoolsMap["weight"].(int); ok && v != 0 {
							mainPoolWeight.Weight = helper.IntUint64(v)
						}

						mainAddressPool.AddressPools = append(mainAddressPool.AddressPools, &mainPoolWeight)
					}
				}

				if v, ok := fallbackAddressPoolSetMap["main_address_pool_id"].(int); ok && v != 0 {
					mainAddressPool.MainAddressPoolId = helper.IntUint64(v)
				}

				if v, ok := fallbackAddressPoolSetMap["min_survive_num"].(int); ok && v != 0 {
					mainAddressPool.MinSurviveNum = helper.IntUint64(v)
				}

				if v, ok := fallbackAddressPoolSetMap["traffic_strategy"].(string); ok && v != "" {
					mainAddressPool.TrafficStrategy = helper.String(v)
				}

				request.FallbackAddressPoolSet = append(request.FallbackAddressPoolSet, &mainAddressPool)
			}
		}

		if v, ok := d.GetOk("strategy_name"); ok {
			request.StrategyName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("is_enabled"); ok {
			request.IsEnabled = helper.String(v.(string))
		}

		if v, ok := d.GetOk("keep_domain_records"); ok {
			request.KeepDomainRecords = helper.String(v.(string))
		}

		if v, ok := d.GetOk("switch_pool_type"); ok {
			request.SwitchPoolType = helper.String(v.(string))
		}

		request.InstanceId = &instanceId
		request.StrategyId = helper.StrToUint64Point(strategyId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().ModifyStrategyWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil || result.Response.Msg == nil {
				return resource.NonRetryableError(fmt.Errorf("Modify igtm strategy failed, Response is nil."))
			}

			if *result.Response.Msg == "success" {
				return nil
			} else {
				return resource.NonRetryableError(fmt.Errorf("Modify igtm strategy failed, Msg is %s.", *result.Response.Msg))
			}
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update igtm strategy failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudIgtmStrategyRead(d, meta)
}

func resourceTencentCloudIgtmStrategyDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_strategy.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = igtmv20231024.NewDeleteStrategyRequest()
	)

	idSplit := strings.Split(d.Id(), tccommon.FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}

	instanceId := idSplit[0]
	strategyId := idSplit[1]

	request.InstanceId = &instanceId
	request.StrategyId = helper.StrToUint64Point(strategyId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().DeleteStrategyWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Msg == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete igtm strategy failed, Response is nil."))
		}

		if *result.Response.Msg == "success" {
			return nil
		} else {
			return resource.NonRetryableError(fmt.Errorf("Delete igtm strategy failed, Msg is %s.", *result.Response.Msg))
		}
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete igtm strategy failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
