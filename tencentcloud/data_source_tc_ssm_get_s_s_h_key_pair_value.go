/*
Use this data source to query detailed information of ssm get_s_s_h_key_pair_value

Example Usage

```hcl
data "tencentcloud_ssm_get_s_s_h_key_pair_value" "get_s_s_h_key_pair_value" {
  secret_name = ""
  s_s_h_key_id = ""
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

func dataSourceTencentCloudSsmGetSSHKeyPairValue() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSsmGetSSHKeyPairValueRead,
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Secret name.",
			},

			"s_s_h_key_id": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "The key pair ID is the unique identifier of the key pair in the cloud server.",
			},

			"s_s_h_key_i_d": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "SSH key pair ID.",
			},

			"s_s_h_key_name": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "The name of the SSH key pair. Users can modify the name of the key pair in the CVM console.",
			},

			"public_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Public key plain text, encoded using base64.",
			},

			"private_key": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Private key plain text, encoded using base64.",
			},

			"project_i_d": {
				Computed:    true,
				Type:        schema.TypeInt,
				Description: "The project ID to which this key pair belongs.",
			},

			"s_s_h_key_description": {
				Computed:    true,
				Type:        schema.TypeString,
				Description: "Description of the SSH key pair. Users can modify the description information of the key pair in the CVM console.",
			},

			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
		},
	}
}

func dataSourceTencentCloudSsmGetSSHKeyPairValueRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssm_get_s_s_h_key_pair_value.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	paramMap := make(map[string]interface{})
	if v, ok := d.GetOk("secret_name"); ok {
		paramMap["SecretName"] = helper.String(v.(string))
	}

	if v, ok := d.GetOk("s_s_h_key_id"); ok {
		paramMap["SSHKeyId"] = helper.String(v.(string))
	}

	service := SsmService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := service.DescribeSsmGetSSHKeyPairValueByFilter(ctx, paramMap)
		if e != nil {
			return retryError(e)
		}
		sSHKeyID = result
		return nil
	})
	if err != nil {
		return err
	}

	ids := make([]string, 0, len(sSHKeyID))
	if sSHKeyID != nil {
		_ = d.Set("s_s_h_key_i_d", sSHKeyID)
	}

	if sSHKeyName != nil {
		_ = d.Set("s_s_h_key_name", sSHKeyName)
	}

	if publicKey != nil {
		_ = d.Set("public_key", publicKey)
	}

	if privateKey != nil {
		_ = d.Set("private_key", privateKey)
	}

	if projectID != nil {
		_ = d.Set("project_i_d", projectID)
	}

	if sSHKeyDescription != nil {
		_ = d.Set("s_s_h_key_description", sSHKeyDescription)
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
