/*
Use this data source to query table groups of the TcaplusDB cluster.

Example Usage

```hcl
data "tencentcloud_tcaplus_tablegroups" "null" {
  cluster_id = "19162256624"
}
data "tencentcloud_tcaplus_tablegroups" "id" {
  cluster_id    = "19162256624"
  tablegroup_id = "19162256624:1"
}
data "tencentcloud_tcaplus_tablegroups" "name" {
  cluster_id      = "19162256624"
  tablegroup_name = "test"
}
data "tencentcloud_tcaplus_tablegroups" "all" {
  cluster_id      = "19162256624"
  tablegroup_id   = "19162256624:1"
  tablegroup_name = "test"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTencentCloudTcaplusTableGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcaplusTableGroupsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the TcaplusDB cluster to be query.",
			},
			"tablegroup_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Id of the table group to be query.",
			},
			"tablegroup_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the table group to be query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "File for saving results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of table group. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tablegroup_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the table group.",
						},
						"tablegroup_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the table group.",
						},
						"table_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of tables.",
						},
						"total_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total storage size (MB).",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the table group..",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTcaplusTableGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcaplus_tablegroups.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcaplusService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	clusterId := d.Get("cluster_id").(string)
	groupId := d.Get("tablegroup_id").(string)
	groupName := d.Get("tablegroup_name").(string)

	groups, err := service.DescribeGroups(ctx, clusterId, groupId, groupName)
	if err != nil {
		groups, err = service.DescribeGroups(ctx, clusterId, groupId, groupName)
	}

	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(groups))

	for _, group := range groups {
		listItem := make(map[string]interface{})
		listItem["tablegroup_name"] = group.TableGroupName
		listItem["tablegroup_id"] = fmt.Sprintf("%s:%s", clusterId, *group.TableGroupId)
		listItem["table_count"] = group.TableCount
		listItem["total_size"] = group.TotalSize
		listItem["create_time"] = group.CreatedTime
		list = append(list, listItem)
	}

	d.SetId("group." + clusterId + "." + groupId + "." + groupName)
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
