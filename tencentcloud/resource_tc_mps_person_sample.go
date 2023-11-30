package tencentcloud

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mps "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mps/v20190612"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
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
				Description: "Material application scene, optional value:1. Recognition.Face: used for content recognition 2. Review.Face: used for inappropriate content identification 3. All: contains all of the above, equivalent to 1+2.",
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

	personId = *response.Response.Person.PersonId
	d.SetId(personId)

	return resourceTencentCloudMpsPersonSampleRead(d, meta)
}

func resourceTencentCloudMpsPersonSampleRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_person_sample.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}

	personId := d.Id()

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

	if personSample.UsageSet != nil {
		_ = d.Set("usages", personSample.UsageSet)
	}

	if personSample.Description != nil {
		_ = d.Set("description", personSample.Description)
	}

	if personSample.FaceInfoSet != nil {
		faceContents := []*string{}
		for _, faceInfo := range personSample.FaceInfoSet {
			url := faceInfo.Url
			res, err := http.Get(*url)
			if err != nil {
				return err
			}
			content, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return err
			}
			base64Encode := base64.StdEncoding.EncodeToString(content)
			faceContents = append(faceContents, &base64Encode)
		}
		_ = d.Set("face_contents", faceContents)
	}

	return nil
}

func resourceTencentCloudMpsPersonSampleUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_person_sample.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	request := mps.NewModifyPersonSampleRequest()

	personId := d.Id()

	needChange := false
	request.PersonId = &personId

	mutableArgs := []string{"name", "usages", "description", "face_contents"}

	for _, v := range mutableArgs {
		if d.HasChange(v) {
			needChange = true
			break
		}
	}

	if needChange {

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
			operationInfo := mps.AiSampleFaceOperation{}
			for i := range faceContentsSet {
				faceContents := faceContentsSet[i].(string)
				operationInfo.FaceContents = append(operationInfo.FaceContents, &faceContents)
			}
			operationInfo.Type = helper.String("reset")
			request.FaceOperationInfo = &operationInfo
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
	}

	return resourceTencentCloudMpsPersonSampleRead(d, meta)
}

func resourceTencentCloudMpsPersonSampleDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_mps_person_sample.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := MpsService{client: meta.(*TencentCloudClient).apiV3Conn}
	personId := d.Id()

	if err := service.DeleteMpsPersonSampleById(ctx, personId); err != nil {
		return err
	}

	return nil
}
