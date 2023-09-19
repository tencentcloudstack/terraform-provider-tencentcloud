/*
Provides a resource to create a waf saas instance

Example Usage

Create a basic waf premium saas instance

```hcl
resource "tencentcloud_waf_saas_instance" "example" {
  goods_category   = "premium_saas"
  instance_name    = "tf-example-saas-waf"
}
```

Create a complete waf ultimate_saas instance(Chinese Mainland)

```hcl
resource "tencentcloud_waf_saas_instance" "example" {
  goods_category   = "ultimate_saas"
  instance_name    = "tf-example-saas-waf"
  time_span        = 1
  time_unit        = "m"
  auto_renew_flag  = 1
  elastic_mode     = 1
  is_cn_mainland   = 1
  real_region      = "gz"
  domain_pkg_count = 3
  qps_pkg_count    = 3
}
```

Create a complete waf ultimate_saas instance(Non Chinese Mainland)

```hcl
resource "tencentcloud_waf_saas_instance" "example" {
  goods_category   = "ultimate_saas"
  instance_name    = "tf-example-saas-waf"
  time_span        = 1
  time_unit        = "m"
  auto_renew_flag  = 1
  elastic_mode     = 1
  is_cn_mainland   = 0
  real_region      = "sg"
  domain_pkg_count = 3
  qps_pkg_count    = 3
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudWafSaasInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafSaasInstanceCreate,
		Read:   resourceTencentCloudWafSaasInstanceRead,
		Update: resourceTencentCloudWafSaasInstanceUpdate,
		Delete: resourceTencentCloudWafSaasInstanceDelete,

		Schema: map[string]*schema.Schema{
			"goods_category": {
				Required:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue(WAF_CATEGORY_SAAS),
				Description:  "Billing order parameters. support premium_saas, enterprise_saas, ultimate_saas.",
			},
			"time_span": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerMin(1),
				Default:      1,
				Description:  "Time interval.",
			},
			"time_unit": {
				Optional:     true,
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue(TIME_UNIT),
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
				ValidateFunc: validateAllowedIntValue(AUTO_RENEW_FLAG),
				Description:  "Auto renew flag, 1: enable, 0: disable.",
			},
			"elastic_mode": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      ELASTIC_MODE_0,
				ValidateFunc: validateAllowedIntValue(ELASTIC_MODE),
				Description:  "Is elastic billing enabled, 1: enable, 0: disable.",
			},
			"is_cn_mainland": {
				Optional:     true,
				Type:         schema.TypeInt,
				Default:      MAINLAND_1,
				ValidateFunc: validateAllowedIntValue(MAINLAND),
				Description:  "Chinese Mainland or not, 1: Chinese Mainland, 0: Non Chinese Mainland.",
			},
			"real_region": {
				Optional:     true,
				Type:         schema.TypeString,
				Default:      SAAS_REAL_REGION_NON_MAINLAND_HK,
				ValidateFunc: validateAllowedStringValue(SAAS_REAL_REGIONS),
				Description:  "region. If `is_cn_mainland` is 1, support: gz, sh, bj, cd (Means: GuangZhou, ShangHai, BeiJing, ChengDu); If `is_cn_mainland` is 0, support: hk, sg, th, kr, in, de, ca, use, sao, usw, jkt (Means: HongKong, Singapore, Bandkok, Seoul, Mumbai, Frankfurt, Toronto, Virginia, SaoPaulo, SiliconValley, Jakarta).",
			},
			"domain_pkg_count": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerMin(1),
				Description:  "Domain extension package count.",
			},
			"qps_pkg_count": {
				Optional:     true,
				Type:         schema.TypeInt,
				ValidateFunc: validateIntegerMin(1),
				Description:  "QPS extension package count.",
			},
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
			"api_security": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "waf instance api security status.",
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
	defer logElapsed("resource.tencentcloud_waf_saas_instance.create")()
	defer inconsistentCheck(d, meta)()

	var (
		logId         = getLogId(contextNil)
		request       = waf.NewGenerateDealsAndPayNewRequest()
		response      = waf.NewGenerateDealsAndPayNewResponse()
		instanceId    string
		mainlandMode  int
		readRegion    string
		realRegionInt int64
	)

	goods := []*waf.GoodNews{}

	// make main instance
	instanceGood := waf.GoodNews{}
	instanceGoodDetail := new(waf.GoodsDetailNew)
	instanceGood.GoodsNum = helper.IntInt64(1)
	if v, ok := d.GetOk("goods_category"); ok {
		goodsCategory := v.(string)
		goodsCategoryId := int64(WAF_CATEGORY_ID_SAAS[goodsCategory])
		subProductCode := SUB_PRODUCT_CODE_SAAS[goodsCategory]
		labelTypes := LABEL_TYPES_SAAS[goodsCategory]
		pid := int64(PID_CLB[goodsCategory])
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
		readRegion = v.(string)
	}

	if v, ok := d.GetOkExists("is_cn_mainland"); ok {
		mainlandMode = v.(int)
		if mainlandMode == MAINLAND_1 {
			if !IsContains(SAAS_REAL_REGION_MAINLAND, readRegion) {
				return fmt.Errorf("If `is_cn_mainland` is 1, parameter `real_region` is not legal")
			}

			instanceGood.RegionId = helper.IntInt64(REGION_ID_MAINLAND)
			realRegionInt = int64(SAAS_REAL_REGION_MAINLAND_ID_MAP[readRegion])
			instanceGoodDetail.RealRegion = &realRegionInt

		} else {
			if !IsContains(SAAS_REAL_REGION_NON_MAINLAND, readRegion) {
				return fmt.Errorf("If `is_cn_mainland` is 0, parameter `real_region` is not legal")
			}

			instanceGood.RegionId = helper.IntInt64(REGION_ID_NON_MAINLAND)
			realRegionInt = int64(SAAS_REAL_REGION_NON_MAINLAND_ID_MAP[readRegion])
			instanceGoodDetail.RealRegion = &realRegionInt
		}
	}

	instanceGood.GoodsDetail = instanceGoodDetail
	goods = append(goods, &instanceGood)

	// make domain pkg
	if v, ok := d.GetOkExists("domain_pkg_count"); ok {
		domainPkgGood := waf.GoodNews{}
		domainPkgGoodDetail := new(waf.GoodsDetailNew)
		domainPkgGood.GoodsCategoryId = helper.IntInt64(DOMIAN_CATEGORY_ID_SAAS)
		domainPkgGood.GoodsNum = helper.IntInt64(1)
		domainPkgGoodDetail.SubProductCode = helper.String(DOMAIN_SUB_PRODUCT_CODE_SAAS)
		domainPkgGoodDetail.Pid = helper.IntInt64(DOMAIN_PID_SAAS)
		domainPkgGoodDetail.LabelTypes = helper.Strings([]string{DOMAIN_LABEL_TYPE_SAAS})
		domainPkgGoodDetail.LabelCounts = []*int64{helper.IntInt64(v.(int))}

		if v, ok := d.GetOkExists("time_span"); ok {
			domainPkgGoodDetail.TimeSpan = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("time_unit"); ok {
			domainPkgGoodDetail.TimeUnit = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("auto_renew_flag"); ok {
			domainPkgGoodDetail.AutoRenewFlag = helper.IntInt64(v.(int))
		}

		if mainlandMode == MAINLAND_1 {
			domainPkgGood.RegionId = helper.IntInt64(REGION_ID_MAINLAND)
		} else {
			domainPkgGood.RegionId = helper.IntInt64(REGION_ID_NON_MAINLAND)
		}

		domainPkgGoodDetail.RealRegion = &realRegionInt
		domainPkgGood.GoodsDetail = domainPkgGoodDetail
		goods = append(goods, &domainPkgGood)
	}

	// make qps pkg
	if v, ok := d.GetOkExists("qps_pkg_count"); ok {
		qpsPkgGood := waf.GoodNews{}
		qpsPkgGoodDetail := new(waf.GoodsDetailNew)
		qpsPkgGood.GoodsCategoryId = helper.IntInt64(QPS_CATEGORY_ID_SAAS)
		qpsPkgGood.GoodsNum = helper.IntInt64(1)
		qpsPkgGoodDetail.SubProductCode = helper.String(QPS_SUB_PRODUCT_CODE_SAAS)
		qpsPkgGoodDetail.Pid = helper.IntInt64(QPS_PID_SAAS)
		qpsPkgGoodDetail.LabelTypes = helper.Strings([]string{QPS_LABEL_TYPE_SAAS})
		qpsPkgGoodDetail.LabelCounts = []*int64{helper.IntInt64(v.(int) * 1000)}

		if v, ok := d.GetOkExists("time_span"); ok {
			qpsPkgGoodDetail.TimeSpan = helper.IntInt64(v.(int))
		}

		if v, ok := d.GetOk("time_unit"); ok {
			qpsPkgGoodDetail.TimeUnit = helper.String(v.(string))
		}

		if v, ok := d.GetOkExists("auto_renew_flag"); ok {
			qpsPkgGoodDetail.AutoRenewFlag = helper.IntInt64(v.(int))
		}

		if mainlandMode == MAINLAND_1 {
			qpsPkgGood.RegionId = helper.IntInt64(REGION_ID_MAINLAND)
		} else {
			qpsPkgGood.RegionId = helper.IntInt64(REGION_ID_NON_MAINLAND)
		}

		qpsPkgGoodDetail.RealRegion = &realRegionInt
		qpsPkgGood.GoodsDetail = qpsPkgGoodDetail
		goods = append(goods, &qpsPkgGood)
	}

	request.Goods = goods
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().GenerateDealsAndPayNew(request)
		if e != nil {
			return retryError(e)
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
			modifyInstanceElasticModeRequest := waf.NewModifyInstanceElasticModeRequest()
			modifyInstanceElasticModeRequest.InstanceId = &instanceId
			modifyInstanceElasticModeRequest.Mode = helper.IntInt64(elasticMode)
			err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyInstanceElasticMode(modifyInstanceElasticModeRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyInstanceElasticModeRequest.GetAction(), modifyInstanceElasticModeRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update waf saas instance elastic mode failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudWafClbInstanceRead(d, meta)
}

func resourceTencentCloudWafSaasInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_saas_instance.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = WafService{client: meta.(*TencentCloudClient).apiV3Conn}
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

	if instanceInfo.Region != nil {
		_ = d.Set("real_region", instanceInfo.Region)
	}

	if instanceInfo.DomainPkg != nil {
		_ = d.Set("domain_pkg_count", instanceInfo.DomainPkg.Count)
	}

	if instanceInfo.QPS != nil {
		_ = d.Set("qps_pkg_count", instanceInfo.QPS.Count)
	}

	if instanceInfo.Edition != nil {
		_ = d.Set("edition", instanceInfo.Edition)
	}

	if instanceInfo.BeginTime != nil {
		_ = d.Set("begin_time", instanceInfo.BeginTime)
	}

	if instanceInfo.ValidTime != nil {
		_ = d.Set("valid_time", instanceInfo.ValidTime)
	}

	if instanceInfo.APISecurity != nil {
		_ = d.Set("api_security", instanceInfo.APISecurity)
	}

	if instanceInfo.Status != nil {
		_ = d.Set("status", instanceInfo.Status)
	}

	return nil
}

func resourceTencentCloudWafSaasInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_saas_instance.update")()
	defer inconsistentCheck(d, meta)()

	var (
		logId                            = getLogId(contextNil)
		modifyInstanceNameRequest        = waf.NewModifyInstanceNameRequest()
		modifyInstanceRenewFlagRequest   = waf.NewModifyInstanceRenewFlagRequest()
		modifyInstanceElasticModeRequest = waf.NewModifyInstanceElasticModeRequest()
		instanceId                       = d.Id()
	)

	immutableArgs := []string{"goods_category", "time_span", "time_unit", "is_cn_mainland", "domain_pkg_count", "qps_pkg_count"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("instance_name") {
		if v, ok := d.GetOkExists("instance_name"); ok {
			modifyInstanceNameRequest.InstanceID = &instanceId
			modifyInstanceNameRequest.InstanceName = helper.String(v.(string))
			modifyInstanceNameRequest.Edition = helper.String("sparta-waf")
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyInstanceName(modifyInstanceNameRequest)
				if e != nil {
					return retryError(e)
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
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyInstanceRenewFlag(modifyInstanceRenewFlagRequest)
				if e != nil {
					return retryError(e)
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
			modifyInstanceElasticModeRequest.InstanceId = &instanceId
			modifyInstanceElasticModeRequest.Mode = helper.IntInt64(v.(int))
			err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
				result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyInstanceElasticMode(modifyInstanceElasticModeRequest)
				if e != nil {
					return retryError(e)
				} else {
					log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, modifyInstanceElasticModeRequest.GetAction(), modifyInstanceElasticModeRequest.ToJsonString(), result.ToJsonString())
				}

				return nil
			})

			if err != nil {
				log.Printf("[CRITAL]%s update waf saas instance elastic mode failed, reason:%+v", logId, err)
				return err
			}
		}
	}

	return resourceTencentCloudWafSaasInstanceRead(d, meta)
}

func resourceTencentCloudWafSaasInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_saas_instance.delete")()
	defer inconsistentCheck(d, meta)()

	return fmt.Errorf("tencentcloud waf saas instance not supported delete, please contact the work order for processing")
}
