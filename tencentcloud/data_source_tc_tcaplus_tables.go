/*
Use this data source to query tcaplus tables

Example Usage

```hcl
data "tencentcloud_tcaplus_tables" "null" {
  cluster_id = "19162256624"
}

data "tencentcloud_tcaplus_tables" "group" {
  cluster_id = "19162256624"
  group_id   = "19162256624:3"
}

data "tencentcloud_tcaplus_tables" "name" {
  cluster_id = "19162256624"
  group_id   = "19162256624:3"
  table_name = "guagua"
}

data "tencentcloud_tcaplus_tables" "id" {
  cluster_id = "19162256624"
  table_id   = "tcaplus-faa65eb7"
}
data "tencentcloud_tcaplus_tables" "all" {
  cluster_id = "19162256624"
  group_id   = "19162256624:3"
  table_id   = "tcaplus-faa65eb7"
  table_name = "guagua"
}
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceTencentCloudTcaplusTables() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTcaplusTablesRead,
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
			"table_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Table id to be query.",
			},
			"table_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Table name to be query.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of tcaplus groups. Each element contains the following attributes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Group of this table belongs.",
						},
						"table_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Id of this table.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of this table.",
						},
						"table_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of this table.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of this table.",
						},
						"idl_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Idl file id for this table.",
						},
						"table_idl_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of this table idl.",
						},
						"reserved_read_qps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Table reserved read QPS.",
						},
						"reserved_write_qps": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Table reserved write QPS.",
						},
						"reserved_volume": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Table reserved capacity(GB).",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Create time of the tcaplus table.",
						},
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Show if this table  create error.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Status of this table.",
						},
						"table_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Size of this table.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudTcaplusTablesRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_tcaplus_tables.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := TcaplusService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	clusterId := d.Get("cluster_id").(string)
	groupId := d.Get("group_id").(string)
	tableId := d.Get("table_id").(string)
	tableName := d.Get("table_name").(string)

	tables, err := service.DescribeTables(ctx, clusterId, groupId, tableId, tableName)
	if err != nil {
		tables, err = service.DescribeTables(ctx, clusterId, groupId, tableId, tableName)
	}
	if err != nil {
		return err
	}

	list := make([]map[string]interface{}, 0, len(tables))

	for _, tableInfo := range tables {

		listItem := make(map[string]interface{})

		if tableInfo.IdlFiles != nil && len(tableInfo.IdlFiles) > 0 {
			idlFile := tableInfo.IdlFiles[0]

			var tcaplusIdlId TcaplusIdlId
			tcaplusIdlId.ClusterId = clusterId
			tcaplusIdlId.FileName = *idlFile.FileName
			tcaplusIdlId.FileType = *idlFile.FileType

			tcaplusIdlId.FileExtType = *idlFile.FileExtType
			tcaplusIdlId.FileSize = *idlFile.FileSize
			tcaplusIdlId.FileId = *idlFile.FileId
			id, err := json.Marshal(tcaplusIdlId)

			if err != nil {
				return fmt.Errorf("format idl id fail,%s", err.Error())
			}
			listItem["idl_id"] = string(id)
		}

		if tableInfo.Error != nil && tableInfo.Error.Message != nil {
			listItem["error"] = *tableInfo.Error.Message
		} else {
			listItem["error"] = ""
		}
		if tableInfo.TableGroupId != nil {
			listItem["group_id"] = fmt.Sprintf("%s:%s", clusterId, *tableInfo.TableGroupId)
		}
		if tableInfo.TableInstanceId != nil {
			listItem["table_id"] = *tableInfo.TableInstanceId
		}
		if tableInfo.TableName != nil {
			listItem["table_name"] = *tableInfo.TableName
		}
		if tableInfo.TableType != nil {
			listItem["table_type"] = *tableInfo.TableType
		}
		if tableInfo.Memo != nil {
			listItem["description"] = *tableInfo.Memo
		}
		if tableInfo.TableIdlType != nil {
			listItem["table_idl_type"] = *tableInfo.TableIdlType
		}
		if tableInfo.ReservedReadQps != nil {
			listItem["reserved_read_qps"] = *tableInfo.ReservedReadQps
		}
		if tableInfo.ReservedWriteQps != nil {
			listItem["reserved_write_qps"] = *tableInfo.ReservedWriteQps
		}
		if tableInfo.ReservedVolume != nil {
			listItem["reserved_volume"] = *tableInfo.ReservedVolume
		}
		if tableInfo.CreatedTime != nil {
			listItem["create_time"] = *tableInfo.CreatedTime
		}
		if tableInfo.Status != nil {
			listItem["status"] = *tableInfo.Status
		}
		if tableInfo.TableSize != nil {
			listItem["table_size"] = *tableInfo.TableSize
		}
		list = append(list, listItem)
	}

	d.SetId("table." + clusterId + "." + groupId + "." + tableId + "." + tableName)
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
