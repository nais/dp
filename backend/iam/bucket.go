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

func ChangeBucketAccessControl(bucketName, member string, start, end time.Time) error {

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := getBucketPolicy(client, bucketName); err != nil {
		log.Fatal(err)
	}

	expression := getCondition(start, end)

	if err := addUser(client, bucketName, member, expression); err != nil {
		log.Fatal(err)
	}

	return nil
}

func getCondition(start, end time.Time) string {

	startString := start.String()
	endString := end.String()
	var expression string
	if len(startString) > 0 {
		expression = "request.time > timestamp('" + startString + "')"
	}
	if len(endString) > 0 {
		if len(startString) > 0 {
			expression = expression + " && request.time < timestamp('" + endString + "')"
		} else {
			expression = "request.time < timestamp('" + endString + "')"
		}

	}
	return expression
}

func getBucketPolicy(client *storage.Client, bucketName string) (*iam.Policy3, error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	policy, err := client.Bucket(bucketName).IAM().V3().Policy(ctx)
	if err != nil {
		return nil, err
	}

	for _, binding := range policy.Bindings {
		log.Printf("%q: %q (condition: %v)", binding.Role, binding.Members, binding.Condition)
	}
	return policy, nil
}

func addUser(client *storage.Client, bucketName, member, expression string) error {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	bucket := client.Bucket(bucketName)
	policy, err := bucket.IAM().V3().Policy(ctx)
	if err != nil {
		return err
	}
	userMember := "user:" + member
	policy.Bindings = append(policy.Bindings, &iampb.Binding{
		Role:    "roles/storage.objectViewer",
		Members: []string{userMember},
		Condition: &expr.Expr{
			Title:      "Conditional access",
			Expression: expression,
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
