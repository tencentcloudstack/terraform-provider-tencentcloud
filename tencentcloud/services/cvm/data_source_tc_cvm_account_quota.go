package cvm

import (
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCvmAccountQuota() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCvmAccountQuotaRead,

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Filter by availability zone, such as ap-guangzhou-3.",
			},
			"quota_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter by quota type. Valid values: PostPaidQuotaSet, PrePaidQuotaSet, SpotPaidQuotaSet, ImageQuotaSet, DisasterRecoverGroupQuotaSet.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// Computed
			"app_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "User AppId.",
			},
			"account_quota_overview": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Account quota overview.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Region.",
						},
						"account_quota": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Account quota details.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"post_paid_quota_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Post-paid quota list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Availability zone.",
												},
												"total_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Total quota.",
												},
												"used_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Used quota.",
												},
												"remaining_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Remaining quota.",
												},
											},
										},
									},
									"pre_paid_quota_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Pre-paid quota list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Availability zone.",
												},
												"total_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Total quota.",
												},
												"used_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Used quota.",
												},
												"remaining_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Remaining quota.",
												},
												"once_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Single purchase quota.",
												},
											},
										},
									},
									"spot_paid_quota_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Spot instance quota list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"zone": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Availability zone.",
												},
												"total_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Total quota.",
												},
												"used_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Used quota.",
												},
												"remaining_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Remaining quota.",
												},
											},
										},
									},
									"image_quota_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Image quota list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"total_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Total quota.",
												},
												"used_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Used quota.",
												},
											},
										},
									},
									"disaster_recover_group_quota_set": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Disaster recover group quota list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"group_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Group quota.",
												},
												"current_num": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Current number of groups.",
												},
												"cvm_in_host_group_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum instances in host group.",
												},
												"cvm_in_switch_group_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum instances in switch group.",
												},
												"cvm_in_rack_group_quota": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "Maximum instances in rack group.",
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
		},
	}
}

func dataSourceTencentCloudCvmAccountQuotaRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cvm_account_quota.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	request := cvm.NewDescribeAccountQuotaRequest()

	// Build filters
	var filters []*cvm.Filter

	if v, ok := d.GetOk("zone"); ok {
		zoneSet := v.(*schema.Set).List()
		if len(zoneSet) > 0 {
			filter := &cvm.Filter{
				Name:   helper.String("zone"),
				Values: helper.InterfacesStringsPoint(zoneSet),
			}
			filters = append(filters, filter)
		}
	}

	if v, ok := d.GetOk("quota_type"); ok {
		filter := &cvm.Filter{
			Name:   helper.String("quota-type"),
			Values: []*string{helper.String(v.(string))},
		}
		filters = append(filters, filter)
	}

	if len(filters) > 0 {
		request.Filters = filters
	}

	var response *cvm.DescribeAccountQuotaResponse
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().DescribeAccountQuota(request)
		if e != nil {
			return tccommon.RetryError(e, tccommon.InternalError)
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s read CVM account quota failed, reason:%+v", logId, err)
		return err
	}

	if response == nil || response.Response == nil {
		d.SetId("")
		return nil
	}

	// Set app_id
	if response.Response.AppId != nil {
		_ = d.Set("app_id", int(*response.Response.AppId))
	}

	// Set account_quota_overview
	if response.Response.AccountQuotaOverview != nil {
		accountQuotaOverviewList := []map[string]interface{}{
			flattenAccountQuotaOverview(response.Response.AccountQuotaOverview),
		}
		_ = d.Set("account_quota_overview", accountQuotaOverviewList)
	}

	d.SetId(helper.DataResourceIdsHash([]string{helper.UInt64ToStr(*response.Response.AppId)}))

	// Save to file if specified
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), map[string]interface{}{
			"app_id":                 response.Response.AppId,
			"account_quota_overview": response.Response.AccountQuotaOverview,
		}); e != nil {
			return e
		}
	}

	return nil
}

func flattenAccountQuotaOverview(overview *cvm.AccountQuotaOverview) map[string]interface{} {
	result := make(map[string]interface{})

	if overview.Region != nil {
		result["region"] = *overview.Region
	}

	if overview.AccountQuota != nil {
		accountQuotaList := []map[string]interface{}{
			flattenAccountQuota(overview.AccountQuota),
		}
		result["account_quota"] = accountQuotaList
	}

	return result
}

