/*
Provides a resource to create a waf instance

Example Usage

```hcl
resource "tencentcloud_waf_instance" "instance" {
  goods {
		goods_num =
		goods_detail {
			time_span =
			time_unit = ""
			sub_product_code = ""
			pid =
			instance_name = ""
			auto_renew_flag =
			real_region =
			label_types =
			label_counts =
			cur_deadline = ""
			instance_id = ""
		}
		goods_category_id =
		region_id =

  }
}
```

Import

waf instance can be imported using the id, e.g.

```
terraform import tencentcloud_waf_instance.instance instance_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudWafInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudWafInstanceCreate,
		Read:   resourceTencentCloudWafInstanceRead,
		Update: resourceTencentCloudWafInstanceUpdate,
		Delete: resourceTencentCloudWafInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"goods": {
				Required:    true,
				Type:        schema.TypeList,
				Description: "Billing order parameters.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"goods_num": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Goods quantity.",
						},
						"goods_detail": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Required:    true,
							Description: "Goods detail info.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time_span": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Time intervalNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"time_unit": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Time unit, support m、y、dNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"sub_product_code": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Sub product labels, and must be filled when new purchase or renew. when changing configurations, it should be placed in oldConfig and newConfig.saaswaf premium edition : sp_wsm_waf_premiumsaaswaf enterprise edition : sp_wsm_waf_enterprisesaaswaf ultimate edition : sp_wsm_waf_ultimateclbwaf premium edition : sp_wsm_waf_premium_clbclbwaf enterprise edition : sp_wsm_waf_enterprise_clbclbwaf ultimate edition : sp_wsm_waf_ultimate_clbNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"pid": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Pid of waf product (corresponding to a pricing formula), can be used to query the pricing model through pid billing.saaswaf premium edition : 1000827saaswaf enterprise edition : 1000830saaswaf ultimate edition : 1000832clbwaf premium edition : 1001150clbwaf enterprise edition : 1001152clbwaf ultimate edition : 1001154Note: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"instance_name": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Waf instance nameNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"auto_renew_flag": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Auto renew flag, 1:enable, 0:disableNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"real_region": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "Region infoNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"label_types": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Optional:    true,
										Description: "Label array of billing itemsaaswaf premium edition : sv_wsm_waf_package_premiumsaaswaf enterprise edition : sv_wsm_waf_package_enterprisesaaswaf ultimate edition : sv_wsm_waf_package_ultimatesaaswaf premium edition and non chinese mainland : sv_wsm_waf_package_premium_intlsaaswaf enterprise edition and non chinese mainland : sv_wsm_waf_package_enterprise_intlsaaswaf ultimate edition and non chinese mainland : sv_wsm_waf_package_ultimate _intlsaaswaf business expansion package : sv_wsm_waf_qps_epsaaswaf domain expansion package : sv_wsm_waf_domainclbwaf premium edition : sv_wsm_waf_package_premium_clbclbwaf enterprise edition : sv_wsm_waf_package_enterprise_clbclbwaf ultimate edition : sv_wsm_waf_package_ultimate_clbclbwaf premium edition and non chinese mainland : sv_wsm_waf_package_premium_clb_intlclbwaf enterprise edition and non chinese mainland : sv_wsm_waf_package_premium_clb_intlclbwaf ultimate edition and non chinese mainland : sv_wsm_waf_package_ultimate_clb _intlclbwaf business expansion package : sv_wsm_waf_qps_ep_clbclbwaf domain expansion package : sv_wsm_waf_domain_clbNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"label_counts": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Optional:    true,
										Description: "Label count of billing item, generally corresponds one-to-one with SvLabelTypeNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"cur_deadline": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Instance expiration time, used when changing configurationNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
									"instance_id": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Purchase bot or API security for the instanceNote: This field may return null, indicating that a valid value cannot be obtained.",
									},
								},
							},
						},
						"goods_category_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Order type ID, used to uniquely identify one scenario of a business (three scenarios in total: new purchase, configuration change, renewal)saaswaf premium edition: 102375(new purchase), 102376(renewal), 102377(configuration change)saaswaf enterprise edition : 102378(new purchase), 102379(renewal), 102380(configuration change)saaswaf ultimate edition : 102369(new purchase), 102370(renewal), 102371(configuration change)clbwaf premium edition: 101198(new purchase), 101199(renewal), 101200(configuration change)clbwaf enterprise edition: 101204(new purchase), 101205(renewal), 101206(configuration change)clbwaf ultimate edition: 101201(new purchase), 101202(renewal), 101203(configuration change)Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Instance region id1 represents purchasing resources from mainland China, and 2 represents purchasing of non Chinese Mainland resourcesNote: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudWafInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_instance.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = waf.NewGenerateDealsAndPayNewRequest()
		response   = waf.NewGenerateDealsAndPayNewResponse()
		instanceId string
	)
	if v, ok := d.GetOk("goods"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			goodNews := waf.GoodNews{}
			if v, ok := dMap["goods_num"]; ok {
				goodNews.GoodsNum = helper.IntInt64(v.(int))
			}
			if goodsDetailMap, ok := helper.InterfaceToMap(dMap, "goods_detail"); ok {
				goodsDetailNew := waf.GoodsDetailNew{}
				if v, ok := goodsDetailMap["time_span"]; ok {
					goodsDetailNew.TimeSpan = helper.IntInt64(v.(int))
				}
				if v, ok := goodsDetailMap["time_unit"]; ok {
					goodsDetailNew.TimeUnit = helper.String(v.(string))
				}
				if v, ok := goodsDetailMap["sub_product_code"]; ok {
					goodsDetailNew.SubProductCode = helper.String(v.(string))
				}
				if v, ok := goodsDetailMap["pid"]; ok {
					goodsDetailNew.Pid = helper.IntInt64(v.(int))
				}
				if v, ok := goodsDetailMap["instance_name"]; ok {
					goodsDetailNew.InstanceName = helper.String(v.(string))
				}
				if v, ok := goodsDetailMap["auto_renew_flag"]; ok {
					goodsDetailNew.AutoRenewFlag = helper.IntInt64(v.(int))
				}
				if v, ok := goodsDetailMap["real_region"]; ok {
					goodsDetailNew.RealRegion = helper.IntInt64(v.(int))
				}
				if v, ok := goodsDetailMap["label_types"]; ok {
					labelTypesSet := v.(*schema.Set).List()
					for i := range labelTypesSet {
						labelTypes := labelTypesSet[i].(string)
						goodsDetailNew.LabelTypes = append(goodsDetailNew.LabelTypes, &labelTypes)
					}
				}
				if v, ok := goodsDetailMap["label_counts"]; ok {
					labelCountsSet := v.(*schema.Set).List()
					for i := range labelCountsSet {
						labelCounts := labelCountsSet[i].(int)
						goodsDetailNew.LabelCounts = append(goodsDetailNew.LabelCounts, helper.IntInt64(labelCounts))
					}
				}
				if v, ok := goodsDetailMap["cur_deadline"]; ok {
					goodsDetailNew.CurDeadline = helper.String(v.(string))
				}
				if v, ok := goodsDetailMap["instance_id"]; ok {
					goodsDetailNew.InstanceId = helper.String(v.(string))
				}
				goodNews.GoodsDetail = &goodsDetailNew
			}
			if v, ok := dMap["goods_category_id"]; ok {
				goodNews.GoodsCategoryId = helper.IntInt64(v.(int))
			}
			if v, ok := dMap["region_id"]; ok {
				goodNews.RegionId = helper.IntInt64(v.(int))
			}
			request.Goods = append(request.Goods, &goodNews)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().GenerateDealsAndPayNew(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create waf instance failed, reason:%+v", logId, err)
		return err
	}

	instanceId = *response.Response.InstanceId
	d.SetId(instanceId)

	return resourceTencentCloudWafInstanceRead(d, meta)
}

func resourceTencentCloudWafInstanceRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_instance.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}

	instanceId := d.Id()

	instance, err := service.DescribeWafInstanceById(ctx, instanceId)
	if err != nil {
		return err
	}

	if instance == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `WafInstance` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if instance.Goods != nil {
		goodsList := []interface{}{}
		for _, goods := range instance.Goods {
			goodsMap := map[string]interface{}{}

			if instance.Goods.GoodsNum != nil {
				goodsMap["goods_num"] = instance.Goods.GoodsNum
			}

			if instance.Goods.GoodsDetail != nil {
				goodsDetailMap := map[string]interface{}{}

				if instance.Goods.GoodsDetail.TimeSpan != nil {
					goodsDetailMap["time_span"] = instance.Goods.GoodsDetail.TimeSpan
				}

				if instance.Goods.GoodsDetail.TimeUnit != nil {
					goodsDetailMap["time_unit"] = instance.Goods.GoodsDetail.TimeUnit
				}

				if instance.Goods.GoodsDetail.SubProductCode != nil {
					goodsDetailMap["sub_product_code"] = instance.Goods.GoodsDetail.SubProductCode
				}

				if instance.Goods.GoodsDetail.Pid != nil {
					goodsDetailMap["pid"] = instance.Goods.GoodsDetail.Pid
				}

				if instance.Goods.GoodsDetail.InstanceName != nil {
					goodsDetailMap["instance_name"] = instance.Goods.GoodsDetail.InstanceName
				}

				if instance.Goods.GoodsDetail.AutoRenewFlag != nil {
					goodsDetailMap["auto_renew_flag"] = instance.Goods.GoodsDetail.AutoRenewFlag
				}

				if instance.Goods.GoodsDetail.RealRegion != nil {
					goodsDetailMap["real_region"] = instance.Goods.GoodsDetail.RealRegion
				}

				if instance.Goods.GoodsDetail.LabelTypes != nil {
					goodsDetailMap["label_types"] = instance.Goods.GoodsDetail.LabelTypes
				}

				if instance.Goods.GoodsDetail.LabelCounts != nil {
					goodsDetailMap["label_counts"] = instance.Goods.GoodsDetail.LabelCounts
				}

				if instance.Goods.GoodsDetail.CurDeadline != nil {
					goodsDetailMap["cur_deadline"] = instance.Goods.GoodsDetail.CurDeadline
				}

				if instance.Goods.GoodsDetail.InstanceId != nil {
					goodsDetailMap["instance_id"] = instance.Goods.GoodsDetail.InstanceId
				}

				goodsMap["goods_detail"] = []interface{}{goodsDetailMap}
			}

			if instance.Goods.GoodsCategoryId != nil {
				goodsMap["goods_category_id"] = instance.Goods.GoodsCategoryId
			}

			if instance.Goods.RegionId != nil {
				goodsMap["region_id"] = instance.Goods.RegionId
			}

			goodsList = append(goodsList, goodsMap)
		}

		_ = d.Set("goods", goodsList)

	}

	return nil
}

func resourceTencentCloudWafInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_instance.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		modifyInstanceNameRequest  = waf.NewModifyInstanceNameRequest()
		modifyInstanceNameResponse = waf.NewModifyInstanceNameResponse()
	)

	instanceId := d.Id()

	request.InstanceId = &instanceId

	immutableArgs := []string{"goods"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseWafClient().ModifyInstanceName(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update waf instance failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudWafInstanceRead(d, meta)
}

func resourceTencentCloudWafInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_waf_instance.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := WafService{client: meta.(*TencentCloudClient).apiV3Conn}
	instanceId := d.Id()

	if err := service.DeleteWafInstanceById(ctx, instanceId); err != nil {
		return err
	}

	return nil
}
