package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudGaapGroupAndStatisticsProxy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapGroupAndStatisticsProxyRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Project Id.",
			},

			"group_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Channel group information that can be counted.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Channel Group ID.",
						},
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Channel Group name.",
						},
						"proxy_set": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Channel list in the proxy group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"proxy_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Channel Id.",
									},
									"proxy_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Channel name.",
									},
									"listener_list": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: "listener list.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"listener_id": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "listener Id.",
												},
												"listener_name": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "listener name.",
												},
												"port": {
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "listened port.",
												},
												"protocol": {
													Type:        schema.TypeString,
													Computed:    true,
													Description: "Listener protocol type.",
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

func dataSourceTencentCloudGaapGroupAndStatisticsProxyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_gaap_group_and_statistics_proxy.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ProjectId"] = helper.IntUint64(v.(int))
	}

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

	var groupSet []*gaap.GroupStatisticsInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapGroupAndStatisticsProxyByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		groupSet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(groupSet))
	tmpList := make([]map[string]interface{}, 0, len(groupSet))

	if groupSet != nil {
		for _, groupStatisticsInfo := range groupSet {
			groupStatisticsInfoMap := map[string]interface{}{}

			if groupStatisticsInfo.GroupId != nil {
				groupStatisticsInfoMap["group_id"] = groupStatisticsInfo.GroupId
				ids = append(ids, *groupStatisticsInfo.GroupId)
			}

			if groupStatisticsInfo.GroupName != nil {
				groupStatisticsInfoMap["group_name"] = groupStatisticsInfo.GroupName
			}

			if groupStatisticsInfo.ProxySet != nil {
				proxySetList := []interface{}{}
				for _, proxySet := range groupStatisticsInfo.ProxySet {
					proxySetMap := map[string]interface{}{}

					if proxySet.ProxyId != nil {
						proxySetMap["proxy_id"] = proxySet.ProxyId
					}

					if proxySet.ProxyName != nil {
						proxySetMap["proxy_name"] = proxySet.ProxyName
					}

					if proxySet.ListenerList != nil {
						listenerListList := []interface{}{}
						for _, listenerList := range proxySet.ListenerList {
							listenerListMap := map[string]interface{}{}

							if listenerList.ListenerId != nil {
								listenerListMap["listener_id"] = listenerList.ListenerId
							}

							if listenerList.ListenerName != nil {
								listenerListMap["listener_name"] = listenerList.ListenerName
							}

							if listenerList.Port != nil {
								listenerListMap["port"] = listenerList.Port
							}

							if listenerList.Protocol != nil {
								listenerListMap["protocol"] = listenerList.Protocol
							}

							listenerListList = append(listenerListList, listenerListMap)
						}

						proxySetMap["listener_list"] = listenerListList
					}

					proxySetList = append(proxySetList, proxySetMap)
				}

				groupStatisticsInfoMap["proxy_set"] = proxySetList
			}

			tmpList = append(tmpList, groupStatisticsInfoMap)
		}

		_ = d.Set("group_set", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
