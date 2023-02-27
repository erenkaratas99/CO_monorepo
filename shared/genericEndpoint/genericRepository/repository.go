package genericRepository

import (
	"CustomerOrderMonoRepo/shared/helpers"
	"context"
	"fmt"
	"github.com/loov/hrtime"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"
)

type Repository struct {
	collection *mongo.Collection
}

func NewGenericRepository(mc *mongo.Collection) *Repository {
	return &Repository{collection: mc}
}

func (r *Repository) GenericGetRepo(FSF *helpers.FSF, getTotalCount bool) ([]*interface{}, *int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := options.Find().SetLimit(8).SetSkip(0)
	projection := FSF.GetProjection()
	if projection != nil {
		opts = opts.SetProjection(projection)
	}
	filter := FSF.GetFilter()
	sort := FSF.GetSortCondition()
	if sort != nil {
		opts = opts.SetSort(sort)
	}
	var wg sync.WaitGroup
	//var order []*entities.Order
	var order []*interface{}
	var err error
	var totalCountG int
	wg.Add(1)
	go func() {
		defer wg.Done()
		start_tc := hrtime.Now()
		cur, err := r.collection.Find(ctx, filter, opts)
		if err != nil {
			return
		}
		defer cur.Close(ctx)
		err = cur.All(ctx, &order)
		elapsed := hrtime.Since(start_tc).Milliseconds()
		fmt.Println("elapsed : ", elapsed)
	}()

	if getTotalCount {
		totalCount64, _ := r.collection.CountDocuments(ctx, filter)
		totalCountG = int(totalCount64)
	}
	wg.Wait()
	if err != nil {
		return nil, nil, err
	}
	return order, &totalCountG, nil
}
