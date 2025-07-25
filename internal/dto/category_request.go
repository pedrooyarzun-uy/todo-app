package dto

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type DeleteCategoryRequest struct {
	Id int `json:"id" binding:"required"`
}

type ModifyCategoryRequest struct {
	Id   int    `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
}
