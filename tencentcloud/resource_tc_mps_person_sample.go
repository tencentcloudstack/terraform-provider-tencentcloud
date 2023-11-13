/*
Provides a resource to create a mps person_sample

Example Usage

```hcl
resource "tencentcloud_mps_person_sample" "person_sample" {
  name = &lt;nil&gt;
  usages = &lt;nil&gt;
  description = &lt;nil&gt;
  face_contents = &lt;nil&gt;
}
```

Import

mps person_sample can be imported using the id, e.g.

```
terraform import tencentcloud_mps_person_sample.person_sample person_sample_id
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
)

func resourceTencentCloudMpsPersonSample() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudMpsPersonSampleCreate,
		Read:   resourceTencentCloudMpsPersonSampleRead,
		Update: resourceTencentCloudMpsPersonSampleUpdate,
		Delete: resourceTencentCloudMpsPersonSampleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Material name, length limit: 20 characters.",
			},

			"usages": {
				Required: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Material application scene, optional value:1. Recognition: used for content recognition, equivalent to Recognition.Face.2. Review: used for inappropriate content identification, equivalent to Review.Face.3. All: contains all of the above, equivalent to 1+2.",
			},

			"description": {
				Optional:    true,
				Type:        schema.TypeString,
				Description: "Material description, length limit: 1024 characters.",
			},

			"face_contents": {
				Optional: true,
				Type:     schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Material image [Base64](https://tools.ietf.org/html/rfc4648) encoded string only supports jpeg and png image formats. Array length limit: 5 images.Note: The picture must be a single portrait with clearer facial features, with a pixel size of not less than 200*200.",
			},
		},
	}
}

func resourceTencentCloudMpsPersonSampleCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_person_sample.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	var (
		request  = mps.NewCreatePersonSampleRequest()
		response = mps.NewCreatePersonSampleResponse()
		personId string
	)
	if v, ok := d.GetOk("name"); ok {
		request.Name = helper.String(v.(string))
	}

	if v, ok := d.GetOk("usages"); ok {
		usagesSet := v.(*schema.Set).List()
		for i := range usagesSet {
			usages := usagesSet[i].(string)
			request.Usages = append(request.Usages, &usages)
		}
	}

	if v, ok := d.GetOk("description"); ok {
		request.Description = helper.String(v.(string))
	}

	if v, ok := d.GetOk("face_contents"); ok {
		faceContentsSet := v.(*schema.Set).List()
		for i := range faceContentsSet {
			faceContents := faceContentsSet[i].(string)
			request.FaceContents = append(request.FaceContents, &faceContents)
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().CreatePersonSample(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create mps personSample failed, reason:%+v", logId, err)
		return err
	}

	personId = *response.Response.PersonId
	d.SetId(personId)

	return resourceTencentCloudMpsPersonSampleRead(d, meta)
}

func resourceTencentCloudMpsPersonSampleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_person_sample.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	personSampleId := d.Id()

	personSample, err := service.DescribeMpsPersonSampleById(ctx, personId)
	if err != nil {
		return err
	}

	if personSample == nil {
		d.SetId("")
		log.Printf("[WARN]%s resource `MpsPersonSample` [%s] not found, please check if it has been deleted.\n", logId, d.Id())
		return nil
	}

	if personSample.Name != nil {
		_ = d.Set("name", personSample.Name)
	}

	if personSample.Usages != nil {
		_ = d.Set("usages", personSample.Usages)
	}

	if personSample.Description != nil {
		_ = d.Set("description", personSample.Description)
	}

	if personSample.FaceContents != nil {
		_ = d.Set("face_contents", personSample.FaceContents)
	}

	return nil
}

func resourceTencentCloudMpsPersonSampleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_person_sample.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyPersonSampleRequest()

	personSampleId := d.Id()

	request.PersonId = &personId

	immutableArgs := []string{"name", "usages", "description", "face_contents"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("name") {
		if v, ok := d.GetOk("name"); ok {
			request.Name = helper.String(v.(string))
		}
	}

	if d.HasChange("usages") {
		if v, ok := d.GetOk("usages"); ok {
			usagesSet := v.(*schema.Set).List()
			for i := range usagesSet {
				usages := usagesSet[i].(string)
				request.Usages = append(request.Usages, &usages)
			}
		}
	}

	if d.HasChange("description") {
		if v, ok := d.GetOk("description"); ok {
			request.Description = helper.String(v.(string))
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseMpsClient().ModifyPersonSample(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update mps personSample failed, reason:%+v", logId, err)
		return err
	}

	return resourceTencentCloudMpsPersonSampleRead(d, meta)
}

func resourceTencentCloudMpsPersonSampleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_person_sample.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	personSampleId := d.Id()

	if err := service.DeleteMpsPersonSampleById(ctx, personId); err != nil {
		return err
	}

	return nil
}
