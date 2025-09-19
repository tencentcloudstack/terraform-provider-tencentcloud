package tag

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudTagKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudTagKeysRead,
		Schema: map[string]*schema.Schema{
			"create_uin": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Creator `Uin`. If not specified, `Uin` is only used as the query condition.",
			},

			"show_project": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Whether to show project. Allow values: 0: no, 1: yes.",
			},

			"category": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Tag type. Valid values: Custom: custom tag; System: system tag; All: all tags. Default value: All.",
			},

			"tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Tag list.",
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

func dataSourceTencentCloudTagKeysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_tag_keys.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId   = tccommon.GetLogId(nil)
		ctx     = tccommon.NewResourceLifeCycleHandleFuncContext(context.Background(), logId, d, meta)
		service = TagService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOkExists("create_uin"); ok {
		paramMap["CreateUin"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOkExists("show_project"); ok {
		paramMap["ShowProject"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("category"); ok {
		paramMap["Category"] = helper.String(v.(string))
	}

	var respData []*string
	reqErr := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeTagKeysByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		respData = result
		return nil
	})

	if reqErr != nil {
		return reqErr
	}

	if respData != nil {
		_ = d.Set("tags", respData)
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
