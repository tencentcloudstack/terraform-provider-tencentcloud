package tcaplusdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

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

func ResourceTencentCloudTcaplusIdl() *schema.Resource {
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
				ValidateFunc: tccommon.ValidateAllowedStringValue(TCAPLUS_IDL_TYPES),
				Description:  "Type of the IDL file. Valid values are PROTO and TDR.",
			},
			"file_ext_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: tccommon.ValidateAllowedStringValue(TCAPLUS_FILE_EXT_TYPES),
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
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_idl.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

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
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			idlId, parseTableInfos, err = tcaplusService.VerifyIdlFiles(ctx, tcaplusIdlId, groupId, fileContent)
			if err != nil {
				return tccommon.RetryError(err)
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
	defer tccommon.LogElapsed("resource.tencentcloud_tcaplus_idl.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var tcaplusIdlId TcaplusIdlId

	if err := json.Unmarshal([]byte(d.Id()), &tcaplusIdlId); err != nil {
		return fmt.Errorf("idl id is broken,%s", err.Error())
	}

	parseTableInfos, err := tcaplusService.DesOldIdlFiles(ctx, tcaplusIdlId)
	if err != nil {
		err = resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			parseTableInfos, err = tcaplusService.DesOldIdlFiles(ctx, tcaplusIdlId)
			if err != nil {
				return tccommon.RetryError(err)
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

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	tcaplusService := TcaplusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var tcaplusIdlId TcaplusIdlId

	if err := json.Unmarshal([]byte(d.Id()), &tcaplusIdlId); err != nil {
		return fmt.Errorf("idl id is broken,%s", err.Error())
	}

	err := tcaplusService.DeleteIdlFiles(ctx, tcaplusIdlId)

	if err != nil {
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			err = tcaplusService.DeleteIdlFiles(ctx, tcaplusIdlId)
			if err != nil {
				return tccommon.RetryError(err)
			}
			return nil
		})
	}
	return err
}
