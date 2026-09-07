package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoEdgeKvList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoEdgeKvListRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Site ID.",
			},

			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "EdgeKV namespace name.",
			},

			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key name prefix filter. Only keys starting with the specified prefix are returned. If not set, all keys are returned.",
			},

			"cursor": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Cursor position for traversal. Do not set this field for the first query; for subsequent queries, set it to the cursor returned by the previous query. After Read completes, this field is populated with the cursor from the last API response.",
			},

			"keys": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of key names in the specified namespace.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoEdgeKvListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_edge_kv_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TeoService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("zone_id"); ok {
		paramMap["zone_id"] = v.(string)
	}

	if v, ok := d.GetOk("namespace"); ok {
		paramMap["namespace"] = v.(string)
	}

	if v, ok := d.GetOk("prefix"); ok {
		paramMap["prefix"] = v.(string)
	}

	if v, ok := d.GetOk("cursor"); ok {
		paramMap["cursor"] = v.(string)
	}

	keysList, lastCursor, err := service.DescribeTeoEdgeKvListByFilter(ctx, paramMap)
	if err != nil {
		return err
	}

	if lastCursor != nil {
		_ = d.Set("cursor", *lastCursor)
	}

	if keysList != nil {
		_ = d.Set("keys", keysList)
	}

	d.SetId(helper.BuildToken())

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
