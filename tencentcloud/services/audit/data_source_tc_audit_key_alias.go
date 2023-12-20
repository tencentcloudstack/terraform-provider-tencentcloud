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

func DataSourceTencentCloudAuditKeyAlias() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudAuditKeyAliasRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"audit_key_alias_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of available key alias supported by audit.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key ID.",
						},
						"key_alias": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Key alias.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudAuditKeyAliasRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_audit_cmq_regions.read")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
	auditService := AuditService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	region := d.Get("region").(string)
	var keyAlias []*audit.KeyMetadata
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		keyAlias, errRet = auditService.DescribeKeyAlias(ctx, region)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	keyList := make([]map[string]interface{}, 0, len(keyAlias))
	ids := make([]string, 0, len(keyAlias))
	for _, keyData := range keyAlias {
		mapping := map[string]interface{}{
			"key_id":    keyData.KeyId,
			"key_alias": keyData.Alias,
		}
		keyList = append(keyList, mapping)
		ids = append(ids, *keyData.KeyId)
	}
	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("audit_key_alias_list", keyList)
	if err != nil {
		log.Printf("[CRITAL]%s audit read key alias list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), keyList); e != nil {
			return e
		}
	}
	return nil
}
