package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/ankeshnirala/sqlscan"
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
	rows, err := s.sqlStore.GetRecords(ctx, procedures.SP_GETRECORDS, tables.USERS, ctxVal.SearchColumn, ctxVal.SearchTerm, ctxVal.SortColumn, ctxVal.SortOrder, ctxVal.Offset, ctxVal.Limit)
	if err != nil {
		return err
	}

	var users []types.ProductSummarized
	if err := sqlscan.Rows(&users, rows); err != nil {
		return err
	}

	return commfunc.WriteJSON(w, http.StatusOK, users)
}
