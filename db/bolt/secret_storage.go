package bolt

import (
	"fmt"

	"github.com/semaphoreui/semaphore/db"
)

func (d *BoltDb) GetSecretStorages(projectID int) ([]db.SecretStorage, error) {
	return []db.SecretStorage{}, nil
}

func (d *BoltDb) CreateSecretStorage(storage db.SecretStorage) (db.SecretStorage, error) {
	return db.SecretStorage{}, fmt.Errorf("secret storage is not supported in BoltDB, please use SQL database")
}

func (d *BoltDb) GetSecretStorage(projectID int, storageID int) (db.SecretStorage, error) {
	return db.SecretStorage{}, fmt.Errorf("secret storage is not supported in BoltDB, please use SQL database")
}

func (d *BoltDb) DeleteSecretStorage(projectID int, storageID int) error {
	return fmt.Errorf("secret storage is not supported in BoltDB, please use SQL database")
}

func (d *BoltDb) UpdateSecretStorage(storage db.SecretStorage) error {
	return fmt.Errorf("secret storage is not supported in BoltDB, please use SQL database")
}

func (d *BoltDb) GetSecretStorageRefs(projectID int, storageID int) (db.ObjectReferrers, error) {
	return d.getObjectRefs(projectID, db.SecretStorageProps, storageID)
}
