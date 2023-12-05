package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	emr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/emr/v20190103"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudEmrCvmQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEmrCvmQuotaRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "EMR cluster ID.",
			},

			"zone_id": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Zone ID.",
			},

			"post_paid_quota_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Postpaid quota list Note: This field may return null, indicating that no valid value can be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"used_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used quota Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"remaining_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Residual quota Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"total_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total quota Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available area Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},

			"spot_paid_quota_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Biding instance quota list Note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"used_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Used quota Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"remaining_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Residual quota Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"total_quota": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total quota Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"zone": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Available area Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
					},
				},
			},

			"eks_quota_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Eks quota Note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The specifications of the marketable resource are as follows: `TASK`, `CORE`, `MASTER`, `ROUTER`.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cpu cores.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Memory quantity (unit: GB).",
						},
						"number": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Specifies the maximum number of resources that can be applied for.",
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

func dataSourceTencentCloudEmrCvmQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_emr_cvm_quota.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var clusterId string

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("cluster_id"); ok {
		clusterId = v.(string)
		paramMap["ClusterId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOkExists("zone_id"); ok {
		paramMap["ZoneId"] = helper.IntInt64(v.(int))
	}

	service := EMRService{client: meta.(*TencentCloudClient).apiV3Conn}

	var cvmQuota *emr.DescribeCvmQuotaResponseParams
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEmrCvmQuotaByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		cvmQuota = result
		return nil
	})
	if err != nil {
		return err
	}

	//ids := make([]string, 0, len(postPaidQuotaSet))
	tmpList := make([]map[string]interface{}, 0)

	if cvmQuota.PostPaidQuotaSet != nil {
		tmpList := make([]map[string]interface{}, 0, len(cvmQuota.PostPaidQuotaSet))

		for _, quotaEntity := range cvmQuota.PostPaidQuotaSet {
			quotaEntityMap := map[string]interface{}{}

			if quotaEntity.UsedQuota != nil {
				quotaEntityMap["used_quota"] = quotaEntity.UsedQuota
			}

			if quotaEntity.RemainingQuota != nil {
				quotaEntityMap["remaining_quota"] = quotaEntity.RemainingQuota
			}

			if quotaEntity.TotalQuota != nil {
				quotaEntityMap["total_quota"] = quotaEntity.TotalQuota
			}

			if quotaEntity.Zone != nil {
				quotaEntityMap["zone"] = quotaEntity.Zone
			}

			tmpList = append(tmpList, quotaEntityMap)
		}
		_ = d.Set("post_paid_quota_set", tmpList)
	}

	if cvmQuota.SpotPaidQuotaSet != nil {
		tmpList := make([]map[string]interface{}, 0, len(cvmQuota.SpotPaidQuotaSet))

		for _, quotaEntity := range cvmQuota.SpotPaidQuotaSet {
			quotaEntityMap := map[string]interface{}{}

			if quotaEntity.UsedQuota != nil {
				quotaEntityMap["used_quota"] = quotaEntity.UsedQuota
			}

			if quotaEntity.RemainingQuota != nil {
				quotaEntityMap["remaining_quota"] = quotaEntity.RemainingQuota
			}

			if quotaEntity.TotalQuota != nil {
				quotaEntityMap["total_quota"] = quotaEntity.TotalQuota
			}

			if quotaEntity.Zone != nil {
				quotaEntityMap["zone"] = quotaEntity.Zone
			}
			tmpList = append(tmpList, quotaEntityMap)
		}

		_ = d.Set("spot_paid_quota_set", tmpList)
	}

	if cvmQuota.EksQuotaSet != nil {
		tmpList := make([]map[string]interface{}, 0, len(cvmQuota.EksQuotaSet))

		for _, podSaleSpec := range cvmQuota.EksQuotaSet {
			podSaleSpecMap := map[string]interface{}{}

			if podSaleSpec.NodeType != nil {
				podSaleSpecMap["node_type"] = podSaleSpec.NodeType
			}

			if podSaleSpec.Cpu != nil {
				podSaleSpecMap["cpu"] = podSaleSpec.Cpu
			}

			if podSaleSpec.Memory != nil {
				podSaleSpecMap["memory"] = podSaleSpec.Memory
			}

			if podSaleSpec.Number != nil {
				podSaleSpecMap["number"] = podSaleSpec.Number
			}

			tmpList = append(tmpList, podSaleSpecMap)
		}

		_ = d.Set("eks_quota_set", tmpList)
	}

	d.SetId(clusterId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
