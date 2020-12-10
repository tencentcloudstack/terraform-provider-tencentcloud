/*
Use this data source to query the list of SQL Server readonly groups.

Example Usage

```hcl
data "tencentcloud_sqlserver_readonly_groups" "master" {
  master_instance_id           = "mssql-3cdq7kx5"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSqlserverReadonlyGroups() *schema.Resource {

	return &schema.Resource{
		Read: dataSourceTencentSqlserverReadonlyGroups,
		Schema: map[string]*schema.Schema{
			"master_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Master SQL Server instance ID.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to store results.",
			},
			// Computed values
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of SQL Server readonly group. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the readonly group.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the readonly group.",
						},
						"master_instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Master instance id.",
						},
						"max_delay_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Maximum delay time of the readonly instances.",
						},
						"is_offline_delay": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicate whether to offline delayed readonly instances.",
						},
						"min_instances": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Minimum readonly instances that stays in the group.",
						},
						"vip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Virtual IP address of the readonly group.",
						},
						"vport": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Virtual port of the readonly group.",
						},
						"readonly_instance_set": {
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Readonly instance ID set of the readonly group.",
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Status of the readonly group. `1` for running, `5` for applying.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentSqlserverReadonlyGroups(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_sqlserver_readonly_groups.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	instanceId := d.Get("master_instance_id").(string)
	sqlserverService := SqlserverService{client: meta.(*TencentCloudClient).apiV3Conn}

	groupList, err := sqlserverService.DescribeReadonlyGroupList(ctx, instanceId)

	if err != nil {
		return fmt.Errorf("api[DescribeReadOnlyGroupList]fail, return %s", err.Error())
	}

	var list []map[string]interface{}
	var ids = make([]string, len(groupList))

	for _, item := range groupList {
		roSet := make([]string, 0)
		for _, v := range item.ReadOnlyInstanceSet {
			roSet = append(roSet, *v.InstanceId)
		}
		mapping := map[string]interface{}{
			"name":                  item.ReadOnlyGroupName,
			"vip":                   item.Vip,
			"vport":                 item.Vport,
			"is_offline_delay":      item.IsOfflineDelay,
			"max_delay_time":        item.ReadOnlyMaxDelayTime,
			"min_instances":         item.MinReadOnlyInGroup,
			"status":                item.Status,
			"master_instance_id":    item.MasterInstanceId,
			"id":                    item.ReadOnlyGroupId,
			"readonly_instance_set": roSet,
		}
		list = append(list, mapping)
		ids = append(ids, *item.ReadOnlyGroupId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("list", list); e != nil {
		log.Printf("[CRITAL]%s provider set list fail, reason:%s\n", logId, e.Error())
		return e
	}
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		return writeToFile(output.(string), list)
	}

	return nil
}
