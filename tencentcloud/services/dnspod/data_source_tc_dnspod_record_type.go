package dnspod

import (
	"context"
	"strconv"
	"time"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudDnspodRecordType() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudDnspodRecordTypeRead,
		Schema: map[string]*schema.Schema{
			"domain_grade": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Domain level. + Old packages: D_FREE, D_PLUS, D_EXTRA, D_EXPERT, D_ULTRA correspond to free package, personal luxury, enterprise 1, enterprise 2, enterprise 3. + New packages: DP_FREE, DP_PLUS, DP_EXTRA, DP_EXPERT, DP_ULTRA correspond to new free, personal professional, enterprise basic, enterprise standard, enterprise flagship.",
			},

			"type_list": {
				Computed: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Record type list.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudDnspodRecordTypeRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_dnspod_record_type.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)

	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("domain_grade"); ok {
		paramMap["DomainGrade"] = helper.String(v.(string))
	}

	service := DnspodService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}

	var typeList []*string

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeDnspodRecordTypeByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}
		typeList = result
		return nil
	})
	if err != nil {
		return err
	}

	// ids := make([]string, 0, len(typeList))
	if typeList != nil {
		_ = d.Set("type_list", typeList)
	}

	// d.SetId(helper.DataResourceIdsHash(ids))
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), typeList); e != nil {
			return e
		}
	}
	return nil
}
