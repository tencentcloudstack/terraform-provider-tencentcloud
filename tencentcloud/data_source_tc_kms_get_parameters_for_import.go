/*
Use this data source to query detailed information of kms get_parameters_for_import

Example Usage

```hcl
data "tencentcloud_kms_get_parameters_for_import" "example" {
  key_id             = "786aea8c-4aec-11ee-b601-525400281a45"
  wrapping_algorithm = "RSAES_OAEP_SHA_1"
  wrapping_key_spec  = "RSA_2048"
}
```
*/
package tencentcloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	kms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/kms/v20190118"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKmsGetParametersForImport() *schema.Resource {
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
	defer logElapsed("data_source.tencentcloud_kms_get_parameters_for_import.read")()
	defer inconsistentCheck(d, meta)()

	var (
		logId                  = getLogId(contextNil)
		ctx                    = context.WithValue(context.TODO(), logIdKey, logId)
		service                = KmsService{client: meta.(*TencentCloudClient).apiV3Conn}
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

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsGetParametersForImportByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
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
		if e := writeToFile(output.(string), d); e != nil {
			return e
		}
	}

	return nil
}
