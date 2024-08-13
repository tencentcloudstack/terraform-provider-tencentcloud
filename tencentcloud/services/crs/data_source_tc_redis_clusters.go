package crs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	redis "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/redis/v20180412"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudRedisClusters() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudRedisClustersRead,
		Schema: map[string]*schema.Schema{
			"redis_cluster_ids": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Redis Cluster Ids.",
			},
			"status": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Cluster status: 1- In process, 2- Running, 3- Isolated.",
			},
			"project_ids": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Project Ids.",
			},
			"auto_renew_flag": {
				Optional:    true,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Renewal mode: 0- default state (manual renewal); 1- Automatic renewal; 2- Clearly stating that automatic renewal is not allowed.",
			},
			"cluster_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Cluster name.",
			},
			"dedicated_cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Dedicated cluster Id.",
			},
			// computed
			"resources": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: ".",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"app_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "User's Appid.",
						},
						"region_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Region Id.",
						},
						"zone_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "zone Id.",
						},
						"redis_cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Redis Cluster Id.",
						},
						"pay_mode": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Billing mode, 1-annual and monthly package, 0-quantity based billing.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project Id.",
						},
						"auto_renew_flag": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Renewal mode: 0- default state (manual renewal); 1- Automatic renewal; 2- Clearly stating that automatic renewal is not allowed.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster name.",
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance create time.",
						},
						"end_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance expiration time.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cluster status: 1- In process, 2- Running, 3- Isolated.",
						},
						"base_bundles": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Basic Control Resource Package.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_bundle_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource bundle name.",
									},
									"available_memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Saleable memory, unit: GB.",
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Resource bundle count.",
									},
								},
							},
						},
						"resource_bundles": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "List of Resource Packages.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_bundle_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Resource bundle name.",
									},
									"available_memory": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Saleable memory, unit: GB.",
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Resource bundle count.",
									},
								},
							},
						},
						"dedicated_cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated cluster Id.",
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

func dataSourceTencentCloudRedisClustersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_redis_clusters.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(tccommon.ContextNil)
		ctx     = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service = RedisService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("redis_cluster_ids"); ok {
		redisClusterIdsSet := v.(*schema.Set).List()
		paramMap["RedisClusterIds"] = helper.InterfacesStringsPoint(redisClusterIdsSet)
	}

	if v, ok := d.GetOk("status"); ok {
		statusList := make([]*int64, 0)
		statusSet := v.(*schema.Set).List()
		for i := range statusSet {
			status := statusSet[i].(int)
			statusList = append(statusList, helper.IntInt64(status))
		}

		paramMap["Status"] = statusList
	}

	if v, ok := d.GetOk("project_ids"); ok {
		projectIdsList := make([]*int64, 0)
		projectIdsSet := v.(*schema.Set).List()
		for i := range projectIdsSet {
			projectIds := projectIdsSet[i].(int)
			projectIdsList = append(projectIdsList, helper.IntInt64(projectIds))
		}

		paramMap["ProjectIds"] = projectIdsList
	}

	if v, ok := d.GetOk("auto_renew_flag"); ok {
		autoRenewFlagList := make([]*int64, 0)
		autoRenewFlagSet := v.(*schema.Set).List()
		for i := range autoRenewFlagSet {
			autoRenewFlag := autoRenewFlagSet[i].(int)
			autoRenewFlagList = append(autoRenewFlagList, helper.IntInt64(autoRenewFlag))
		}

		paramMap["AutoRenewFlag"] = autoRenewFlagList
	}

	if v, ok := d.GetOk("cluster_name"); ok {
		paramMap["ClusterName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("search_key"); ok {
		paramMap["SearchKey"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		paramMap["DedicatedClusterId"] = helper.String(v.(string))
	}

	var resources []*redis.CDCResource
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeRedisClustersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		resources = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(resources))
	tmpList := make([]map[string]interface{}, 0, len(resources))
	if resources != nil {
		for _, cDCResource := range resources {
			cDCResourceMap := map[string]interface{}{}
			if cDCResource.AppId != nil {
				cDCResourceMap["app_id"] = cDCResource.AppId
			}

			if cDCResource.RegionId != nil {
				cDCResourceMap["region_id"] = cDCResource.RegionId
			}

			if cDCResource.ZoneId != nil {
				cDCResourceMap["zone_id"] = cDCResource.ZoneId
			}

			if cDCResource.RedisClusterId != nil {
				cDCResourceMap["redis_cluster_id"] = cDCResource.RedisClusterId
			}

			if cDCResource.PayMode != nil {
				cDCResourceMap["pay_mode"] = cDCResource.PayMode
			}

			if cDCResource.ProjectId != nil {
				cDCResourceMap["project_id"] = cDCResource.ProjectId
			}

			if cDCResource.AutoRenewFlag != nil {
				cDCResourceMap["auto_renew_flag"] = cDCResource.AutoRenewFlag
			}

			if cDCResource.ClusterName != nil {
				cDCResourceMap["cluster_name"] = cDCResource.ClusterName
			}

			if cDCResource.StartTime != nil {
				cDCResourceMap["start_time"] = cDCResource.StartTime
			}

			if cDCResource.EndTime != nil {
				cDCResourceMap["end_time"] = cDCResource.EndTime
			}

			if cDCResource.Status != nil {
				cDCResourceMap["status"] = cDCResource.Status
			}

			if cDCResource.BaseBundles != nil {
				baseBundlesList := []interface{}{}
				for _, baseBundles := range cDCResource.BaseBundles {
					baseBundlesMap := map[string]interface{}{}
					if baseBundles.ResourceBundleName != nil {
						baseBundlesMap["resource_bundle_name"] = baseBundles.ResourceBundleName
					}

					if baseBundles.AvailableMemory != nil {
						baseBundlesMap["available_memory"] = baseBundles.AvailableMemory
					}

					if baseBundles.Count != nil {
						baseBundlesMap["count"] = baseBundles.Count
					}

					baseBundlesList = append(baseBundlesList, baseBundlesMap)
				}

				cDCResourceMap["base_bundles"] = baseBundlesList
			}

			if cDCResource.ResourceBundles != nil {
				resourceBundlesList := []interface{}{}
				for _, resourceBundles := range cDCResource.ResourceBundles {
					resourceBundlesMap := map[string]interface{}{}
					if resourceBundles.ResourceBundleName != nil {
						resourceBundlesMap["resource_bundle_name"] = resourceBundles.ResourceBundleName
					}

					if resourceBundles.AvailableMemory != nil {
						resourceBundlesMap["available_memory"] = resourceBundles.AvailableMemory
					}

					if resourceBundles.Count != nil {
						resourceBundlesMap["count"] = resourceBundles.Count
					}

					resourceBundlesList = append(resourceBundlesList, resourceBundlesMap)
				}

				cDCResourceMap["resource_bundles"] = resourceBundlesList
			}

			if cDCResource.DedicatedClusterId != nil {
				cDCResourceMap["dedicated_cluster_id"] = cDCResource.DedicatedClusterId
			}

			ids = append(ids, *cDCResource.DedicatedClusterId)
			tmpList = append(tmpList, cDCResourceMap)
		}

		_ = d.Set("resources", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
