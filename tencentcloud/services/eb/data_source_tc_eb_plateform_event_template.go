package eb

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudEbPlateformEventTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudEbPlateformEventTemplateRead,
		Schema: map[string]*schema.Schema{
			"event_type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Platform product event type.",
			},

			"event_template": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Platform product event template.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudEbPlateformEventTemplateRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_eb_plateform_event_template.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	var eventType string
	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("event_type"); ok {
		eventType = v.(string)
		paramMap["EventType"] = helper.String(v.(string))
	}

	service := EbService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var eventTemplate *string
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeEbPlateformEventTemplateByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		eventTemplate = result
		return nil
	})
	if err != nil {
		return err
	}

	if eventTemplate != nil {
		_ = d.Set("event_template", eventTemplate)
	}

	d.SetId(helper.DataResourceIdsHash([]string{eventType, *eventTemplate}))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), eventTemplate); e != nil {
			return e
		}
	}
	return nil
}
