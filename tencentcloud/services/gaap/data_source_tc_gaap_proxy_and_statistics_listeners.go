package gaap

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudGaapProxyAndStatisticsListeners() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudGaapProxyAndStatisticsListenersRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Required:    true,
				Type:        schema.TypeInt,
				Description: "Project Id.",
			},

			"proxy_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "proxy information that can be counted.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"proxy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Proxy Id.",
						},
						"proxy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Proxy Name.",
						},
						"listener_list": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Listener List.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener Id.",
									},
									"listener_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Listener Name.",
									},
									"port": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "listerned port.",
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

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudGaapProxyAndStatisticsListenersRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_gaap_proxy_and_statistics_listeners.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("project_id"); v != nil {
		paramMap["ProjectId"] = helper.IntUint64(v.(int))
	}

	service := GaapService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var proxySet []*gaap.ProxySimpleInfo

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeGaapProxyAndStatisticsListenersByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		proxySet = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(proxySet))
	tmpList := make([]map[string]interface{}, 0, len(proxySet))

	if proxySet != nil {
		for _, proxySimpleInfo := range proxySet {
			proxySimpleInfoMap := map[string]interface{}{}

			if proxySimpleInfo.ProxyId != nil {
				proxySimpleInfoMap["proxy_id"] = proxySimpleInfo.ProxyId
			}

			if proxySimpleInfo.ProxyName != nil {
				proxySimpleInfoMap["proxy_name"] = proxySimpleInfo.ProxyName
			}

			if proxySimpleInfo.ListenerList != nil {
				listenerListList := []interface{}{}
				for _, listenerList := range proxySimpleInfo.ListenerList {
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

				proxySimpleInfoMap["listener_list"] = listenerListList
			}

			ids = append(ids, *proxySimpleInfo.ProxyId)
			tmpList = append(tmpList, proxySimpleInfoMap)
		}

		_ = d.Set("proxy_set", tmpList)
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
