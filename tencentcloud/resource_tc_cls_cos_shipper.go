/*
Provides a resource to create a cls cos shipper.

Example Usage

```hcl
resource "tencentcloud_cls_cos_shipper" "shipper" {
  bucket       = "preset-scf-bucket-1308919341"
  interval     = 300
  max_size     = 200
  partition    = "/%Y/%m/%d/%H/"
  prefix       = "ap-guangzhou-fffsasad-1649734752"
  shipper_name = "ap-guangzhou-fffsasad-1649734752"
  topic_id     = "4d07fba0-b93e-4e0b-9a7f-d58542560bbb"

  compress {
    format = "lzop"
  }

  content {
    format = "json"

    json {
      enable_tag  = true
      meta_fields = [
        "__FILENAME__",
        "__SOURCE__",
        "__TIMESTAMP__",
      ]
    }
  }
}
```

Import

cls cos shipper can be imported using the id, e.g.

```
$ terraform import tencentcloud_cls_cos_shipper.shipper 5d1b7b2a-c163-4c48-bb01-9ee00584d761
```
*/
package tencentcloud

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

func resourceTencentCloudClsCosShipper() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClsCosShipperCreate,
		Read:   resourceTencentCloudClsCosShipperRead,
		Delete: resourceTencentCloudClsCosShipperDelete,
		Update: resourceTencentCloudClsCosShipperUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"topic_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the log topic to which the shipping rule to be created belongs.",
			},
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Destination bucket in the shipping rule to be created.",
			},
			"prefix": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Prefix of the shipping directory in the shipping rule to be created.",
			},
			"shipper_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Shipping rule name.",
			},
			"interval": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Shipping time interval in seconds. Default value: 300. Value range: 300~900.",
			},
			"max_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "Maximum size of a file to be shipped, in MB. Default value: 256. Value range: 100~256.",
			},
			"filter_rules": {
				Type:     schema.TypeList,
				Optional: true,
				Description: "Filter rules for shipped logs. Only logs matching the rules can be shipped. All rules are in the AND relationship, and up to five rules can be added. " +
					"If the array is empty, no filtering will be performed, and all logs will be shipped.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter rule key.",
						},
						"regex": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter rule.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Filter rule value.",
						},
					},
				},
			},
			"partition": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Partition rule of shipped log, which can be represented in strftime time format.",
			},
			"compress": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Compression format. Valid values: gzip, lzop, none (no compression).",
						},
					},
				},
				Description: "Compression configuration of shipped log.",
			},
			"content": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Format configuration of shipped log content.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"format": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Content format. Valid values: json, csv.",
						},
						"csv": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"print_key": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether to print key on the first row of the CSV file.",
									},
									"keys": {
										Type:        schema.TypeSet,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Names of keys.Note: this field may return null, indicating that no valid values can be obtained.",
									},
									"delimiter": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Field delimiter.",
									},
									"escape_char": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Field delimiter.",
									},
									"non_existing_field": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Content used to populate non-existing fields.",
									},
								},
							},
							Description: "CSV format content description.Note: this field may return null, indicating that no valid values can be obtained.",
						},
						"json": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enable_tag": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Enablement flag.",
									},
									"meta_fields": {
										Type:        schema.TypeSet,
										Required:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "Metadata information list\nNote: this field may return null, indicating that no valid values can be obtained..",
									},
								},
							},
							Description: "JSON format content description.Note: this field may return null, indicating that no valid values can be obtained.",
						},
					},
				},
			},
		},
	}
}

