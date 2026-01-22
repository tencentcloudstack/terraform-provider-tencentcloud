package teo

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudTeoPlan() *schema.Resource {
	return &schema.Resource{
		Create: ResourceTencentCloudTeoPlanCreate,
		Read:   ResourceTencentCloudTeoPlanRead,
		Update: ResourceTencentCloudTeoPlanUpdate,
		Delete: ResourceTencentCloudTeoPlanDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"plan_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"personal", "basic", "standard", "enterprise"}),
				Description:  "The subscription package type, the possible values are: `personal`: personal package, prepaid package; `basic`: basic package, prepaid package; `standard`: standard package, prepaid package; `enterprise`: enterprise package, postpaid package.",
			},

			"prepaid_plan_param": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Subscription prepaid package parameters. When PlanType is personal, basic, or standard, this parameter is optional and is used to enter the subscription duration of the package and whether to enable automatic renewal. If this parameter is not filled in, the default subscription duration is 1 month and automatic renewal is not enabled.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36}),
							Description:  "The subscription period of the prepaid package, in months, with possible values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36. If not filled in, the default value 1 is used.",
						},
						"renew_flag": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: tccommon.ValidateAllowedStringValue([]string{"on", "off"}),
							Description:  "The automatic renewal flag of the prepaid package, the values are: `on`: turn on automatic renewal; `off`: do not turn on automatic renewal. If not filled in, the default value off is used. When automatic renewal occurs, the default renewal period is 1 month.",
						},
					},
				},
			},

			// computed
			"plan_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Plan ID.",
			},

			"area": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Service area, possible values are: <li>mainland: Mainland China; </li><li>overseas: Worldwide (excluding Mainland China); </li><li>global: Worldwide (including Mainland China). </li>.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Package status, the values are: <li>normal: normal status; </li><li>expiring-soon: about to expire; </li><li>expired: expired; </li><li>isolated: isolated; </li><li>overdue-isolated: overdue isolated. </li>.",
			},

			"pay_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Payment type, possible values: <li>0: post-payment; </li><li>1: pre-payment. </li>.",
			},

			"enabled_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the package takes effect.",
			},

			"expired_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration date of the package.",
			},
		},
	}
}

func ResourceTencentCloudTeoPlanCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_plan.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request  = teov20220901.NewCreatePlanRequest()
		response = teov20220901.NewCreatePlanResponse()
		planId   string
	)

	if v, ok := d.GetOk("plan_type"); ok {
		request.PlanType = helper.String(v.(string))
	}

	if dMap, ok := helper.InterfacesHeadMap(d, "prepaid_plan_param"); ok {
		prepaidPlanParam := teov20220901.PrepaidPlanParam{}
		if v, ok := dMap["period"].(int); ok && v != 0 {
			prepaidPlanParam.Period = helper.IntInt64(v)
		}

		if v, ok := dMap["renew_flag"].(string); ok && v != "" {
			prepaidPlanParam.RenewFlag = helper.String(v)
		}

		request.PrepaidPlanParam = &prepaidPlanParam
	}

	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().CreatePlanWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Create teo plan failed, Response is nil."))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create teo function failed, reason:%+v", logId, err)
		return err
	}

	if response.Response.PlanId == nil {
		return fmt.Errorf("PlanId is nil.")
	}

	planId = *response.Response.PlanId
	d.SetId(planId)
	return ResourceTencentCloudTeoPlanRead(d, meta)
}

func ResourceTencentCloudTeoPlanRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_plan.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		planId  = d.Id()
	)

	respData, err := service.DescribeTeoPlansById(ctx, planId)
	if err != nil {
		return err
	}

	if respData == nil {
		log.Printf("[WARN]%s resource `tencentcloud_teo_plan` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		d.SetId("")
		return nil
	}

	if respData.PlanType != nil {
		_ = d.Set("plan_type", respData.PlanType)
	}

	if respData.PlanId != nil {
		_ = d.Set("plan_id", respData.PlanId)
	}

	if respData.Area != nil {
		_ = d.Set("area", respData.Area)
	}

	if respData.Status != nil {
		_ = d.Set("status", respData.Status)
	}

	if respData.PayMode != nil {
		_ = d.Set("pay_mode", respData.PayMode)
	}

	if respData.EnabledTime != nil {
		_ = d.Set("enabled_time", respData.EnabledTime)
	}

	if respData.ExpiredTime != nil {
		_ = d.Set("expired_time", respData.ExpiredTime)
	}

	return nil
}

func ResourceTencentCloudTeoPlanUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_plan.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId  = tccommon.GetLogId(tccommon.ContextNil)
		ctx    = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		planId = d.Id()
	)

	if d.HasChange("plan_type") {
		request := teov20220901.NewUpgradePlanRequest()
		if v, ok := d.GetOk("plan_type"); ok {
			request.PlanType = helper.String(v.(string))
		}

		request.PlanId = &planId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().UpgradePlanWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("prepaid_plan_param.0.period") {
		request := teov20220901.NewRenewPlanRequest()
		if v, ok := d.GetOk("period"); ok {
			request.Period = helper.IntInt64(v.(int))
		}

		request.PlanId = &planId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().RenewPlanWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	if d.HasChange("prepaid_plan_param.0.renew_flag") {
		request := teov20220901.NewModifyPlanRequest()
		if v, ok := d.GetOk("renew_flag"); ok {
			request.RenewFlag = &teov20220901.RenewFlag{
				Switch: helper.String(v.(string)),
			}
		}

		request.PlanId = &planId
		err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().ModifyPlanWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			} else {
				log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return ResourceTencentCloudTeoPlanRead(d, meta)
}

func ResourceTencentCloudTeoPlanDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_teo_plan.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		request = teov20220901.NewDestroyPlanRequest()
		planId  = d.Id()
	)

	request.PlanId = &planId
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client().DestroyPlanWithContext(ctx, request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s delete teo plan failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
