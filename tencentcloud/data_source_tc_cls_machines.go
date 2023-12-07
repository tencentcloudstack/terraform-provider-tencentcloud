package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudClsMachines() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudClsMachinesRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Group id.",
			},

			"machines": {
				Computed:    true,
				Type:        schema.TypeList,
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
	defer logElapsed("data_source.tencentcloud_cls_machines.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("group_id"); ok {
		paramMap["GroupId"] = helper.String(v.(string))
	}

	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	var machines []*cls.MachineInfo

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeClsMachinesByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		machines = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(machines))
	tmpList := make([]map[string]interface{}, 0, len(machines))

	if machines != nil {
		for _, machineInfo := range machines {
			machineInfoMap := map[string]interface{}{}

			if machineInfo.Ip != nil {
				machineInfoMap["ip"] = machineInfo.Ip
			}

			if machineInfo.Status != nil {
				machineInfoMap["status"] = machineInfo.Status
			}

			if machineInfo.OfflineTime != nil {
				machineInfoMap["offline_time"] = machineInfo.OfflineTime
			}

			if machineInfo.AutoUpdate != nil {
				machineInfoMap["auto_update"] = machineInfo.AutoUpdate
			}

			if machineInfo.Version != nil {
				machineInfoMap["version"] = machineInfo.Version
			}

			if machineInfo.UpdateStatus != nil {
				machineInfoMap["update_status"] = machineInfo.UpdateStatus
			}

			if machineInfo.ErrCode != nil {
				machineInfoMap["err_code"] = machineInfo.ErrCode
			}

			if machineInfo.ErrMsg != nil {
				machineInfoMap["err_msg"] = machineInfo.ErrMsg
			}

			ids = append(ids, *machineInfo.Ip)
			tmpList = append(tmpList, machineInfoMap)
		}

		_ = d.Set("machines", tmpList)
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
