package models

import "mime/multipart"

// import "mime/multipart"

// //	type ProfileStruct struct {
// //		UserId       int     `json:"id,omitempty" form:"id,omitempty" db:"user_id"`
// //		Firstname    *string `json:"firstname" form:"firstname" db:"firstname"`
// //		Email        string  `json:"email" form:"email" db:"email"`
// //		Lastname     *string `json:"lastname" form:"lastname" db:"lastname"`
// //		PhoneNumber  *string `json:"phonenumber,omitempty" form:"phonenumber" db:"phone_number"`
// //		Point        *int    `json:"point,omitempty" form:"point,omitempty" db:"point"`
// //		IdMember     *int    `json:"idmember,omitempty" form:"idmember,omitempty" db:"id_member"`
// //		ProfileImage *string `json:"profileimage,omitempty" form:"profileimage,omitempty" db:"profile_image"`
// //	}
// type Profile struct {
// 	User_Id      int    `db:"users_id" json:"users_id"`
// 	Firstname    string `db:"first_name" json:"first_name"`
// 	Lastname     string `db:"last_name" json:"last_name"`
// 	PhoneNumber  string `db:"phone_number" json:"phone_number"`
// 	ProfileImage string `db:"Profile_image" json:"Profile_image"`
// 	Point        int    `db:"point" json:"point"`
// }

// type ProfileBody struct {
// 	Profile
// 	ProfileIMage *multipart.FileHeader
// }

type Profile struct {
	User_Id      int    `db:"users_id" json:"users_id"`
	Firstname    string `db:"first_name" json:"first_name"`
	Lastname     string `db:"last_name" json:"last_name"`
	Email        string `form:"email" json:"email"`
	PhoneNumber  string `db:"phone_number" json:"phone_number"`
	ProfileImage string `db:"profile_image" json:"profile_image"`
	Password     string `json:"password,omitempty"`
	Point        int    `db:"point" json:"point"`
}

type ProfileBody struct {
	Profile
	ProfileImage *multipart.FileHeader `form:"profile_image"`
}

type UpdateProfileRequest struct {
	Firstname    *string               `form:"first_name" json:"first_name"`
	Lastname     *string               `form:"last_name" json:"last_name"`
	Email        *string               `form:"email" json:"email"`
	PhoneNumber  *string               `form:"phone_number" json:"phone_number"`
	OldPassword  *string               `form:"old_password" json:"old_password"`
	Password     *string               `form:"password" json:"password"`
	ProfileImage *multipart.FileHeader `form:"profile_image" json:"profile_image"`
}

type ResponseUpdateProfile struct {
	Message string
	Status  string
	Data    interface{} `json:"data,omitempty"`
}
