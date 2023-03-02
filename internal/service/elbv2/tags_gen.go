// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package elbv2

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// ListTags lists elbv2 service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(ctx context.Context, conn elbv2iface.ELBV2API, identifier string) (tftags.KeyValueTags, error) {
	input := &elbv2.DescribeTagsInput{
		ResourceArns: aws.StringSlice([]string{identifier}),
	}

	output, err := conn.DescribeTagsWithContext(ctx, input)

	if err != nil {
		return tftags.New(ctx, nil), err
	}

	return KeyValueTags(ctx, output.TagDescriptions[0].Tags), nil
}

// []*SERVICE.Tag handling

// Tags returns elbv2 service tags.
func Tags(tags tftags.KeyValueTags) []*elbv2.Tag {
	result := make([]*elbv2.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &elbv2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from elbv2 service tags.
func KeyValueTags(ctx context.Context, tags []*elbv2.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// UpdateTags updates elbv2 service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.

func UpdateTags(ctx context.Context, conn elbv2iface.ELBV2API, identifier string, oldTagsMap, newTagsMap any) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &elbv2.RemoveTagsInput{
			ResourceArns: aws.StringSlice([]string{identifier}),
			TagKeys:      aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.RemoveTagsWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &elbv2.AddTagsInput{
			ResourceArns: aws.StringSlice([]string{identifier}),
			Tags:         Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.AddTagsWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags, newTags any) error {
	return UpdateTags(ctx, meta.(*conns.AWSClient).ELBV2Conn(), identifier, oldTags, newTags)
}