func resourceTencentCloudClsCosShipperCreate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_cos_shipper.create")()

	logId := getLogId(contextNil)

	var (
		request  = cls.NewCreateShipperRequest()
		response *cls.CreateShipperResponse
	)

	if v, ok := d.GetOk("topic_id"); ok {
		request.TopicId = helper.String(v.(string))
	}

	if v, ok := d.GetOk("bucket"); ok {
		request.Bucket = helper.String(v.(string))
	}

	if v, ok := d.GetOk("prefix"); ok {
		request.Prefix = helper.String(v.(string))
	}

	if v, ok := d.GetOk("shipper_name"); ok {
		request.ShipperName = helper.String(v.(string))
	}

	if v, ok := d.GetOk("interval"); ok {
		request.Interval = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("max_size"); ok {
		request.MaxSize = helper.IntUint64(v.(int))
	}

	if v, ok := d.GetOk("filter_rules"); ok {
		filterRules := make([]*cls.FilterRuleInfo, 0, 10)
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			filterRule := cls.FilterRuleInfo{}
			if v, ok := dMap["key"]; ok {
				filterRule.Key = helper.String(v.(string))
			}
			if v, ok := dMap["regex"]; ok {
				filterRule.Regex = helper.String(v.(string))
			}
			if v, ok := dMap["value"]; ok {
				filterRule.Value = helper.String(v.(string))
			}
			filterRules = append(filterRules, &filterRule)
		}
		request.FilterRules = filterRules
	}

	if v, ok := d.GetOk("partition"); ok {
		request.Partition = helper.String(v.(string))
	}

	if v, ok := d.GetOk("compress"); ok {
		compresses := make([]*cls.CompressInfo, 0, 10)
		if len(v.([]interface{})) != 1 {
			return fmt.Errorf("need only one compress.")
		}
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			compress := cls.CompressInfo{}
			if v, ok := dMap["format"]; ok {
				compress.Format = helper.String(v.(string))
			}
			compresses = append(compresses, &compress)
		}
		request.Compress = compresses[0]
	}

	if v, ok := d.GetOk("content"); ok {
		contents := make([]*cls.ContentInfo, 0, 10)
		if len(v.([]interface{})) != 1 {
			return fmt.Errorf("need only one content.")
		}
		for _, item := range v.([]interface{}) {
			dMap := item.(map[string]interface{})
			content := cls.ContentInfo{}
			if v, ok := dMap["format"]; ok {
				content.Format = helper.String(v.(string))
			}
			if v, ok := dMap["csv"]; ok {
				if len(v.([]interface{})) == 1 {
					csv := v.([]interface{})[0].(map[string]interface{})
					csvInfo := cls.CsvInfo{}
					csvInfo.PrintKey = helper.Bool(csv["print_key"].(bool))
					keys := csv["keys"].(*schema.Set).List()
					for _, key := range keys {
						csvInfo.Keys = append(csvInfo.Keys, helper.String(key.(string)))
					}
					csvInfo.Delimiter = helper.String(csv["delimiter"].(string))
					csvInfo.EscapeChar = helper.String(csv["escape_char"].(string))
					csvInfo.NonExistingField = helper.String(csv["non_existing_field"].(string))
					content.Csv = &csvInfo
				}
			}
			if v, ok := dMap["json"]; ok {
				if len(v.([]interface{})) == 1 {

					json := v.([]interface{})[0].(map[string]interface{})
					jsonInfo := cls.JsonInfo{}
					jsonInfo.EnableTag = helper.Bool(json["enable_tag"].(bool))
					metaFields := json["meta_fields"].(*schema.Set).List()
					for _, metaField := range metaFields {
						jsonInfo.MetaFields = append(jsonInfo.MetaFields, helper.String(metaField.(string)))
					}
					content.Json = &jsonInfo
				}
			}
			contents = append(contents, &content)
		}
		request.Content = contents[0]
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().CreateShipper(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		response = result
		return nil
	})

	if err != nil {
		log.Printf("[CRITAL]%s create cls cos shipper failed, reason:%+v", logId, err)
		return err
	}

	id := *response.Response.ShipperId
	d.SetId(id)
	return resourceTencentCloudClsCosShipperRead(d, meta)
}

func resourceTencentCloudClsCosShipperRead(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_cos_shipper.read")()
	defer inconsistentCheck(d, meta)()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}

	id := d.Id()

	shipper, err := service.DescribeClsCosShipperById(ctx, id)

	if err != nil {
		return err
	}

	if shipper == nil {
		d.SetId("")
		return fmt.Errorf("resource `Shipper` %s does not exist", id)
	}

	_ = d.Set("topic_id", shipper.TopicId)
	_ = d.Set("bucket", shipper.Bucket)
	_ = d.Set("prefix", shipper.Prefix)
	_ = d.Set("shipper_name", shipper.ShipperName)
	if shipper.Interval != nil {
		_ = d.Set("interval", shipper.Interval)
	}
	if shipper.MaxSize != nil {
		_ = d.Set("max_size", shipper.MaxSize)
	}

	if shipper.FilterRules != nil {
		filterRules := make([]interface{}, 0, 100)
		for _, v := range shipper.FilterRules {
			filterRule := map[string]interface{}{
				"key":   v.Key,
				"regex": v.Regex,
				"value": v.Value,
			}
			filterRules = append(filterRules, filterRule)
		}
		_ = d.Set("filter_rules", filterRules)
	}

	if shipper.Partition != nil {
		_ = d.Set("partition", shipper.Partition)
	}

	if shipper.Compress != nil {
		compress := map[string]interface{}{
			"format": shipper.Compress.Format,
		}
		_ = d.Set("compress", []interface{}{compress})
	}

	if shipper.Content != nil {
		content := map[string]interface{}{
			"format": shipper.Content.Format,
		}
		if shipper.Content.Csv != nil {
			csv := map[string]interface{}{
				"print_key":          shipper.Content.Csv.PrintKey,
				"keys":               shipper.Content.Csv.Keys,
				"delimiter":          shipper.Content.Csv.Delimiter,
				"escape_char":        shipper.Content.Csv.EscapeChar,
				"non_existing_field": shipper.Content.Csv.NonExistingField,
			}
			content["csv"] = []interface{}{csv}
		}
		if shipper.Content.Json != nil {
			json := map[string]interface{}{
				"enable_tag":  shipper.Content.Json.EnableTag,
				"meta_fields": shipper.Content.Json.MetaFields,
			}
			content["json"] = []interface{}{json}
		}
		_ = d.Set("content", []interface{}{content})
	}
	return nil
}

