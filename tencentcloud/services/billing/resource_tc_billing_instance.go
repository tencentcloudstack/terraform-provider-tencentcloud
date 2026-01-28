package billing

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	billingv20180709 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/billing/v20180709"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudBillingInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudBillingInstanceCreate,
		Read:   resourceTencentCloudBillingInstanceRead,
		Update: resourceTencentCloudBillingInstanceUpdate,
		Delete: resourceTencentCloudBillingInstanceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"product_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Product code.",
			},

			"sub_product_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sub-product code.",
			},

			"region_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region code.",
			},

			"zone_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Availability zone code.",
			},

			"pay_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Payment mode. Available values: PrePay: upfront charge.",
			},

			"parameter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Product detailed information.",
			},

			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Project id, default value is 0.",
			},

			"period": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Purchase duration, max number is 36, default value is 1.",
			},

			"period_unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Purchase duration unit. valid values: \nm: month,\ny: year. \ndefault value is: m.",
			},

			"renew_flag": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Auto-renewal flag. valid values: NOTIFY_AND_MANUAL_RENEW: manually renew, NOTIFY_AND_AUTO_RENEW: automatically renew, DISABLE_NOTIFY_AND_MANUAL_RENEW: renewal is disabled. \ndefault value is NOTIFY_AND_MANUAL_RENEW.",
			},

			// computed
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance id.",
			},
		},
	}
}

func resourceTencentCloudBillingInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_billing_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = billingv20180709.NewCreateInstanceRequest()
		response   = billingv20180709.NewCreateInstanceResponse()
		instanceId string
	)

	// get current date string
	startDate, endDate := getDateRange()

	if v, ok := d.GetOk("product_code"); ok {
		request.ProductCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sub_product_code"); ok {
		request.SubProductCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("region_code"); ok {
		request.RegionCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("zone_code"); ok {
		request.ZoneCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("pay_mode"); ok {
		request.PayMode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("parameter"); ok {
		request.Parameter = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("project_id"); ok {
		request.ProjectId = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOkExists("period"); ok {
		request.Period = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("period_unit"); ok {
		request.PeriodUnit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("renew_flag"); ok {
		request.RenewFlag = helper.String(v.(string))
	}

	tmpUUID := helper.BuildUUID()
	request.ClientToken = &tmpUUID
	request.Quantity = helper.IntInt64(1)
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBillingV20180709Client().CreateInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create billing instance failed, Response is nil."))
		}

		// if null, need to retry, until get InstanceIdList
		if len(result.Response.InstanceIdList) == 0 {
			return resource.RetryableError(fmt.Errorf("InstanceIdList is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create billing instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	instanceId = *response.Response.InstanceIdList[0]
	d.SetId(instanceId)

	if response.Response.OrderId == nil {
		return fmt.Errorf("OrderId is nil.")
	}

	orderId := *response.Response.OrderId

	// wait
	waitReq := billingv20180709.NewDescribeDealsByCondRequest()
	waitReq.StartTime = &startDate
	waitReq.EndTime = &endDate
	waitReq.Offset = helper.IntInt64(0)
	waitReq.Limit = helper.IntInt64(1)
	waitReq.OrderId = &orderId
	waitReq.StatusSet = helper.Int64Slice2Int64PointerSlice([]int64{4, 6})
	reqErr = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBillingV20180709Client().DescribeDealsByCondWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe deals by cond failed, Response is nil."))
		}

		if len(result.Response.Deals) == 0 {
			return resource.RetryableError(fmt.Errorf("Deals is nil."))
		}

		deal := result.Response.Deals[0]
		if deal.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe deals by cond failed, Status is nil."))
		}

		if *deal.Status == 4 {
			return nil
		}

		if *deal.Status == 6 {
			return resource.NonRetryableError(fmt.Errorf("Create billing instance failed, Status is 6."))
		}

		return resource.RetryableError(fmt.Errorf("Billing instance is still creating. Status is %d.", *deal.Status))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s create billing instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return resourceTencentCloudBillingInstanceRead(d, meta)
}

func resourceTencentCloudBillingInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_billing_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = BillingService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	respData, err := service.DescribeBillingInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_billing_instance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.ProductCode != nil {
		_ = d.Set("product_code", respData.ProductCode)
	}

	if respData.SubProductCode != nil {
		_ = d.Set("sub_product_code", respData.SubProductCode)
	}

	if respData.RegionCode != nil {
		_ = d.Set("region_code", respData.RegionCode)
	}

	if respData.RenewPeriodUnit != nil {
		_ = d.Set("period_unit", respData.RenewPeriodUnit)
	}

	if respData.RenewFlag != nil {
		_ = d.Set("renew_flag", respData.RenewFlag)
	}

	if respData.InstanceId != nil {
		_ = d.Set("instance_id", respData.InstanceId)
	}

	return nil
}

func resourceTencentCloudBillingInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_billing_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		instanceId = d.Id()
	)

	immutableArgs := []string{"product_code", "sub_product_code", "region_code", "zone_code", "pay_mode", "parameter", "project_id", "renew_flag"}
	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	needChange := false
	mutableArgs := []string{"period", "period_unit"}
	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {
		request := billingv20180709.NewRenewInstanceRequest()
		response := billingv20180709.NewRenewInstanceResponse()

		// get current date string
		startDate, endDate := getDateRange()

		if v, ok := d.GetOk("product_code"); ok {
			request.ProductCode = helper.String(v.(string))
		}

		if v, ok := d.GetOk("sub_product_code"); ok {
			request.SubProductCode = helper.String(v.(string))
		}

		if v, ok := d.GetOk("region_code"); ok {
			request.RegionCode = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("period"); ok {
			request.Period = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("period_unit"); ok {
			request.PeriodUnit = helper.String(v.(string))
		}

		tmpUUID := helper.BuildUUID()
		request.ClientToken = &tmpUUID
		request.InstanceId = &instanceId
		reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBillingV20180709Client().RenewInstanceWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Renew instance failed, Response is nil."))
			}

			response = result
			return nil
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update billing instance failed, reason:%+v", logId, reqErr)
			return reqErr
		}

		if len(response.Response.OrderIdList) == 0 {
			return fmt.Errorf("OrderIdList is nil.")
		}

		orderId := *response.Response.OrderIdList[0]

		// wait
		waitReq := billingv20180709.NewDescribeDealsByCondRequest()
		waitReq.StartTime = &startDate
		waitReq.EndTime = &endDate
		waitReq.Offset = helper.IntInt64(0)
		waitReq.Limit = helper.IntInt64(1)
		waitReq.OrderId = &orderId
		waitReq.StatusSet = helper.Int64Slice2Int64PointerSlice([]int64{4, 6})
		reqErr = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBillingV20180709Client().DescribeDealsByCondWithContext(ctx, waitReq)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
			}

			if result == nil || result.Response == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe deals by cond failed, Response is nil."))
			}

			if len(result.Response.Deals) == 0 {
				return resource.RetryableError(fmt.Errorf("Deals is nil."))
			}

			deal := result.Response.Deals[0]
			if deal.Status == nil {
				return resource.NonRetryableError(fmt.Errorf("Describe deals by cond failed, Status is nil."))
			}

			if *deal.Status == 4 {
				return nil
			}

			if *deal.Status == 6 {
				return resource.NonRetryableError(fmt.Errorf("Update billing instance failed, Status is 6."))
			}

			return resource.RetryableError(fmt.Errorf("Billing instance is still updating. Status is %d.", *deal.Status))
		})

		if reqErr != nil {
			log.Printf("[CRITAL]%s update billing instance failed, reason:%+v", logId, reqErr)
			return reqErr
		}
	}

	return resourceTencentCloudBillingInstanceRead(d, meta)
}

func resourceTencentCloudBillingInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_billing_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request    = billingv20180709.NewRefundInstanceRequest()
		response   = billingv20180709.NewRefundInstanceResponse()
		instanceId = d.Id()
	)

	// get current date string
	startDate, endDate := getDateRange()

	if v, ok := d.GetOk("product_code"); ok {
		request.ProductCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("sub_product_code"); ok {
		request.SubProductCode = helper.String(v.(string))
	}

	if v, ok := d.GetOk("region_code"); ok {
		request.RegionCode = helper.String(v.(string))
	}

	tmpUUID := helper.BuildUUID()
	request.ClientToken = &tmpUUID
	request.InstanceId = &instanceId
	reqErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBillingV20180709Client().RefundInstanceWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Refund billing instance failed, Response is nil."))
		}

		response = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete billing instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	if len(response.Response.OrderIdList) == 0 {
		return fmt.Errorf("OrderIdList is nil.")
	}

	orderId := *response.Response.OrderIdList[0]

	// wait
	waitReq := billingv20180709.NewDescribeDealsByCondRequest()
	waitReq.StartTime = &startDate
	waitReq.EndTime = &endDate
	waitReq.Offset = helper.IntInt64(0)
	waitReq.Limit = helper.IntInt64(1)
	waitReq.OrderId = &orderId
	waitReq.StatusSet = helper.Int64Slice2Int64PointerSlice([]int64{6, 7})
	reqErr = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseBillingV20180709Client().DescribeDealsByCondWithContext(ctx, waitReq)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, waitReq.GetAction(), waitReq.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe deals by cond failed, Response is nil."))
		}

		if len(result.Response.Deals) == 0 {
			return resource.RetryableError(fmt.Errorf("Deals is nil."))
		}

		deal := result.Response.Deals[0]
		if deal.Status == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe deals by cond failed, Status is nil."))
		}

		if *deal.Status == 6 {
			return nil
		}

		if *deal.Status == 7 {
			return resource.NonRetryableError(fmt.Errorf("Delete billing instance failed, Status is 6."))
		}

		return resource.RetryableError(fmt.Errorf("Billing instance is still deleting. Status is %d.", *deal.Status))
	})

	if reqErr != nil {
		log.Printf("[CRITAL]%s delete billing instance failed, reason:%+v", logId, reqErr)
		return reqErr
	}

	return nil
}

func getDateRange() (startDate, endDate string) {
	now := time.Now()
	startDate = now.AddDate(0, 0, -1).Format("2006-01-02") + " 00:00:00"
	endDate = now.AddDate(0, 0, 1).Format("2006-01-02") + " 00:00:00"
	return
}
