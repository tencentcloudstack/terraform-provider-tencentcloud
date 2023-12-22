package waf

import (
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svccwp "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cwp"

	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func ResourceTencentCloudWafSaasInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafSaasInstanceCreate,
		Read:   resourceTencentCloudWafSaasInstanceRead,
		Update: resourceTencentCloudWafSaasInstanceUpdate,
		Delete: resourceTencentCloudWafSaasInstanceDelete,

		Schema: map[string]*schema.Schema{
			"goods_category": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(WAF_CATEGORY_SAAS),
				Description:  "Billing order parameters. support premium_saas, enterprise_saas, ultimate_saas.",
			},
			"time_span": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerMin(1),
				Default:      1,
				Description:  "Time interval.",
			},
			"time_unit": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: tccommon.ValidateAllowedStringValue(TIME_UNIT),
				Default:      "m",
				Description:  "Time unit, support d, m, y. d: day, m: month, y: year.",
			},
			"instance_name": {
				Optional:    true,
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Waf instance name.",
			},
			"auto_renew_flag": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      AUTO_RENEW_FLAG_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(AUTO_RENEW_FLAG),
				Description:  "Auto renew flag, 1: enable, 0: disable.",
			},
			"elastic_mode": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      ELASTIC_MODE_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(ELASTIC_MODE),
				Description:  "Is elastic billing enabled, 1: enable, 0: disable.",
			},
			"qps_limit": {
				Optional:     true,
				Computed:     true,
				Type:         schema.TypeInt,
				ValidateFunc: tccommon.ValidateIntegerMin(10000),
				Description:  "QPS Limit, Minimum setting 10000. Only `elastic_mode` is 1, can be set.",
			},
			"real_region": {
				Optional:     true,
				Type:         schema.TypeString,
				Default:      SAAS_REAL_REGION_MAINLAND_GZ,
				ValidateFunc: tccommon.ValidateAllowedStringValue(SAAS_REAL_REGIONS),
				Description:  "region. If Region is `ap-guangzhou`, support: gz, sh, bj, cd (Means: GuangZhou, ShangHai, BeiJing, ChengDu); If Region is `ap-seoul`, support: hk, sg, th, kr, in, de, ca, use, sao, usw, jkt (Means: HongKong, Singapore, Bandkok, Seoul, Mumbai, Frankfurt, Toronto, Virginia, SaoPaulo, SiliconValley, Jakarta).",
			},
			"bot_management": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      BOT_MANAGEMENT_STATUS_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(BOT_MANAGEMENT_STATUS),
				Description:  "Whether to purchase Bot management, 1: yes, 0: no. Default is 0.",
			},
			"api_security": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      API_SECURITY_STATUS_0,
				ValidateFunc: tccommon.ValidateAllowedIntValue(API_SECURITY_STATUS),
				Description:  "Whether to purchase API Security, 1: yes, 0: no. Default is 0.",
			},
			//"domain_pkg_count": {
			//	Optional:     true,
			//	Type:         schema.TypeInt,
			//	ValidateFunc: tccommon.ValidateIntegerMin(1),
			//	Description:  "Domain extension package count.",
			//},
			//"qps_pkg_count": {
			//	Optional:     true,
			//	Type:         schema.TypeInt,
			//	ValidateFunc: tccommon.ValidateIntegerMin(1),
			//	Description:  "QPS extension package count.",
			//},
			// computed
			"instance_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "waf instance id.",
			},
			"edition": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "waf instance edition, clb or saas.",
			},
			"begin_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "waf instance start time.",
			},
			"valid_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "waf instance valid time.",
			},
			"status": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "waf instance status.",
			},
		},
	}
}

func resourceTencentCloudWafSaasInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_saas_instance.create")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId         = tccommon.GetLogId(tccommon.ContextNil)
		request       = waf.NewGenerateDealsAndPayNewRequest()
		response      = waf.NewGenerateDealsAndPayNewResponse()
		client        = meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		instanceId    string
		mainlandMode  int
		realRegion    string
		realRegionInt int64
		goodsCategory string
	)

	region := client.Region
	if region == REGION_GZ {
		mainlandMode = REGION_ID_MAINLAND

	} else if region == REGION_KR {
		mainlandMode = REGION_ID_NON_MAINLAND

	} else {
		return fmt.Errorf("Region only supports `ap-guangzhou` and `ap-seoul`.")
	}

	goods := []*waf.GoodNews{}

	// make main instance
	instanceGood := waf.GoodNews{}
	instanceGoodDetail := new(waf.GoodsDetailNew)
	instanceGood.GoodsNum = helper.IntInt64(1)
	if v, ok := d.GetOk("goods_category"); ok {
		goodsCategory = v.(string)
		goodsCategoryId := int64(WAF_CATEGORY_ID_SAAS[goodsCategory])
		subProductCode := SUB_PRODUCT_CODE_SAAS[goodsCategory]
		labelTypes := LABEL_TYPES_SAAS[goodsCategory]
		pid := int64(PID_SAAS[goodsCategory])
		labelCounts := int64(1)

		instanceGood.GoodsCategoryId = &goodsCategoryId
		instanceGoodDetail.SubProductCode = &subProductCode
		instanceGoodDetail.Pid = &pid
		instanceGoodDetail.LabelTypes = helper.Strings([]string{labelTypes})
		instanceGoodDetail.LabelCounts = []*int64{&labelCounts}
	}

	if v, ok := d.GetOkExists("time_span"); ok {
		instanceGoodDetail.TimeSpan = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("time_unit"); ok {
		instanceGoodDetail.TimeUnit = helper.String(v.(string))
	}

	if v, ok := d.GetOk("instance_name"); ok {
		instanceGoodDetail.InstanceName = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
		instanceGoodDetail.AutoRenewFlag = helper.IntInt64(v.(int))
	}

	if v, ok := d.GetOk("real_region"); ok {
		realRegion = v.(string)
	}

	if mainlandMode == REGION_ID_MAINLAND {
		if !tccommon.IsContains(SAAS_REAL_REGION_MAINLAND, realRegion) {
			return fmt.Errorf("If Region is `ap-guangzhou`, parameter `real_region` is not legal")
		}

		instanceGood.RegionId = helper.IntInt64(mainlandMode)
		realRegionInt = int64(SAAS_REAL_REGION_MAINLAND_ID_MAP[realRegion])
		instanceGoodDetail.RealRegion = &realRegionInt

	} else {
		if !tccommon.IsContains(SAAS_REAL_REGION_NON_MAINLAND, realRegion) {
			return fmt.Errorf("If Region is `ap-seoul`, parameter `real_region` is not legal")
		}

		instanceGood.RegionId = helper.IntInt64(mainlandMode)
		realRegionInt = int64(SAAS_REAL_REGION_NON_MAINLAND_ID_MAP[realRegion])
		instanceGoodDetail.RealRegion = &realRegionInt
	}

	instanceGood.GoodsDetail = instanceGoodDetail
	goods = append(goods, &instanceGood)

	// bot management
	if v, ok := d.GetOkExists("bot_management"); ok {
		if v.(int) == 1 {
			botManagementGood := new(waf.GoodNews)
			botManagementDetail := new(waf.GoodsDetailNew)
			botManagementGood.GoodsCategoryId = helper.IntInt64(BOT_MANAGEMENT_CATEGORY_ID_SAAS)
			botManagementGood.GoodsNum = helper.IntInt64(1)
			botManagementDetail.SubProductCode = helper.String(BOT_MANAGEMENT_SUB_PRODUCT_CODE_SAAS)
			botManagementDetail.Pid = helper.IntInt64(BOT_MANAGEMENT_PID_SAAS)
			botManagementDetail.LabelTypes = helper.Strings([]string{BOT_MANAGEMENT_LABEL_TYPES_SAAS[mainlandMode]})
			botManagementDetail.LabelCounts = []*int64{helper.IntInt64(1)}

			if v, ok := d.GetOkExists("time_span"); ok {
				botManagementDetail.TimeSpan = helper.IntInt64(v.(int))
			}

			if v, ok := d.GetOk("time_unit"); ok {
				botManagementDetail.TimeUnit = helper.String(v.(string))
			}

			botManagementGood.RegionId = helper.IntInt64(mainlandMode)
			botManagementGood.GoodsDetail = botManagementDetail
			goods = append(goods, botManagementGood)
		}
	}

	// api security
	if v, ok := d.GetOkExists("api_security"); ok {
		if v.(int) == 1 {
			apiSecurityGood := new(waf.GoodNews)
			apiSecurityDetail := new(waf.GoodsDetailNew)
			apiSecurityGood.GoodsCategoryId = helper.IntInt64(API_SECURITY_CATEGORY_ID_SAAS)
			apiSecurityGood.GoodsNum = helper.IntInt64(1)
			apiSecurityDetail.SubProductCode = helper.String(API_SECURITY_SUB_PRODUCT_CODE_SAAS)
			apiSecurityDetail.Pid = helper.IntInt64(API_SECURITY_PID_SAAS)
			if mainlandMode == svccwp.REGION_ID_1 {
				apiSecurityDetail.LabelTypes = helper.Strings([]string{API_SECURITY_LABEL_TYPES_SAAS_REGION1[goodsCategory]})
			} else {
				apiSecurityDetail.LabelTypes = helper.Strings([]string{API_SECURITY_LABEL_TYPES_SAAS_REGION9[goodsCategory]})
			}
			apiSecurityDetail.LabelCounts = []*int64{helper.IntInt64(1)}

			if v, ok := d.GetOkExists("time_span"); ok {
				apiSecurityDetail.TimeSpan = helper.IntInt64(v.(int))
			}

			if v, ok := d.GetOk("time_unit"); ok {
				apiSecurityDetail.TimeUnit = helper.String(v.(string))
			}

			apiSecurityGood.RegionId = helper.IntInt64(mainlandMode)
			apiSecurityGood.GoodsDetail = apiSecurityDetail
			goods = append(goods, apiSecurityGood)
		}
	}

	// make domain pkg
	//if v, ok := d.GetOkExists("domain_pkg_count"); ok {
	//	domainPkgGood := waf.GoodNews{}
	//	domainPkgGoodDetail := new(waf.GoodsDetailNew)
	//	domainPkgGood.GoodsCategoryId = helper.IntInt64(DOMIAN_CATEGORY_ID_SAAS)
	//	domainPkgGood.GoodsNum = helper.IntInt64(1)
	//	domainPkgGoodDetail.SubProductCode = helper.String(DOMAIN_SUB_PRODUCT_CODE_SAAS)
	//	domainPkgGoodDetail.Pid = helper.IntInt64(DOMAIN_PID_SAAS)
	//	domainPkgGoodDetail.LabelTypes = helper.Strings([]string{DOMAIN_LABEL_TYPE_SAAS})
	//	domainPkgGoodDetail.LabelCounts = []*int64{helper.IntInt64(v.(int))}
	//
	//	if v, ok := d.GetOkExists("time_span"); ok {
	//		domainPkgGoodDetail.TimeSpan = helper.IntInt64(v.(int))
	//	}
	//
	//	if v, ok := d.GetOk("time_unit"); ok {
	//		domainPkgGoodDetail.TimeUnit = helper.String(v.(string))
	//	}
	//
	//	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
	//		domainPkgGoodDetail.AutoRenewFlag = helper.IntInt64(v.(int))
	//	}
	//
	//	domainPkgGood.RegionId = helper.IntInt64(mainlandMode)
	//	domainPkgGoodDetail.RealRegion = &realRegionInt
	//	domainPkgGood.GoodsDetail = domainPkgGoodDetail
	//	goods = append(goods, &domainPkgGood)
	//}
	//
	//// make qps pkg
	//if v, ok := d.GetOkExists("qps_pkg_count"); ok {
	//	qpsPkgGood := waf.GoodNews{}
	//	qpsPkgGoodDetail := new(waf.GoodsDetailNew)
	//	qpsPkgGood.GoodsCategoryId = helper.IntInt64(QPS_CATEGORY_ID_SAAS)
	//	qpsPkgGood.GoodsNum = helper.IntInt64(1)
	//	qpsPkgGoodDetail.SubProductCode = helper.String(QPS_SUB_PRODUCT_CODE_SAAS)
	//	qpsPkgGoodDetail.Pid = helper.IntInt64(QPS_PID_SAAS)
	//	qpsPkgGoodDetail.LabelTypes = helper.Strings([]string{QPS_LABEL_TYPE_SAAS})
	//	qpsPkgGoodDetail.LabelCounts = []*int64{helper.IntInt64(v.(int) * 1000)}
	//
	//	if v, ok := d.GetOkExists("time_span"); ok {
	//		qpsPkgGoodDetail.TimeSpan = helper.IntInt64(v.(int))
	//	}
	//
	//	if v, ok := d.GetOk("time_unit"); ok {
	//		qpsPkgGoodDetail.TimeUnit = helper.String(v.(string))
	//	}
	//
	//	if v, ok := d.GetOkExists("auto_renew_flag"); ok {
	//		qpsPkgGoodDetail.AutoRenewFlag = helper.IntInt64(v.(int))
	//	}
	//
	//	qpsPkgGood.RegionId = helper.IntInt64(mainlandMode)
	//	qpsPkgGoodDetail.RealRegion = &realRegionInt
	//	qpsPkgGood.GoodsDetail = qpsPkgGoodDetail
	//	goods = append(goods, &qpsPkgGood)
	//}

	request.Goods = goods
	err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().GenerateDealsAndPayNew(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if *result.Response.Status == 0 || *result.Response.InstanceId == "" {
			return resource.NonRetryableError(fmt.Errorf("create waf saas instance status error: %s", *result.Response.ReturnMessage))
		}

		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create waf saas instance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	// set elastic mode
	if v, ok := d.GetOkExists("elastic_mode"); ok {
		elasticMode := v.(int)
		if elasticMode == ELASTIC_MODE_1 {
			newSwitchElasticModeRequest := waf.NewSwitchElasticModeRequest()
			newSwitchElasticModeRequest.InstanceID = &instanceId
			newSwitchElasticModeRequest.Mode = helper.IntInt64(elasticMode)
			newSwitchElasticModeRequest.Edition = helper.String(EDITION_SAAS)
			err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().SwitchElasticMode(newSwitchElasticModeRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, newSwitchElasticModeRequest.GetAction(), newSwitchElasticModeRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update waf saas instance elastic mode failed, reason:%+v", logId, err)
				return err
			}

			// set qpsLimit
			if v, ok = d.GetOkExists("qps_limit"); ok {
				qpsLimit := v.(int)
				modifyInstanceQpsLimitRequest := waf.NewModifyInstanceQpsLimitRequest()
				modifyInstanceQpsLimitRequest.InstanceId = &instanceId
				modifyInstanceQpsLimitRequest.QpsLimit = helper.IntInt64(qpsLimit)
				err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyInstanceQpsLimit(modifyInstanceQpsLimitRequest)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyInstanceQpsLimitRequest.GetAction(), modifyInstanceQpsLimitRequest.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s update waf clb instance qpsLimit failed, reason:%+v", logId, err)
					return err
				}
			}
		} else {
			if _, ok = d.GetOkExists("qps_limit"); ok {
				return fmt.Errorf("If `elastic_mode` is 0, not support set `qps_limit`.")
			}
		}
	}

	return resourceTencentCloudWafSaasInstanceRead(d, meta)
}

func resourceTencentCloudWafSaasInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_saas_instance.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Id()
	)

	instanceInfo, err := service.DescribeWafInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instanceInfo == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instanceInfo.InstanceId != nil {
		_ = d.Set("instance_id", instanceInfo.InstanceId)
	}

	if instanceInfo.InstanceName != nil {
		_ = d.Set("instance_name", instanceInfo.InstanceName)
	}

	if instanceInfo.RenewFlag != nil {
		_ = d.Set("auto_renew_flag", instanceInfo.RenewFlag)
	}

	if instanceInfo.Mode != nil {
		_ = d.Set("elastic_mode", instanceInfo.Mode)
	}

	if instanceInfo.ElasticBilling != nil {
		_ = d.Set("qps_limit", instanceInfo.ElasticBilling)
	}

	if instanceInfo.Region != nil {
		_ = d.Set("real_region", instanceInfo.Region)
	}

	if instanceInfo.BotPkg != nil {
		_ = d.Set("bot_management", instanceInfo.BotPkg.Status)
	} else {
		_ = d.Set("bot_management", 0)
	}

	if instanceInfo.ApiPkg != nil {
		_ = d.Set("api_security", instanceInfo.ApiPkg.Status)
	} else {
		_ = d.Set("api_security", 0)
	}

	//if instanceInfo.DomainPkg != nil {
	//	_ = d.Set("domain_pkg_count", instanceInfo.DomainPkg.Count)
	//}
	//
	//if instanceInfo.QPS != nil {
	//	qpsCount := *instanceInfo.QPS.Count / 1000
	//	_ = d.Set("qps_pkg_count", qpsCount)
	//}

	if instanceInfo.Edition != nil {
		_ = d.Set("edition", instanceInfo.Edition)
	}

	if instanceInfo.BeginTime != nil {
		_ = d.Set("begin_time", instanceInfo.BeginTime)
	}

	if instanceInfo.ValidTime != nil {
		_ = d.Set("valid_time", instanceInfo.ValidTime)
	}

	if instanceInfo.Status != nil {
		_ = d.Set("status", instanceInfo.Status)
	}

	return nil
}

func resourceTencentCloudWafSaasInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_saas_instance.update")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                          = tccommon.GetLogId(tccommon.ContextNil)
		modifyInstanceNameRequest      = waf.NewModifyInstanceNameRequest()
		modifyInstanceRenewFlagRequest = waf.NewModifyInstanceRenewFlagRequest()
		newSwitchElasticModeRequest    = waf.NewSwitchElasticModeRequest()
		instanceId                     = d.Id()
		elasticMode                    int
	)

	immutableArgs := []string{"goods_category", "time_span", "time_unit", "domain_pkg_count", "qps_pkg_count", "bot_management", "api_security"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_name") {
		if v, ok := d.GetOkExists("instance_name"); ok {
			modifyInstanceNameRequest.InstanceID = &instanceId
			modifyInstanceNameRequest.InstanceName = helper.String(v.(string))
			modifyInstanceNameRequest.Edition = helper.String(EDITION_SAAS)
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyInstanceName(modifyInstanceNameRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyInstanceNameRequest.GetAction(), modifyInstanceNameRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update waf saas instance name failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if d.HasChange("auto_renew_flag") {
		if v, ok := d.GetOkExists("auto_renew_flag"); ok {
			modifyInstanceRenewFlagRequest.InstanceId = &instanceId
			modifyInstanceRenewFlagRequest.RenewFlag = helper.IntInt64(v.(int))
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyInstanceRenewFlag(modifyInstanceRenewFlagRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyInstanceRenewFlagRequest.GetAction(), modifyInstanceRenewFlagRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update waf saas instance auto renew flag failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if d.HasChange("elastic_mode") {
		if v, ok := d.GetOkExists("elastic_mode"); ok {
			newSwitchElasticModeRequest.InstanceID = &instanceId
			newSwitchElasticModeRequest.Mode = helper.IntInt64(v.(int))
			newSwitchElasticModeRequest.Edition = helper.String(EDITION_SAAS)
			err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
				result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().SwitchElasticMode(newSwitchElasticModeRequest)
				if e != nil {
					return tccommon.RetryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, newSwitchElasticModeRequest.GetAction(), newSwitchElasticModeRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update waf saas instance elastic mode failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	if v, ok := d.GetOkExists("elastic_mode"); ok {
		elasticMode = v.(int)
	}

	if elasticMode == ELASTIC_MODE_1 {
		if d.HasChange("qps_limit") {
			if v, ok := d.GetOkExists("qps_limit"); ok {
				qpsLimit := v.(int)
				modifyInstanceQpsLimitRequest := waf.NewModifyInstanceQpsLimitRequest()
				modifyInstanceQpsLimitRequest.InstanceId = &instanceId
				modifyInstanceQpsLimitRequest.QpsLimit = helper.IntInt64(qpsLimit)
				err := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
					result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseWafClient().ModifyInstanceQpsLimit(modifyInstanceQpsLimitRequest)
					if e != nil {
						return tccommon.RetryError(e)
					} else {
						log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyInstanceQpsLimitRequest.GetAction(), modifyInstanceQpsLimitRequest.ToJsonString(), result.ToJsonString())
					}

					return nil
				})

				if err != nil {
					log.Printf("[CRITAL]%s update waf clb instance qpsLimit failed, reason:%+v", logId, err)
					return err
				}
			}
		}
	} else {
		if _, ok := d.GetOkExists("qps_limit"); ok {
			return fmt.Errorf("If `elastic_mode` is 0, not support set `qps_limit`.")
		}
	}

	return resourceTencentCloudWafSaasInstanceRead(d, meta)
}

func resourceTencentCloudWafSaasInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_waf_saas_instance.delete")()
	defer tccommon.InconsistentCheck(d, meta)()

	return fmt.Errorf("tencentcloud waf saas instance not supported delete, please contact the work order for processing")
}
