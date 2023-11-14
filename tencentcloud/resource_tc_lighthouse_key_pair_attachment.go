/*
Provides a resource to create a lighthouse key_pair_attachment

Example Usage

```hcl
resource "tencentcloud_lighthouse_key_pair_attachment" "key_pair_attachment" {
  key_ids =
  instance_ids =
}
```

Import

lighthouse key_pair_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_key_pair_attachment.key_pair_attachment key_pair_attachment_id
```
*/
package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	lighthouse "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/lighthouse/v20200324"
	"log"
	"strings"
	"time"
)

func resourceTencentCloudLighthouseKeyPairAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudLighthouseKeyPairAttachmentCreate,
		Read:   resourceTencentCloudLighthouseKeyPairAttachmentRead,
		Delete: resourceTencentCloudLighthouseKeyPairAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"key_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Key pair ID list. Each request can contain up to 100 key pairs.",
			},

			"instance_ids": {
				Required: true,
				ForceNew: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Instance ID list. Each request can contain up to 100 instances at a time.",
			},
		},
	}
}

func resourceTencentCloudLighthouseKeyPairAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_key_pair_attachment.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request    = lighthouse.NewAssociateInstancesKeyPairsRequest()
		response   = lighthouse.NewAssociateInstancesKeyPairsResponse()
		keyIds     string
		instanceId string
	)
	if v, ok := d.GetOk("key_ids"); ok {
		keyIdsSet := v.(*schema.Set).List()
		for i := range keyIdsSet {
			keyIds := keyIdsSet[i].(string)
			request.KeyIds = append(request.KeyIds, &keyIds)
		}
	}

	if v, ok := d.GetOk("instance_ids"); ok {
		instanceIdsSet := v.(*schema.Set).List()
		for i := range instanceIdsSet {
			instanceIds := instanceIdsSet[i].(string)
			request.InstanceIds = append(request.InstanceIds, &instanceIds)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseLighthouseClient().AssociateInstancesKeyPairs(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create lighthouse keyPairAttachment failed, reason:%+v", logId, err)
		return err
	}

	keyIds = *response.Response.KeyIds
	d.SetId(strings.Join([]string{keyIds, instanceId}, FILED_SP))

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseKeyPairAttachmentStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return resourceTencentCloudLighthouseKeyPairAttachmentRead(d, meta)
}

func resourceTencentCloudLighthouseKeyPairAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_key_pair_attachment.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	keyIds := idSplit[0]
	instanceId := idSplit[1]

	keyPairAttachment, err := service.DescribeLighthouseKeyPairAttachmentById(ctx, keyIds, instanceId)
	if err != nil {
		return err
	}

	if keyPairAttachment == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `LighthouseKeyPairAttachment` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if keyPairAttachment.KeyIds != nil {
		_ = d.Set("key_ids", keyPairAttachment.KeyIds)
	}

	if keyPairAttachment.InstanceIds != nil {
		_ = d.Set("instance_ids", keyPairAttachment.InstanceIds)
	}

	return nil
}

func resourceTencentCloudLighthouseKeyPairAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_lighthouse_key_pair_attachment.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}
	idSplit := strings.Split(d.Id(), FILED_SP)
	if len(idSplit) != 2 {
		return fmt.Errorf("id is broken,%s", d.Id())
	}
	keyIds := idSplit[0]
	instanceId := idSplit[1]

	if err := service.DeleteLighthouseKeyPairAttachmentById(ctx, keyIds, instanceId); err != nil {
		return err
	}

	service := LighthouseService{client: meta.(*TencentCloudClient).apiV3Conn}

	conf := BuildStateChangeConf([]string{}, []string{"SUCCESS"}, 20*readRetryTimeout, time.Second, service.LighthouseKeyPairAttachmentStateRefreshFunc(d.Id(), []string{}))

	if _, e := conf.WaitForState(); e != nil {
		return e
	}

	return nil
}
