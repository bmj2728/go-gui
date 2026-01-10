package catdb

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"

	"github.com/bmj2728/catfetch/pkg/shared/metadata"
	"go.etcd.io/bbolt"
)

const (
	bucketCats    = "cats"
	dbKeyMetadata = "metadata"
	dbKeyData     = "data"
)

type CatDB struct {
	db *bbolt.DB
}

func OpenDB(path string) (*CatDB, error) {
	db, err := bbolt.Open(path, 0644, nil)
	if err != nil {
		return nil, err
	}
	cdb := &CatDB{db: db}
	err = cdb.InitDB()
	if err != nil {
		return nil, err
	}

	return &CatDB{db: db}, nil
}

func (c *CatDB) Close() error {
	fmt.Println("Closing CatDB")
	fmt.Println(c.DB().Stats())

	return c.db.Close()
}

func (c *CatDB) DB() *bbolt.DB {
	return c.db
}

func (c *CatDB) InitDB() error {
	return c.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketCats))
		if err != nil {
			return err
		}
		return nil
	})
}

func (c *CatDB) AddCatVersion(metadata *metadata.CatMetadata, catData []byte) (string, string, error) {
	catId := metadata.ID
	versionId, err := hashURL(metadata.URL)
	if err != nil {
		return "", "", err
	}

	meta, err := json.Marshal(metadata)
	if err != nil {
		return "", "", err
	}

	err = c.db.Update(func(tx *bbolt.Tx) error {

		b := tx.Bucket([]byte(bucketCats))
		cat, err := b.CreateBucketIfNotExists([]byte(catId))
		if err != nil {
			return err
		}
		version, err := cat.CreateBucketIfNotExists([]byte(versionId))
		if err != nil {
			return err
		}
		err = version.Put([]byte(dbKeyMetadata), meta)
		if err != nil {
			return err
		}
		err = version.Put([]byte(dbKeyData), catData)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", "", err
	}
	return catId, versionId, nil
}

func hashURL(url string) (string, error) {
	h := fnv.New64a()
	_, err := h.Write([]byte(url))
	if err != nil {
		log.Default().Println(err)
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum64()), nil
}