func flattenAccountQuota(quota *cvm.AccountQuota) map[string]interface{} {
	result := make(map[string]interface{})

	if quota.PostPaidQuotaSet != nil && len(quota.PostPaidQuotaSet) > 0 {
		result["post_paid_quota_set"] = flattenPostPaidQuotaSet(quota.PostPaidQuotaSet)
	}

	if quota.PrePaidQuotaSet != nil && len(quota.PrePaidQuotaSet) > 0 {
		result["pre_paid_quota_set"] = flattenPrePaidQuotaSet(quota.PrePaidQuotaSet)
	}

	if quota.SpotPaidQuotaSet != nil && len(quota.SpotPaidQuotaSet) > 0 {
		result["spot_paid_quota_set"] = flattenSpotPaidQuotaSet(quota.SpotPaidQuotaSet)
	}

	if quota.ImageQuotaSet != nil && len(quota.ImageQuotaSet) > 0 {
		result["image_quota_set"] = flattenImageQuotaSet(quota.ImageQuotaSet)
	}

	if quota.DisasterRecoverGroupQuotaSet != nil && len(quota.DisasterRecoverGroupQuotaSet) > 0 {
		result["disaster_recover_group_quota_set"] = flattenDisasterRecoverGroupQuotaSet(quota.DisasterRecoverGroupQuotaSet)
	}

	return result
}

func flattenPostPaidQuotaSet(quotaSet []*cvm.PostPaidQuota) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(quotaSet))

	for _, quota := range quotaSet {
		m := make(map[string]interface{})

		if quota.Zone != nil {
			m["zone"] = *quota.Zone
		}
		if quota.TotalQuota != nil {
			m["total_quota"] = int(*quota.TotalQuota)
		}
		if quota.UsedQuota != nil {
			m["used_quota"] = int(*quota.UsedQuota)
		}
		if quota.RemainingQuota != nil {
			m["remaining_quota"] = int(*quota.RemainingQuota)
		}

		result = append(result, m)
	}

	return result
}

func flattenPrePaidQuotaSet(quotaSet []*cvm.PrePaidQuota) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(quotaSet))

	for _, quota := range quotaSet {
		m := make(map[string]interface{})

		if quota.Zone != nil {
			m["zone"] = *quota.Zone
		}
		if quota.TotalQuota != nil {
			m["total_quota"] = int(*quota.TotalQuota)
		}
		if quota.UsedQuota != nil {
			m["used_quota"] = int(*quota.UsedQuota)
		}
		if quota.RemainingQuota != nil {
			m["remaining_quota"] = int(*quota.RemainingQuota)
		}
		if quota.OnceQuota != nil {
			m["once_quota"] = int(*quota.OnceQuota)
		}

		result = append(result, m)
	}

	return result
}

func flattenSpotPaidQuotaSet(quotaSet []*cvm.SpotPaidQuota) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(quotaSet))

	for _, quota := range quotaSet {
		m := make(map[string]interface{})

		if quota.Zone != nil {
			m["zone"] = *quota.Zone
		}
		if quota.TotalQuota != nil {
			m["total_quota"] = int(*quota.TotalQuota)
		}
		if quota.UsedQuota != nil {
			m["used_quota"] = int(*quota.UsedQuota)
		}
		if quota.RemainingQuota != nil {
			m["remaining_quota"] = int(*quota.RemainingQuota)
		}

		result = append(result, m)
	}

	return result
}

func flattenImageQuotaSet(quotaSet []*cvm.ImageQuota) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(quotaSet))

	for _, quota := range quotaSet {
		m := make(map[string]interface{})

		if quota.TotalQuota != nil {
			m["total_quota"] = int(*quota.TotalQuota)
		}
		if quota.UsedQuota != nil {
			m["used_quota"] = int(*quota.UsedQuota)
		}

		result = append(result, m)
	}

	return result
}

func flattenDisasterRecoverGroupQuotaSet(quotaSet []*cvm.DisasterRecoverGroupQuota) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(quotaSet))

	for _, quota := range quotaSet {
		m := make(map[string]interface{})

		if quota.GroupQuota != nil {
			m["group_quota"] = int(*quota.GroupQuota)
		}
		if quota.CurrentNum != nil {
			m["current_num"] = int(*quota.CurrentNum)
		}
		if quota.CvmInHostGroupQuota != nil {
			m["cvm_in_host_group_quota"] = int(*quota.CvmInHostGroupQuota)
		}
		if quota.CvmInSwitchGroupQuota != nil {
			m["cvm_in_switch_group_quota"] = int(*quota.CvmInSwitchGroupQuota)
		}
		if quota.CvmInRackGroupQuota != nil {
			m["cvm_in_rack_group_quota"] = int(*quota.CvmInRackGroupQuota)
		}

		result = append(result, m)
	}

	return result
}
