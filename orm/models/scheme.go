package orm

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName    string `gorm:"column:firstname;type:varchar;not null;size:100",json:"firstName"`
	LastName     string `gorm:"column:lastname;type:varchar;not null;size:100",json:"lastName"`
	Email        string `gorm:"column:email;type:text;not null;uniqueIndex",json:"email"`
	Password     string `gorm:"column:password;not null;type:text",json:"-"`
	ResetToken   string `gorm:"column:reset_token;type:text;comment:Used for resetting password",json:"-"`
	DateOfBirth  string `gorm:"column:date_of_birth;not null;type:varchar;size:11",json:"dateOfBirth"`
	ProfileImage string `gorm:"column:profile_image;not null;type:text",json:"profileImage"`
}

func (User) TableName() string {
	return "users"
}
