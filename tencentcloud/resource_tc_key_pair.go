/*
Provides a key pair resource.

Example Usage

```hcl
resource "tencentcloud_key_pair" "foo" {
	key_name   = "terraform_test"
}

resource "tencentcloud_key_pair" "foo1" {
  key_name   = "terraform_test"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQDjd8fTnp7Dcuj4mLaQxf9Zs/ORgUL9fQxRCNKkPgP1paTy1I513maMX126i36Lxxl3+FUB52oVbo/FgwlIfX8hyCnv8MCxqnuSDozf1CD0/wRYHcTWAtgHQHBPCC2nJtod6cVC3kB18KeV4U7zsxmwFeBIxojMOOmcOBuh7+trRw=="
}
```

Import

Key pair can be imported using the id, e.g.

```
$ terraform import tencentcloud_key_pair.foo skey-17634f05
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func resourceTencentCloudKeyPair() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudKeyPairCreate,
		Read:   resourceTencentCloudKeyPairRead,
		Update: resourceTencentCloudKeyPairUpdate,
		Delete: resourceTencentCloudKeyPairDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateKeyPairName,
				Description:  "The key pair's name. It is the only in one TencentCloud account.",
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					switch value := v.(type) {
					case string:
						publicKey := value
						split := strings.Split(value, " ")
						if len(split) > 2 {
							publicKey = strings.Join(split[0:2], " ")
						}
						return strings.TrimSpace(publicKey)
					default:
						return ""
					}
				},
				Description: "You can import an existing public key and using TencentCloud key pair to manage it.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				ForceNew:    true,
				Description: "Specifys to which project the key pair belongs.",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags of the key pair.",
			},
		},
	}
}

func cvmCreateKeyPair(ctx context.Context, d *schema.ResourceData, meta interface{}) (keyId string, err error) {
	logId := getLogId(ctx)
	request := cvm.NewCreateKeyPairRequest()
	response := cvm.NewCreateKeyPairResponse()
	request.KeyName = helper.String(d.Get("key_name").(string))
	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))

	innerErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().CreateKeyPair(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if innerErr != nil {
		log.Printf("[CRITAL]%s create cvm keyPair by import failed, reason:%+v", logId, err)
		err = innerErr
		return
	}
	if response == nil || response.Response == nil || response.Response.KeyPair == nil {
		err = fmt.Errorf("Response is nil")
		return
	}

	keyId = *response.Response.KeyPair.KeyId
	return
}

func cvmCreateKeyPairByImportPublicKey(ctx context.Context, d *schema.ResourceData, meta interface{}) (keyId string, err error) {
	logId := getLogId(ctx)
	request := cvm.NewImportKeyPairRequest()
	response := cvm.NewImportKeyPairResponse()
	request.KeyName = helper.String(d.Get("key_name").(string))
	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
	request.PublicKey = helper.String(d.Get("public_key").(string))

	innerErr := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCvmClient().ImportKeyPair(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if innerErr != nil {
		log.Printf("[CRITAL]%s create cvm keyPair by import failed, reason:%+v", logId, err)
		err = innerErr
		return
	}
	if response == nil || response.Response == nil {
		err = fmt.Errorf("Response is nil")
		return
	}

	keyId = *response.Response.KeyId
	return
}

func resourceTencentCloudKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_key_pair.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	var (
		keyId string
		err   error
	)

	if _, ok := d.GetOk("public_key"); ok {
		keyId, err = cvmCreateKeyPairByImportPublicKey(ctx, d, meta)
	} else {
		keyId, err = cvmCreateKeyPair(ctx, d, meta)
	}
	if err != nil {
		return err
	}
	d.SetId(keyId)

	if tags := helper.GetTags(d, "tags"); len(tags) > 0 {
		tcClient := meta.(*TencentCloudClient).apiV3Conn
		tagService := &TagService{client: tcClient}
		resourceName := BuildTagResourceName("cvm", "keypair", tcClient.Region, keyId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudKeyPairRead(d, meta)
}

func resourceTencentCloudKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_key_pair.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	keyId := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var keyPair *cvm.KeyPair
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		keyPair, errRet = cvmService.DescribeKeyPairById(ctx, keyId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if keyPair == nil {
		d.SetId("")
		return nil
	}

	_ = d.Set("key_name", keyPair.KeyName)
	_ = d.Set("project_id", keyPair.ProjectId)
	if keyPair.PublicKey != nil {
		publicKey := *keyPair.PublicKey
		split := strings.Split(publicKey, " ")
		if len(split) > 2 {
			publicKey = strings.Join(split[0:2], " ")
		}
		_ = d.Set("public_key", publicKey)
	}

	client := meta.(*TencentCloudClient).apiV3Conn
	tagService := TagService{client}

	tags, err := tagService.DescribeResourceTags(ctx, "cvm", "keypair", client.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_key_pair.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	keyId := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	if d.HasChange("key_name") {
		keyName := d.Get("key_name").(string)
		err := cvmService.ModifyKeyPairName(ctx, keyId, keyName)
		if err != nil {
			return err
		}
	}

	if d.HasChange("tags") {
		oldInterface, newInterface := d.GetChange("tags")
		replaceTags, deleteTags := diffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := TagService{
			client: meta.(*TencentCloudClient).apiV3Conn,
		}
		region := meta.(*TencentCloudClient).apiV3Conn.Region
		resourceName := BuildTagResourceName("cvm", "keypair", region, keyId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudKeyPairRead(d, meta)
}

func resourceTencentCloudKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_key_pair.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	keyId := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var keyPair *cvm.KeyPair
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		keyPair, errRet = cvmService.DescribeKeyPairById(ctx, keyId)
		if errRet != nil {
			return retryError(errRet, InternalError)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if keyPair == nil {
		d.SetId("")
		return nil
	}

	if len(keyPair.AssociatedInstanceIds) > 0 {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := cvmService.UnbindKeyPair(ctx, []*string{&keyId}, keyPair.AssociatedInstanceIds)
			if errRet != nil {
				if sdkErr, ok := errRet.(*errors.TencentCloudSDKError); ok {
					if sdkErr.Code == CVM_NOT_FOUND_ERROR {
						return nil
					}
				}
				return retryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		errRet := cvmService.DeleteKeyPair(ctx, keyId)
		if errRet != nil {
			return retryError(errRet, KYE_PAIR_INVALID_ERROR, KEY_PAIR_NOT_SUPPORT_ERROR)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
