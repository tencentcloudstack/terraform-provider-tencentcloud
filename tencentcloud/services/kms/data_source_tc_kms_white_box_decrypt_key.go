package kms

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKmsWhiteBoxDecryptKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsWhiteBoxDecryptKeyRead,
		Schema: map[string]*schema.Schema{
			"key_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Globally unique identifier for the white box key.",
			},
			"decrypt_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "White box decryption key, base64 encoded.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudKmsWhiteBoxDecryptKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kms_white_box_decrypt_key.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId              = tccommon.GetLogId(tccommon.ContextNil)
		ctx                = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service            = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		whiteBoxDecryptKey *kms.DescribeWhiteBoxDecryptKeyResponseParams
		keyId              string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("key_id"); ok {
		paramMap["KeyId"] = helper.String(v.(string))
		keyId = v.(string)
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsWhiteBoxDecryptKeyByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		whiteBoxDecryptKey = result
		return nil
	})

	if err != nil {
		return err
	}

	if whiteBoxDecryptKey.DecryptKey != nil {
		_ = d.Set("decrypt_key", whiteBoxDecryptKey.DecryptKey)
	}

	d.SetId(keyId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := tccommon.WriteToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
