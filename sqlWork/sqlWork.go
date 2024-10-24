package sqlWork

import (
	"uuu/lib"

	"gorm.io/gorm"
)

/*
// Get the first record ordered by primary key
db.First(&user)
// SELECT * FROM users ORDER BY id LIMIT 1;

// Get one record, no specified order
db.Take(&user)
// SELECT * FROM users LIMIT 1;

// Get last record, ordered by primary key desc
db.Last(&user)
// SELECT * FROM users ORDER BY id DESC LIMIT 1;
*/

func GetUserByEmail(db *gorm.DB, name string) ([]lib.User, error) {

	us := make([]lib.User, 0)
	err := db.Where("email LIKE ?", "%"+name+"%").Find(&us).Error
	return us, err
}

func GetUserByName(db *gorm.DB, name string) (lib.User, error) {
	var user lib.User
	err := db.Where("name = ?", name).Take(&user).Error
	return user, err
}

func GetUserNew(db *gorm.DB, userId int) (lib.User, error) {
	var user lib.User
	err := db.Take(&user, userId).Error
	return user, err
}

func GetUserListNew(db *gorm.DB) ([]lib.User, error) {
	us := make([]lib.User, 0)
	// err:= db.Preload("Images").Find(&us).Error
	err := db.Find(&us).Error
	return us, err
}

func CreateUserNew(db *gorm.DB, user lib.User) (uint, error) {
	err := db.Create(&user).Error
	return user.Id, err
}

func RemoveUserNew(db *gorm.DB, id uint) (int, error) {
	command := db.Delete(&lib.User{}, id)
	return int(command.RowsAffected), command.Error
}
