package cvm

import (
	"context"
	"fmt"
	"log"
	"strings"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	svctag "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tag"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func ResourceTencentCloudKeyPair() *schema.Resource {
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
				ValidateFunc: tccommon.ValidateKeyPairName,
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
	logId := tccommon.GetLogId(ctx)
	request := cvm.NewCreateKeyPairRequest()
	response := cvm.NewCreateKeyPairResponse()
	request.KeyName = helper.String(d.Get("key_name").(string))
	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))

	innerErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().CreateKeyPair(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	logId := tccommon.GetLogId(ctx)
	request := cvm.NewImportKeyPairRequest()
	response := cvm.NewImportKeyPairResponse()
	request.KeyName = helper.String(d.Get("key_name").(string))
	request.ProjectId = helper.IntInt64(d.Get("project_id").(int))
	request.PublicKey = helper.String(d.Get("public_key").(string))

	innerErr := resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		result, e := meta.(tccommon.ProviderMeta).GetAPIV3Conn().UseCvmClient().ImportKeyPair(request)
		if e != nil {
			return tccommon.RetryError(e)
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
	defer tccommon.LogElapsed("resource.tencentcloud_key_pair.create")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)
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
		tcClient := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
		tagService := svctag.NewTagService(tcClient)
		resourceName := tccommon.BuildTagResourceName("cvm", "keypair", tcClient.Region, keyId)
		if err := tagService.ModifyTags(ctx, resourceName, tags, nil); err != nil {
			return err
		}
	}

	return resourceTencentCloudKeyPairRead(d, meta)
}

func resourceTencentCloudKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_key_pair.read")()
	defer tccommon.InconsistentCheck(d, meta)()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	keyId := d.Id()
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}
	var keyPair *cvm.KeyPair
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		keyPair, errRet = cvmService.DescribeKeyPairById(ctx, keyId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
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

	client := meta.(tccommon.ProviderMeta).GetAPIV3Conn()
	tagService := svctag.NewTagService(client)

	tags, err := tagService.DescribeResourceTags(ctx, "cvm", "keypair", client.Region, d.Id())
	if err != nil {
		return err
	}
	_ = d.Set("tags", tags)

	return nil
}

func resourceTencentCloudKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_key_pair.update")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	keyId := d.Id()
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
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
		replaceTags, deleteTags := svctag.DiffTags(oldInterface.(map[string]interface{}), newInterface.(map[string]interface{}))
		tagService := svctag.NewTagService(meta.(tccommon.ProviderMeta).GetAPIV3Conn())
		region := meta.(tccommon.ProviderMeta).GetAPIV3Conn().Region
		resourceName := tccommon.BuildTagResourceName("cvm", "keypair", region, keyId)
		err := tagService.ModifyTags(ctx, resourceName, replaceTags, deleteTags)
		if err != nil {
			return err
		}
	}

	return resourceTencentCloudKeyPairRead(d, meta)
}

func resourceTencentCloudKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	defer tccommon.LogElapsed("resource.tencentcloud_key_pair.delete")()

	logId := tccommon.GetLogId(tccommon.ContextNil)
	ctx := context.WithValue(context.TODO(), tccommon.LogIdKey, logId)

	keyId := d.Id()
	cvmService := CvmService{
		client: meta.(tccommon.ProviderMeta).GetAPIV3Conn(),
	}

	var keyPair *cvm.KeyPair
	var errRet error
	err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
		keyPair, errRet = cvmService.DescribeKeyPairById(ctx, keyId)
		if errRet != nil {
			return tccommon.RetryError(errRet, tccommon.InternalError)
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
		err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
			errRet := cvmService.UnbindKeyPair(ctx, []*string{&keyId}, keyPair.AssociatedInstanceIds)
			if errRet != nil {
				if sdkErr, ok := errRet.(*errors.TencentCloudSDKError); ok {
					if sdkErr.Code == CVM_NOT_FOUND_ERROR {
						return nil
					}
				}
				return tccommon.RetryError(errRet)
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	err = resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		errRet := cvmService.DeleteKeyPair(ctx, keyId)
		if errRet != nil {
			return tccommon.RetryError(errRet, KYE_PAIR_INVALID_ERROR, KEY_PAIR_NOT_SUPPORT_ERROR)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
