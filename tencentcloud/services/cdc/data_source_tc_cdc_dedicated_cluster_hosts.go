package cdc

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdc/v20201214"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudCdcDedicatedClusterHosts() *schema.Resource {
	return &schema.Resource{
		Read: DataSourceTencentCloudCdcDedicatedClusterHostsRead,
		Schema: map[string]*schema.Schema{
			"dedicated_cluster_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Dedicated Cluster ID.",
			},
			// computed
			"host_info_set": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Dedicated Cluster Host Info.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Host Ip (Deprecated).",
						},
						"service_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Service Type.",
						},
						"host_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Host Status.",
						},
						"host_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Host Type.",
						},
						"cpu_available": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dedicated Cluster Host CPU Available Count.",
						},
						"cpu_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dedicated Cluster Host CPU Total Count.",
						},
						"mem_available": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dedicated Cluster Host Memory Available Count (GB).",
						},
						"mem_total": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Dedicated Cluster Host Memory Total Count (GB).",
						},
						"run_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Host Run Time.",
						},
						"expire_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Host Expire Time.",
						},
						"host_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Dedicated Cluster Host ID.",
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

func DataSourceTencentCloudCdcDedicatedClusterHostsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cdc_dedicated_cluster_hosts.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = CdcService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		hostInfoSet []*cdc.HostInfo
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("dedicated_cluster_id"); ok {
		paramMap["DedicatedClusterId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeCdcHostByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		hostInfoSet = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(hostInfoSet))
	tmpList := make([]map[string]interface{}, 0, len(hostInfoSet))

	if hostInfoSet != nil {
		for _, hostInfo := range hostInfoSet {
			hostInfoMap := map[string]interface{}{}

			if hostInfo.HostIp != nil {
				hostInfoMap["host_ip"] = hostInfo.HostIp
			}

			if hostInfo.ServiceType != nil {
				hostInfoMap["service_type"] = hostInfo.ServiceType
			}

			if hostInfo.HostStatus != nil {
				hostInfoMap["host_status"] = hostInfo.HostStatus
			}

			if hostInfo.HostType != nil {
				hostInfoMap["host_type"] = hostInfo.HostType
			}

			if hostInfo.CpuAvailable != nil {
				hostInfoMap["cpu_available"] = hostInfo.CpuAvailable
			}

			if hostInfo.CpuTotal != nil {
				hostInfoMap["cpu_total"] = hostInfo.CpuTotal
			}

			if hostInfo.MemAvailable != nil {
				hostInfoMap["mem_available"] = hostInfo.MemAvailable
			}

			if hostInfo.MemTotal != nil {
				hostInfoMap["mem_total"] = hostInfo.MemTotal
			}

			if hostInfo.RunTime != nil {
				hostInfoMap["run_time"] = hostInfo.RunTime
			}

			if hostInfo.ExpireTime != nil {
				hostInfoMap["expire_time"] = hostInfo.ExpireTime
			}

			if hostInfo.HostId != nil {
				hostInfoMap["host_id"] = hostInfo.HostId
			}

			ids = append(ids, *hostInfo.HostId)
			tmpList = append(tmpList, hostInfoMap)
		}

		_ = d.Set("host_info_set", tmpList)
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
