/*
Use this resource to create TcaplusDB IDL file.

Example Usage

Create a tcaplus database idl file

The file will be with a specified cluster and tablegroup.

```hcl
locals {
  vpc_id    = data.tencentcloud_vpc_subnets.vpc.instance_list.0.vpc_id
  subnet_id = data.tencentcloud_vpc_subnets.vpc.instance_list.0.subnet_id
}

variable "availability_zone" {
  default = "ap-guangzhou-3"
}

data "tencentcloud_vpc_subnets" "vpc" {
  is_default        = true
  availability_zone = var.availability_zone
}

resource "tencentcloud_tcaplus_cluster" "example" {
  idl_type                 = "PROTO"
  cluster_name             = "tf_example_tcaplus_cluster"
  vpc_id                   = local.vpc_id
  subnet_id                = local.subnet_id
  password                 = "your_pw_123111"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_tablegroup" "example" {
  cluster_id      = tencentcloud_tcaplus_cluster.example.id
  tablegroup_name = "tf_example_group_name"
}

resource "tencentcloud_tcaplus_idl" "main" {
  cluster_id    = tencentcloud_tcaplus_cluster.example.id
  tablegroup_id = tencentcloud_tcaplus_tablegroup.example.id
  file_name     = "tf_example_tcaplus_idl"
  file_type     = "PROTO"
  file_ext_type = "proto"
  file_content  = <<EOF
    syntax = "proto2";
    package myTcaplusTable;
    import "tcaplusservice.optionv1.proto";
    message tb_online {
        option(tcaplusservice.tcaplus_primary_key) = "uin,name,region";
        required int64 uin = 1;
        required string name = 2;
        required int32 region = 3;
        required int32 gamesvrid = 4;
        optional int32 logintime = 5 [default = 1];
        repeated int64 lockid = 6 [packed = true];
        optional bool is_available = 7 [default = false];
        optional pay_info pay = 8;
    }

    message pay_info {
        required int64 pay_id = 1;
        optional uint64 total_money = 2;
        optional uint64 pay_times = 3;
        optional pay_auth_info auth = 4;
        message pay_auth_info {
            required string pay_keys = 1;
            optional int64 update_time = 2;
        }
    }
    EOF
}
```
*/
package tencentcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TcaplusIdlId struct {
	ClusterId   string
	FileExtType string
	FileId      int64
	FileName    string
	FileSize    int64
	FileType    string
}

