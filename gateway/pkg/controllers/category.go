package controllers

import (
	"context"
	pb "github.com/JMURv/e-commerce/api/pb/item"
	"github.com/JMURv/e-commerce/gateway/pkg/utils"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func ListCategory(w http.ResponseWriter, r *http.Request) {
	cli := pb.NewCategoryServiceClient(itemConn)
	categories, _ := cli.GetAllCategories(context.Background(), &pb.EmptyRequest{})

	utils.OkResponse(w, http.StatusOK, categories)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	cli := pb.NewCategoryServiceClient(itemConn)
	c, err := cli.GetCategoryByID(context.Background(), &pb.GetCategoryByIDRequest{CategoryId: categoryID})
	if err != nil {
		utils.ErrResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.OkResponse(w, http.StatusOK, c)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	newCategory := &pb.CreateCategoryRequest{}
	utils.ParseBody(r, newCategory)

	cli := pb.NewCategoryServiceClient(itemConn)
	c, err := cli.CreateCategory(context.Background(), newCategory)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.OkResponse(w, http.StatusCreated, c)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	newData := &pb.UpdateCategoryRequest{CategoryId: categoryID}
	utils.ParseBody(r, newData)

	cli := pb.NewCategoryServiceClient(itemConn)
	c, err := cli.UpdateCategory(context.Background(), newData)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	utils.OkResponse(w, http.StatusOK, c)
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 32)
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	cli := pb.NewCategoryServiceClient(itemConn)
	_, err = cli.DeleteCategory(context.Background(), &pb.DeleteCategoryRequest{CategoryId: categoryID})
	if err != nil {
		utils.ErrResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
