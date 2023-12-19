package cbs

import (
	"context"
	"log"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/pkg/errors"
	tag "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tag/v20180813"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/ratelimit"
)

//internal version: replace tagFmt begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
//internal version: replace tagFmt end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

type TagService struct {
	client *connectivity.TencentCloudClient
}

func (me *TagService) ModifyTags(ctx context.Context, resourceName string, replaceTags map[string]string, deleteKeys []string) error {
	request := tag.NewModifyResourceTagsRequest()
	request.Resource = &resourceName
	if len(replaceTags) > 0 {
		request.ReplaceTags = make([]*tag.Tag, 0, len(replaceTags))
		for k, v := range replaceTags {
			key := k
			value := v
			replaceTag := &tag.Tag{
				TagKey:   &key,
				TagValue: &value,
			}
			request.ReplaceTags = append(request.ReplaceTags, replaceTag)
		}
	}
	if len(deleteKeys) > 0 {
		request.DeleteTags = make([]*tag.TagKeyObject, 0, len(deleteKeys))
		for _, v := range deleteKeys {
			key := v
			deleteKey := &tag.TagKeyObject{
				TagKey: &key,
			}
			request.DeleteTags = append(request.DeleteTags, deleteKey)
		}
	}

	return resource.Retry(tccommon.WriteRetryTimeout, func() *resource.RetryError {
		ratelimit.Check(request.GetAction())

		if _, err := me.client.UseTagClient().ModifyResourceTags(request); err != nil {
			return tccommon.RetryError(errors.WithStack(err))
		}

		return nil
	})
}

func (me *TagService) DescribeResourceTags(ctx context.Context, serviceType, resourceType, region, resourceId string) (tags map[string]string, err error) {
	request := tag.NewDescribeResourceTagsByResourceIdsRequest()
	request.ServiceType = &serviceType
	request.ResourcePrefix = &resourceType
	request.ResourceRegion = &region
	request.ResourceIds = []*string{&resourceId}
	request.Limit = helper.IntUint64(DESCRIBE_TAGS_LIMIT)

	var offset uint64
	request.Offset = &offset

	// for run loop at least once
	count := DESCRIBE_TAGS_LIMIT
	for count == DESCRIBE_TAGS_LIMIT {
		if err := resource.Retry(tccommon.ReadRetryTimeout, func() *resource.RetryError {
			ratelimit.Check(request.GetAction())

			response, err := me.client.UseTagClient().DescribeResourceTagsByResourceIds(request)
			if err != nil {
				count = 0

				return tccommon.RetryError(errors.WithStack(err))
			}

			allTags := response.Response.Tags
			count = len(allTags)

			for _, t := range allTags {
				if *t.ResourceId != resourceId {
					continue
				}
				if tags == nil {
					tags = make(map[string]string)
				}

				tags[*t.TagKey] = *t.TagValue
			}

			return nil
		}); err != nil {
			return nil, err
		}

		offset += uint64(count)
	}

	return
}

//internal version: replace waitTag begin, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.
//internal version: replace waitTag end, please do not modify this annotation and refrain from inserting any code between the beginning and end lines of the annotation.

func diffTags(oldTags, newTags map[string]interface{}) (replaceTags map[string]string, deleteTags []string) {
	replaceTags = make(map[string]string)
	deleteTags = make([]string, 0)
	for k, v := range newTags {
		_, ok := oldTags[k]
		if !ok || oldTags[k].(string) != v.(string) {
			replaceTags[k] = v.(string)
		}
	}
	for k := range oldTags {
		_, ok := newTags[k]
		if !ok {
			deleteTags = append(deleteTags, k)
		}
	}
	return
}

