/*
Provides a resource to create a ci hot_link

Example Usage

```hcl
resource "tencentcloud_ci_hot_link" "hot_link" {
  bucket = "terraform-ci-xxxxxx"
  hot_link {
		url = &lt;nil&gt;
		type = &lt;nil&gt;

  }
}
```

Import

ci hot_link can be imported using the id, e.g.

```
terraform import tencentcloud_ci_hot_link.hot_link hot_link_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	ci "github.com/tencentyun/cos-go-sdk-v5"
	"log"
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
				Description: ".",
			},

			"hot_link": {
				Required:    true,
				Description: "Bucket name.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Domain address.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Anti-leech type, `white` is whitelist, `black` is blacklist, `off` is closed.",
						},
					},
				},
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

	hotLinkId := d.Id()

	hotLink, err := service.DescribeCiHotLinkById(ctx, bucket)
	if err != nil {
		return err
	}

	if hotLink == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `CiHotLink` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if hotLink.bucket != nil {
		_ = d.Set("bucket", hotLink.bucket)
	}

	if hotLink.hotLink != nil {
	}

	return nil
}

func resourceTencentCloudCiHotLinkUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_hot_link.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := ci.NewSetHotLinkRequest()

	hotLinkId := d.Id()

	request.Bucket = &bucket

	immutableArgs := []string{"bucket", "hot_link"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseCiClient().SetHotLink(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update ci hotLink failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudCiHotLinkRead(d, meta)
}

func resourceTencentCloudCiHotLinkDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_ci_hot_link.delete")()
	defer inconsistentCheck(d, meta)()

	return nil
}
