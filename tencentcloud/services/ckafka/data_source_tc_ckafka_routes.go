package ckafka

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ckafka "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ckafka/v20190819"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCkafkaRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudCkafkaRoutesRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the ckafka instance.",
			},

			"route_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Route id, used to query a specific route.",
			},

			"main_route_flag": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to display the main route. When set to true, the main route created at instance creation will be additionally returned.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			"routers": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of ckafka routes. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Instance access type. 0: PLAINTEXT, 1: SASL_PLAINTEXT, 2: SSL, 3: SASL_SSL.",
						},
						"route_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Route id.",
						},
						"vip_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Routing network type. 3: vpc routing, 7: internal support routing, 1: public network routing.",
						},
						"vip_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Virtual IP list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Virtual IP.",
									},
									"vport": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Virtual port.",
									},
								},
							},
						},
						"domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Domain.",
						},
						"domain_port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Domain port.",
						},
						"delete_timestamp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Delete timestamp.",
						},
						"subnet": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet id.",
						},
						"broker_vip_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Virtual IP list (1 to 1 broker nodes).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vip": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Virtual IP.",
									},
									"vport": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Virtual port.",
									},
								},
							},
						},
						"vpc_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Vpc id.",
						},
						"note": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Remark.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Route status. 1: creating, 2: created, 3: create failed, 4: deleting, 6: delete failed.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudCkafkaRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_ckafka_routes.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId      = tccommon.GetLogId(nil)
		ctx        = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service    = CkafkaService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		instanceId = d.Get("instance_id").(string)
	)

	var routeId *int64
	if v, ok := d.GetOkExists("route_id"); ok {
		routeId = helper.IntInt64(v.(int))
	}

	var mainRouteFlag *bool
	if v, ok := d.GetOkExists("main_route_flag"); ok {
		mainRouteFlag = helper.Bool(v.(bool))
	}

	var respData *ckafka.RouteResponse
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCkafkaRouteByFilter(ctx, instanceId, routeId, mainRouteFlag)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil || len(result.Routers) == 0 {
			log.Printf("[DATASOURCE] read empty, skip SetId, instance_id=%s", instanceId)
			return resource.NonRetryableError(fmt.Errorf("ckafka routes is empty, instance_id=%s", instanceId))
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	routersList := make([]map[string]interface{}, 0, len(respData.Routers))
	for _, route := range respData.Routers {
		routeMap := map[string]interface{}{}
		if route.AccessType != nil {
			routeMap["access_type"] = route.AccessType
		}

		if route.RouteId != nil {
			routeMap["route_id"] = route.RouteId
		}

		if route.VipType != nil {
			routeMap["vip_type"] = route.VipType
		}

		vipList := make([]map[string]interface{}, 0, len(route.VipList))
		for _, vip := range route.VipList {
			vipMap := map[string]interface{}{}
			if vip.Vip != nil {
				vipMap["vip"] = vip.Vip
			}
			if vip.Vport != nil {
				vipMap["vport"] = vip.Vport
			}
			vipList = append(vipList, vipMap)
		}
		routeMap["vip_list"] = vipList

		if route.Domain != nil {
			routeMap["domain"] = route.Domain
		}

		if route.DomainPort != nil {
			routeMap["domain_port"] = route.DomainPort
		}

		if route.DeleteTimestamp != nil {
			routeMap["delete_timestamp"] = route.DeleteTimestamp
		}

		if route.Subnet != nil {
			routeMap["subnet"] = route.Subnet
		}

		brokerVipList := make([]map[string]interface{}, 0, len(route.BrokerVipList))
		for _, brokerVip := range route.BrokerVipList {
			brokerVipMap := map[string]interface{}{}
			if brokerVip.Vip != nil {
				brokerVipMap["vip"] = brokerVip.Vip
			}
			if brokerVip.Vport != nil {
				brokerVipMap["vport"] = brokerVip.Vport
			}
			brokerVipList = append(brokerVipList, brokerVipMap)
		}
		routeMap["broker_vip_list"] = brokerVipList

		if route.VpcId != nil {
			routeMap["vpc_id"] = route.VpcId
		}

		if route.Note != nil {
			routeMap["note"] = route.Note
		}

		if route.Status != nil {
			routeMap["status"] = route.Status
		}

		routersList = append(routersList, routeMap)
	}

	_ = d.Set("routers", routersList)

	d.SetId(helper.BuildToken())
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), routersList); e != nil {
			return e
		}
	}

	return nil
}
