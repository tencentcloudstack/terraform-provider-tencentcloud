/*
Use this data source to query detailed information of kms get_parameters_for_import

Example Usage

```hcl
data "tencentcloud_kms_get_parameters_for_import" "get_parameters_for_import" {
      wrapping_algorithm = "RSAES_OAEP_SHA_1"
  wrapping_key_spec = "RSA_2048"
    }
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKmsGetParametersForImport() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKmsGetParametersForImportRead,
		Schema: map[string]*schema.Schema{
			"key_id": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "CMK unique identifier.",
			},

			"public_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Base64-encoded public key content.",
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

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("wrapping_algorithm"); ok {
		paramMap["WrappingAlgorithm"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("wrapping_key_spec"); ok {
		paramMap["WrappingKeySpec"] = helper.String(v.(string))
	}

	service := KmsService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeKmsGetParametersForImportByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		keyId = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(keyId))
	if keyId != nil {
		_ = d.Set("key_id", keyId)
	}

	if publicKey != nil {
		_ = d.Set("public_key", publicKey)
	}

	if importToken != nil {
		_ = d.Set("import_token", importToken)
	}

	if parametersValidTo != nil {
		_ = d.Set("parameters_valid_to", parametersValidTo)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if e := writeToFile(output.(string)); e != nil {
			return e
		}
	}
	return nil
}
