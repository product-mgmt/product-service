package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ankeshnirala/sqlscan"
	"github.com/gorilla/mux"
	"github.com/product-mgmt/common-service/constants/messages"
	"github.com/product-mgmt/common-service/constants/procedures"
	"github.com/product-mgmt/common-service/constants/tables"
	"github.com/product-mgmt/common-service/types"
	"github.com/product-mgmt/common-service/utils/commfunc"
)

func (s *Storage) GetProductsHandler(w http.ResponseWriter, r *http.Request) error {

	ctxVal := r.Context().Value(types.CTXKey{Key: messages.PAGINATE}).(types.Paginate)

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// fetching all products
	rows, err := s.sqlStore.GetRecords(ctx, procedures.SP_GETRECORDS, tables.PRODUCTS, ctxVal.SearchColumn, ctxVal.SearchTerm, ctxVal.SortColumn, ctxVal.SortOrder, ctxVal.Offset, ctxVal.Limit)
	if err != nil {
		return err
	}

	var users []types.ProductSummarized
	if err := sqlscan.Rows(&users, rows); err != nil {
		return err
	}

	return commfunc.WriteJSON(w, http.StatusOK, users)
}

func (s *Storage) AddProductHandler(w http.ResponseWriter, r *http.Request) error {
	// sync request body data with AddProductRequest
	req := new(types.AddProductRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.sqlStore.GetRecordByArgs(ctx, procedures.SP_GETRECORD, tables.PRODUCTS, "sku", req.SKU)
	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf(messages.ALREDAYEXISTS, req.SKU)
	}

	// create a new product
	newProduct := types.NewProduct(req.Name, req.Description, req.SKU, req.CategoryID, req.Price)

	// add new product in db
	product, err := s.sqlStore.AddReord(ctx, procedures.SP_CREATE_PRODUCTS, newProduct.Name, newProduct.Description, newProduct.SKU, newProduct.CategoryID, newProduct.Price)
	if err != nil {
		return err
	}

	var output types.StandardResponse
	if err := sqlscan.Row(&output, product); err != nil {
		s.logger.Error(fmt.Errorf(messages.PRODUCTCREATING, err.Error()))
		return fmt.Errorf(messages.SOMETHINGWENTWRONG)
	}

	return commfunc.WriteJSON(w, http.StatusOK, output)
}

func (s *Storage) UpdateProductHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Storage) ViewProductHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		return fmt.Errorf(messages.IDREQUIRED)
	}

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.sqlStore.GetRecordByArgs(ctx, procedures.SP_GETRECORD, tables.PRODUCTS, "id", id)
	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf(messages.RECORDNOTFOUND, &id)
	}

	var output types.Product
	if err := sqlscan.Row(&output, rows); err != nil {
		s.logger.Error(fmt.Errorf(messages.PRODUCTFETCHING, err.Error()))
		return fmt.Errorf(messages.SOMETHINGWENTWRONG)
	}

	return commfunc.WriteJSON(w, http.StatusOK, output)
}

func (s *Storage) DeleteProductHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		return fmt.Errorf(messages.IDREQUIRED)
	}

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.sqlStore.GetRecordByArgs(ctx, procedures.SP_GETRECORD, tables.PRODUCTS, "id", id)
	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf(messages.RECORDNOTFOUND, &id)
	}

	// soft delete
	row, err := s.sqlStore.DeleteRecordByArgs(ctx, procedures.SP_SOFTDELETE, tables.PRODUCTS, "id", id)
	if err != nil {
		return err
	}

	var output types.StandardResponse
	if err := sqlscan.Row(&output, row); err != nil {
		s.logger.Error(fmt.Errorf(messages.PRODUCTDELETING, err.Error()))
		return fmt.Errorf(messages.SOMETHINGWENTWRONG)
	}

	return commfunc.WriteJSON(w, http.StatusOK, output)
}
