package controllers

import (
	"context"
	"dev-boiler/models"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type UserController struct {
	db        *gorm.DB
	createLog *mongo.Collection
	updateLog *mongo.Collection
	deleteLog *mongo.Collection
}

func NewUserController(db *gorm.DB, mongo *mongo.Client) (UserController, error) {
	user := models.User{}
	db.AutoMigrate(&user)
	createLog := mongo.Database("test").Collection("LogUserCreate")
	updateLog := mongo.Database("test").Collection("LogUserUpdate")
	deleteLog := mongo.Database("test").Collection("LogUserDelete")
	return UserController{db, createLog, updateLog, deleteLog}, nil
}

func saveFile(file *multipart.FileHeader) (bool, string) {
	src, err := file.Open()
	if err != nil {
		return false, "fail to open file"
	}
	defer src.Close()

	filepath := "statics/assets/" + file.Filename
	dst, err := os.Create(filepath)
	if err != nil {
		return false, "fail to create file path"
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return false, "fail to copy file path"
	}

	return true, filepath
}

func (u UserController) FindUsers(ctx echo.Context) error {
	users := []models.User{}
	u.db.Find(&users)
	return ctx.JSON(http.StatusOK, users)
}

func (u UserController) FindUser(ctx echo.Context) error {
	user := models.User{}
	if err := u.db.Find(&user).Error; err != nil {
		return ctx.String(http.StatusBadRequest, "Not found")
	}

	return ctx.JSON(http.StatusOK, user)
}

func (u UserController) CreateUser(ctx echo.Context) error {
	user := models.User{}
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	success, detail := saveFile(file)
	if !success {
		return ctx.String(http.StatusInternalServerError, detail)
	}

	user.ImageURL = detail
	if err := u.db.Create(&user).Error; err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	u.createLog.InsertOne(context.Background(), &user)
	return ctx.JSON(http.StatusOK, user)
}

func (u UserController) UpdateUser(ctx echo.Context) error {
	id := ctx.Param("id")
	user := models.User{}

	// find user
	result := u.db.Find(&user, id)
	if result.RowsAffected == 0 {
		return ctx.String(http.StatusNotFound, "Not found")
	}

	// check
	if err := ctx.Bind(&user); err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	u.updateLog.InsertOne(context.Background(), &user)
	u.db.Save(&user)
	return ctx.JSON(http.StatusOK, user)
}

func (u UserController) DeleteUser(ctx echo.Context) error {
	id := ctx.Param("id")
	user := models.User{}

	// find user
	result := u.db.Find(&user, id)
	if result.RowsAffected == 0 {
		return ctx.String(http.StatusNotFound, "Not found")
	}

	// delete user
	u.deleteLog.InsertOne(context.Background(), &user)
	u.db.Delete(&user, id)
	return ctx.NoContent(http.StatusOK)
}
