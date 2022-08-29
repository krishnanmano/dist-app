package badgerdb

import (
	"context"
	"dist-app/logger"
	"fmt"
	"os"
	"time"

	"github.com/dgraph-io/badger"
)

const (
	// Default BadgerDB discardRatio. It represents the discard ratio for the
	// BadgerDB GC.
	//
	// Ref: https://godoc.org/github.com/dgraph-io/badger#DB.RunValueLogGC
	badgerDiscardRatio = 0.5

	// Default BadgerDB GC interval
	badgerGCInterval = 10 * time.Minute
)

var (
	// BadgerAlertNamespace defines the alerts BadgerDB namespace.
	BadgerAlertNamespace = []byte("alerts")
	BDBClient            DB
)

// DB defines an embedded key/value store database interface.
type DB interface {
	Get(namespace, key []byte) (value []byte, err error)
	PrefixScan(namespace []byte) (value [][]byte, err error)
	Set(namespace, key, value []byte) error
	Has(namespace, key []byte) (bool, error)
	Delete(namespace, key []byte) error
	Close() error
}

// BadgerDB is a wrapper around a BadgerDB backend database that implements
// the DB interface.
type BadgerDB struct {
	db         *badger.DB
	ctx        context.Context
	cancelFunc context.CancelFunc
	logger     logger.ILogger
}

// NewBadgerDB returns a new initialized BadgerDB database implementing the DB
// interface. If the database cannot be initialized, an error will be returned.
func NewBadgerDB(ctx context.Context, dataDir string) (DB, error) {
	if err := os.MkdirAll(dataDir, 0774); err != nil {
		return nil, err
	}

	opts := badger.DefaultOptions(dataDir)
	opts.SyncWrites = true
	opts.Dir, opts.ValueDir = dataDir, dataDir
	badgerDB, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	bdb := &BadgerDB{
		db:     badgerDB,
		logger: logger.Log, //.With("module", "db"),
	}
	bdb.ctx, bdb.cancelFunc = context.WithCancel(ctx)

	go bdb.runGC()

	BDBClient = bdb
	return bdb, nil
}

// Get implements the DB interface. It attempts to get a value for a given key
// and namespace. If the key does not exist in the provided namespace, an error
// is returned, otherwise the retrieved value.
func (bdb *BadgerDB) Get(namespace, key []byte) (value []byte, err error) {
	err = bdb.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(badgerNamespaceKey(namespace, key))
		if err != nil {
			return err
		}

		fmt.Println("*****: ", item.Version())
		item.UserMeta()

		value, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return value, nil
}

func (bdb *BadgerDB) PrefixScan(namespace []byte) (value [][]byte, err error) {
	valueData := make([][]byte, 0)
	err = bdb.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := namespace
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%s, value=%s\n", k, v)
				valueData = append(valueData, v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return valueData, nil
}

// Set implements the DB interface. It attempts to store a value for a given key
// and namespace. If the key/value pair cannot be saved, an error is returned.
func (bdb *BadgerDB) Set(namespace, key, value []byte) error {
	err := bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Set(badgerNamespaceKey(namespace, key), value)
	})

	if err != nil {
		bdb.logger.Debug(fmt.Sprintf("failed to set key %s for namespace %s: %v", key, namespace, err))
		return err
	}

	return nil
}

func (bdb *BadgerDB) Delete(namespace, key []byte) error {
	err := bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Delete(badgerNamespaceKey(namespace, key))
	})

	if err != nil {
		bdb.logger.Debug(fmt.Sprintf("failed to delete key %s for namespace %s: %v", key, namespace, err))
		return err
	}

	return nil
}

// Has implements the DB interface. It returns a boolean reflecting if the
// datbase has a given key for a namespace or not. An error is only returned if
// an error to Get would be returned that is not of type badgerdb.ErrKeyNotFound.
func (bdb *BadgerDB) Has(namespace, key []byte) (ok bool, err error) {
	_, err = bdb.Get(namespace, key)
	switch err {
	case badger.ErrKeyNotFound:
		ok, err = false, nil
	case nil:
		ok, err = true, nil
	}

	return
}

// Close implements the DB interface. It closes the connection to the underlying
// BadgerDB database as well as invoking the context's cancel function.
func (bdb *BadgerDB) Close() error {
	bdb.cancelFunc()
	return bdb.db.Close()
}

// runGC triggers the garbage collection for the BadgerDB backend database. It
// should be run in a goroutine.
func (bdb *BadgerDB) runGC() {
	ticker := time.NewTicker(badgerGCInterval)
	for {
		select {
		case <-ticker.C:
			err := bdb.db.RunValueLogGC(badgerDiscardRatio)
			if err != nil {
				// don't report error when GC didn't result in any cleanup
				if err == badger.ErrNoRewrite {
					bdb.logger.Debug(fmt.Sprintf("no BadgerDB GC occurred: %v", err))
				} else {
					bdb.logger.Error(fmt.Sprintf("failed to GC BadgerDB: %v", err))
				}
			}

		case <-bdb.ctx.Done():
			return
		}
	}
}

// badgerNamespaceKey returns a composite key used for lookup and storage for a
// given namespace and key.
func badgerNamespaceKey(namespace, key []byte) []byte {
	prefix := []byte(fmt.Sprintf("%s/", namespace))
	return append(prefix, key...)
}
