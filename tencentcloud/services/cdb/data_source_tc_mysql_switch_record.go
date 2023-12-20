package cdb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	cdb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdb/v20170320"
)

func DataSourceTencentCloudMysqlSwitchRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudMysqlSwitchRecordRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Instance ID, the format is: cdb-c1nl9rpv or cdbro-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page.",
			},

			"items": {
				Computed:    true,
				Type:        schema.TypeList,
				Description: "Instance switching record details.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switching time, the format is: 2017-09-03 01:34:31.",
						},
						"switch_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Switch type, possible return values: TRANSFER - data migration; MASTER2SLAVE - master-standby switch; RECOVERY - master-slave recovery.",
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

func dataSourceTencentCloudMysqlSwitchRecordRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_mysql_switch_record.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var instanceId string
	if v, ok := d.GetOk("instance_id"); ok {
		instanceId = v.(string)
	}

	service := MysqlService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	var items []*cdb.DBSwitchInfo
	err := resource.Retry(5*tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeMysqlSwitchRecordById(ctx, instanceId)
		if e != nil {
			return tccommon.RetryError(e)
		}
		items = result
		return nil
	})
	if err != nil {
		return err
	}

	tmpList := make([]map[string]interface{}, 0, len(items))
	if items != nil {
		for _, dBSwitchInfo := range items {
			dBSwitchInfoMap := map[string]interface{}{}

			if dBSwitchInfo.SwitchTime != nil {
				dBSwitchInfoMap["switch_time"] = dBSwitchInfo.SwitchTime
			}

			if dBSwitchInfo.SwitchType != nil {
				dBSwitchInfoMap["switch_type"] = dBSwitchInfo.SwitchType
			}

			tmpList = append(tmpList, dBSwitchInfoMap)
		}

		_ = d.Set("items", tmpList)
	}

	d.SetId(instanceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}
	return nil
}
