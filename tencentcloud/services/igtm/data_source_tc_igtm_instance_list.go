package igtm

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	igtmv20231024 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/igtm/v20231024"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudIgtmInstanceList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudIgtmInstanceListRead,
		Schema: map[string]*schema.Schema{
			"filters": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Filter conditions.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter field name, supported list as follows:\n- InstanceId: IGTM instance ID.\n- Domain: IGTM instance domain.\n- MonitorId: Monitor ID.\n- PoolId: Pool ID. This is a required parameter, not passing it will cause interface query failure.",
						},
						"value": {
							Type:        schema.TypeSet,
							Required:    true,
							Description: "Filter field value.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"fuzzy": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether to enable fuzzy query, only supports filter field name as domain.\nWhen fuzzy query is enabled, maximum Value length is 1, otherwise maximum Value length is 5. (Reserved field, not currently used).",
						},
					},
				},
			},

			"instance_set": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Instance list.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance ID.",
						},
						"instance_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance name.",
						},
						"resource_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Resource ID.",
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Business domain.",
						},
						"access_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cname domain access method\nCUSTOM: Custom access domain\nSYSTEM: System access domain.",
						},
						"access_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access domain.",
						},
						"access_sub_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Access subdomain.",
						},
						"global_ttl": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Global record expiration time.",
						},
						"package_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Package type\nFREE: Free version\nSTANDARD: Standard version\nULTIMATE: Ultimate version.",
						},
						"working_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance running status\nNORMAL: Healthy\nFAULTY: At risk\nDOWN: Down\nUNKNOWN: Unknown.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance status, ENABLED: Normal, DISABLED: Disabled.",
						},
						"is_cname_configured": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether cname access: true accessed; false not accessed.",
						},
						"remark": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark.",
						},
						"strategy_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Strategy count.",
						},
						"address_pool_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bound address pool count.",
						},
						"monitor_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Bound monitor count.",
						},
						"pool_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Address pool ID.",
						},
						"pool_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Address pool name.",
						},
						"created_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance creation time.",
						},
						"updated_on": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance update time.",
						},
					},
				},
			},

			"system_access_enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether system domain access is supported: true supported; false not supported.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudIgtmInstanceListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_igtm_instance_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = IgtmService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*igtmv20231024.ResourceFilter, 0, len(filtersSet))
		for _, item := range filtersSet {
			filtersMap := item.(map[string]interface{})
			resourceFilter := igtmv20231024.ResourceFilter{}
			if v, ok := filtersMap["name"].(string); ok && v != "" {
				resourceFilter.Name = helper.String(v)
			}

			if v, ok := filtersMap["value"]; ok {
				valueSet := v.(*schema.Set).List()
				for i := range valueSet {
					value := valueSet[i].(string)
					resourceFilter.Value = append(resourceFilter.Value, helper.String(value))
				}
			}

			if v, ok := filtersMap["fuzzy"].(bool); ok {
				resourceFilter.Fuzzy = helper.Bool(v)
			}
			tmpSet = append(tmpSet, &resourceFilter)
		}

		paramMap["Filters"] = tmpSet
	}

	var (
		respData            []*igtmv20231024.Instance
		systemAccessEnabled *bool
	)
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, saEnabled, e := service.DescribeIgtmInstanceListByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		systemAccessEnabled = saEnabled
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	instanceSetList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, instanceSet := range respData {
			instanceSetMap := map[string]interface{}{}
			if instanceSet.InstanceId != nil {
				instanceSetMap["instance_id"] = instanceSet.InstanceId
			}

			if instanceSet.InstanceName != nil {
				instanceSetMap["instance_name"] = instanceSet.InstanceName
			}

			if instanceSet.ResourceId != nil {
				instanceSetMap["resource_id"] = instanceSet.ResourceId
			}

			if instanceSet.Domain != nil {
				instanceSetMap["domain"] = instanceSet.Domain
			}

			if instanceSet.AccessType != nil {
				instanceSetMap["access_type"] = instanceSet.AccessType
			}

			if instanceSet.AccessDomain != nil {
				instanceSetMap["access_domain"] = instanceSet.AccessDomain
			}

			if instanceSet.AccessSubDomain != nil {
				instanceSetMap["access_sub_domain"] = instanceSet.AccessSubDomain
			}

			if instanceSet.GlobalTtl != nil {
				instanceSetMap["global_ttl"] = instanceSet.GlobalTtl
			}

			if instanceSet.PackageType != nil {
				instanceSetMap["package_type"] = instanceSet.PackageType
			}

			if instanceSet.WorkingStatus != nil {
				instanceSetMap["working_status"] = instanceSet.WorkingStatus
			}

			if instanceSet.Status != nil {
				instanceSetMap["status"] = instanceSet.Status
			}

			if instanceSet.IsCnameConfigured != nil {
				instanceSetMap["is_cname_configured"] = instanceSet.IsCnameConfigured
			}

			if instanceSet.Remark != nil {
				instanceSetMap["remark"] = instanceSet.Remark
			}

			if instanceSet.StrategyNum != nil {
				instanceSetMap["strategy_num"] = instanceSet.StrategyNum
			}

			if instanceSet.AddressPoolNum != nil {
				instanceSetMap["address_pool_num"] = instanceSet.AddressPoolNum
			}

			if instanceSet.MonitorNum != nil {
				instanceSetMap["monitor_num"] = instanceSet.MonitorNum
			}

			if instanceSet.PoolId != nil {
				instanceSetMap["pool_id"] = instanceSet.PoolId
			}

			if instanceSet.PoolName != nil {
				instanceSetMap["pool_name"] = instanceSet.PoolName
			}

			if instanceSet.CreatedOn != nil {
				instanceSetMap["created_on"] = instanceSet.CreatedOn
			}

			if instanceSet.UpdatedOn != nil {
				instanceSetMap["updated_on"] = instanceSet.UpdatedOn
			}

			instanceSetList = append(instanceSetList, instanceSetMap)
		}

		_ = d.Set("instance_set", instanceSetList)
	}

	if systemAccessEnabled != nil {
		_ = d.Set("system_access_enabled", systemAccessEnabled)
	}

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
