package tencentcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

func TencentMysqlBackupInfosItem() map[string]*schema.Schema {

	return map[string]*schema.Schema{
		"time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"finish_time": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"size": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"backup_id": {
			Type:     schema.TypeInt,
			Computed: true,
		},
		"backup_model": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"intranet_url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"internet_url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"creator": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

/*
data_source data_source_tc_mysql_backup_list{
	instance_id      string  = "Database (master)ID"
    max_number       int64  = "Recent log numbers, default 10,min 0 ,max10000"
    result_output_file string ="Output path"
	list             []TencentMsyqlBackupInfosItem = "Backup information list"
}
struct  TencentMsyqlBackupInfosItem {
    time              string = "Snapshot time for this backup"
    finish_time	      string = "The completion time of the backup task"
	size              int64  = "Backup file size"
	backup_id         int64  = "Backup task ID to be used when deleting the backup file"
	backup_model      string  = "logical - logical cold backup,physical - physical cold backup"
	intranet_url      string = "Intranet download address"
	internet_url      string = "Tnternet download address"
	creator           string = "Backup creator"
}
*/
func dataSourceTencentMysqlBackupList() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceTencentMysqlBackupListRead,
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"max_number": {
				Type:         schema.TypeInt,
				ForceNew:     true,
				Optional:     true,
				Default:      10,
				ValidateFunc: validateIntegerInRange(1, 10000),
			},
			"result_output_file": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			// Computed values
			"list": {Type: schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: TencentMysqlBackupInfosItem(),
				},
			},
		},
	}
}

func dataSourceTencentMysqlBackupListRead(d *schema.ResourceData, meta interface{}) error {

	ctx := context.WithValue(context.TODO(), "logId", GetLogId(nil))

	mysqlService := MysqlService{client: meta.(*TencentCloudClient).apiV3Conn}

	max_number, _ := d.Get("max_number").(int)
	backInfoItems, err := mysqlService.DescribeBackupsByInstanceId(ctx, d.Get("instance_id").(string), int64(max_number))

	if err != nil {
		return fmt.Errorf("api[DescribeBackups]fail, return %s", err.Error())
	}

	var itemShemas []map[string]interface{}
	var ids = make([]string, len(backInfoItems))

	for index, item := range backInfoItems {
		mapping := map[string]interface{}{
			"time":         *item.Date,
			"finish_time":  *item.FinishTime,
			"size":         *item.Size,
			"backup_id":    *item.BackupId,
			"backup_model": *item.Type,
			"intranet_url": strings.Replace(*item.IntranetUrl, "\u0026", "&", -1),
			"internet_url": strings.Replace(*item.InternetUrl, "\u0026", "&", -1),
			"creator":      *item.Creator,
		}
		ids[index] = fmt.Sprintf("%d", *item.BackupId)
		itemShemas = append(itemShemas, mapping)
	}

	d.Set("list", itemShemas)
	d.SetId(dataResourceIdsHash(ids))

	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), itemShemas)
	}
	return nil
}
