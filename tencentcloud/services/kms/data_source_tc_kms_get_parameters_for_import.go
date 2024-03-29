package kms

import (
	"context"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func DataSourceTencentCloudKmsGetParametersForImport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsGetParametersForImportRead,
		Schema: map[string]*schema.Schema{
			"key_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "CMK unique identifier.",
			},
			"wrapping_algorithm": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Specifies the algorithm for encrypting key material, currently supports RSAES_PKCS1_V1_5, RSAES_OAEP_SHA_1, RSAES_OAEP_SHA_256.",
			},
			"wrapping_key_spec": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Specifies the type of encryption key material, currently only supports RSA_2048.",
			},
			"public_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Base64-encoded public key content.",
			},
			"import_token": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The token required for importing key material is used as the parameter of ImportKeyMaterial.",
			},
			"parameters_valid_to": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The validity period of the exported token and public key cannot be imported after this period, and you need to call GetParametersForImport again to obtain it.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudKmsGetParametersForImportRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("data_source.tencentcloud_kms_get_parameters_for_import.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	var (
		logId                  = tccommon.GetLogId(tccommon.ContextNil)
		ctx                    = context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
		service                = KmsService{client: meta.(tccommon.ProviderMeta).GetAPIV3Conn()}
		getParametersForImport *kms.GetParametersForImportResponseParams
		keyId                  string
	)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("key_id"); ok {
		paramMap["KeyId"] = helper.String(v.(string))
		keyId = v.(string)
	}

	if v, ok := d.GetOk("wrapping_algorithm"); ok {
		paramMap["WrappingAlgorithm"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("wrapping_key_spec"); ok {
		paramMap["WrappingKeySpec"] = helper.String(v.(string))
	}

	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsGetParametersForImportByFilter(ctx, paramMap)
		if e != nil {
			return tccommon.RetryError(e)
		}

		getParametersForImport = result
		return nil
	})

	if err != nil {
		return err
	}

	if getParametersForImport.KeyId != nil {
		_ = d.Set("key_id", getParametersForImport.KeyId)
	}

	if getParametersForImport.PublicKey != nil {
		_ = d.Set("public_key", getParametersForImport.PublicKey)
	}

	if getParametersForImport.ImportToken != nil {
		_ = d.Set("import_token", getParametersForImport.ImportToken)
	}

	if getParametersForImport.ParametersValidTo != nil {
		_ = d.Set("parameters_valid_to", getParametersForImport.ParametersValidTo)
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
