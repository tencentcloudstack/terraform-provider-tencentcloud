/*
Use this data source to query detailed information of tse gateway_routes

Example Usage

```hcl
data "tencentcloud_tse_gateway_routes" "gateway_routes" {
  gateway_id = "gateway-xxxxxx"
  service_name = "serviceA"
  route_name = "123"
  filters {
		key = "name"
		value = "123"

  }
  }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tse/v20201207"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudTseGatewayRoutes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTseGatewayRoutesRead,
		Schema: map[string]*schema.Schema{
			"gateway_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Gateway ID.",
			},

			"service_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Service name.",
			},

			"route_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Route name.",
			},

			"filters": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Filter conditions, valid value:name, path, host, method, service, protocol.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter name.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Filter value.",
						},
					},
				},
			},

			"result": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Result.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"route_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Route list.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"i_d": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service ID.",
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service name.",
									},
									"methods": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Method list.",
									},
									"paths": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Path list.",
									},
									"hosts": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Host list.",
									},
									"protocols": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Computed:    true,
										Description: "Protocol list.",
									},
									"preserve_host": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to keep the host when forwarding to the backend.",
									},
									"https_redirect_status_code": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "Https redirection status code.",
									},
									"strip_path": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to strip path when forwarding to the backend.",
									},
									"created_time": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Created time.",
									},
									"force_https": {
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "Whether to enable forced HTTPS, no longer use.",
									},
									"service_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service name.",
									},
									"service_i_d": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Service ID.",
									},
									"destination_ports": {
										Type: schema.TypeSet,
										Elem: &schema.Schema{
											Type: schema.TypeInt,
										},
										Computed:    true,
										Description: "Destination port for Layer 4 matching.",
									},
									"headers": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The headers of route.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"key": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Key of header.",
												},
												"value": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Value of header.",
												},
											},
										},
									},
								},
							},
						},
						"total_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total count.",
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

func dataSourceTencentCloudTseGatewayRoutesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tse_gateway_routes.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("gateway_id"); ok {
		paramMap["GatewayId"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("service_name"); ok {
		paramMap["ServiceName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("route_name"); ok {
		paramMap["RouteName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("filters"); ok {
		filtersSet := v.([]interface{})
		tmpSet := make([]*tse.ListFilter, 0, len(filtersSet))

		for _, item := range filtersSet {
			listFilter := tse.ListFilter{}
			listFilterMap := item.(map[string]interface{})

			if v, ok := listFilterMap["key"]; ok {
				listFilter.Key = helper.String(v.(string))
			}
			if v, ok := listFilterMap["value"]; ok {
				listFilter.Value = helper.String(v.(string))
			}
			tmpSet = append(tmpSet, &listFilter)
		}
		paramMap["filters"] = tmpSet
	}

	service := TseService{client: meta.(*TencentCloudClient).apiV3Conn}

	var result []*tse.KongServiceRouteList

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTseGatewayRoutesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		result = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(result))
	if result != nil {
		kongServiceRouteListMap := map[string]interface{}{}

		if result.RouteList != nil {
			routeListList := []interface{}{}
			for _, routeList := range result.RouteList {
				routeListMap := map[string]interface{}{}

				if routeList.ID != nil {
					routeListMap["i_d"] = routeList.ID
				}

				if routeList.Name != nil {
					routeListMap["name"] = routeList.Name
				}

				if routeList.Methods != nil {
					routeListMap["methods"] = routeList.Methods
				}

				if routeList.Paths != nil {
					routeListMap["paths"] = routeList.Paths
				}

				if routeList.Hosts != nil {
					routeListMap["hosts"] = routeList.Hosts
				}

				if routeList.Protocols != nil {
					routeListMap["protocols"] = routeList.Protocols
				}

				if routeList.PreserveHost != nil {
					routeListMap["preserve_host"] = routeList.PreserveHost
				}

				if routeList.HttpsRedirectStatusCode != nil {
					routeListMap["https_redirect_status_code"] = routeList.HttpsRedirectStatusCode
				}

				if routeList.StripPath != nil {
					routeListMap["strip_path"] = routeList.StripPath
				}

				if routeList.CreatedTime != nil {
					routeListMap["created_time"] = routeList.CreatedTime
				}

				if routeList.ForceHttps != nil {
					routeListMap["force_https"] = routeList.ForceHttps
				}

				if routeList.ServiceName != nil {
					routeListMap["service_name"] = routeList.ServiceName
				}

				if routeList.ServiceID != nil {
					routeListMap["service_i_d"] = routeList.ServiceID
				}

				if routeList.DestinationPorts != nil {
					routeListMap["destination_ports"] = routeList.DestinationPorts
				}

				if routeList.Headers != nil {
					headersList := []interface{}{}
					for _, headers := range routeList.Headers {
						headersMap := map[string]interface{}{}

						if headers.Key != nil {
							headersMap["key"] = headers.Key
						}

						if headers.Value != nil {
							headersMap["value"] = headers.Value
						}

						headersList = append(headersList, headersMap)
					}

					routeListMap["headers"] = []interface{}{headersList}
				}

				routeListList = append(routeListList, routeListMap)
			}

			kongServiceRouteListMap["route_list"] = []interface{}{routeListList}
		}

		if result.TotalCount != nil {
			kongServiceRouteListMap["total_count"] = result.TotalCount
		}

		ids = append(ids, *result.GatewayId)
		_ = d.Set("result", kongServiceRouteListMap)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), kongServiceRouteListMap); e != nil {
			return e
		}
	}
	return nil
}
