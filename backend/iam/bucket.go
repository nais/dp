package iam

import (
	"context"
	"errors"
	"google.golang.org/genproto/googleapis/type/expr"
	"log"
	"time"

	"cloud.google.com/go/iam"
	"cloud.google.com/go/storage"
	iampb "google.golang.org/genproto/googleapis/iam/v1"
)

func UpdateBucketAccessControl(bucketName, member string) error {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := getBucketPolicy(client, bucketName); err != nil {
		log.Fatal(err)
	}

	if err := addUser(client, bucketName, member); err != nil {
		log.Fatal(err)
	}
	return nil
}

func getBucketPolicy(c *storage.Client, bucketName string) (*iam.Policy3, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	policy, err := c.Bucket(bucketName).IAM().V3().Policy(ctx)
	if err != nil {
		return nil, err
	}
	for _, binding := range policy.Bindings {
		log.Printf("%q: %q (condition: %v)", binding.Role, binding.Members, binding.Condition)
	}
	return policy, nil
}

func addUser(c *storage.Client, bucketName, member string) error {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := c.Bucket(bucketName)
	policy, err := bucket.IAM().V3().Policy(ctx)
	if err != nil {
		return err
	}

	policy.Bindings = append(policy.Bindings, &iampb.Binding{
		Role:    "roles/storage.objectViewer",
		Members: []string{member},
		Condition: &expr.Expr{
			Title:       "Expires_2022",
			Description: "Expires at noon on 2022-12-31",
			Expression:  "request.time < timestamp('2022-12-31T12:00:00Z')",
		},
	})
	if err := bucket.IAM().V3().SetPolicy(ctx, policy); err != nil {
		return err
	}
	// NOTE: It may be necessary to retry this operation if IAM policies are
	// being modified concurrently. SetPolicy will return an error if the policy
	// was modified since it was retrieved.
	return nil
}

func RemoveMemberFromBucket(bucketName, bucketMember string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	bucket := client.Bucket(bucketName)
	policy, err := bucket.IAM().V3().Policy(ctx)
	if err != nil {
		return err
	}
	for _, binding := range policy.Bindings {
		// Only remove matching role
		if binding.Role == "roles/storage.objectViewer" {
			// Filter out member.
			i := -1
			for j, member := range binding.Members {
				if member == bucketMember {
					i = j
				}
			}

			if i == -1 {
				return errors.New("no matching binding group found")
			} else {
				binding.Members = append(binding.Members[:i], binding.Members[i+1:]...)
			}
		}
	}
	if err := bucket.IAM().V3().SetPolicy(ctx, policy); err != nil {
		return err
	}
	// NOTE: It may be necessary to retry this operation if IAM policies are
	// being modified concurrently. SetPolicy will return an error if the policy
	// was modified since it was retrieved.
	return nil
}
