package products

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/bootcamp-go/custom-mocking/pkg/store"
	"github.com/stretchr/testify/assert"
)

type MockFunc struct {
	calledMethod string
	retValue     *interface{}
}

func (m *MockFunc) Return(val interface{}) {
	m.retValue = &val
}

type StoreMockWithOrder struct {
	Data      []Product
	Execs     []MockFunc
	execCount int
}

func (s *StoreMockWithOrder) Read(data interface{}) error {
	if s.Execs[s.execCount].calledMethod != "Read" {
		panic(fmt.Sprintf("se esperaba un llamado a %s pero se llamó Read", s.Execs[s.execCount].calledMethod))

	}
	ret := *s.Execs[s.execCount].retValue
	if valueToReturn, ok := ret.(error); ok {
		s.execCount++
		return valueToReturn
	}
	ref := data.(*[]Product)
	*ref = s.Data
	s.execCount++
	return nil
}

func (s *StoreMockWithOrder) Write(data interface{}) error {
	if s.Execs[s.execCount].calledMethod != "Write" {
		panic(fmt.Sprintf("se esperaba un llamado a %s pero se llamó Write", s.Execs[s.execCount].calledMethod))
	}
	ret := *s.Execs[s.execCount].retValue
	if valueToReturn, ok := ret.(error); ok {
		s.execCount++
		return valueToReturn
	}
	s.Data = data.([]Product)
	s.execCount++
	return nil
}

func (s *StoreMockWithOrder) On(methodName string) *MockFunc {
	h := MockFunc{calledMethod: methodName}
	s.Execs = append(s.Execs, h)
	return &s.Execs[len(s.Execs)-1]
}
func TestIntegrationServiceGetAll(t *testing.T) {
	input := []Product{
		{
			ID:    1,
			Name:  "CellPhone",
			Type:  "Tech",
			Count: 3,
			Price: 250,
		}, {
			ID:    2,
			Name:  "Notebook",
			Type:  "Tech",
			Count: 10,
			Price: 1750.5,
		},
	}
	dataJson, _ := json.Marshal(input)
	dbMock := store.Mock{
		Data: dataJson,
	}
	storeStub := store.FileStore{
		FileName: "",
		Mock:     &dbMock,
	}
	myRepo := NewRepository(&storeStub)
	myService := NewService(myRepo)

	result, err := myService.GetAll()

	assert.Equal(t, input, result)
	assert.Nil(t, err)
}

func TestServiceGetAllError(t *testing.T) {
	expectedError := errors.New("error for GetAll")
	dbMock := store.Mock{
		Err: expectedError,
	}
	storeStub := store.FileStore{
		FileName: "",
		Mock:     &dbMock,
	}
	myRepo := NewRepository(&storeStub)
	myService := NewService(myRepo)

	result, err := myService.GetAll()

	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
}

func TestStore(t *testing.T) {
	testProduct := Product{
		Name:  "CellPhone",
		Type:  "Tech",
		Count: 3,
		Price: 52.0,
	}
	encodedData, _ := json.Marshal([]Product{})
	dbMock := store.Mock{Data: encodedData}
	storeStub := store.FileStore{
		FileName: "",
		Mock:     &dbMock,
	}
	myRepo := NewRepository(&storeStub)
	myService := NewService(myRepo)
	result, _ := myService.Store(testProduct.Name, testProduct.Type,
		testProduct.Count, testProduct.Price)
	assert.Equal(t, testProduct.Name, result.Name)
	assert.Equal(t, testProduct.Type, result.Type)
	assert.Equal(t, testProduct.Price, result.Price)
	assert.Equal(t, 1, result.ID)
}

func TestStoreError(t *testing.T) {
	testProduct := Product{
		Name:  "CellPhone",
		Type:  "Tech",
		Count: 3,
		Price: 52.0,
	}
	expectedError := errors.New("error for Storage")
	dbMock := store.Mock{
		Err: expectedError,
	}
	storeStub := store.FileStore{
		FileName: "",
		Mock:     &dbMock,
	}
	myRepo := NewRepository(&storeStub)
	myService := NewService(myRepo)
	result, err := myService.Store(testProduct.Name, testProduct.Type,
		testProduct.Count, testProduct.Price)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, Product{}, result)
}

func TestStoreErrorFromDB(t *testing.T) {
	testProduct := Product{
		Name:  "CellPhone",
		Type:  "Tech",
		Count: 3,
		Price: 52.0,
	}
	storeStub := StoreMockWithOrder{Data: []Product{testProduct}}
	storeStub.On("Read").Return(nil)
	storeStub.On("Read").Return(nil)
	storeStub.On("Write").Return(fmt.Errorf("stub error"))
	myRepo := NewRepository(&storeStub)
	myService := NewService(myRepo)
	_, err := myService.Store(testProduct.Name, testProduct.Type,
		testProduct.Count, testProduct.Price)
	assert.NotNil(t, err)
}
