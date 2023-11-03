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

func (s *Storage) GetInventoryHandler(w http.ResponseWriter, r *http.Request) error {

	ctxVal := r.Context().Value(types.CTXKey{Key: messages.PAGINATE}).(types.Paginate)

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// fetching all products
	rows, err := s.sqlStore.GetRecords(ctx, procedures.SP_GETRECORDS, tables.INVENTORY, ctxVal.SearchColumn, ctxVal.SearchTerm, ctxVal.SortColumn, ctxVal.SortOrder, ctxVal.Offset, ctxVal.Limit)
	if err != nil {
		return err
	}

	var users []types.ProductInventory
	if err := sqlscan.Rows(&users, rows); err != nil {
		return err
	}

	return commfunc.WriteJSON(w, http.StatusOK, users)
}

func (s *Storage) AddInventoryHandler(w http.ResponseWriter, r *http.Request) error {
	// sync request body data with AddProductRequest
	req := new(types.ProductInventoryRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.sqlStore.GetRecordByArgs(ctx, procedures.SP_GETRECORD, tables.PRODUCTS, "product_id", req.ProductId)
	if err != nil {
		return err
	}

	if rows.Next() {
		return fmt.Errorf(messages.RECORDNOTFOUND, req.ProductId)
	}

	// create a new product
	newInventory := types.NewProductInventory(req.ProductId, req.Quantity)

	// add new product in db
	inventory, err := s.sqlStore.AddReord(ctx, procedures.SP_CREATE_PRODUCT_INVENTORY, newInventory.ProductId, newInventory.Quantity)
	if err != nil {
		return err
	}

	var output types.StandardResponse
	if err := sqlscan.Row(&output, inventory); err != nil {
		s.logger.Error(fmt.Errorf(messages.INVENTORYCREATE, err.Error()))
		return fmt.Errorf(messages.SOMETHINGWENTWRONG)
	}

	return commfunc.WriteJSON(w, http.StatusOK, output)
}

func (s *Storage) UpdateInventoryHandler(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *Storage) ViewInventoryHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		return fmt.Errorf(messages.IDREQUIRED)
	}

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.sqlStore.GetRecordByArgs(ctx, procedures.SP_GETRECORD, tables.INVENTORY, "id", id)
	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf(messages.RECORDNOTFOUND, &id)
	}

	var output types.ProductInventory
	if err := sqlscan.Row(&output, rows); err != nil {
		s.logger.Error(fmt.Errorf(messages.INVENTORYFETCH, err.Error()))
		return fmt.Errorf(messages.SOMETHINGWENTWRONG)
	}

	return commfunc.WriteJSON(w, http.StatusOK, output)
}

func (s *Storage) DeleteInventoryHandler(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id := vars["id"]

	if id == "" {
		return fmt.Errorf(messages.IDREQUIRED)
	}

	// create a context to timeout db operation once work end
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.sqlStore.GetRecordByArgs(ctx, procedures.SP_GETRECORD, tables.INVENTORY, "id", id)
	if err != nil {
		return err
	}

	if !rows.Next() {
		return fmt.Errorf(messages.RECORDNOTFOUND, &id)
	}

	// soft delete
	row, err := s.sqlStore.DeleteRecordByArgs(ctx, procedures.SP_SOFTDELETE, tables.INVENTORY, "id", id)
	if err != nil {
		return err
	}

	var output types.StandardResponse
	if err := sqlscan.Row(&output, row); err != nil {
		s.logger.Error(fmt.Errorf(messages.INVENTORYDELETE, err.Error()))
		return fmt.Errorf(messages.SOMETHINGWENTWRONG)
	}

	return commfunc.WriteJSON(w, http.StatusOK, output)
}
