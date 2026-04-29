package repository

import (
	"context"
	"fmt"
	"user_management_system/database"
	"user_management_system/models"

	"github.com/google/uuid"
)

type PaginatedUsers struct {
	TotalCount int64         `json:"total_count"`
	UserItems  []models.User `json:"user_items"`
}

//create the user
func Create(user *models.User) error {
	return database.DB.Create(user).Error
}

//sql command find email
func FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

//sql query find user id
func FindByUserID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := database.DB.WithContext(ctx).Where("id = ?", id).First(&user).Error
	return &user, err
}

//show all user
func GetAll() ([]models.User, error) {
	var users []models.User
	err := database.DB.Find(&users).Error
	return users, err
}

//1 page asign 10 data
func GetAllPaginated(ctx context.Context, page, recordPerPage int) (*PaginatedUsers, error) {
	var users []models.User
	var totalCount int64

	offset := (page - 1) * recordPerPage

	//count give me total pages info
	if err := database.DB.WithContext(ctx).Model(&models.User{}).Count(&totalCount).Error; err != nil {
		return nil, err
	}

	//give me at most 10 records from db
	if err := database.DB.WithContext(ctx).Limit(recordPerPage).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	return &PaginatedUsers{
		TotalCount: totalCount,
		UserItems:  users,
	}, nil
}

func DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	result := database.DB.WithContext(ctx).Delete(&models.User{}, "id = ?", userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}