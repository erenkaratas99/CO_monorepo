package repositories

import (
	"CustomerOrderMonoRepo/OrderService/src/internal/entities"
	"context"
	"github.com/erenkaratas99/COApiCore/pkg/customErrors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"sync"
	"time"
)

type Repository struct {
	Collection *mongo.Collection
}

func NewRepository(mc *mongo.Collection) *Repository {
	return &Repository{Collection: mc}
}

func (r *Repository) CreateOrder(order *entities.Order) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	id := uuid.New()
	timeNow := time.Now().Format(time.RFC3339)
	c := bson.M{
		"_id":             id.String(),
		"customer_id":     order.CustomerID,
		"order_date":      timeNow,
		"order_total":     order.OrderTotal,
		"payment_status":  order.PaymentStatus,
		"shipment_status": order.ShipmentStatus,
		"order_item":      order.OrderItem,
		"address":         order.Address,
		"updated_at":      timeNow,
	}
	res, err := r.Collection.InsertOne(ctx, c)
	if err != nil {
		return nil, err
	}
	insertedId := res.InsertedID.(string)
	return &insertedId, nil
}

func (r *Repository) GetSingleOrder(orderid string, getTotalCount bool) (*entities.Order, *int, error) {
	filter := bson.M{"_id": orderid}
	order := entities.Order{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := r.Collection.FindOne(ctx, filter).Decode(&order)
	if err != nil {
		return nil, nil, customErrors.DocNotFound
	}
	var totalCountG int
	if getTotalCount {
		totalCount64, _ := r.Collection.CountDocuments(ctx, filter)
		totalCountG = int(totalCount64)
	}
	return &order, &totalCountG, nil
}

func (r *Repository) GetAllOrders(l int64, o int64, getTotalCount bool) ([]*entities.OrderResponseModel, *int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	blankFilter := bson.M{}
	opts := options.Find().SetLimit(l).SetSkip(o)
	var err error
	var order []*entities.OrderResponseModel
	var wg sync.WaitGroup
	var totalCountG int
	wg.Add(1)
	go func() {
		defer wg.Done()
		cur, err := r.Collection.Find(ctx, blankFilter, opts)
		if err != nil {
			return
		}
		defer cur.Close(ctx)
		err = cur.All(ctx, &order)
	}()
	if getTotalCount {
		totalCount64, _ := r.Collection.CountDocuments(ctx, blankFilter)
		totalCountG = int(totalCount64)
	}
	wg.Wait()
	if err != nil {
		return nil, nil, err
	}
	return order, &totalCountG, err
}

func (r *Repository) DeleteOrder(orderid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{"_id": orderid}
	res, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return customErrors.NewHTTPError(http.StatusInternalServerError,
			"ServerError", "No order found with given ID.")
	}
	return nil
}

func (r *Repository) UpdateOrder(orderid string, order *entities.Order) (error, *string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{"_id": orderid}
	update := bson.M{
		"$set": bson.M{
			"order_total":     order.OrderTotal,
			"payment_status":  order.PaymentStatus,
			"shipment_status": order.ShipmentStatus,
			"order_item":      order.OrderItem,
			"updated_at":      time.Now().Format(time.RFC3339),
		},
	}
	res, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err, nil
	}
	if res.ModifiedCount == 0 {
		return customErrors.NewHTTPError(http.StatusInternalServerError,
			"ServerErr",
			"No document could had been modified."), nil
	}
	return nil, &order.Id
}

func (r *Repository) GetCustomerOrders(customerID string, l, o int64) ([]*entities.OrderResponseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{"customer_id": customerID}
	opts := options.Find().SetLimit(l).SetSkip(o).SetProjection(bson.D{{"shipment_status", 1}})
	cur, err := r.Collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var orderStats []*entities.OrderResponseModel
	err = cur.All(ctx, &orderStats)
	return orderStats, err
}
