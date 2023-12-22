package dayu

import (
	"context"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	dayu "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dayu/v20180709"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDayuDdosPolicyCases() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDayuDdosPolicyCasesRead,
		Schema: map[string]*schema.Schema{
			"resource_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(DAYU_RESOURCE_TYPE),
				Description:  "Type of the resource that the DDoS policy case works for, valid values are `bgpip`, `bgp`, `bgp-multip` and `net`.",
			},
			"scene_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the DDoS policy case to be query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of DDoS policy cases. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the resource that the DDoS policy case works for.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the DDoS policy case.",
						},
						"platform_types": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "Platform of the DDoS policy case.",
							},
							Computed:    true,
							Description: "Platform set of the DDoS policy case.",
						},
						"app_protocols": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "App protocol of the DDoS policy case.",
							},
							Computed:    true,
							Description: "App protocol set of the DDoS policy case.",
						},
						"app_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "App type of the DDoS policy case.",
						},
						"tcp_start_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start port of the TCP service.",
						},
						"tcp_end_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End port of the TCP service.",
						},
						"udp_start_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Start port of the UDP service.",
						},
						"udp_end_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "End port of the UDP service.",
						},
						"has_abroad": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicate whether the service involves overseas or not.",
						},
						"has_initiate_tcp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicate whether the service actively initiates TCP requests or not.",
						},
						"has_initiate_udp": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicate whether the actively initiate UDP requests or not.",
						},
						"has_vpn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Indicate whether the service involves VPN service or not.",
						},
						"peer_tcp_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The port that actively initiates TCP requests.",
						},
						"peer_udp_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The port that actively initiates UDP requests.",
						},
						"tcp_footprint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The fixed signature of TCP protocol load.",
						},
						"udp_footprint": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The fixed signature of TCP protocol load.",
						},
						"web_api_urls": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type:        schema.TypeString,
								Description: "Web API url.",
							},
							Computed:    true,
							Description: "Web API url set.",
						},
						"min_tcp_package_len": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The minimum length of TCP message package.",
						},
						"max_tcp_package_len": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The max length of TCP message package.",
						},
						"min_udp_package_len": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The minimum length of UDP message package.",
						},
						"max_udp_package_len": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The max length of UDP message package.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the DDoS policy case.",
						},
						"scene_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the DDoS policy case.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudDayuDdosPolicyCasesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dayu_ddos_policy_cases.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := DayuService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	resourceType := d.Get("resource_type").(string)
	sceneId := d.Get("scene_id").(string)

	var ddosPolicyCase dayu.KeyValueRecord
	has := false
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, flag, err := service.DescribeDdosPolicyCase(ctx, resourceType, sceneId)
		if err != nil {
			return tccommon.RetryError(err)
		}
		ddosPolicyCase = result
		has = flag
		return nil
	})

	if err != nil {
		return err
	}

	var list []map[string]interface{}
	var ids []string
	if has {
		list = make([]map[string]interface{}, 0, 1)
		ids = make([]string, 0, 1)
	}
	listItem := make(map[string]interface{})
	for _, record := range ddosPolicyCase.Record {
		key := *record.Key
		if key == "CaseName" {
			listItem["name"] = *record.Value
		}
		if key == "HasInitiateTcp" {

			listItem["has_initiate_tcp"] = *record.Value
		}
		if key == "HasInitiateUdp" {
			listItem["has_initiate_udp"] = *record.Value
		}
		if key == "HasVPN" {
			listItem["has_vpn"] = *record.Value
		}
		if key == "PeerTcpPort" {
			listItem["peer_tcp_port"] = *record.Value
		}
		if key == "PeerUdpPort" {
			listItem["peer_udp_port"] = *record.Value
		}
		if key == "TcpFootprint" {
			listItem["tcp_footprint"] = *record.Value
		}
		if key == "UdpFootprint" {
			listItem["udp_footprint"] = *record.Value
		}
		if key == "HasAbroad" {
			_ = d.Set("has_abroad", *record.Value)
			listItem["has_abroad"] = *record.Value
		}
		if key == "TcpSportStart" {
			listItem["tcp_start_port"] = *record.Value
		}
		if key == "TcpSportEnd" {
			listItem["tcp_end_port"] = *record.Value
		}
		if key == "UdpSportStart" {
			listItem["udp_start_port"] = *record.Value
		}
		if key == "UdpSportEnd" {
			listItem["udp_end_port"] = *record.Value
		}
		if key == "MaxUdpPackageLen" {
			listItem["max_udp_package_len"] = *record.Value
		}
		if key == "MinUdpPackageLen" {
			listItem["min_udp_package_len"] = *record.Value
		}
		if key == "MaxTcpPackageLen" {
			listItem["max_tcp_package_len"] = *record.Value
		}
		if key == "MinTcpPackageLen" {
			listItem["min_tcp_package_len"] = *record.Value
		}
		if key == "AppType" {
			listItem["app_type"] = *record.Value
		}
		if key == "AppProtocols" {
			listItem["app_protocols"] = strings.Split(*record.Value, ";")
		}
		if key == "WebApiUrl" {
			listItem["web_api_urls"] = strings.Split(*record.Value, ";")
		}
		if key == "PlatformTypes" {
			listItem["platform_types"] = strings.Split(*record.Value, ";")
		}
		if key == "Id" {
			listItem["scene_id"] = *record.Value
		}
		if key == "CreateTime" {
			listItem["create_time"] = *record.Value
		}
	}
	list = append(list, listItem)
	ids = append(ids, listItem["scene_id"].(string))

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return tccommon.WriteToFile(output.(string), list)
	}
	return nil

}
