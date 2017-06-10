package db

import (
	"github.com/boltdb/bolt"
	"log"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
)

var acronBucket = []byte("acron")

func createDB() *bolt.DB {
	home, err := homedir.Dir()
	path := home + "/.acron"
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 755)
	}
	fmt.Println("Path: ", path)
	db, err := bolt.Open(path + "/bolt.db", 0644, nil)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Add(key, value []byte) {
	db := createDB()
	defer db.Close()

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists(acronBucket)
		if err != nil {
			return err
		}
		err = bucket.Put(key, value)
		if err != nil {
			return err
		}
		fmt.Printf("%s was added to the db\n", string(key))
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func Get(key []byte) string {
	db := createDB()
	defer db.Close()
	var value string
	// retrieve the data
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(acronBucket)
		if bucket == nil {
			return fmt.Errorf("Bucket %q not found!", acronBucket)
		}

		val := bucket.Get(key)
		value = string(val)

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	return value
}