package cvm

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"log"
	"regexp"
	"strings"

	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
)

func dataSourceTencentCloudKeyPairsReadPostHandleResponse0(ctx context.Context, req map[string]interface{}, resp *cvm.DescribeKeyPairsResponseParams) error {
	d := tccommon.ResourceDataFromContext(ctx)
	logId := tccommon.GetLogId(tccommon.ContextNil)
	var err error
	keyPairs := resp.KeyPairSet
	keyPairList := make([]map[string]interface{}, 0, len(keyPairs))
	ids := make([]string, 0, len(keyPairs))
	keyName := ctx.Value("key_name").(string)
	namePattern, _ := regexp.Compile(keyName)
	for _, keyPair := range keyPairs {
		if match := namePattern.MatchString(*keyPair.KeyName); !match {
			continue
		}
		mapping := map[string]interface{}{
			"key_id":      keyPair.KeyId,
			"key_name":    keyPair.KeyName,
			"project_id":  keyPair.ProjectId,
			"create_time": keyPair.CreatedTime,
		}
		if keyPair.PublicKey != nil {
			publicKey := *keyPair.PublicKey
			split := strings.Split(publicKey, " ")
			publicKey = strings.Join(split[0:len(split)-1], " ")
			mapping["public_key"] = publicKey
		}
		keyPairList = append(keyPairList, mapping)
		ids = append(ids, *keyPair.KeyId)
	}

	d.SetId(helper.DataResourceIdsHash(ids))
	err = d.Set("key_pair_list", keyPairList)
	if err != nil {
		log.Printf("[CRITAL]%s provider set key pair list fail, reason:%s\n ", logId, err.Error())
		return err
	}

	context.WithValue(ctx, "keyPairList", keyPairList)
	return nil
}

func dataSourceTencentCloudKeyPairsReadOutputContent(ctx context.Context) interface{} {
	eipList := ctx.Value("keyPairList").(*schema.Set).List()
	return eipList
}

func dataSourceTencentCloudKeyPairsReadPreRequest0(ctx context.Context, req *cvm.DescribeKeyPairsRequest) error {
	d := tccommon.ResourceDataFromContext(ctx)
	keyId := d.Get("key_id").(string)
	keyName := d.Get("key_name").(string)
	name := keyName
	if keyName != "" {
		if name[0] == '^' {
			name = name[1:]
		}
		length := len(name)
		if length > 0 && name[length-1] == '$' {
			name = name[:length-1]
		}

		pattern := `^[a-zA-Z0-9_]+$`
		if match, _ := regexp.MatchString(pattern, name); !match {
			return fmt.Errorf("key_name only support letters, numbers, and _ : %s", keyName)
		}
	}

	var projectId *int
	if v, ok := d.GetOkExists("project_id"); ok {
		vv := v.(int)
		projectId = &vv
	}

	if keyId != "" {
		req.KeyIds = []*string{&keyId}
	}
	req.Filters = make([]*cvm.Filter, 0)
	if name != "" {
		filter := &cvm.Filter{
			Name:   helper.String("key-name"),
			Values: []*string{&name},
		}
		req.Filters = append(req.Filters, filter)
	}
	if projectId != nil {
		filter := &cvm.Filter{
			Name:   helper.String("project-id"),
			Values: []*string{helper.String(fmt.Sprintf("%d", *projectId))},
		}
		req.Filters = append(req.Filters, filter)
	}

	context.WithValue(ctx, "keyName", keyName)
	return nil
}
