package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cynosdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cynosdb/v20190107"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudCynosdbResourcePackageList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCynosdbResourcePackageListRead,
		Schema: map[string]*schema.Schema{
			"package_id": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Resource Package Unique ID.",
			},
			"package_name": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Resource Package Name.",
			},
			"package_type": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Resource package type CCU - Compute resource package, DISK - Storage resource package.",
			},
			"package_region": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Resource package usage region China - common in mainland China, overseas - common in Hong Kong, Macao, Taiwan, and overseas.",
			},
			"status": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Resource package status creating - creating; Using - In use; Expired - has expired; Normal_ Finish - used up; Apply_ Refund - Applying for a refund; Refund - The fee has been refunded.",
			},
			"order_by": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Sorting conditions supported: startTime - effective time, expireTime - expiration time, packageUsedSpec - usage capacity, and packageTotalSpec - total storage capacity. Arrange in array order;.",
			},
			"order_direction": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Sort by, DESC Descending, ASC Ascending.",
			},
			"resource_package_list": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Resource package details note: This field may return null, indicating that a valid value cannot be obtained.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "AppID note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"package_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource Package Unique ID Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"package_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource package name note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"package_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource package type CCU - Compute resource package, DISK - Store resource package Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"package_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource package is used in China, which is commonly used in mainland China, and in overseas, which is commonly used in Hong Kong, Macao, Taiwan, and overseas. Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource package status creating - creating; Using - In use; Expired - has expired; Normal_ Finish - used up; Apply_ Refund - Applying for a refund; Refund - The fee has been refunded. Note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"package_total_spec": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Attention to the total amount of resource packages: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"package_used_spec": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Resource package usage note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"has_quota": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Resource package usage note: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"bind_instance_infos": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Note for binding instance information: This field may return null, indicating that a valid value cannot be obtained.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance ID.",
									},
									"instance_region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Region of instance.",
									},
									"instance_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Instance type.",
									},
								},
							},
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Effective time: July 1st, 2022 00:00:00 Attention: This field may return null, indicating that a valid value cannot be obtained.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Expiration time: August 1st, 2022 00:00:00 Attention: This field may return null, indicating that a valid value cannot be obtained.",
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

func dataSourceTencentCloudCynosdbResourcePackageListRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_cynosdb_resource_package_list.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId      = getLogId(contextNil)
		ctx        = context.WithValue(context.TODO(), logIdKey, logId)
		service    = CynosdbService{client: meta.(*TencentCloudClient).apiV3Conn}
		detailList []*cynosdb.Package
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("package_id"); ok {
		packageIdSet := v.(*schema.Set).List()
		paramMap["PackageId"] = helper.InterfacesStringsPoint(packageIdSet)
	}

	if v, ok := d.GetOk("package_name"); ok {
		packageNameSet := v.(*schema.Set).List()
		paramMap["PackageName"] = helper.InterfacesStringsPoint(packageNameSet)
	}

	if v, ok := d.GetOk("package_type"); ok {
		packageTypeSet := v.(*schema.Set).List()
		paramMap["PackageType"] = helper.InterfacesStringsPoint(packageTypeSet)
	}

	if v, ok := d.GetOk("package_region"); ok {
		packageRegionSet := v.(*schema.Set).List()
		paramMap["PackageRegion"] = helper.InterfacesStringsPoint(packageRegionSet)
	}

	if v, ok := d.GetOk("status"); ok {
		statusSet := v.(*schema.Set).List()
		paramMap["Status"] = helper.InterfacesStringsPoint(statusSet)
	}

	if v, ok := d.GetOk("order_by"); ok {
		orderBySet := v.(*schema.Set).List()
		paramMap["OrderBy"] = helper.InterfacesStringsPoint(orderBySet)
	}

	if v, ok := d.GetOk("order_direction"); ok {
		paramMap["OrderDirection"] = helper.String(v.(string))
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCynosdbResourcePackageListByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		detailList = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(detailList))

	if detailList != nil {
		tmpList := []interface{}{}
		for _, detail := range detailList {
			detailMap := map[string]interface{}{}
			if detail.AppId != nil {
				detailMap["app_id"] = detail.AppId
			}

			if detail.PackageId != nil {
				detailMap["package_id"] = detail.PackageId
			}

			if detail.PackageName != nil {
				detailMap["package_name"] = detail.PackageName
			}

			if detail.PackageType != nil {
				detailMap["package_type"] = detail.PackageType
			}

			if detail.PackageRegion != nil {
				detailMap["package_region"] = detail.PackageRegion
			}

			if detail.Status != nil {
				detailMap["status"] = detail.Status
			}

			if detail.PackageTotalSpec != nil {
				detailMap["package_total_spec"] = detail.PackageTotalSpec
			}

			if detail.PackageUsedSpec != nil {
				detailMap["package_used_spec"] = detail.PackageUsedSpec
			}

			if detail.HasQuota != nil {
				detailMap["has_quota"] = detail.HasQuota
			}

			if detail.BindInstanceInfos != nil {
				insList := []interface{}{}
				for _, instanceInfo := range detail.BindInstanceInfos {
					insMap := map[string]interface{}{}
					if instanceInfo.InstanceId != nil {
						insMap["instance_id"] = instanceInfo.InstanceId
					}

					if instanceInfo.InstanceRegion != nil {
						insMap["instance_region"] = instanceInfo.InstanceRegion
					}

					if instanceInfo.InstanceType != nil {
						insMap["instance_type"] = instanceInfo.InstanceType
					}
					insList = append(insList, insMap)
				}

				detailMap["bind_instance_infos"] = insList
			}

			if detail.StartTime != nil {
				detailMap["start_time"] = detail.StartTime
			}

			if detail.ExpireTime != nil {
				detailMap["expire_time"] = detail.ExpireTime
			}

			tmpList = append(tmpList, detailMap)
		}

		_ = d.Set("resource_package_list", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
