package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"app/models"
)

func (h *handler) Category(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":
		h.CreateCategory(w, r)
	case "GET":
		var (
			values = r.URL.Query()
			method = values.Get("method")
		)

		if method == "GET_LIST" {
			h.GetListCategory(w, r)
		} else if method == "GET" {
			h.GetByIdCategory(w, r)
		}
	case "PUT":
		h.UpdateCategory(w, r)
	case "DELETE":
		h.DeleteCategory(w, r)
	}

}

func (h *handler) CreateCategory(w http.ResponseWriter, r *http.Request) {

	var createCategory models.CreateCategory

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.handlerResponse(w, "error while read body: "+err.Error(), http.StatusBadRequest, nil)
		return
	}

	err = json.Unmarshal(body, &createCategory)
	if err != nil {
		h.handlerResponse(w, "error while unmarshal body: "+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	id, err := h.strg.Category().Create(&createCategory)
	if err != nil {
		h.handlerResponse(w, "error while storage category create:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	resp, err := h.strg.Category().GetByID(&models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(w, "error while storage category get by id:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Success", http.StatusOK, resp)
}

func (h *handler) GetByIdCategory(w http.ResponseWriter, r *http.Request) {

	var id string = r.URL.Query().Get("id")

	resp, err := h.strg.Category().GetByID(&models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(w, "error while storage category get by id:"+err.Error(), http.StatusInternalServerError, nil)
		return
	}

	h.handlerResponse(w, "Success", http.StatusOK, resp)
}

func (h *handler) GetListCategory(w http.ResponseWriter, r *http.Request) {

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

	resp, err := h.strg.Category().GetList(&models.CategoryGetListRequest{
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

func (h *handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	var UpdateCategory models.UpdateCategory

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.handlerResponse(w, "Error while updating category"+err.Error(), http.StatusBadRequest, nil)
	}

	err = json.Unmarshal(body, &UpdateCategory)
	if err != nil {
		h.handlerResponse(w, "Error while updating category"+err.Error(), http.StatusBadRequest, nil)
	}

	_, err = h.strg.Category().Update(&UpdateCategory)
	if err != nil {
		h.handlerResponse(w, "Error while updating category"+err.Error(), http.StatusBadRequest, nil)
	}

	category, err := h.strg.Category().GetByID(&models.CategoryPrimaryKey{Id: UpdateCategory.Id})
	if err != nil {
		h.handlerResponse(w, "Error while getbyid category"+err.Error(), http.StatusBadRequest, nil)
	}

	resp, err := json.Marshal(category)
	if err != nil {
		h.handlerResponse(w, "Error while marshal category"+err.Error(), http.StatusBadRequest, nil)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (h *handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	var id string = r.URL.Query().Get("id")
	log.Println(id)

	err := h.strg.Category().Delete(&models.CategoryPrimaryKey{Id: id})
	if err != nil {
		h.handlerResponse(w, "Error while deleting category"+err.Error(), http.StatusBadRequest, nil)
	}
	w.WriteHeader(http.StatusOK)
}
