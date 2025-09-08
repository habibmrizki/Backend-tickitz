package models

type UsersStruct struct {
	Id       int    `db:"id,omitempty" json:"id,omitempty"`
	Email    string `db:"email" json:"email" form:"email" binding:"required,email"`
	Password string `db:"password" json:"password" form:"password"`
	Role     string `db:"role,omitempty" json:"role,omitempty"`
}

type Profile struct {
	User_Id      int    `db:"users_id" json:"users_id"`
	Firstname    string `db:"first_name" json:"first_name"`
	Lastname     string `db:"last_name" json:"last_name"`
	PhoneNumber  string `db:"phone_number" json:"phone_number"`
	ProfileImage string `db:"Profile_image" json:"Profile_image"`
	Point        int    `db:"point" json:"point"`
}
