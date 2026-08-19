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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

func DataSourceTencentCloudTeoOriginGroupHealthStatus() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoOriginGroupHealthStatusRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"lb_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Load balancer instance ID.",
			},

			"origin_group_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Origin group ID list. When not specified, the health status of all origin groups under the load balancer is returned by default.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"origin_group_health_status_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Health status list of origin groups under the load balancer.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"origin_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Origin group ID.",
						},
						"origin_health_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The health status of each origin in the origin group, which is comprehensively determined based on the results of all detection regions. If more than half of the regions determine the origin as unhealthy, the corresponding status is unhealthy; otherwise, it is healthy.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"origin": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Origin.",
									},
									"healthy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Origin health status. Valid values: Healthy, Unhealthy, Undetected.",
									},
								},
							},
						},
						"check_region_health_status": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Health status of origins under each health check region.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Health check region, ISO-3166-1 two-letter code.",
									},
									"healthy": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Health status of origins under a single health check region. Valid values: Healthy, Unhealthy, Undetected. When all origins in a single health check region are healthy, the status is healthy; otherwise, it is unhealthy.",
									},
									"origin_health_status": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "Origin health status under the health check region.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"origin": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Origin.",
												},
												"healthy": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Origin health status. Valid values: Healthy, Unhealthy, Undetected.",
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

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoOriginGroupHealthStatusRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_origin_group_health_status.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["ZoneId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("lb_instance_id"); ok {
		paramMap["LBInstanceId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("origin_group_ids"); ok {
		originGroupIdsSet := v.([]interface{})
		tmpSet := make([]*string, 0, len(originGroupIdsSet))
		for _, item := range originGroupIdsSet {
			originGroupId := item.(string)
			tmpSet = append(tmpSet, helper.String(originGroupId))
		}
		paramMap["OriginGroupIds"] = tmpSet
	}

	var respData *teov20220901.DescribeOriginGroupHealthStatusResponseParams
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTeoOriginGroupHealthStatusByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		log.Printf("[DATASOURCE] read empty, skip SetId")
		return reqErr
	}

	originGroupHealthStatusList := make([]map[string]interface{}, 0, len(respData.OriginGroupHealthStatusList))
	if respData.OriginGroupHealthStatusList != nil {
		for _, originGroupHealthStatusDetail := range respData.OriginGroupHealthStatusList {
			originGroupHealthStatusMap := map[string]interface{}{}
			if originGroupHealthStatusDetail.OriginGroupId != nil {
				originGroupHealthStatusMap["origin_group_id"] = originGroupHealthStatusDetail.OriginGroupId
			}

			if originGroupHealthStatusDetail.OriginHealthStatus != nil {
				originHealthStatusList := make([]map[string]interface{}, 0, len(originGroupHealthStatusDetail.OriginHealthStatus))
				for _, originHealthStatus := range originGroupHealthStatusDetail.OriginHealthStatus {
					originHealthStatusMap := map[string]interface{}{}
					if originHealthStatus.Origin != nil {
						originHealthStatusMap["origin"] = originHealthStatus.Origin
					}

					if originHealthStatus.Healthy != nil {
						originHealthStatusMap["healthy"] = originHealthStatus.Healthy
					}

					originHealthStatusList = append(originHealthStatusList, originHealthStatusMap)
				}

				originGroupHealthStatusMap["origin_health_status"] = originHealthStatusList
			}

			if originGroupHealthStatusDetail.CheckRegionHealthStatus != nil {
				checkRegionHealthStatusList := make([]map[string]interface{}, 0, len(originGroupHealthStatusDetail.CheckRegionHealthStatus))
				for _, checkRegionHealthStatus := range originGroupHealthStatusDetail.CheckRegionHealthStatus {
					checkRegionHealthStatusMap := map[string]interface{}{}
					if checkRegionHealthStatus.Region != nil {
						checkRegionHealthStatusMap["region"] = checkRegionHealthStatus.Region
					}

					if checkRegionHealthStatus.Healthy != nil {
						checkRegionHealthStatusMap["healthy"] = checkRegionHealthStatus.Healthy
					}

					if checkRegionHealthStatus.OriginHealthStatus != nil {
						regionOriginHealthStatusList := make([]map[string]interface{}, 0, len(checkRegionHealthStatus.OriginHealthStatus))
						for _, regionOriginHealthStatus := range checkRegionHealthStatus.OriginHealthStatus {
							regionOriginHealthStatusMap := map[string]interface{}{}
							if regionOriginHealthStatus.Origin != nil {
								regionOriginHealthStatusMap["origin"] = regionOriginHealthStatus.Origin
							}

							if regionOriginHealthStatus.Healthy != nil {
								regionOriginHealthStatusMap["healthy"] = regionOriginHealthStatus.Healthy
							}

							regionOriginHealthStatusList = append(regionOriginHealthStatusList, regionOriginHealthStatusMap)
						}

						checkRegionHealthStatusMap["origin_health_status"] = regionOriginHealthStatusList
					}

					checkRegionHealthStatusList = append(checkRegionHealthStatusList, checkRegionHealthStatusMap)
				}

				originGroupHealthStatusMap["check_region_health_status"] = checkRegionHealthStatusList
			}

			originGroupHealthStatusList = append(originGroupHealthStatusList, originGroupHealthStatusMap)
		}
	}

	_ = d.Set("origin_group_health_status_list", originGroupHealthStatusList)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), originGroupHealthStatusList); e != nil {
			return e
		}
	}

	return nil
}

func (me *TeoService) DescribeTeoOriginGroupHealthStatusByFilter(ctx context.Context, param map[string]interface{}) (ret *teov20220901.DescribeOriginGroupHealthStatusResponseParams, errRet error) {
	var (
		logId    = tccommon.GetLogId(ctx)
		request  = teov20220901.NewDescribeOriginGroupHealthStatusRequest()
		response = teov20220901.NewDescribeOriginGroupHealthStatusResponse()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "ZoneId" {
			request.ZoneId = v.(*string)
		}

		if k == "LBInstanceId" {
			request.LBInstanceId = v.(*string)
		}

		if k == "OriginGroupIds" {
			request.OriginGroupIds = v.([]*string)
		}
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())
		result, e := me.client.UseTeoV20220901Client().DescribeOriginGroupHealthStatus(request)
		if e != nil {
			return tccommon.RetryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}

		if result == nil || result.Response == nil {
			return resource.NonRetryableError(fmt.Errorf("Describe teo origin group health status failed, Response is nil."))
		}

		if len(result.Response.OriginGroupHealthStatusList) == 0 {
			return resource.NonRetryableError(fmt.Errorf("Describe teo origin group health status failed, OriginGroupHealthStatusList is empty."))
		}

		response = result
		return nil
	})

	if err != nil {
		errRet = err
		return
	}

	ret = response.Response
	return
}
