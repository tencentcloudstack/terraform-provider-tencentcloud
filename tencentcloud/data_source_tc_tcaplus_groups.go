/*
Use this data source to query tcaplus table groups

Example Usage

```hcl
data "tencentcloud_tcaplus_groups" "null" {
  cluster_id = "19162256624"
}
data "tencentcloud_tcaplus_groups" "id" {
  cluster_id = "19162256624"
  group_id   = "19162256624:1"
}
data "tencentcloud_tcaplus_groups" "name" {
  cluster_id = "19162256624"
  group_name = "test"
}
data "tencentcloud_tcaplus_groups" "all" {
  cluster_id = "19162256624"
  group_id   = "19162256624:1"
  group_name = "test"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTencentCloudTcaplusGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcaplusGroupsRead,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Id of the tcaplus cluster to be query.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group id to be query.",
			},
			"group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Group name to be query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of tcaplus table groups. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the tcaplus group.",
						},
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of the tcaplus group.",
						},
						"table_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of tables.",
						},
						"total_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The total storage(MB).",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the tcaplus group.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTcaplusGroupsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcaplus_groups.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcaplusService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	clusterId := d.Get("cluster_id").(string)
	groupId := d.Get("group_id").(string)
	groupName := d.Get("group_name").(string)

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
		listItem["group_name"] = group.TableGroupName
		listItem["group_id"] = fmt.Sprintf("%s:%s", clusterId, *group.TableGroupId)
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
