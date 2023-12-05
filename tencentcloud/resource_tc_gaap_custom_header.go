package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gaap "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/gaap/v20180529"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudGaapCustomHeader() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudGaapCustomHeaderCreate,
		Read:   resourceTencentCloudGaapCustomHeaderRead,
		Update: resourceTencentCloudGaapCustomHeaderUpdate,
		Delete: resourceTencentCloudGaapCustomHeaderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"rule_id": {
				Required:    true,
				Type:        schema.TypeString,
				Description: "Rule id.",
			},

			"headers": {
				Optional:    true,
				Type:        schema.TypeList,
				Description: "Headers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"header_name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Header name.",
						},
						"header_value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Header value.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudGaapCustomHeaderCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_custom_header.create")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	var ruleId string
	if v, ok := d.GetOk("rule_id"); ok {
		ruleId = v.(string)
	}
	headers := make([]*gaap.HttpHeaderParam, 0)
	if v, ok := d.GetOk("headers"); ok {
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			httpHeaderParam := &gaap.HttpHeaderParam{}
			if v, ok := dMap["header_name"]; ok {
				httpHeaderParam.HeaderName = helper.String(v.(string))
			}
			if v, ok := dMap["header_value"]; ok {
				httpHeaderParam.HeaderValue = helper.String(v.(string))
			}
			headers = append(headers, httpHeaderParam)
		}
	}

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		e := service.CreateCustomHeader(ctx, ruleId, headers)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create gaap customHeader failed, reason:%+v", logId, err)
		return err
	}

	d.SetId(ruleId)

	return resourceTencentCloudGaapCustomHeaderRead(d, meta)
}

func resourceTencentCloudGaapCustomHeaderRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_custom_header.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)

	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

	ruleId := d.Id()
	headers, err := service.DescribeGaapCustomHeader(ctx, ruleId)
	if err != nil {
		return err
	}

	_ = d.Set("rule_id", ruleId)

	headersList := []interface{}{}
	for _, header := range headers {
		headersMap := map[string]interface{}{}

		if header.HeaderName != nil {
			headersMap["header_name"] = header.HeaderName
		}

		if header.HeaderValue != nil {
			headersMap["header_value"] = header.HeaderValue
		}

		headersList = append(headersList, headersMap)
	}

	_ = d.Set("headers", headersList)

	return nil
}

func resourceTencentCloudGaapCustomHeaderUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_custom_header.update")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	ruleId := d.Id()

	immutableArgs := []string{"rule_id"}

	for _, v := range immutableArgs {
		if d.HasChange(v) {
			return fmt.Errorf("argument `%s` cannot be changed", v)
		}
	}

	if d.HasChange("headers") {
		headers := make([]*gaap.HttpHeaderParam, 0)
		if v, ok := d.GetOk("headers"); ok {
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				httpHeaderParam := &gaap.HttpHeaderParam{}
				if v, ok := dMap["header_name"]; ok {
					httpHeaderParam.HeaderName = helper.String(v.(string))
				}
				if v, ok := dMap["header_value"]; ok {
					httpHeaderParam.HeaderValue = helper.String(v.(string))
				}
				headers = append(headers, httpHeaderParam)
			}
		}
		service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}

		err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
			e := service.CreateCustomHeader(ctx, ruleId, headers)
			if e != nil {
				return retryError(e)
			}
			return nil
		})
		if err != nil {
			log.Printf("[CRITAL]%s create gaap customHeader failed, reason:%+v", logId, err)
			return err
		}
	}

	return resourceTencentCloudGaapCustomHeaderRead(d, meta)
}

func resourceTencentCloudGaapCustomHeaderDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_gaap_custom_header.delete")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)

	service := GaapService{client: meta.(*TencentCloudClient).apiV3Conn}
	ruleId := d.Id()

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		headers := make([]*gaap.HttpHeaderParam, 0)
		e := service.CreateCustomHeader(ctx, ruleId, headers)
		if e != nil {
			return retryError(e)
		}
		return nil
	})
	if err != nil {
		log.Printf("[CRITAL]%s create gaap customHeader failed, reason:%+v", logId, err)
		return err
	}

	return nil
}
