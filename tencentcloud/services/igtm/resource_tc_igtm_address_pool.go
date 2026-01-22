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

func ResourceTencentCloudIgtmAddressPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudIgtmAddressPoolCreate,
		Read:   resourceTencentCloudIgtmAddressPoolRead,
		Update: resourceTencentCloudIgtmAddressPoolUpdate,
		Delete: resourceTencentCloudIgtmAddressPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Address pool name, duplicates are not allowed.",
			},

			"traffic_strategy": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Traffic strategy: WEIGHT for load balancing, ALL for resolving all healthy addresses.",
			},

			"address_set": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "Address list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"addr": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Address value: only supports IPv4, IPv6, and domain name formats.\nLoopback addresses, reserved addresses, internal addresses, and Tencent reserved network segments are not supported.",
						},
						"is_enable": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Whether to enable: DISABLED for disabled, ENABLED for enabled.",
						},
						"address_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Address ID.",
						},
						"location": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Address name.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OK for normal, DOWN for failure, WARN for risk, UNKNOWN for probing, UNMONITORED for unknown.",
						},
						"weight": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Weight, required when traffic strategy is WEIGHT; range 1-100.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Modification time.",
						},
					},
				},
			},

			"monitor_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Monitor ID.",
			},

			// computed
			"pool_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Address pool ID.",
			},
		},
	}
}

func resourceTencentCloudIgtmAddressPoolCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_address_pool.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = igtmv20231024.NewCreateAddressPoolRequest()
		response = igtmv20231024.NewCreateAddressPoolResponse()
		poolId   string
	)

	if v, ok := d.GetOk("pool_name"); ok {
		request.PoolName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("traffic_strategy"); ok {
		request.TrafficStrategy = helper.String(v.(string))
	}

	if v, ok := d.GetOk("address_set"); ok {
		for _, item := range v.(*schema.Set).List() {
			addressSetMap := item.(map[string]interface{})
			address := igtmv20231024.Address{}
			if v, ok := addressSetMap["addr"].(string); ok && v != "" {
				address.Addr = helper.String(v)
			}

			if v, ok := addressSetMap["is_enable"].(string); ok && v != "" {
				address.IsEnable = helper.String(v)
			}

			if v, ok := addressSetMap["location"].(string); ok && v != "" {
				address.Location = helper.String(v)
			}

			if v, ok := addressSetMap["weight"].(int); ok {
				address.Weight = helper.IntUint64(v)
			}

			request.AddressSet = append(request.AddressSet, &address)
		}
	}

	if v, ok := d.GetOkExists("monitor_id"); ok {
		request.MonitorId = helper.IntUint64(v.(int))
	}

	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().CreateAddressPoolWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create igtm address pool failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create igtm address pool failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if response.Response.AddressPoolId == nil {
		return fmt.Errorf("AddressPoolId is nil.")
	}

	poolId = helper.UInt64ToStr(*response.Response.AddressPoolId)
	d.SetId(poolId)
	return resourceTencentCloudIgtmAddressPoolRead(d, meta)
}

func resourceTencentCloudIgtmAddressPoolRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_address_pool.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		poolId  = d.Id()
	)

	respData, err := service.DescribeIgtmAddressPoolById(ctx, poolId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_igtm_address_pool` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.AddressPool != nil {
		if respData.AddressPool.PoolName != nil {
			_ = d.Set("pool_name", respData.AddressPool.PoolName)
		}

		if respData.AddressPool.TrafficStrategy != nil {
			_ = d.Set("traffic_strategy", respData.AddressPool.TrafficStrategy)
		}

		if respData.AddressPool.MonitorId != nil {
			_ = d.Set("monitor_id", respData.AddressPool.MonitorId)
		}

		if respData.AddressPool.PoolId != nil {
			_ = d.Set("pool_id", respData.AddressPool.PoolId)
		}
	}

	if respData.AddressSet != nil && len(respData.AddressSet) > 0 {
		addressSetList := make([]map[string]interface{}, 0, len(respData.AddressSet))
		for _, addressSet := range respData.AddressSet {
			addressSetMap := map[string]interface{}{}
			if addressSet.Addr != nil {
				addressSetMap["addr"] = addressSet.Addr
			}

			if addressSet.IsEnable != nil {
				addressSetMap["is_enable"] = addressSet.IsEnable
			}

			if addressSet.AddressId != nil {
				addressSetMap["address_id"] = addressSet.AddressId
			}

			if addressSet.Location != nil {
				addressSetMap["location"] = addressSet.Location
			}

			if addressSet.Status != nil {
				addressSetMap["status"] = addressSet.Status
			}

			if addressSet.Weight != nil {
				addressSetMap["weight"] = addressSet.Weight
			}

			if addressSet.CreatedOn != nil {
				addressSetMap["created_on"] = addressSet.CreatedOn
			}

			if addressSet.UpdatedOn != nil {
				addressSetMap["updated_on"] = addressSet.UpdatedOn
			}

			addressSetList = append(addressSetList, addressSetMap)
		}

		_ = d.Set("address_set", addressSetList)
	}

	return nil
}

func resourceTencentCloudIgtmAddressPoolUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_address_pool.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		poolId = d.Id()
	)

	needChange := false
	mutableArgs := []string{"pool_name", "traffic_strategy", "monitor_id", "address_set"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := igtmv20231024.NewModifyAddressPoolRequest()
		if v, ok := d.GetOk("pool_name"); ok {
			request.PoolName = helper.String(v.(string))
		}

		if v, ok := d.GetOk("traffic_strategy"); ok {
			request.TrafficStrategy = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("monitor_id"); ok {
			request.MonitorId = helper.IntUint64(v.(int))
		}

		if v, ok := d.GetOk("address_set"); ok {
			for _, item := range v.(*schema.Set).List() {
				addressSetMap := item.(map[string]interface{})
				address := igtmv20231024.Address{}
				if v, ok := addressSetMap["addr"].(string); ok && v != "" {
					address.Addr = helper.String(v)
				}

				if v, ok := addressSetMap["is_enable"].(string); ok && v != "" {
					address.IsEnable = helper.String(v)
				}

				if v, ok := addressSetMap["address_id"].(int); ok && v != 0 {
					address.AddressId = helper.IntUint64(v)
				}

				if v, ok := addressSetMap["location"].(string); ok && v != "" {
					address.Location = helper.String(v)
				}

				if v, ok := addressSetMap["weight"].(int); ok {
					address.Weight = helper.IntUint64(v)
				}

				request.AddressSet = append(request.AddressSet, &address)
			}
		}

		request.PoolId = helper.StrToUint64Point(poolId)
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().ModifyAddressPoolWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update igtm address pool failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudIgtmAddressPoolRead(d, meta)
}

func resourceTencentCloudIgtmAddressPoolDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_igtm_address_pool.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = igtmv20231024.NewDeleteAddressPoolRequest()
		poolId  = d.Id()
	)

	request.PoolId = helper.StrToUint64Point(poolId)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseIgtmV20231024Client().DeleteAddressPoolWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil || result.Response.Msg == nil {
			return resource.NonRetryableError(fmt.Errorf("Delete igtm address pool failed, Response is nil."))
		}

		if *result.Response.Msg == "success" {
			return nil
		} else {
			return resource.NonRetryableError(fmt.Errorf("Delete igtm address pool failed, Msg is %s.", *result.Response.Msg))
		}
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete igtm address pool failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}