func (me *TagService) DescribeProjectById(ctx context.Context, projectId uint64) (project *tag.Project, disable *uint64, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tag.NewDescribeProjectsRequest()

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	// query enable project
	request.AllList = helper.Uint64(0)

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	instances := make([]*tag.Project, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTagClient().DescribeProjects(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Projects) < 1 {
			break
		}
		instances = append(instances, response.Response.Projects...)
		if len(response.Response.Projects) < int(limit) {
			break
		}
		offset += limit
	}

	for _, instance := range instances {
		if *instance.ProjectId == projectId {
			project = instance
			disable = helper.Uint64(0)
			break
		}
	}

	if project != nil {
		return
	}

	// query all project
	offset = 0
	limit = 20

	request.AllList = helper.Uint64(1)
	instances = make([]*tag.Project, 0)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTagClient().DescribeProjects(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Projects) < 1 {
			break
		}
		instances = append(instances, response.Response.Projects...)
		if len(response.Response.Projects) < int(limit) {
			break
		}
		offset += limit
	}

	for _, instance := range instances {
		if *instance.ProjectId == projectId {
			project = instance
			disable = helper.Uint64(1)
			break
		}
	}

	return
}

func (me *TagService) DisableProjectById(ctx context.Context, projectId uint64) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tag.NewUpdateProjectRequest()
	request.ProjectId = &projectId
	request.Disable = helper.Int64(1)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTagClient().UpdateProject(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}

func (me *TagService) DescribeProjects(ctx context.Context, param map[string]interface{}) (project []*tag.Project, errRet error) {
	var (
		logId   = tccommon.GetLogId(ctx)
		request = tag.NewDescribeProjectsRequest()
	)

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	for k, v := range param {
		if k == "AllList" {
			request.AllList = v.(*uint64)
		}
	}

	ratelimit.Check(request.GetAction())

	var (
		offset uint64 = 0
		limit  uint64 = 20
	)
	for {
		request.Offset = &offset
		request.Limit = &limit
		response, err := me.client.UseTagClient().DescribeProjects(request)
		if err != nil {
			errRet = err
			return
		}
		log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

		if response == nil || len(response.Response.Projects) < 1 {
			break
		}
		project = append(project, response.Response.Projects...)
		if len(response.Response.Projects) < int(limit) {
			break
		}

		offset += limit
	}

	return
}

func (me *TagService) DescribeTagResourceById(ctx context.Context, tagKey string, tagValue string) (tagRes *tag.Tag, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tag.NewGetTagsRequest()
	request.TagKeys = []*string{&tagKey}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTagClient().GetTags(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.Tags) < 1 {
		return
	}
	for _, v := range response.Response.Tags {
		if *v.TagKey == tagKey && *v.TagValue == tagValue {
			tagRes = v
		}
	}
	return
}

func (me *TagService) DeleteTagResourceById(ctx context.Context, tagKey string, tagValue string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tag.NewDeleteTagRequest()
	request.TagKey = &tagKey
	request.TagValue = &tagValue

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTagClient().DeleteTag(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
func (me *TagService) DescribeTagTagAttachmentById(ctx context.Context, tagKey string, tagValue string, resource string) (resourceTag *tag.ResourceTagMapping, errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tag.NewGetResourcesRequest()
	request.ResourceList = []*string{&resource}

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTagClient().GetResources(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	if response == nil || response.Response == nil || len(response.Response.ResourceTagMappingList) < 1 {
		return
	}
	for _, resourceTagMap := range response.Response.ResourceTagMappingList {
		if *resourceTagMap.Resource == resource {
			for _, v := range resourceTagMap.Tags {
				if *v.TagKey == tagKey && *v.TagValue == tagValue {
					resourceTag = &tag.ResourceTagMapping{
						Resource: &resource,
						Tags: []*tag.Tag{
							{TagKey: v.TagKey, TagValue: &tagValue},
						},
					}
				}
			}
		}
	}
	return
}

func (me *TagService) DeleteTagTagAttachmentById(ctx context.Context, tagKey string, resource string) (errRet error) {
	logId := tccommon.GetLogId(ctx)

	request := tag.NewDeleteResourceTagRequest()
	request.TagKey = &tagKey
	request.Resource = &resource

	defer func() {
		if errRet != nil {
			log.Printf("[CRITAL]%s api[%s] fail, request body [%s], reason[%s]\n", logId, request.GetAction(), request.ToJsonString(), errRet.Error())
		}
	}()

	ratelimit.Check(request.GetAction())

	response, err := me.client.UseTagClient().DeleteResourceTag(request)
	if err != nil {
		errRet = err
		return
	}
	log.Printf("[DEBUG]%s api[%s] success, request body [%s], response body [%s]\n", logId, request.GetAction(), request.ToJsonString(), response.ToJsonString())

	return
}
