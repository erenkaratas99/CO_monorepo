package helpers

import (
	sharedentities "CustomerOrderMonoRepo/shared/entities"
	"encoding/json"
	"fmt"
	"github.com/erenkaratas99/COApiCore/pkg/customErrors"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"net/http"
	"strconv"
)

type ExactFilter struct {
	Key []interface{}
	Val [][]interface{}
}

type MatchFilter struct {
	Key []interface{}
	Val [][]interface{}
}

type Field struct {
	DoInclude interface{}
	Field     []interface{}
}

type Sort struct {
	Value  interface{}
	HiToLo interface{}
}

type FSF struct {
	ExactFilter ExactFilter
	MatchFilter MatchFilter
	Field       *Field
	SortCond    Sort
}

func ReadFieldsFromJSON(path string) (*interface{}, error) {
	data, err := ioutil.ReadFile(path)
	var castData interface{}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &castData)
	if err != nil {
		return nil, err
	}
	return &castData, nil
}

func (cf *CheckFields) FieldsCompatibilityCheck(fsf *FSF) error {
	var fieldsAppended []interface{}
	for _, fieldRequested := range fsf.ExactFilter.Key {
		fieldsAppended = append(fieldsAppended, fieldRequested)
	}
	for _, fieldRequested := range fsf.MatchFilter.Key {
		fieldsAppended = append(fieldsAppended, fieldRequested)
	}
	isField := ListIterator(fieldsAppended, cf.Fields)
	if !isField {
		return customErrors.NewHTTPError(400, "FieldErr",
			"Requested field for filtering not could have been found.")
	}
	isFieldExactComp := ListIterator(fsf.ExactFilter.Key, cf.ExactFields)
	if !isFieldExactComp {
		return customErrors.NewHTTPError(http.StatusBadRequest, "FieldErr",
			"Requested field for exact search is not compatible for that.")
	}
	isFieldMatchComp := ListIterator(fsf.MatchFilter.Key, cf.MatchFields)
	if !isFieldMatchComp {
		return customErrors.NewHTTPError(http.StatusBadRequest, "FieldErr",
			"Requested field for match filtering is not compatible for that.")
	}
	var valTempArr []interface{}
	if fsf.SortCond.Value != nil {
		valTempArr = append(valTempArr, fsf.SortCond.Value)
		isFieldSortable := ListIterator(valTempArr, cf.SortableFields)
		if !isFieldSortable {
			return customErrors.NewHTTPError(http.StatusBadRequest, "FieldErr",
				"Requested field for sorting is not compatible for that.")
		}
	}
	return nil
}
func ListIterator(arr1 []interface{}, arr2 []string) bool {
	if arr1 == nil {
		return true
	}
	var notExist []string
	for _, val1 := range arr1 {
		val1 := fmt.Sprintf("%v", val1)
		match := false
		for _, val2 := range arr2 {
			if val1 == val2 {
				match = true
				break
			}
		}
		if !match {
			notExist = append(notExist, val1)
		}
	}
	if len(notExist) > 0 {
		return false
	}
	return true
}

func NewExactFilter(jsonMap sharedentities.JsonMap) *ExactFilter {
	filters := jsonMap.ExactFilters
	if filters == nil {
		return &ExactFilter{
			Key: nil,
			Val: nil,
		}
	}
	var valueArr [][]interface{}
	var keyArr []interface{}
	for k, v := range filters {
		valueArr = append(valueArr, v)
		keyArr = append(keyArr, k)
	}
	return &ExactFilter{
		Val: valueArr,
		Key: keyArr,
	}
}
func NewMatchFilter(jsonMap sharedentities.JsonMap) *MatchFilter {
	filters := jsonMap.Match
	if filters == nil {
		return &MatchFilter{
			Key: nil,
			Val: nil,
		}
	}
	var valueArr [][]interface{}
	var keyArr []interface{}
	for k, v := range filters {
		valueArr = append(valueArr, v)
		keyArr = append(keyArr, k)
	}
	return &MatchFilter{
		Val: valueArr,
		Key: keyArr,
	}
}

//func NewField(jsonMap entities.JsonMap) *Field {
//	fields := jsonMap.Fields
//	if fields == nil || fields["fields"] == nil {
//		return &Field{
//			DoInclude: nil,
//			Field:     nil,
//		}
//	}
//	doInclude, _ := strconv.Atoi(fields["fields"][0].(string))
//	var fieldArr []interface{}
//	for i := 1; i < len(fields["fields"]); i++ {
//		fieldArr = append(fieldArr, fmt.Sprintf("%v", fields["fields"][i]))
//	}
//	return &Field{
//		DoInclude: doInclude,
//		Field:     fieldArr,
//	}
//
//}
func NewField(jsonMap sharedentities.JsonMap) *Field {
	fields := jsonMap.Fields
	if fields == nil || fields["fields"] == nil {
		return &Field{
			DoInclude: nil,
			Field:     nil,
		}
	}
	doInclude, _ := strconv.Atoi(fields["fields"][0].(string))
	var fieldArr []interface{}
	for i := 1; i < len(fields["fields"]); i++ {
		val, ok := fields["fields"][i].(string)
		if !ok {
			return nil
		}
		fieldArr = append(fieldArr, val)
	}
	return &Field{
		DoInclude: doInclude,
		Field:     fieldArr,
	}
}