func resourceTencentCloudTcaplusIdl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcaplusIdlCreate,
		Read:   resourceTencentCloudTcaplusIdlRead,
		Delete: resourceTencentCloudTcaplusIdlDelete,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the TcaplusDB cluster to which the table group belongs.",
			},
			"tablegroup_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the table group to which the IDL file belongs.",
			},
			"file_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the IDL file.",
			},
			"file_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(TCAPLUS_IDL_TYPES),
				Description:  "Type of the IDL file. Valid values are PROTO and TDR.",
			},
			"file_ext_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(TCAPLUS_FILE_EXT_TYPES),
				Description:  "File ext type of the IDL file. If `file_type` is `PROTO`, `file_ext_type` must be 'proto'; If `file_type` is `TDR`, `file_ext_type` must be 'xml'.",
			},
			"file_content": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IDL file content of the TcaplusDB table.",
			},

			// Computed values.
			"table_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Table info of the IDL.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Error messages for creating IDL file.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the TcaplusDB table.",
						},
						"key_fields": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Primary key fields of the TcaplusDB table.",
						},
						"sum_key_field_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total size of primary key field of the TcaplusDB table.",
						},
						"value_fields": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Non-primary key fields of the TcaplusDB table.",
						},
						"sum_value_field_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Total size of non-primary key fields of the TcaplusDB table.",
						},
						"index_key_set": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Index key set of the TcaplusDB table.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudTcaplusIdlCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_idl.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	var tcaplusIdlId TcaplusIdlId
	tcaplusIdlId.ClusterId = d.Get("cluster_id").(string)
	tcaplusIdlId.FileName = d.Get("file_name").(string)
	tcaplusIdlId.FileType = d.Get("file_type").(string)
	tcaplusIdlId.FileExtType = d.Get("file_ext_type").(string)

	groupId := d.Get("tablegroup_id").(string)
	fileContent := d.Get("file_content").(string)

	items := strings.Split(groupId, ":")
	if len(items) != 2 {
		return fmt.Errorf("group id is broken,%s", groupId)
	}
	groupId = items[1]

	matchExtTypes := FileExtTypeMatch[tcaplusIdlId.FileType]
	if matchExtTypes == nil || !matchExtTypes[tcaplusIdlId.FileExtType] {
		return fmt.Errorf("file_ext_type %s not match file_type %s",
			tcaplusIdlId.FileExtType, tcaplusIdlId.FileType)
	}

	tcaplusIdlId.FileSize = int64(len(fileContent))

	idlId, parseTableInfos, err := tcaplusService.VerifyIdlFiles(ctx, tcaplusIdlId, groupId, fileContent)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			idlId, parseTableInfos, err = tcaplusService.VerifyIdlFiles(ctx, tcaplusIdlId, groupId, fileContent)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	if len(parseTableInfos) == 0 {
		return fmt.Errorf("Illegal idl. no tables defined in this file were found")
	}

	tableInfos := make([]map[string]interface{}, 0, len(parseTableInfos))

	for _, tableInfo := range parseTableInfos {
		var infoMap = map[string]interface{}{}
		if tableInfo.Error == nil {
			infoMap["error"] = ""
		} else {
			infoMap["error"] = *tableInfo.Error.Message
			tableInfos = append(tableInfos, infoMap)
			continue
		}
		infoMap["key_fields"] = tableInfo.KeyFields
		infoMap["sum_key_field_size"] = tableInfo.SumKeyFieldSize
		infoMap["value_fields"] = tableInfo.ValueFields
		infoMap["sum_value_field_size"] = tableInfo.SumValueFieldSize
		infoMap["index_key_set"] = tableInfo.IndexKeySet
		infoMap["table_name"] = tableInfo.TableName
		tableInfos = append(tableInfos, infoMap)
	}

	tcaplusIdlId.FileId = idlId
	_ = d.Set("table_infos", tableInfos)
	id, err := json.Marshal(tcaplusIdlId)
	if err != nil {
		return fmt.Errorf("format idl id fail,%s", err.Error())
	}
	d.SetId(string(id))
	return resourceTencentCloudTcaplusIdlRead(d, meta)
}

func resourceTencentCloudTcaplusIdlRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_tcaplus_idl.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	var tcaplusIdlId TcaplusIdlId

	if err := json.Unmarshal([]byte(d.Id()), &tcaplusIdlId); err != nil {
		return fmt.Errorf("idl id is broken,%s", err.Error())
	}

	parseTableInfos, err := tcaplusService.DesOldIdlFiles(ctx, tcaplusIdlId)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			parseTableInfos, err = tcaplusService.DesOldIdlFiles(ctx, tcaplusIdlId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	if err != nil {
		return err
	}

	if len(parseTableInfos) == 0 {
		d.SetId("")
		return nil
	}
	tableInfos := make([]map[string]interface{}, 0, len(parseTableInfos))

	for _, tableInfo := range parseTableInfos {
		var infoMap = map[string]interface{}{}
		if tableInfo.Error == nil {
			infoMap["error"] = ""
		} else {
			infoMap["error"] = *tableInfo.Error.Message
			tableInfos = append(tableInfos, infoMap)
			continue
		}
		infoMap["key_fields"] = tableInfo.KeyFields
		infoMap["sum_key_field_size"] = tableInfo.SumKeyFieldSize
		infoMap["value_fields"] = tableInfo.ValueFields
		infoMap["sum_value_field_size"] = tableInfo.SumValueFieldSize
		infoMap["index_key_set"] = tableInfo.IndexKeySet
		infoMap["table_name"] = tableInfo.TableName
		tableInfos = append(tableInfos, infoMap)
	}

	return d.Set("table_infos", tableInfos)
}

func resourceTencentCloudTcaplusIdlDelete(d *schema.ResourceData, meta interface{}) error {

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(*TencentCloudClient).apiV3Conn}

	var tcaplusIdlId TcaplusIdlId

	if err := json.Unmarshal([]byte(d.Id()), &tcaplusIdlId); err != nil {
		return fmt.Errorf("idl id is broken,%s", err.Error())
	}

	err := tcaplusService.DeleteIdlFiles(ctx, tcaplusIdlId)

	if err != nil {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			err = tcaplusService.DeleteIdlFiles(ctx, tcaplusIdlId)
			if err != nil {
				return retryError(err)
			}
			return nil
		})
	}
	return err
}
