/*
Provide a resource to create a CLB instance.

Example Usage

```hcl
resource "tencentcloud_clb_rewrite" "rewrite" {
  clb_id                = "lb-p7olt9e5"
  source_listener_id    = "lbl-jc1dx6ju#lb-p7olt9e5"
  target_listener_id    = "lbl-asj1hzuo#lb-p7olt9e5"
  rewrite_source_loc_id = "loc-ft8fmngv#lbl-jc1dx6ju#lb-p7olt9e5"
  rewrite_target_loc_id = "loc-4xxr2cy7#lbl-asj1hzuo#lb-p7olt9e5"
}
```

Import

CLB instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_rewrite.rewrite loc-ft8fmngv#loc-4xxr2cy7#lbl-jc1dx6ju#lbl-asj1hzuo#lb-p7olt9e5
```
*/
package tencentcloud

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	clb "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/clb/v20180317"
)

func resourceTencentCloudClbRewrite() *schema.Resource {
	return &schema.Resource{
		Create: resourceTencentCloudClbRewriteCreate,
		Read:   resourceTencentCloudClbRewriteRead,
		Delete: resourceTencentCloudClbRewriteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"clb_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Id of CLB instance.",
			},
			"source_listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Id of source listener. ",
			},
			"target_listener_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Id of source listener. ",
			},
			"rewrite_source_loc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Id of rule id of source listener. ",
			},
			"rewrite_target_loc_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "Id of rule id of target listener. ",
			},
		},
	}
}

func resourceTencentCloudClbRewriteCreate(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_clb_rewrite.create")()

	//暂时不支持auto
	clbId := d.Get("clb_id").(string)
	sourceListenerId := strings.Split(d.Get("source_listener_id").(string), "#")[0]
	targertListenerId := strings.Split(d.Get("target_listener_id").(string), "#")[0]
	sourceLocId := strings.Split(d.Get("rewrite_source_loc_id").(string), "#")[0]
	targetLocId := strings.Split(d.Get("rewrite_target_loc_id").(string), "#")[0]

	request := clb.NewManualRewriteRequest()

	request.LoadBalancerId = stringToPointer(clbId)
	request.SourceListenerId = stringToPointer(sourceListenerId)
	request.TargetListenerId = stringToPointer(targertListenerId)

	var rewriteInfo clb.RewriteLocationMap
	rewriteInfo.SourceLocationId = stringToPointer(sourceLocId)
	rewriteInfo.TargetLocationId = stringToPointer(targetLocId)
	request.RewriteInfos = []*clb.RewriteLocationMap{&rewriteInfo}
	response, err := meta.(*TencentCloudClient).apiV3Conn.UseClbClient().ManualRewrite(request)
	if err != nil {
		log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n",
			logId, request.GetAction(), request.ToJsonString(), err.Error())
		return err
	} else {
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n",
			logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())
		requestId := *response.Response.RequestId
		retryErr := retrySet(requestId, meta.(*TencentCloudClient).apiV3Conn.UseClbClient())
		if retryErr != nil {
			return retryErr
		}
	}
	d.SetId(sourceLocId + "#" + targetLocId + "#" + sourceListenerId + "#" + targertListenerId + "#" + clbId)
	return resourceTencentCloudClbRewriteRead(d, meta)
}

func resourceTencentCloudClbRewriteRead(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_clb_rewrite.read")()
	ctx := context.WithValue(context.TODO(), "logId", logId)

	rewriteId := d.Id()

	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	instance, err := clbService.DescribeRewriteInfoById(ctx, rewriteId)
	if err != nil {
		return err
	}
	//遍历得到符合target的重定向

	d.Set("clb_id", (*instance)["clb_id"])
	d.Set("source_listener_id", (*instance)["source_listener_id"]+"#"+(*instance)["clb_id"])
	d.Set("target_listener_id", (*instance)["target_listener_id"]+"#"+(*instance)["clb_id"])
	d.Set("rewrite_source_loc_id", (*instance)["rewrite_source_loc_id"]+"#"+(*instance)["source_listener_id"]+"#"+(*instance)["clb_id"])
	d.Set("rewrite_target_loc_id", (*instance)["rewrite_target_loc_id"]+"#"+(*instance)["target_listener_id"]+"#"+(*instance)["clb_id"])

	return nil
}

func resourceTencentCloudClbRewriteDelete(d *schema.ResourceData, meta interface{}) error {
	logId := GetLogId(nil)
	defer LogElapsed(logId + "resource.tencentcloud_clb_rewrite.delete")()
	ctx := context.WithValue(context.TODO(), "logId", logId)

	clbId := d.Id()
	clbService := ClbService{
		client: meta.(*TencentCloudClient).apiV3Conn,
	}
	err := clbService.DeleteRewriteInfoById(ctx, clbId)
	if err != nil {
		log.Printf("[CRITAL]%s reason[%s]\n", logId, err.Error())
		return err
	}
	return nil
}
