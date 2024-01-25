// Code generated by iacg; DO NOT EDIT.
package cls

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudClsMachines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsMachinesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Group id.",
			},

			"machines": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Info of Machines.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ip of machine.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "status of machine.",
						},
						"offline_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "offline time of machine.",
						},
						"auto_update": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "if open auto update flag.",
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "current machine version.",
						},
						"update_status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "machine update status.",
						},
						"err_code": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "code of update operation.",
						},
						"err_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "msg of update operation.",
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

func dataSourceTencentCloudClsMachinesRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_cls_machines.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	service := ClsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	var respData []*cls.MachineInfo
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClsMachinesByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		respData = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(respData))
	machinesList := make([]map[string]interface{}, 0, len(respData))
	if respData != nil {
		for _, machines := range respData {
			machinesMap := map[string]interface{}{}

			if machines.Ip != nil {
				machinesMap["ip"] = machines.Ip
			}
			ip := *machines.Ip

			if machines.Status != nil {
				machinesMap["status"] = machines.Status
			}

			if machines.OfflineTime != nil {
				machinesMap["offline_time"] = machines.OfflineTime
			}

			if machines.AutoUpdate != nil {
				machinesMap["auto_update"] = machines.AutoUpdate
			}

			if machines.Version != nil {
				machinesMap["version"] = machines.Version
			}

			if machines.UpdateStatus != nil {
				machinesMap["update_status"] = machines.UpdateStatus
			}

			if machines.ErrCode != nil {
				machinesMap["err_code"] = machines.ErrCode
			}

			if machines.ErrMsg != nil {
				machinesMap["err_msg"] = machines.ErrMsg
			}

			ids = append(ids, ip)
			machinesList = append(machinesList, machinesMap)
		}

		_ = d.Set("machines", machinesList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), machinesList); e != nil {
			return e
		}
	}

	return nil
}
