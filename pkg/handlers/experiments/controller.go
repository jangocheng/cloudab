package experiments

import (
	"context"
	"fmt"
	"time"

	"github.com/caicloud/nirvana/log"
	v1 "github.com/dyweb/cloudab/pkg/apis/v1"
	"github.com/dyweb/cloudab/pkg/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Controller struct {
	collection string
}

func New() Controller {
	return Controller{
		collection: "experiments",
	}
}

func (c Controller) CreateExperiment(
	ctx context.Context, exp *v1.Experiment) (*v1.Experiment, error) {
	cctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	logger := log.DefaultLogger()
	logger.V(log.LevelDebug).Infof(
		"Creating the experiment %s in collection %s", exp.Name, c.collection)

	if _, err := store.DB.Collection(c.collection).InsertOne(cctx, exp); err != nil {
		return nil, err
	}
	return exp, nil
}

// ListExperiments returns all experiments.
func (c Controller) ListExperiments(ctx context.Context, count int) ([]v1.Experiment, error) {
	cctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	var exp v1.Experiment
	exps, err := store.DB.Collection(c.collection).Find(cctx, exp)
	if err != nil {
		return nil, err
	}

	results := &[]v1.Experiment{}
	if err := exps.All(cctx, results); err != nil {
		return nil, err
	}
	return *results, nil
}

// GetExperiment return an experiment by id.
func (c Controller) GetExperiment(ctx context.Context, id primitive.ObjectID) (*v1.Experiment, error) {
	cctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	exp := store.DB.Collection(c.collection).FindOne(cctx, bson.D{{"_id", id}})
	if exp == nil {
		return nil, fmt.Errorf("Failed to find the experiment with id %s", id.Hex())
	}
	var result v1.Experiment
	if err := exp.Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}
