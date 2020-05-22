/*
Use this resource to create tcaplus idl file

Example Usage

```hcl
resource "tencentcloud_tcaplus_application" "test" {
  idl_type                 = "PROTO"
  app_name                 = "tf_tcaplus_app_test"
  vpc_id                   = "vpc-7k6gzox6"
  subnet_id                = "subnet-akwgvfa3"
  password                 = "1qaA2k1wgvfa3ZZZ"
  old_password_expire_last = 3600
}

resource "tencentcloud_tcaplus_zone" "zone" {
  app_id    = tencentcloud_tcaplus_application.test.id
  zone_name = "tf_test_zone_name"
}

resource "tencentcloud_tcaplus_idl" "main" {
  app_id        = tencentcloud_tcaplus_application.test.id
  zone_id       = tencentcloud_tcaplus_zone.zone.id
  file_name     = "tf_idl_test"
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

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type TcaplusIdlId struct {
	ApplicationId string
	FileExtType   string
	FileId        int64
	FileName      string
	FileSize      int64
	FileType      string
}

func resourceTencentCloudTcaplusIdl() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudTcaplusIdlCreate,
		Read:   resourceTencentCloudTcaplusIdlRead,
		Delete: resourceTencentCloudTcaplusIdlDelete,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Application id of the idl belongs..",
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Zone of this idl belongs.",
			},
			"file_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of this idl file.",
			},
			"file_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(TCAPLUS_IDL_TYPES),
				Description:  "Type of this idl file, Valid values are " + strings.Join(TCAPLUS_IDL_TYPES, ",") + ".",
			},
			"file_ext_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue(TCAPLUS_FILE_EXT_TYPES),
				Description:  "File ext type of this idl file. if `file_type` is PROTO  `file_ext_type` must be 'proto',if `file_type` is TDR  `file_ext_type` must be 'xml',if `file_type` is MIX  `file_ext_type` must be 'xml' or 'proto'.",
			},
			"file_content": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Idl file content.",
			},

			// Computed values.
			"table_infos": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Table infos in this idl.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Show if this table  error.",
						},
						"table_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of this table.",
						},
						"key_fields": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key fields of this table.",
						},
						"sum_key_field_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Key fields size of this table.",
						},
						"value_fields": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value fields of this table.",
						},
						"sum_value_field_size": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Value fields size of this table.",
						},
						"index_key_set": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Index key set of this table.",
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
	tcaplusIdlId.ApplicationId = d.Get("app_id").(string)
	tcaplusIdlId.FileName = d.Get("file_name").(string)
	tcaplusIdlId.FileType = d.Get("file_type").(string)
	tcaplusIdlId.FileExtType = d.Get("file_ext_type").(string)

	zoneId := d.Get("zone_id").(string)
	fileContent := d.Get("file_content").(string)

	items := strings.Split(zoneId, ":")
	if len(items) != 2 {
		return fmt.Errorf("zone id is broken,%s", zoneId)
	}
	zoneId = items[1]

	matchExtTypes := FileExtTypeMatch[tcaplusIdlId.FileType]
	if matchExtTypes == nil || !matchExtTypes[tcaplusIdlId.FileExtType] {
		return fmt.Errorf("file_ext_type %s not match file_type %s",
			tcaplusIdlId.FileExtType, tcaplusIdlId.FileType)
	}

	tcaplusIdlId.FileSize = int64(len(fileContent))

	idlId, parseTableInfos, err := tcaplusService.VerifyIdlFiles(ctx, tcaplusIdlId, zoneId, fileContent)
	if err != nil {
		err = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			idlId, parseTableInfos, err = tcaplusService.VerifyIdlFiles(ctx, tcaplusIdlId, zoneId, fileContent)
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
		infoMap["key_fields"] = *tableInfo.KeyFields
		infoMap["sum_key_field_size"] = *tableInfo.SumKeyFieldSize
		infoMap["value_fields"] = *tableInfo.ValueFields
		infoMap["sum_value_field_size"] = *tableInfo.SumValueFieldSize
		infoMap["index_key_set"] = *tableInfo.IndexKeySet
		infoMap["table_name"] = *tableInfo.TableName
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
		infoMap["key_fields"] = *tableInfo.KeyFields
		infoMap["sum_key_field_size"] = *tableInfo.SumKeyFieldSize
		infoMap["value_fields"] = *tableInfo.ValueFields
		infoMap["sum_value_field_size"] = *tableInfo.SumValueFieldSize
		infoMap["index_key_set"] = *tableInfo.IndexKeySet
		infoMap["table_name"] = *tableInfo.TableName
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
