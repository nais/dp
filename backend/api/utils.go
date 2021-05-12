package api

import (
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	log "github.com/sirupsen/logrus"
)

func (a *api) createUpdates(dp DataProduct, existingDp DataProduct) ([]firestore.Update, error) {
	var updates []firestore.Update
	newAccess := existingDp.Access

	if len(dp.Name) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "name",
			Value: dp.Name,
		})
	}
	if len(dp.Description) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "description",
			Value: dp.Description,
		})
	}
	if len(dp.Datastore) > 0 {
		if errs := ValidateDatastore(dp.Datastore[0]); errs != nil {
			return nil, errs
		}
		updates = append(updates, firestore.Update{
			Path:  "datastore",
			Value: dp.Datastore,
		})
	}
	if len(dp.Owner) > 0 {
		updates = append(updates, firestore.Update{
			Path:  "owner",
			Value: dp.Owner,
		})
	}
	if len(dp.Access) > 0 {
		for _, access := range dp.Access {
			if errs := a.validate.Struct(access); errs != nil {
				return nil, errs
			}
		}
		newAccess = append(newAccess, dp.Access...)
		updates = append(updates, firestore.Update{
			Path:  "access",
			Value: newAccess,
		})
	}

	return updates, nil
}

func ValidateDatastore(store map[string]string) error {
	datastoreType := store["type"]
	if len(datastoreType) == 0 {
		return fmt.Errorf("no type defined")
	}

	switch datastoreType {
	case BucketType:
		return hasKeys(store, "project_id", "bucket_id")
	case BigQueryType:
		return hasKeys(store, "dataset_id", "project_id", "resource_id")
	}

	return fmt.Errorf("unknown datastore type: %v", datastoreType)
}

func hasKeys(m map[string]string, keys ...string) error {
	for _, k := range keys {
		if _, found := m[k]; !found {
			return fmt.Errorf("missing key: %v", k)
		}
	}
	return nil
}

func respondf(w http.ResponseWriter, statusCode int, format string, args ...interface{}) {
	w.WriteHeader(statusCode)

	if _, wErr := w.Write([]byte(fmt.Sprintf(format, args...))); wErr != nil {
		log.Errorf("unable to write response: %v", wErr)
	}
}

func documentToProduct(d *firestore.DocumentSnapshot) (DataProductResponse, error) {
	var dpr DataProductResponse
	var dp DataProduct

	if err := d.DataTo(&dp); err != nil {
		return dpr, err
	}

	if dp.Access == nil {
		dp.Access = make([]*AccessEntry, 0)
	}
	dpr.ID = d.Ref.ID
	dpr.Updated = d.UpdateTime
	dpr.Created = d.CreateTime
	dpr.DataProduct = dp

	return dpr, nil
}
