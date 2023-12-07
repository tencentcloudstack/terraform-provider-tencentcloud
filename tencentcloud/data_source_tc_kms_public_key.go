package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKmsPublicKey() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsPublicKeyRead,
		Schema: map[string]*schema.Schema{
			"key_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "CMK unique identifier.",
			},
			"public_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Base64-encoded public key content.",
			},
			"public_key_pem": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Public key content in PEM format.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudKmsPublicKeyRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_kms_public_key.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId     = getLogId(contextNil)
		ctx       = context.WithValue(context.TODO(), logIdKey, logId)
		service   = KmsService{client: meta.(*TencentCloudClient).apiV3Conn}
		publicKey *kms.GetPublicKeyResponseParams
		keyId     string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("key_id"); ok {
		paramMap["KeyId"] = helper.String(v.(string))
		keyId = v.(string)
	}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsPublicKeyByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}

		publicKey = result
		return nil
	})

	if err != nil {
		return err
	}

	if publicKey.KeyId != nil {
		_ = d.Set("key_id", publicKey.KeyId)
	}

	if publicKey.PublicKey != nil {
		_ = d.Set("public_key", publicKey.PublicKey)
	}

	if publicKey.PublicKeyPem != nil {
		_ = d.Set("public_key_pem", publicKey.PublicKeyPem)
	}

	d.SetId(keyId)
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
