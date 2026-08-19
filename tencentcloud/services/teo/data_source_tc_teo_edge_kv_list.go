package teo

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
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

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.Background()
	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseTeoV20220901Client()

	request := teov20220901.NewEdgeKVListRequest()

	if v, ok := d.GetOk("zone_id"); ok {
		request.ZoneId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("namespace"); ok {
		request.Namespace = helper.String(v.(string))
	}

	if v, ok := d.GetOk("prefix"); ok {
		request.Prefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("cursor"); ok {
		request.Cursor = helper.String(v.(string))
	}

	var keysList []*string
	var lastCursor *string
	request.Limit = helper.IntInt64(1000)

	for {
		ratelimit.Check(request.GetAction())
		var response *teov20220901.EdgeKVListResponse
		err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			resp, e := client.EdgeKVListWithContext(ctx, request)
			if e != nil {
				return tccommon.RetryError(e)
			}
			if resp == nil || resp.Response == nil {
				log.Printf("[DATASOURCE] read empty, skip SetId, teo_edge_kv_list zone_id=%s namespace=%s", d.Get("zone_id"), d.Get("namespace"))
				return resource.NonRetryableError(errors.New("teo_edge_kv_list EdgeKVList response is nil"))
			}
			response = resp
			return nil
		})
		if err != nil {
			return err
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response.Response.Keys != nil {
			keysList = append(keysList, response.Response.Keys...)
		}

		if response.Response.Cursor != nil {
			lastCursor = response.Response.Cursor
			// Empty string cursor means traversal is complete.
			if strings.TrimSpace(*response.Response.Cursor) == "" {
				break
			}
			request.Cursor = response.Response.Cursor
		} else {
			// nil cursor means traversal is complete.
			break
		}
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
