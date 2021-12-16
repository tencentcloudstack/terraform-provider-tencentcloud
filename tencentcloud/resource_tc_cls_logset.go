package tencentcloud

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"time"
)

func resourceTencentCloudClsLogSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsLogSetCreate,
		Read:   resourceTencentCloudClsLogRead,
		Update: resourceTencentCloudClsLogUpdate,
		Delete: resourceTencentCloudClsLogDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(1 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"logset_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateNotEmpty,
				Description:  "Name of the LogSet.",
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    10,
				Description: "The label key value pair of the binding for the log set. A maximum of 10 tag key value pairs are supported. ",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag Key",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Tag value",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClsLogSetCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("source.tencentcloud_cls_logset.create")()
	logId := getLogId(contextNil)
	config := meta.(*TencentCloudClient).apiV3Conn
	request := cls.NewCreateLogsetRequest()
	request.LogsetName = common.StringPtr(d.Get("logset_name").(string))
	if v, ok := d.GetOk("tags"); ok {
		if tags, ok := v.([]interface{}); ok {
			request.Tags = make([]*cls.Tag, 0, len(tags))
			for _, value := range tags {
				if tag, ok := value.(map[string]interface{}); ok {
					tagGet := cls.Tag{
						Key:   helper.String(tag["key"].(string)),
						Value: helper.String(tag["value"].(string)),
					}
					request.Tags = append(request.Tags, &tagGet)
				}
			}
		}
	}
	var logset *cls.CreateLogsetResponse
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		response, e := config.UseClsClient().CreateLogset(request)
		if e != nil {
			return retryError(e)
		}
		logset = response
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create CLS Logset failed, reason:%+v", logId, err)
		return err
	}

	id := logset.Response.LogsetId
	d.SetId(*id)
	return resourceTencentCloudClsLogRead(d, meta)
}
func resourceTencentCloudClsLogRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("source.tencentcloud_cls_logset.read")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	clsService := ClsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var logset *cls.LogsetInfo
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		result, e := clsService.DescribeClsLogSetById(ctx, d.Id())
		if e != nil {
			return retryError(e, InternalError)
		}
		logset = result
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s read CLS logsets failed, reason:%s\n", logId, err.Error())
		return err
	}
	if logset == nil || logset.LogsetId == nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("logset_name", *logset.LogsetName)
	tags := make([]map[string]interface{}, 0, len(logset.Tags))
	for _, tag := range logset.Tags {
		mapping := map[string]interface{}{
			"key":   *tag.Key,
			"value": *tag.Value,
		}
		tags = append(tags, mapping)
	}
	if logset.Tags != nil {
		_ = d.Set("tags", tags)
	}
	return nil
}
func resourceTencentCloudClsLogDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("source.tencentcloud_cls_logset.delete")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	request := cls.NewDeleteLogsetRequest()
	request.LogsetId = common.StringPtr(d.Id())
	clsService := ClsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := resource.Retry(readRetryTimeout, func() *resource.RetryError {
		e := clsService.DeleteClsLogSet(ctx, d.Id())
		if e != nil {
			return retryError(e, InternalError)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s delete CLS logsets failed, reason:%s\n", logId, err.Error())
		return err
	}
	return resourceTencentCloudClsLogRead(d, meta)
}

func resourceTencentCloudClsLogUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("source.tencentcloud_cls_logset.update")()
	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	request := cls.NewModifyLogsetRequest()
	request.LogsetId = common.StringPtr(d.Id())
	if v, ok := d.GetOk("logset_name"); ok {
		request.LogsetName = common.StringPtr(v.(string))
	}
	if v, ok := d.GetOk("tags"); ok {
		if tags, ok := v.([]interface{}); ok {
			request.Tags = make([]*cls.Tag, 0, len(tags))
			for _, value := range tags {
				if tag, ok := value.(map[string]interface{}); ok {
					tagGet := cls.Tag{
						Key:   helper.String(tag["key"].(string)),
						Value: helper.String(tag["value"].(string)),
					}
					request.Tags = append(request.Tags, &tagGet)
				}
			}
		}
	}
	clsService := ClsService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	var readRetryTimeoutTrue = 5 * time.Second
	err := resource.Retry(readRetryTimeoutTrue, func() *resource.RetryError {
		e := clsService.UpdateClsLogSet(ctx, request)
		if e != nil {
			return retryError(e, InternalError)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s update CLS logsets failed, reason:%s\n", logId, err.Error())
		return err
	}
	return resourceTencentCloudClsLogRead(d, meta)
}
