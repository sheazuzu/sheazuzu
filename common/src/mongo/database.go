/*
 * database.go Created on 31.08.2021Copyright (C) 2021 Volkswagen AG, All rights reserved.
 */

// Package mongo provides constants and functions for the creation of a mongo database and the handling of the connection to it.
package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"sheazuzu/common/src/utils"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Database struct {
	Config   *Config
	Client   *mongo.Client
	Database *mongo.Database
	Logger   *zap.SugaredLogger
}

// NewMongoDatabase creates a new mongo database with the passed options and returns a pointer to it.
func NewMongoDatabase(config *Config, logger *zap.SugaredLogger, setOptions ...func(opts *options.ClientOptions)) (*Database, error) {

	logger.Debug("Connecting to mongo db ...")

	opts := options.Client().ApplyURI(config.URI)
	if config.UseSSL {
		logger.Debug("Using SSL certificates ...")
		err := configureSSL(config, opts)
		if err != nil {
			return nil, err
		}
	}

	for _, opt := range setOptions {
		opt(opts)
	}

	client, err := mongo.NewClient(opts)
	if err != nil {
		return nil, errors.Wrap(err, "Could not create mongo client")
	}

	return &Database{Config: config, Client: client, Logger: logger}, nil
}

// Connect initializes the database client and establishes the connection.
func (db *Database) Connect(ctx context.Context) error {

	err := db.Client.Connect(ctx)
	if err != nil {
		return err
	}
	db.Database = db.Client.Database(db.Config.Database)

	return nil
}

// Disconnect disconnects from the mongo db.
func (db *Database) Disconnect(ctx context.Context) {

	db.Logger.Debug("Disconnecting from mongo db ...")
	err := db.Client.Disconnect(ctx)
	if err != nil {
		db.Logger.Errorf("Could not disconnect from mongo db: %s", err.Error())
	}
	db.Logger.Debug("Disconnect done.")
}

// Checks if an index of a given name exists in a collection of the database
func (db *Database) existIndex(ctx context.Context, collectionName string, name string) (bool, error) {

	// access the collection of the given name
	collection := db.Database.Collection(collectionName)

	// there is a special view in the collection containing the indexes
	// we have to iterate through them and check each name
	indexView := collection.Indexes()
	batchSize := int32(10)
	maxTime := time.Duration(db.Config.Timeout) * time.Second

	cur, err := indexView.List(ctx, &options.ListIndexesOptions{
		BatchSize: &batchSize,
		MaxTime:   &maxTime,
	})
	if err != nil {
		return false, errors.Wrapf(err, "could not list the indexes in collection '%s'", collectionName)
	}
	defer CloseCursor(cur, ctx)

	for cur.Next(ctx) {
		var indexDoc bson.D
		err := cur.Decode(&indexDoc)
		if err != nil {
			return false, errors.Wrapf(err, "Failed to decode an index model in collection '%s'", collectionName)
		}
		index := indexDoc.Map()

		if indexName, ok := index["name"]; ok && strings.EqualFold(name, indexName.(string)) {
			return true, nil
		}
	}

	// not found
	return false, nil
}

// creates an index in a collection
func (db *Database) createIndex(ctx context.Context, collectionName string, index *mongo.IndexModel) error {

	collection := db.Database.Collection(collectionName)
	indexView := collection.Indexes()
	maxTime := time.Duration(db.Config.Timeout) * time.Second
	_, err := indexView.CreateOne(ctx, *index, &options.CreateIndexesOptions{MaxTime: &maxTime})
	if err != nil {
		return err
	}
	return nil
}

// InstallIndex installs an index for the given collection, if it does not exist yet.
// An already existing index will not be overwritten.
func (db *Database) InstallIndex(collectionName string, name string, keys bson.D) error {

	db.Logger.Infof("Installing mongo db index '%s' in collection '%s'...", name, collectionName)

	ctx, _ := context.WithTimeout(context.Background(), time.Duration(db.Config.Timeout)*time.Second)
	exists, err := db.existIndex(ctx, collectionName, name)
	if err != nil {
		return errors.Wrapf(err, "Failed to check index '%s' in collection '%s'", name, collectionName)
	}

	if !exists {
		db.Logger.Debugf("The index '%s' in collection '%s' does not exist, creating it ...", name, collectionName)
		index := &mongo.IndexModel{
			Keys: keys,
			Options: &options.IndexOptions{
				Name:       utils.ToStringPtr(name),
				Background: utils.ToBoolPtr(true), // create the index in the background to avoid any blocking
			},
		}

		// new context with new timeout for this second operation
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(db.Config.Timeout)*time.Second)
		err = db.createIndex(ctx, collectionName, index)
		if err != nil {
			return errors.Wrapf(err, "Failed to create index '%s' in collection '%s'", name, collectionName)
		}
	}

	return nil
}

func configureSSL(config *Config, opts *options.ClientOptions) error {
	// User provided both certificates, so its a client and a key cert
	if config.SSLClientCertFile != "" && config.SSLClientKeyFile != "" {
		cert, err := tls.LoadX509KeyPair(config.SSLClientCertFile, config.SSLClientKeyFile)
		if err != nil {
			return errors.Wrap(err, "Error creating tls config")
		}

		opts.SetTLSConfig(&tls.Config{
			Certificates:       []tls.Certificate{cert},
			InsecureSkipVerify: true,
		})
	} else {
		// User provided only a client certificate
		rootPEM, err := ioutil.ReadFile(config.SSLClientCertFile)
		if err != nil {
			return errors.Wrap(err, "could not read ssl client cert")
		}

		roots := x509.NewCertPool()
		ok := roots.AppendCertsFromPEM(rootPEM)
		if !ok {
			return fmt.Errorf("could not parse certificate")
		}

		opts.SetTLSConfig(&tls.Config{
			ClientCAs:          roots,
			InsecureSkipVerify: true,
		})
	}

	return nil
}

// CloseCursor closes the cursor and ignores the error. Can be used with defer.
func CloseCursor(cur *mongo.Cursor, ctx context.Context) {
	_ = cur.Close(ctx)
}
