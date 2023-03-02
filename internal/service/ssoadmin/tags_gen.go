// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package ssoadmin

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssoadmin"
	"github.com/aws/aws-sdk-go/service/ssoadmin/ssoadminiface"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
)

// ListTags lists ssoadmin service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func ListTags(ctx context.Context, conn ssoadminiface.SSOAdminAPI, identifier, resourceType string) (tftags.KeyValueTags, error) {
	input := &ssoadmin.ListTagsForResourceInput{
		ResourceArn: aws.String(identifier),
		InstanceArn: aws.String(resourceType),
	}

	output, err := conn.ListTagsForResourceWithContext(ctx, input)

	if err != nil {
		return tftags.New(ctx, nil), err
	}

	return KeyValueTags(ctx, output.Tags), nil
}

// []*SERVICE.Tag handling

// Tags returns ssoadmin service tags.
func Tags(tags tftags.KeyValueTags) []*ssoadmin.Tag {
	result := make([]*ssoadmin.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := &ssoadmin.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return result
}

// KeyValueTags creates tftags.KeyValueTags from ssoadmin service tags.
func KeyValueTags(ctx context.Context, tags []*ssoadmin.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.StringValue(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// UpdateTags updates ssoadmin service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.

func UpdateTags(ctx context.Context, conn ssoadminiface.SSOAdminAPI, identifier, resourceType string, oldTagsMap, newTagsMap any) error {
	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		input := &ssoadmin.UntagResourceInput{
			ResourceArn: aws.String(identifier),
			InstanceArn: aws.String(resourceType),
			TagKeys:     aws.StringSlice(removedTags.IgnoreAWS().Keys()),
		}

		_, err := conn.UntagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		input := &ssoadmin.TagResourceInput{
			ResourceArn: aws.String(identifier),
			InstanceArn: aws.String(resourceType),
			Tags:        Tags(updatedTags.IgnoreAWS()),
		}

		_, err := conn.TagResourceWithContext(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, resourceType string, oldTags, newTags any) error {
	return UpdateTags(ctx, meta.(*conns.AWSClient).SSOAdminConn(), identifier, resourceType, oldTags, newTags)
}
