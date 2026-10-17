package cos

import (
	"reflect"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCosBucketLifecycleRuleSchema(t *testing.T) {
	rules := ResourceTencentCloudCosBucket().Schema["lifecycle_rules"].Elem.(*schema.Resource).Schema

	if got := rules["status"].Default; got != s3.ExpirationStatusEnabled {
		t.Fatalf("status default = %v, want %q", got, s3.ExpirationStatusEnabled)
	}
	if _, errors := rules["status"].ValidateFunc("invalid", "status"); len(errors) == 0 {
		t.Fatal("status validation accepted an invalid value")
	}
	if got := rules["filter_tags"].Type; got != schema.TypeMap {
		t.Fatalf("filter_tags type = %v, want TypeMap", got)
	}
}

func TestCosBucketLifecycleDataSourceRuleSchema(t *testing.T) {
	bucketList := DataSourceTencentCloudCosBuckets().Schema["bucket_list"].Elem.(*schema.Resource).Schema
	rules := bucketList["lifecycle_rules"].Elem.(*schema.Resource).Schema

	if !rules["status"].Computed || rules["status"].Type != schema.TypeString {
		t.Fatalf("status schema = %#v, want computed string", rules["status"])
	}
	if !rules["filter_tags"].Computed || rules["filter_tags"].Type != schema.TypeMap {
		t.Fatalf("filter_tags schema = %#v, want computed map", rules["filter_tags"])
	}
}

func TestValidateCosBucketLifecycleFilterTags(t *testing.T) {
	tenTags := map[string]interface{}{
		"tag0": "0", "tag1": "1", "tag2": "2", "tag3": "3", "tag4": "4",
		"tag5": "5", "tag6": "6", "tag7": "7", "tag8": "8", "tag9": "9",
	}
	elevenTags := make(map[string]interface{}, len(tenTags)+1)
	for key, value := range tenTags {
		elevenTags[key] = value
	}
	elevenTags["tag10"] = "10"

	tests := map[string]struct {
		tags    map[string]interface{}
		wantErr bool
	}{
		"ten tags":         {tags: tenTags},
		"eleven tags":      {tags: elevenTags, wantErr: true},
		"128 byte key":     {tags: map[string]interface{}{strings.Repeat("k", 128): "value"}},
		"129 byte key":     {tags: map[string]interface{}{strings.Repeat("k", 129): "value"}, wantErr: true},
		"256 byte value":   {tags: map[string]interface{}{"key": strings.Repeat("v", 256)}},
		"257 byte value":   {tags: map[string]interface{}{"key": strings.Repeat("v", 257)}, wantErr: true},
		"empty key":        {tags: map[string]interface{}{"": "value"}, wantErr: true},
		"non-string value": {tags: map[string]interface{}{"key": 1}, wantErr: true},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			_, errors := validateCosBucketLifecycleFilterTags(test.tags, "filter_tags")
			if gotErr := len(errors) > 0; gotErr != test.wantErr {
				t.Fatalf("validateCosBucketLifecycleFilterTags() errors = %v, wantErr %t", errors, test.wantErr)
			}
		})
	}
}

func TestValidateCosBucketLifecycleRules(t *testing.T) {
	abortUploads := schema.NewSet(schema.HashString, []interface{}{"configured"})
	tests := map[string]struct {
		rules   []interface{}
		wantErr bool
	}{
		"tags only": {
			rules: []interface{}{map[string]interface{}{"filter_tags": map[string]interface{}{"env": "prod"}}},
		},
		"abort only": {
			rules: []interface{}{map[string]interface{}{"abort_incomplete_multipart_upload": abortUploads}},
		},
		"tags and abort": {
			rules: []interface{}{map[string]interface{}{
				"filter_tags":                       map[string]interface{}{"env": "prod"},
				"abort_incomplete_multipart_upload": abortUploads,
			}},
			wantErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if gotErr := validateCosBucketLifecycleRules(test.rules) != nil; gotErr != test.wantErr {
				t.Fatalf("validateCosBucketLifecycleRules() error = %t, wantErr %t", gotErr, test.wantErr)
			}
		})
	}
}

func TestExpandCosBucketLifecycleRuleFilter(t *testing.T) {
	tests := map[string]struct {
		prefix string
		tags   map[string]interface{}
		want   *s3.LifecycleRuleFilter
	}{
		"all objects": {
			want: &s3.LifecycleRuleFilter{Prefix: aws.String("")},
		},
		"prefix only": {
			prefix: "logs/",
			want:   &s3.LifecycleRuleFilter{Prefix: aws.String("logs/")},
		},
		"single tag only": {
			tags: map[string]interface{}{"env": "prod"},
			want: &s3.LifecycleRuleFilter{Tag: &s3.Tag{Key: aws.String("env"), Value: aws.String("prod")}},
		},
		"prefix and tag": {
			prefix: "logs/",
			tags:   map[string]interface{}{"env": "prod"},
			want: &s3.LifecycleRuleFilter{And: &s3.LifecycleRuleAndOperator{
				Prefix: aws.String("logs/"),
				Tags:   []*s3.Tag{{Key: aws.String("env"), Value: aws.String("prod")}},
			}},
		},
		"multiple tags are stable": {
			tags: map[string]interface{}{"team": "db", "env": "prod"},
			want: &s3.LifecycleRuleFilter{And: &s3.LifecycleRuleAndOperator{Tags: []*s3.Tag{
				{Key: aws.String("env"), Value: aws.String("prod")},
				{Key: aws.String("team"), Value: aws.String("db")},
			}}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if got := expandCosBucketLifecycleRuleFilter(test.prefix, test.tags); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("expandCosBucketLifecycleRuleFilter() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestFlattenCosBucketLifecycleRuleFilter(t *testing.T) {
	tests := map[string]struct {
		filter     *s3.LifecycleRuleFilter
		wantPrefix string
		wantTags   map[string]interface{}
	}{
		"nil": {
			wantTags: map[string]interface{}{},
		},
		"prefix": {
			filter:     &s3.LifecycleRuleFilter{Prefix: aws.String("logs/")},
			wantPrefix: "logs/",
			wantTags:   map[string]interface{}{},
		},
		"tag": {
			filter:   &s3.LifecycleRuleFilter{Tag: &s3.Tag{Key: aws.String("env"), Value: aws.String("prod")}},
			wantTags: map[string]interface{}{"env": "prod"},
		},
		"and": {
			filter: &s3.LifecycleRuleFilter{And: &s3.LifecycleRuleAndOperator{
				Prefix: aws.String("logs/"),
				Tags: []*s3.Tag{
					{Key: aws.String("env"), Value: aws.String("prod")},
					{Key: aws.String("team"), Value: aws.String("db")},
				},
			}},
			wantPrefix: "logs/",
			wantTags:   map[string]interface{}{"env": "prod", "team": "db"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			prefix, tags := flattenCosBucketLifecycleRuleFilter(test.filter)
			if prefix != test.wantPrefix || !reflect.DeepEqual(tags, test.wantTags) {
				t.Fatalf("flattenCosBucketLifecycleRuleFilter() = (%q, %#v), want (%q, %#v)", prefix, tags, test.wantPrefix, test.wantTags)
			}
		})
	}
}
