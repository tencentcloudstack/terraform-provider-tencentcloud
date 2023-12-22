package waf

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	waf "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/waf/v20180125"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudWafInstanceQpsLimit() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudWafInstanceQpsLimitRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Unique ID of Instance.",
			},
			"type": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Instance type.",
			},
			"qps_data": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Qps info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"elastic_billing_default": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Elastic qps default value.",
						},
						"elastic_billing_min": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum elastic qps.",
						},
						"elastic_billing_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum elastic qps.",
						},
						"qps_extend_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum qps of extend package.",
						},
						"qps_extend_intl_max": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum qps of extend package for overseas.",
						},
					},
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudWafInstanceQpsLimitRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_waf_instance_qps_limit.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(tccommon.ContextNil)
		ctx        = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service    = WafService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		qpsData    *waf.QpsData
		instanceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("instance_id"); ok {
		paramMap["InstanceId"] = helper.String(v.(string))
		instanceId = v.(string)
	}

	if v, ok := d.GetOk("type"); ok {
		paramMap["Type"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeWafInstanceQpsLimitByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		qpsData = result
		return nil
	})

	if err != nil {
		return err
	}

	if qpsData != nil {
		tmqList := []interface{}{}
		qpsDataMap := map[string]interface{}{}

		if qpsData.ElasticBillingDefault != nil {
			qpsDataMap["elastic_billing_default"] = qpsData.ElasticBillingDefault
		}

		if qpsData.ElasticBillingMin != nil {
			qpsDataMap["elastic_billing_min"] = qpsData.ElasticBillingMin
		}

		if qpsData.ElasticBillingMax != nil {
			qpsDataMap["elastic_billing_max"] = qpsData.ElasticBillingMax
		}

		if qpsData.QPSExtendMax != nil {
			qpsDataMap["qps_extend_max"] = qpsData.QPSExtendMax
		}

		if qpsData.QPSExtendIntlMax != nil {
			qpsDataMap["qps_extend_intl_max"] = qpsData.QPSExtendIntlMax
		}

		tmqList = append(tmqList, qpsDataMap)
		_ = d.Set("qps_data", tmqList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
