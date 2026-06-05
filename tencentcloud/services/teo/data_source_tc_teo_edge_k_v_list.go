package teo

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTeoEdgeKVList() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTeoEdgeKVListRead,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Zone ID.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Namespace name.",
			},
			"prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Key name prefix filter. Only returns keys that start with the specified prefix, length 1-512 characters.",
			},
			"cursor": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Cursor position. Identifies the starting position of the current query for traversing large amounts of data. Leave empty for the first query to start from the beginning; for subsequent queries, fill in the Cursor value returned from the previous response to continue traversal.",
			},
			"keys": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of key names.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudTeoEdgeKVListRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_teo_edge_k_v_list.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = NewTeoService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
	)

	zoneId := d.Get("zone_id").(string)
	namespace := d.Get("namespace").(string)

	var prefix *string
	if v, ok := d.GetOk("prefix"); ok {
		tmp := v.(string)
		prefix = &tmp
	}

	var cursor *string
	if v, ok := d.GetOk("cursor"); ok {
		tmp := v.(string)
		cursor = &tmp
	}

	keys, nextCursor, err := service.DescribeTeoEdgeKVList(ctx, zoneId, namespace, prefix, cursor)
	if err != nil {
		return err
	}

	keysList := make([]string, 0, len(keys))
	if keys != nil {
		for _, key := range keys {
			if key != nil {
				keysList = append(keysList, *key)
			}
		}
	}

	_ = d.Set("keys", keysList)

	if nextCursor != nil {
		_ = d.Set("cursor", *nextCursor)
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
