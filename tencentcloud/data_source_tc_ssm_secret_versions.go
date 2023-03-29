/*
Use this data source to query detailed information of SSM secret version
Example Usage
```hcl

data "tencentcloud_ssm_secret_versions" "foo" {
  secret_name = "test"
  version_id = "v1"
}
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	sdkError "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudSsmSecretVersions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudSsmSecretVersionsRead,
		Schema: map[string]*schema.Schema{
			"secret_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Secret name used to filter result.",
			},
			"version_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "VersionId used to filter result.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},
			"secret_version_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of SSM secret versions. When secret status is `Disabled`, this field will not update anymore.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of secret.",
						},
						"secret_binary": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The base64-encoded binary secret.",
						},
						"secret_string": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The string text of secret.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudSsmSecretVersionsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_ssm_secret_versions.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	ssmService := SsmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	secretName := d.Get("secret_name").(string)
	var outErr, inErr error
	var secretInfo *SecretInfo
	outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
		secretInfo, inErr = ssmService.DescribeSecretByName(ctx, secretName)
		if inErr != nil {
			return retryError(inErr)
		}
		return nil
	})
	if outErr != nil {
		sdkErr, ok := outErr.(*sdkError.TencentCloudSDKError)
		if ok && sdkErr.Code == SSMResourceNotFound {
			d.SetId("")
			log.Printf("[WARN]%s resource `secretInfo` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
			return nil
		}

		log.Printf("[CRITAL]%s read SSM secret failed, reason:%+v", logId, outErr)
		return outErr
	}
	if secretInfo.status != SSM_STATUS_ENABLED {
		log.Printf("[CRITAL]%s read SSM secret version failed, reason: secret status is not Enabled", logId)
		return nil
	}
	var secretVersionInfos []*SecretVersionInfo
	var versionIds []string
	if v, ok := d.GetOk("version_id"); ok {
		versionIds = append(versionIds, v.(string))
	} else {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			versionIds, inErr = ssmService.DescribeSecretVersionIdsByName(ctx, secretName)
			if inErr != nil {
				return retryError(inErr)
			}
			return nil
		})
		if outErr != nil {
			log.Printf("[CRITAL]%s read SSM secret versionId list failed, reason:%+v", logId, outErr)
			return outErr
		}
	}

	for _, versionId := range versionIds {
		outErr = resource.Retry(readRetryTimeout, func() *resource.RetryError {
			secretVersionInfo, inErr := ssmService.DescribeSecretVersion(ctx, secretName, versionId)
			if inErr != nil {
				return retryError(inErr)
			}
			secretVersionInfos = append(secretVersionInfos, secretVersionInfo)
			return nil
		})
		if outErr != nil {
			log.Printf("[CRITAL]%s read SSM secret version failed, reason:%+v", logId, outErr)
			return outErr
		}
	}

	var secretVersionList []map[string]interface{}
	var ids []string
	for _, secretVersionInfo := range secretVersionInfos {
		mapping := map[string]interface{}{
			"version_id":    secretVersionInfo.versionId,
			"secret_binary": secretVersionInfo.secretBinary,
			"secret_string": secretVersionInfo.secretString,
		}

		secretVersionList = append(secretVersionList, mapping)
		ids = append(ids, strings.Join([]string{secretVersionInfo.secretName, secretVersionInfo.versionId}, FILED_SP))
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	if e := d.Set("secret_version_list", secretVersionList); e != nil {
		log.Printf("[CRITAL]%s provider set SSM secret version list fail, reason:%+v", logId, e)
		return e
	}
	if output, ok := d.GetOk("result_output_file"); ok && output.(string) != "" {
		return writeToFile(output.(string), secretVersionList)
	}
	return nil
}
