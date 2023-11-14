package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"app/models"
)

func (h *handler) Product(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		h.CreateProduct(w, r)
	case "GET":
		var (
			values = r.URL.Query()
			method = values.Get("method")
		)

		if method == "GET_LIST" {
			h.GetListProduct(w, r)
		} else if method == "GET" {
			h.GetByIdProduct(w, r)
		}
	case "PUT":
		h.UpdateProduct(w, r)
	case "DELETE":
		h.DeleteProduct(w, r)
	}
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var productCreate models.ProductCreate
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.handlerResponse(w, "error while read product body: "+err.Error(), http.StatusBadRequest, nil)
		return
	}

	err = json.Unmarshal(body, &productCreate)
	if err != nil {
		h.handlerResponse(w, "error while unmarshal product body: "+err.Error(), http.StatusInternalServerError, nil)
		return
	}
	id, err := h.strg.Product().Create(&productCreate)
	if err != nil {
		h.handlerResponse(w, "error while storage product create:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	resp, err := h.strg.Product().GetByID(&models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(w, "error while storage product get by id:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Success", http.StatusOK, resp)
}

func (h *handler) GetByIdProduct(w http.ResponseWriter, r *http.Request) {

	var id string = r.URL.Query().Get("id")

	resp, err := h.strg.Product().GetByID(&models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(w, "error while storage category get by id:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Success", http.StatusOK, resp)
}

func (h *handler) GetListProduct(w http.ResponseWriter, r *http.Request) {

	var (
		offsetStr = r.URL.Query().Get("offset")
		limitStr  = r.URL.Query().Get("limit")
		search    = r.URL.Query().Get("search")
	)

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		h.handlerResponse(w, "error while offset: "+err.Error(), http.StatusBadRequest, nil)
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		h.handlerResponse(w, "error while limit: "+err.Error(), http.StatusBadRequest, nil)
		return
	}

	resp, err := h.strg.Product().GetList(&models.ProductGetListRequest{
		Offset: offset,
		Limit:  limit,
		Search: search,
	})
	if err != nil {
		h.handlerResponse(w, "error while storage category get list:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Success", http.StatusOK, resp)
}

func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	var UpdateProduct models.ProductUpdate

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.handlerResponse(w, "Error while updating prdocut"+err.Error(), http.StatusBadRequest, nil)
	}

	err = json.Unmarshal(body, &UpdateProduct)
	if err != nil {
		h.handlerResponse(w, "Error while updating product"+err.Error(), http.StatusBadRequest, nil)
	}

	_, err = h.strg.Product().Update(&UpdateProduct)
	if err != nil {
		h.handlerResponse(w, "Error while updating product"+err.Error(), http.StatusBadRequest, nil)
	}

	product, err := h.strg.Product().GetByID(&models.ProductPrimaryKey{Id: UpdateProduct.Id})
	if err != nil {
		h.handlerResponse(w, "Error while getbyid product"+err.Error(), http.StatusBadRequest, nil)
	}

	resp, err := json.Marshal(product)
	if err != nil {
		h.handlerResponse(w, "Error while marshal product"+err.Error(), http.StatusBadRequest, nil)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	var id string = r.URL.Query().Get("id")
	log.Println(id)

	err := h.strg.Product().Delete(&models.ProductPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(w, "Error while deleting product"+err.Error(), http.StatusBadRequest, nil)
	}
	w.WriteHeader(http.StatusOK)
}
