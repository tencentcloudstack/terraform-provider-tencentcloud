/*
Provides a resource to create a ci hot_link

Example Usage

```hcl
resource "tencentcloud_ci_hot_link" "hot_link" {
	bucket = "terraform-ci-xxxxxx"
	url = ["10.0.0.1", "10.0.0.2"]
	type = "white"
}
```

Import

ci hot_link can be imported using the bucket, e.g.

```
terraform import tencentcloud_ci_hot_link.hot_link terraform-ci-xxxxxx
```
*/
package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	ci "github.com/tencentyun/cos-go-sdk-v5"
)

func resourceTencentCloudCiHotLink() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudCiHotLinkCreate,
		Read:   resourceTencentCloudCiHotLinkRead,
		Update: resourceTencentCloudCiHotLinkUpdate,
		Delete: resourceTencentCloudCiHotLinkDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "bucket name.",
			},

			"url": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "domain address.",
			},

			"type": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Anti-leech type, `white` is whitelist, `black` is blacklist.",
			},
		},
	}
}

func resourceTencentCloudCiHotLinkCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_hot_link.create")()
	defer inconsistentCheck(d, meta)()

	var bucket string
	if v, ok := d.GetOk("bucket"); ok {
		bucket = v.(string)
	} else {
		return errors.New("get bucket failed!")
	}

	d.SetId(bucket)

	return resourceTencentCloudCiHotLinkUpdate(d, meta)
}

func resourceTencentCloudCiHotLinkRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_hot_link.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := CiService{client: meta.(*TencentCloudClient).apiV3Conn}

	bucket := d.Id()

	hotLink, err := service.DescribeCiHotLinkById(ctx, bucket)
	if err != nil {
		return err
	}

	if hotLink == nil {
		d.SetId("")
		return fmt.Errorf("resource `hotLink` %s does not exist", bucket)
	}

	_ = d.Set("bucket", bucket)

	if len(hotLink.Url) > 0 {
		_ = d.Set("url", hotLink.Url)
	}

	if hotLink.Type != "" {
		_ = d.Set("type", hotLink.Type)
	}

	return nil
}

func resourceTencentCloudCiHotLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_hot_link.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Id()

	var hotLinkType string
	var url []string
	if v, _ := d.GetOk("type"); v != nil {
		hotLinkType = v.(string)
	}

	if v, ok := d.GetOk("url"); ok {
		urlList := v.(*schema.Set).List()
		for _, v := range urlList {
			url = append(url, v.(string))
		}
	}

	ciClient := meta.(*TencentCloudClient).apiV3Conn.UsePicClient(bucket)
	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := ciClient.CI.SetHotLink(ctx, &ci.HotLinkOptions{
			Type: hotLinkType,
			Url:  url,
		})
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response status [%s]\n", logId, "SetHotLink", bucket, result.Status)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create ci hotLink failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiHotLinkRead(d, meta)
}

func resourceTencentCloudCiHotLinkDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_hot_link.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	bucket := d.Id()

	ciClient := meta.(*TencentCloudClient).apiV3Conn.UsePicClient(bucket)
	_, err := ciClient.CI.SetHotLink(ctx, &ci.HotLinkOptions{
		Type: "off",
		Url:  []string{},
	})
	if err != nil {
		return err
	}

	return nil
}
