package kms

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKmsListKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsListKeysRead,
		Schema: map[string]*schema.Schema{
			"role": {
				Optional:    true,
				Type:        schema.TypeInt,
				Description: "Filter based on the creator role. The default value is 0, which indicates the cmk created by the user himself, and 1, which indicates the cmk automatically created by authorizing other cloud products.",
			},
			"hsm_cluster_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "HSM cluster ID (only valid for KMS exclusive/managed service instances).",
			},
			"keys": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of KMS keys.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of CMK.",
						},
					},
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

func dataSourceTencentCloudKmsListKeysRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kms_list_keys.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId    = tccommon.GetLogId(tccommon.ContextNil)
		ctx      = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service  = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		listKeys []*kms.Key
	)

	paramMap := make(map[string]interface{})
	if v, _ := d.GetOk("role"); v != nil {
		paramMap["Role"] = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("hsm_cluster_id"); ok {
		paramMap["HsmClusterId"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsListKeysByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		listKeys = result
		return nil
	})

	if err != nil {
		return err
	}

	ids := make([]string, 0, len(listKeys))
	tmpList := make([]map[string]interface{}, 0, len(listKeys))

	if listKeys != nil {
		for _, key := range listKeys {
			mapping := map[string]interface{}{
				"key_id": key.KeyId,
			}

			tmpList = append(tmpList, mapping)
			ids = append(ids, *key.KeyId)
		}

		_ = d.Set("keys", tmpList)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), tmpList); e != nil {
			return e
		}
	}

	return nil
}