func resourceTencentCloudClsCosShipperUpdate(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_cos_shipper.update")()
	logId := getLogId(contextNil)
	request := cls.NewModifyShipperRequest()

	request.ShipperId = helper.String(d.Id())

	if d.HasChange("bucket") {
		if v, ok := d.GetOk("bucket"); ok {
			request.Bucket = helper.String(v.(string))
		}
	}

	if d.HasChange("prefix") {
		if v, ok := d.GetOk("prefix"); ok {
			request.Prefix = helper.String(v.(string))
		}
	}

	if d.HasChange("shipper_name") {
		if v, ok := d.GetOk("shipper_name"); ok {
			request.ShipperName = helper.String(v.(string))
		}
	}

	if d.HasChange("interval") {
		if v, ok := d.GetOk("interval"); ok {
			request.Interval = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("max_size") {
		if v, ok := d.GetOk("max_size"); ok {
			request.MaxSize = helper.IntUint64(v.(int))
		}
	}

	if d.HasChange("filter_rules") {
		if v, ok := d.GetOk("filter_rules"); ok {
			filterRules := make([]*cls.FilterRuleInfo, 0, 10)
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				filterRule := cls.FilterRuleInfo{}
				if v, ok := dMap["key"]; ok {
					filterRule.Key = helper.String(v.(string))
				}
				if v, ok := dMap["regex"]; ok {
					filterRule.Regex = helper.String(v.(string))
				}
				if v, ok := dMap["value"]; ok {
					filterRule.Value = helper.String(v.(string))
				}
				filterRules = append(filterRules, &filterRule)
			}
			request.FilterRules = filterRules
		}
	}

	if d.HasChange("partition") {
		if v, ok := d.GetOk("partition"); ok {
			request.Partition = helper.String(v.(string))
		}
	}

	if d.HasChange("compress") {
		if v, ok := d.GetOk("compress"); ok {
			compresses := make([]*cls.CompressInfo, 0, 10)
			if len(v.([]interface{})) != 1 {
				return fmt.Errorf("need only one compress.")
			}
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				compress := cls.CompressInfo{}
				if v, ok := dMap["format"]; ok {
					compress.Format = helper.String(v.(string))
				}
				compresses = append(compresses, &compress)
			}
			request.Compress = compresses[0]
		}
	}

	if d.HasChange("content") {
		if v, ok := d.GetOk("content"); ok {
			contents := make([]*cls.ContentInfo, 0, 10)
			if len(v.([]interface{})) != 1 {
				return fmt.Errorf("need only one content.")
			}
			for _, item := range v.([]interface{}) {
				dMap := item.(map[string]interface{})
				content := cls.ContentInfo{}
				if v, ok := dMap["format"]; ok {
					content.Format = helper.String(v.(string))
				}
				if v, ok := dMap["csv"]; ok {
					if len(v.([]interface{})) == 1 {
						csv := v.([]interface{})[0].(map[string]interface{})
						csvInfo := cls.CsvInfo{}
						csvInfo.PrintKey = helper.Bool(csv["print_key"].(bool))
						keys := csv["keys"].(*schema.Set).List()
						for _, key := range keys {
							csvInfo.Keys = append(csvInfo.Keys, helper.String(key.(string)))
						}
						csvInfo.Delimiter = helper.String(csv["delimiter"].(string))
						csvInfo.EscapeChar = helper.String(csv["escape_char"].(string))
						csvInfo.NonExistingField = helper.String(csv["non_existing_field"].(string))
						content.Csv = &csvInfo
					}
				}
				if v, ok := dMap["json"]; ok {
					if len(v.([]interface{})) == 1 {

						json := v.([]interface{})[0].(map[string]interface{})
						jsonInfo := cls.JsonInfo{}
						jsonInfo.EnableTag = helper.Bool(json["enable_tag"].(bool))
						metaFields := json["meta_fields"].(*schema.Set).List()
						for _, metaField := range metaFields {
							jsonInfo.MetaFields = append(jsonInfo.MetaFields, helper.String(metaField.(string)))
						}
						content.Json = &jsonInfo
					}
				}
				contents = append(contents, &content)
			}
			request.Content = contents[0]
		}
	}

	err := resource.Retry(writeRetryTimeout, func() *resource.RetryError {
		result, e := meta.(*TencentCloudClient).apiV3Conn.UseClsClient().ModifyShipper(request)
		if e != nil {
			return retryError(e)
		} else {
			log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
				logId, request.GetAction(), request.ToJsonString(), result.ToJsonString())
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceTencentCloudClsCosShipperRead(d, meta)
}

func resourceTencentCloudClsCosShipperDelete(d *schema.ResourceData, meta interface{}) error {
	defer logElapsed("resource.tencentcloud_cls_cos_shipper.delete")()

	logId := getLogId(contextNil)
	ctx := context.WithValue(context.TODO(), logIdKey, logId)
	service := ClsService{client: meta.(*TencentCloudClient).apiV3Conn}
	id := d.Id()

	if err := service.DeleteClsCosShipper(ctx, id); err != nil {
		return err
	}

	return nil
}
