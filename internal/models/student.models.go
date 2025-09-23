package models

import "mime/multipart"

type Student struct {
	Id       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"nama_siswa"`
	Role     string `db:"role" json:"peran,omitempty"`
	Password string `db:"password" json:"password,omitempty"`
	Images   string `db:"images" json:"image"`
}

type StudentBody struct {
	Student
	Images *multipart.FileHeader `form:"image"`
}

// type StudentAuth struct {
// 	Name     string `json:"nama_siswa" binding:"required"`
// 	Password string `json:"password" binding:"required,min=4"`
// }
