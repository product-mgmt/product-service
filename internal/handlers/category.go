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

func (s *Storage) GetCategoryHandler(w http.ResponseWriter, r *http.Request) error {

	ctxVal := r.Context().Value(types.CTXKey{Key: messages.PAGINATE}).(types.Paginate)

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// fetching all products
	rows, err := s.sqlStore.GetRecords(ctx, procedures.SP_GETRECORDS, tables.CATEGORY, ctxVal.SearchColumn, ctxVal.SearchTerm, ctxVal.SortColumn, ctxVal.SortOrder, ctxVal.Offset, ctxVal.Limit)
	if err != nil {
		return err
	}

	var category []types.ProductCategory
	if err := sqlscan.Rows(category, rows); err != nil {
		return err
	}

	return commfunc.WriteJSON(w, http.StatusOK, category)
}

func (s *Storage) AddCategoryHandler(w http.ResponseWriter, r *http.Request) error {
	// sync request body data with ProductCategoryRequest
	req := new(types.ProductCategoryRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// create a new category
	category := types.NewProductCategory(req.Name, req.Description)

	// add new product in db
	createdCategory, err := s.sqlStore.AddReord(ctx, procedures.SP_CREATE_PRODUCT_CATEGORY, category.Name, category.Description)
	if err != nil {
		return err
	}

	var output types.StandardResponse
	if err := sqlscan.Row(&output, createdCategory); err != nil {
		s.logger.Error(fmt.Errorf(messages.CATEGORYCREATE, err.Error()))
		return fmt.Errorf(messages.SOMETHINGWENTWRONG)
	}

	return commfunc.WriteJSON(w, http.StatusOK, output)
}

func (s *Storage) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Storage) ViewCategoryHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		return fmt.Errorf(messages.IDREQUIRED)
	}

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.sqlStore.GetRecordByArgs(ctx, procedures.SP_GETRECORD, tables.CATEGORY, "id", id)
	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf(messages.RECORDNOTFOUND, &id)
	}

	var output types.ProductCategory
	if err := sqlscan.Row(&output, rows); err != nil {
		s.logger.Error(fmt.Errorf(messages.CATEGORYFETCH, err.Error()))
		return fmt.Errorf(messages.SOMETHINGWENTWRONG)
	}

	return commfunc.WriteJSON(w, http.StatusOK, output)
}

func (s *Storage) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		return fmt.Errorf(messages.IDREQUIRED)
	}

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.sqlStore.GetRecordByArgs(ctx, procedures.SP_GETRECORD, tables.CATEGORY, "id", id)
	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf(messages.RECORDNOTFOUND, &id)
	}

	// soft delete
	row, err := s.sqlStore.DeleteRecordByArgs(ctx, procedures.SP_SOFTDELETE, tables.CATEGORY, "id", id)
	if err != nil {
		return err
	}

	var output types.StandardResponse
	if err := sqlscan.Row(&output, row); err != nil {
		s.logger.Error(fmt.Errorf(messages.CATEGORYDELETE, err.Error()))
		return fmt.Errorf(messages.SOMETHINGWENTWRONG)
	}

	return commfunc.WriteJSON(w, http.StatusOK, output)
}
