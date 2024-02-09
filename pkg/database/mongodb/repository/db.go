package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

var ctx context.Context
var db *mongo.Client
var organizationCollection *mongo.Collection
