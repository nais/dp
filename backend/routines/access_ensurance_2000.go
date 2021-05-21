package routines

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/nais/dp/backend/api"
	"github.com/nais/dp/backend/config"
	"github.com/nais/dp/backend/iam"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/iterator"
)

const AccessEnsurance2000 = "AccessEnsurance2000"

func EnsureAccess(ctx context.Context, cfg config.Config, client *firestore.Client, updateFrequency time.Duration) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Debugf("Checking access...")
			if err := ensureAccesses(ctx, cfg, client); err != nil {
				log.Errorf("Checking access: %v", err)
			}
			ticker.Reset(updateFrequency)
		case <-ctx.Done():
			return
		}
	}
}

func ensureAccesses(ctx context.Context, cfg config.Config, client *firestore.Client) error {
	dpc := client.Collection(cfg.Firestore.DataproductsCollection)

	iter := dpc.Documents(ctx)
	defer iter.Stop()
	for {
		document, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Iterating documents: %v", err)
		}

		if err := checkAccess(ctx, cfg, client, document); err != nil {
			return err
		}
	}
	return nil
}

func checkAccess(ctx context.Context, cfg config.Config, client *firestore.Client, snapshot *firestore.DocumentSnapshot) error {
	dataproduct, err := api.DocumentToProductResponse(snapshot)
	if err != nil {
		return err
	}
	if len(dataproduct.DataProduct.Datastore) == 0 {
		// we have no access to check here
		return nil
	}
	datastore := dataproduct.DataProduct.Datastore[0]
	toDelete := make([]string, 0)

	for subject, expiry := range dataproduct.DataProduct.Access {
		if expiry.Before(time.Now()) {
			log.Infof("Access expired, removing %v from %v", subject, datastore["type"])
			if err := iam.RemoveDatastoreAccess(ctx, datastore, subject); err != nil {
				return err
			}

			update := api.AccessUpdate{}
			update.Delete(AccessEnsurance2000, dataproduct.ID, subject)
			api.UpdateHistory(ctx, client, cfg.Firestore.AccessUpdatesCollection, update)
			toDelete = append(toDelete, subject)
		} else {
			access, err := iam.CheckDatastoreAccess(ctx, datastore, subject)
			if err != nil {
				return err
			}
			if !access {
				log.Infof("Access state out of sync with Google %v, giving access to %v", datastore["type"], subject)
				accessMap := map[string]time.Time{subject: expiry}
				if err := iam.UpdateDatastoreAccess(ctx, datastore, accessMap); err != nil {
					return err
				}

				update := api.AccessUpdate{}
				update.Grant(AccessEnsurance2000, dataproduct.ID, subject, expiry)
				api.UpdateHistory(ctx, client, cfg.Firestore.AccessUpdatesCollection, update)
			}
		}
	}

	if len(toDelete) > 0 {
		for _, subject := range toDelete {
			delete(dataproduct.DataProduct.Access, subject)
		}
		snapshot.Ref.Update(ctx, []firestore.Update{{
			Path:  "access",
			Value: dataproduct.DataProduct.Access,
		}})
	}

	update := api.AccessUpdate{}
	update.Verify(AccessEnsurance2000, dataproduct.ID)
	api.UpdateHistory(ctx, client, cfg.Firestore.AccessUpdatesCollection, update)
	return nil
}
