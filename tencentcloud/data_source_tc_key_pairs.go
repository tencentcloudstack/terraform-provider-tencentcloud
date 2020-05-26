/*
Use this data source to query key pairs.

Example Usage

```hcl
data "tencentcloud_key_pairs" "foo" {
  key_id = "skey-ie97i3ml"
}

data "tencentcloud_key_pairs" "name" {
  key_name = "^test$"
}
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"github.com/terraform-providers/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func dataSourceTencentCloudKeyPairs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTencentCloudKeyPairsRead,

		Schema: map[string]*schema.Schema{
			"key_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key_name", "project_id"},
				Description:   "ID of the key pair to be queried.",
			},
			"key_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"key_id"},
				Description:   "Name of the key pair to be queried. Support regular expression search, only `^` and `$` are supported.",
			},
			"project_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"key_id"},
				Description:   "Project id of the key pair to be queried.",
			},
			"result_output_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Used to save results.",
			},

			// computed
			"key_pair_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "An information list of key pair. Each element contains the following attributes:",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of the key pair.",
						},
						"key_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the key pair.",
						},
						"project_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Project id of the key pair.",
						},
						"public_key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "public key of the key pair.",
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation time of the key pair.",
						},
					},
				},
			},
		},
	}
}

func dataSourceTencentCloudKeyPairsRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("data_source.tencentcloud_key_pairs.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	keyId := d.Get("key_id").(string)
	keyName := d.Get("key_name").(string)
	name := keyName
	if keyName != "" {
		if name[0] == '^' {
			name = name[1:]
		}
		length := len(name)
		if length > 0 && name[length-1] == '$' {
			name = name[:length-1]
		}

		pattern := `^[a-zA-Z0-9_]+$`
		if match, _ := regexp.MatchString(pattern, name); !match {
			return fmt.Errorf("key_name only support letters, numbers, and _ : %s", keyName)
		}
	}

	var projectId *int
	if v, ok := d.GetOkExists("project_id"); ok {
		vv := v.(int)
		projectId = &vv
	}

	var keyPairs []*cvm.KeyPair
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		keyPairs, errRet = cvmService.DescribeKeyPairByFilter(ctx, keyId, name, projectId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}

	keyPairList := make([]map[string]interface{}, 0, len(keyPairs))
	ids := make([]string, 0, len(keyPairs))
	namePattern, _ := regexp.Compile(keyName)
	for _, keyPair := range keyPairs {
		if match := namePattern.MatchString(*keyPair.KeyName); !match {
			continue
		}
		mapping := map[string]interface{}{
			"key_id":      keyPair.KeyId,
			"key_name":    keyPair.KeyName,
			"project_id":  keyPair.ProjectId,
			"create_time": keyPair.CreatedTime,
		}
		if keyPair.PublicKey != nil {
			publicKey := *keyPair.PublicKey
			split := strings.Split(publicKey, " ")
			publicKey = strings.Join(split[0:len(split)-1], " ")
			mapping["public_key"] = publicKey
		}
		keyPairList = append(keyPairList, mapping)
		ids = append(ids, *keyPair.KeyId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("key_pair_list", keyPairList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set key pair list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	output, ok := d.GetOk("result_output_file")
	if ok && output.(string) != "" {
		if err := writeToFile(output.(string), keyPairList); err != nil {
			return err
		}
	}
	return nil
}
