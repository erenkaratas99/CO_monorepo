package repositories

import (
	"CustomerOrderMonoRepo/CustomerService/src/internal/entities"
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
func (r *Repository) InsertCustomer(customerReq *entities.CustomerRequestModel) (*string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	id := uuid.New()
	customerId := id.String()
	customerReq.Address.CustomerId = customerId
	timeNow := time.Now().Format(time.RFC3339)
	c := bson.M{
		"_id":        customerId,
		"first_name": customerReq.FirstName,
		"last_name":  customerReq.LastName,
		"email":      customerReq.Email,
		"phone":      customerReq.Phone,
		"address":    customerReq.Address,
		"created_at": timeNow,
		"updated_at": timeNow,
	}
	res, err := r.Collection.InsertOne(ctx, c)
	if err != nil {
		return nil, err
	}
	insertedId := res.InsertedID.(string)
	return &insertedId, nil
}

func (r *Repository) GetSingleCustomer(id string, getTotalCount bool) (*entities.CustomerResponseModel, *int, error) {
	filter := bson.M{"_id": id}
	customer := entities.CustomerResponseModel{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err := r.Collection.FindOne(ctx, filter).Decode(&customer)
	if err != nil {
		return nil, nil, customErrors.DocNotFound
	}
	var totalCountG int
	if getTotalCount {
		totalCount64, _ := r.Collection.CountDocuments(ctx, filter)
		totalCountG = int(totalCount64)
	}
	return &customer, &totalCountG, nil
}

func (r *Repository) GetAllCustomers(l, o int64, getTotalCount bool) ([]*entities.CustomerResponseModel, *int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{}
	opts := options.Find().SetLimit(l).SetSkip(o)
	var err error
	var customer []*entities.CustomerResponseModel
	var wg sync.WaitGroup
	var totalCountG int
	wg.Add(1)
	go func() {
		defer wg.Done()
		cur, err := r.Collection.Find(ctx, filter, opts)
		if err != nil {
			return
		}
		defer cur.Close(ctx)
		err = cur.All(ctx, &customer)
	}()
	if getTotalCount {
		totalCount64, _ := r.Collection.CountDocuments(ctx, filter)
		totalCountG = int(totalCount64)
	}
	wg.Wait()
	if err != nil {
		return nil, nil, err
	}
	return customer, &totalCountG, nil
}

func (r *Repository) UpdateCustomer(c *entities.Customer) (error, *string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{"_id": c.Id}
	uat := time.Now().Format(time.RFC3339)
	update := bson.M{
		"$set": bson.M{
			"first_name": c.FirstName,
			"last_name":  c.LastName,
			"email":      c.Email,
			"phone":      c.Phone,
			"address": bson.M{
				"customer_id":  c.Id,
				"address_name": c.Address.AddressName,
				"address_line": c.Address.AddressLine,
				"city":         c.Address.City,
				"country":      c.Address.Country,
				"city_code":    c.Address.CityCode,
			},
			"updatedAt": uat,
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
	return nil, &c.Id
}

func (r *Repository) DeleteCustomer(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	filter := bson.M{"_id": id}
	res, err := r.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return customErrors.NewHTTPError(http.StatusBadRequest,
			"IdErr", "No customer found with given ID.")
	}
	return nil
}

func (r *Repository) GetAddressOfCustomer(id string) (*entities.CustomerAddressResponseModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	opts := options.FindOne().SetProjection(bson.D{{"address", 1}})
	filter := bson.M{"_id": id}
	addressResp := entities.CustomerAddressResponseModel{}
	err := r.Collection.FindOne(ctx, filter, opts).Decode(&addressResp)
	if err != nil {
		return nil, err
	}
	return &addressResp, nil
}
