/*
Provides a key pair resource.

Example Usage

```hcl
resource "tencentcloud_key_pair" "foo" {
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
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "You can import an existing public key and using TencentCloud key pair to manage it.",
			},
			"project_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				ForceNew:    true,
				Description: "Specifys to which project the key pair belongs.",
			},
		},
	}
}

func resourceTencentCloudKeyPairCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_key_pair.create")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	keyName := d.Get("key_name").(string)
	publicKey := d.Get("public_key").(string)
	projectId := d.Get("project_id").(int)

	keyId := ""
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		id, err := cvmService.CreateKeyPair(ctx, keyName, publicKey, int64(projectId))
		if err != nil {
			return retryError(err)
		}
		keyId = id
		return nil
	})
	if err != nil {
		return err
	}
	d.SetId(keyId)

	return resourceTencentCloudKeyPairRead(d, meta)
}

func resourceTencentCloudKeyPairRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_key_pair.read")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	keyId := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var keyPair *cvm.KeyPair
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		keyPair, errRet = cvmService.DescribeKeyPairById(ctx, keyId)
		if errRet != nil {
			if errRet.Error() == "key pair id not found" {
				d.SetId("")
				return nil
			}
			return retryError(errRet, "InternalError")
		}
		return nil
	})
	if err != nil {
		return err
	}

	d.Set("key_name", keyPair.KeyName)
	d.Set("project_id", keyPair.ProjectId)
	if keyPair.PublicKey != nil {
		publicKey := *keyPair.PublicKey
		split := strings.Split(publicKey, " ")
		publicKey = strings.Join(split[0:len(split)-1], " ")
		d.Set("public_key", publicKey)
	}

	return nil
}

func resourceTencentCloudKeyPairUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_key_pair.update")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

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

	return resourceTencentCloudKeyPairRead(d, meta)
}

func resourceTencentCloudKeyPairDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_key_pair.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), "logId", logId)

	keyId := d.Id()
	cvmService := CvmService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}

	var keyPair *cvm.KeyPair
	var errRet error
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		keyPair, errRet = cvmService.DescribeKeyPairById(ctx, keyId)
		if errRet != nil {
			if errRet.Error() == "key pair id not found" {
				keyId = ""
				return nil
			}
			return retryError(errRet, "InternalError")
		}
		return nil
	})
	if err != nil {
		return err
	}
	if keyId == "" {
		return nil
	}

	if len(keyPair.AssociatedInstanceIds) > 0 {
		err = resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			errRet := cvmService.UnbindKeyPair(ctx, keyId, keyPair.AssociatedInstanceIds)
			if errRet != nil {
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
			return retryError(errRet)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
