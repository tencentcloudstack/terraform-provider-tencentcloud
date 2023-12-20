package audit

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	audit "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cloudaudit/v20190319"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudAuditCosRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAuditCosRegionsRead,

		Schema: map[string]*schema.Schema{
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"audit_cos_region_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of available regions supported by audit cos.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cos_region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cos region.",
						},
						"cos_region_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cos region chinese name.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAuditCosRegionsRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_audit_cos_regions.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	auditService := AuditService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var regions []*audit.CosRegionInfo
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		regions, errRet = auditService.DescribeAuditCosRegions(ctx)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	regionList := make([]map[string]interface{}, 0, len(regions))
	ids := make([]string, 0, len(regions))
	for _, region := range regions {
		mapping := map[string]interface{}{
			"cos_region":      region.CosRegion,
			"cos_region_name": region.CosRegionName,
		}
		regionList = append(regionList, mapping)
		ids = append(ids, *region.CosRegion)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("audit_cos_region_list", regionList)
	if err != nil {
		log.Printf("[CRITAL]%s audit cos read regions list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), regionList); e != nil {
			return e
		}
	}
	return nil
}