func NewSortCondition(jsonMap sharedentities.JsonMap) *Sort {
	sortCond := jsonMap.SortCond
	if len(sortCond) == 0 {
		return &Sort{
			Value:  nil,
			HiToLo: nil,
		}
	}
	field := fmt.Sprintf("%v", sortCond["field"])
	hitoloVal := fmt.Sprintf("%v", sortCond["hitolo"])
	var hitolo bool
	if hitoloVal == "-1" {
		hitolo = true
	} else {
		hitolo = false
	}
	return &Sort{Value: field, HiToLo: hitolo}
}

func NewFSF(jsonMap sharedentities.JsonMap) *FSF {
	ef := NewExactFilter(jsonMap)
	mf := NewMatchFilter(jsonMap)
	fields := NewField(jsonMap)
	sc := NewSortCondition(jsonMap)

	return &FSF{
		ExactFilter: *ef,
		MatchFilter: *mf,
		Field:       fields,
		SortCond:    *sc,
	}
}

func ListRevisedIterator(arr1 []interface{}, arr2 []string) bool {
	if arr1 == nil {
		return true
	}
	arr2Map := make(map[string]bool)
	for _, val := range arr2 {
		arr2Map[val] = true
	}
	for _, val1 := range arr1 {
		val1 := fmt.Sprintf("%v", val1)
		if _, ok := arr2Map[val1]; !ok {
			return false
		}
	}
	return true
}

func (fsf *FSF) GetProjection() *bson.D {
	var projection bson.D
	if fsf.Field.Field != nil && fsf.Field.DoInclude != nil {
		for i := 0; i < len(fsf.Field.Field); i++ {
			projection = append(projection, bson.E{fmt.Sprintf("%v", fsf.Field.Field[i]), fsf.Field.DoInclude})
		}
	} else {
		return nil
	}
	return &projection
}

func (fsf *FSF) GetFilter() *bson.D {
	matchFilterArrPerKey := make(map[string]bson.A)
	var filterExArr bson.A
	var matchFilterArr bson.A
	var filter bson.D
	exFilter := fsf.ExactFilter
	matchFilter := fsf.MatchFilter
	if exFilter.Key != nil && exFilter.Val != nil {
		for i := 0; i < len(exFilter.Key); i++ {
			for j := 0; j < len(exFilter.Val[i]); j++ {
				filterExArr = append(filterExArr, bson.D{{fmt.Sprintf("%v", exFilter.Key[i]),
					fmt.Sprintf("%v", exFilter.Val[i][j])}})
			}
		}

	}
	var strKeys []string
	if matchFilter.Key != nil && matchFilter.Val != nil {
		for i := 0; i < len(matchFilter.Key); i++ {
			temp := fmt.Sprintf("%v", matchFilter.Key[i])
			strKeys = append(strKeys, temp)
			for j := 0; j < len(matchFilter.Val[i]); j++ {
				matchFilterArrPerKey[temp] = append(matchFilterArrPerKey[temp], bson.M{fmt.Sprintf("%v", matchFilter.Key[i]): bson.M{"$gte": fmt.Sprintf("%v", matchFilter.Val[i][j]), "$lt": fmt.Sprintf("%v", matchFilter.Val[i][j]) + "\uffff"}})
			}
		}
	}
	//bson.D{{fmt.Sprintf("%v", matchFilter.Key[i]),
	//					primitive.Regex{Pattern: fmt.Sprintf("%v", matchFilter.Val[i][j]), Options: "i"}}}
	//filter :=
	for _, filters := range matchFilterArrPerKey {
		matchFilterArr = append(matchFilterArr, bson.D{{fmt.Sprintf("$or"), filters}})
	}

	var filterValArr bson.A
	if matchFilterArr != nil {
		filterValArr = append(filterValArr, bson.D{{"$and", matchFilterArr}})
	}
	if filterExArr != nil {
		filterValArr = append(filterValArr, bson.D{{"$and", filterExArr}})
	}
	filter = bson.D{{"$and", filterValArr}}
	if filterValArr == nil {
		filter = bson.D{{}}
	}
	return &filter
}

func (fsf *FSF) GetSortCondition() *bson.D {
	var sort bson.D
	if fsf.SortCond.Value != nil && fsf.SortCond.HiToLo != nil {
		if fsf.SortCond.HiToLo == true {
			sort = bson.D{{fmt.Sprintf("%v", fsf.SortCond.Value), -1}}
		} else {
			sort = bson.D{{fmt.Sprintf("%v", fsf.SortCond.Value), 1}}
		}
	} else {
		return nil
	}
	return &sort
}

func IsSubset(arr2 []interface{}, arr1 []string) bool {
	if arr2 == nil {
		return true
	}
	set := make(map[string]bool)
	for _, v := range arr1 {
		set[v] = true
	}
	for _, v := range arr2 {
		if _, ok := set[fmt.Sprint(v)]; !ok {
			return false
		}
	}
	return true
}
