package oceanus

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	oceanus "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/oceanus/v20190422"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudOceanusMetaTable() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudOceanusMetaTableRead,
		Schema: map[string]*schema.Schema{
			"work_space_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Unique identifier of the space.",
			},
			"catalog": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Catalog name.",
			},
			"database": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Database name.",
			},
			"table": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Table name.",
			},
			// computed
			"serial_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Unique identifier of the metadata table.",
			},
			"ddl": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Table creation statement, encoded in Base64.For example,Q1JFQVRFIFRBQkxFIGRhdGFnZW5fc291cmNlX3RhYmxlICggCiAgICBpZCBJTlQsIAogICAgbmFtZSBTVFJJTkcgCikgV0lUSCAoCidjb25uZWN0b3InPSdkYXRhZ2VuJywKJ3Jvd3MtcGVyLXNlY29uZCcgPSAnMScKKTs=.",
			},
			"create_time": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Scene time.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudOceanusMetaTableRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_oceanus_meta_table.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId       = tccommon.GetLogId(tccommon.ContextNil)
		ctx         = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service     = OceanusService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		MetaTable   *oceanus.GetMetaTableResponseParams
		workSpaceId string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("work_space_id"); ok {
		paramMap["WorkSpaceId"] = helper.String(v.(string))
		workSpaceId = v.(string)
	}

	if v, ok := d.GetOk("catalog"); ok {
		paramMap["Catalog"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("database"); ok {
		paramMap["Database"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("table"); ok {
		paramMap["Table"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeOceanusMetaTableByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		if result == nil {
			e = fmt.Errorf("MetaTable not exists")
			return resource.NonRetryableError(e)
		}

		MetaTable = result
		return nil
	})

	if err != nil {
		return err
	}

	_ = d.Set("work_space_id", workSpaceId)

	if MetaTable.Catalog != nil {
		_ = d.Set("catalog", MetaTable.Catalog)
	}

	if MetaTable.Database != nil {
		_ = d.Set("database", MetaTable.Database)
	}

	if MetaTable.Table != nil {
		_ = d.Set("table", MetaTable.Table)
	}

	if MetaTable.SerialId != nil {
		_ = d.Set("serial_id", MetaTable.SerialId)
	}

	if MetaTable.DDL != nil {
		_ = d.Set("ddl", MetaTable.DDL)
	}

	if MetaTable.CreateTime != nil {
		_ = d.Set("create_time", MetaTable.CreateTime)
	}

	d.SetId(workSpaceId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
