// package handlers

// import (
// 	"errors"
// 	"log"
// 	"net/http"
// 	"strings"

// 	"github.com/gin-gonic/gin"
// 	"github.com/habibmrizki/back-end-tickitz/internal/models"
// 	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
// 	"github.com/habibmrizki/back-end-tickitz/pkg"
// )

// // Register
// // @summary                 Register user
// // @router                  /users/register [post]
// // @Description             Register with email and password for access application
// // @Tags                    Users
// // @Param                   user body models.UsersStruct true "Input email and password"
// // @accept                  json
// // @produce                 json
// // @failure                 500 {object} models.Response
// // @failure                 400 {object} models.Response
// // @failure                 409 {object} models.Response
// // @success                 201 {object} models.Response
// type UsersHandler struct {
// 	usersRepo *repositories.UserRepository
// }

// func NewUsersHandlers(usersRepo *repositories.UserRepository) *UsersHandler {
// 	return &UsersHandler{usersRepo: usersRepo}
// }

// func (u *UsersHandler) UserRegister(ctx *gin.Context) {
// 	// deklarasi body dari input user
// 	newDataUser := models.UsersStruct{}

// 	// binding data
// 	// membaca request dari input user dari JSON sekaligus melakukan verifikasi, jika format json tidak sesuai dengan format yang ada didalam struct maka akan terjadi error
// 	if err := ctx.ShouldBindJSON(&newDataUser); err != nil {
// 		log.Println("[ERROR] :", err.Error())

// 		// error jika format email salah
// 		if strings.Contains(err.Error(), "Error:Field validation for 'Email'") {
// 			ctx.JSON(http.StatusBadRequest, models.Response{
// 				Status:  "salah",
// 				Message: "incorrect email format",
// 			})
// 			return
// 		}

// 		// error jika panjang karakter password kurang dari 8 karakter
// 		if strings.Contains(err.Error(), "Error:Field validation for 'Password'") {
// 			ctx.JSON(http.StatusBadRequest, models.Response{
// 				Status:  "salah",
// 				Message: "password length must be at least 8 characters",
// 			})
// 			return
// 		}

// 		// error yang lainnya
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "salah",
// 			Message: "invalid data sent",
// 		})
// 		return
// 	}

// 	// hash password, mengkonversi password normal menjadi bentuk lain yang sulit dibaca
// 	hash := pkg.NewHashConfig()
// 	// gunakan konfigurasi yang direkomendasikan
// 	hash.UseRecommended()
// 	hashedPass, err := hash.GenHash(newDataUser.Password)

// 	// error jika gagal mengkonversi password menjadi hash
// 	if err != nil {
// 		log.Println("[ERROR] : ", err.Error())
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "salah",
// 			Message: "hash failed",
// 		})
// 		return
// 	}

// 	// default value untuk role
// 	role := "user"

// 	// eksekusi fungsi repository register user
// 	cmd, err := u.usersRepo.UserRegister(ctx.Request.Context(), newDataUser.Email, hashedPass, role)

// 	// lakukan error handling jika terjadi kesalahan dalam menjalankan query / user sudah terdaftar sebelumnya
// 	if err != nil {
// 		log.Println("[ERROR] : ", err.Error())
// 		ctx.JSON(http.StatusConflict, models.Response{
// 			Status:  "salah",
// 			Message: "user already registered",
// 		})
// 		return
// 	}

// 	// cek apakah perintah berhasil manambahkan data di database
// 	if cmd.RowsAffected() == 0 {
// 		log.Println("[ERROR] : ", errors.New("query failed, did not change the data in the database"))
// 		log.Println("Query failed, did not change the data in the database")
// 		return
// 	}

//		// return jika server berhasil memberikan response
//		ctx.JSON(http.StatusCreated, models.Response{
//			Status:  "berhasil", // Perbaikan: Ganti "salah" menjadi "berhasil"
//			Message: "successfully create an account",
//		})
//	}
package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/back-end-tickitz/internal/models"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/habibmrizki/back-end-tickitz/pkg"
)

type UsersHandler struct {
	usersRepo *repositories.UserRepository
}

func NewUsersHandlers(usersRepo *repositories.UserRepository) *UsersHandler {
	return &UsersHandler{usersRepo: usersRepo}
}

// Register
// @summary                 Register user
// @router                  /users/register [post]
// @Description             Register with email and password for access application
// @Tags                    Users
// @Param                   user body models.UsersStruct true "Input email and password"
// @accept                  json
// @produce                 json
// @failure                 500 {object} models.Response
// @failure                 400 {object} models.Response
// @failure                 409 {object} models.Response
// @success                 201 {object} models.Response
func (u *UsersHandler) UserRegister(ctx *gin.Context) {
	// deklarasi body dari input user
	newDataUser := models.UsersStruct{}

	// binding data
	// membaca request dari input user dari JSON sekaligus melakukan verifikasi, jika format json tidak sesuai dengan format yang ada didalam struct maka akan terjadi error
	if err := ctx.ShouldBindJSON(&newDataUser); err != nil {
		log.Println("[ERROR] :", err.Error())

		// error jika format email salah
		// if strings.Contains(err.Error(), "Error:Field validation for 'Email'") {
		// 	ctx.JSON(http.StatusBadRequest, models.Response{
		// 		Status:  "salah",
		// 		Message: "incorrect email format",
		// 	})
		// 	return
		// }

		// error jika panjang karakter password kurang dari 8 karakter
		// if strings.Contains(err.Error(), "Error:Field validation for 'Password'") {
		// 	ctx.JSON(http.StatusBadRequest, models.Response{
		// 		Status:  "salah",
		// 		Message: "password length must be at least 8 characters",
		// 	})
		// 	return
		// }

		// error yang lainnya
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "salah",
			Message: "invalid data sent",
		})
		return
	}

	// / hash password, mengkonversi password normal menjadi bentuk lain yang sulit dibaca
	hash := pkg.NewHashConfig()
	hash.UseRecommended()
	hashedPass, err := hash.GenHash(newDataUser.Password)

	// error jika gagal mengkonversi password menjadi hash
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "salah",
			Message: "hash failed",
		})
		return
	}

	// default value untuk role
	role := "user"

	// eksekusi fungsi repository register user
	cmd, err := u.usersRepo.UserRegister(ctx.Request.Context(), newDataUser.Email, hashedPass, role)

	// lakukan error handling jika terjadi kesalahan dalam menjalankan query / user sudah terdaftar sebelumnya
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusConflict, models.Response{
			Status:  "salah",
			Message: "user already registered",
		})
		return
	}

	// cek apakah perintah berhasil manambahkan data di database
	if cmd.RowsAffected() == 0 {
		log.Println("[ERROR] : ", errors.New("query failed, did not change the data in the database"))
		return
	}

	// return jika server berhasil memberikan response
	ctx.JSON(http.StatusCreated, models.Response{
		Status:  "berhasil",
		Message: "successfully create an account",
	})
}

// UserLogin
// @summary                 Login user
// @router                  /users/login [post]
// @Description             Login with email and password
// @Tags                    Users
// @Param                   user body models.UsersStruct true "Input email and password"
// @accept                  json
// @produce                 json
// @failure                 400 {object} models.Response
// @failure                 401 {object} models.Response
// @failure                 500 {object} models.Response
// @success                 200 {object} models.Response
func (u *UsersHandler) UserLogin(ctx *gin.Context) {
	var userData models.UsersStruct

	if err := ctx.ShouldBindJSON(&userData); err != nil {
		log.Println("[ERROR] :", err.Error())
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid data sent",
		})
		return
	}

	// Ambil data user dari database
	user, err := u.usersRepo.GetUserByEmail(ctx.Request.Context(), userData.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Status:  "gagal",
			Message: "invalid email or password",
		})
		return
	}

	// Verifikasi password
	hash := pkg.NewHashConfig()
	if ok, _ := hash.CompareHashAndPassword(userData.Password, user.Password); !ok {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Status:  "gagal",
			Message: "invalid email or password",
		})
		return
	}

	claim := pkg.NewJWTClaims(user.Id, string(user.Role))

	token, err := claim.GenToken()
	if err != nil {
		log.Println(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token"})
		return
	}

	// Jika berhasil, kirim respons sukses
	ctx.JSON(http.StatusOK, models.ResponseProfile{
		Status:  "berhasil",
		Message: "login successful",
		Data:    gin.H{"id": user.Id, "token": token}, // Add the token here
	})
}

// GetProfileById
// @summary                 Get profile by ID
// @router                  /users/{id} [get]
// @Description             Get user profile by user ID
// @Tags                    Users
// @Param                   id path int true "User ID"
// @accept                  json
// @produce                 json
// @success                 200 {object} models.Response
// @failure                 400 {object} models.Response
// @failure                 404 {object} models.Response
func (u *UsersHandler) GetProfileById(ctx *gin.Context) {
	paramId := ctx.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid user ID",
		})
		return
	}

	profile, err := u.usersRepo.GetProfileById(ctx.Request.Context(), id)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusNotFound, models.Response{
			Status:  "gagal",
			Message: "profile not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.ResponseProfile{
		Status:  "berhasil",
		Message: "successfully get user profile",
		Data:    profile,
	})
}

// UpdateProfile
// @summary                 Update user profile
// @router                  /users/{id} [put]
// @Description             Update a user's profile information by ID
// @Tags                    Users
// @Param                   id path int true "User ID"
// @Param                   profile body models.Profile true "Updated profile information"
// @accept                  json
// @produce                 json
// @success                 200 {object} models.Response
// @failure                 400 {object} models.Response
// @failure                 404 {object} models.Response
// @failure                 500 {object} models.Response
func (u *UsersHandler) UpdateProfile(ctx *gin.Context) {
	// Ambil ID dari parameter URL
	paramId := ctx.Param("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid user ID",
		})
		return
	}

	// Binding data dari body request
	var updatedProfile models.Profile
	if err := ctx.ShouldBindJSON(&updatedProfile); err != nil {
		log.Println("[ERROR] : ", err.Error())
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid data sent",
		})
		return
	}

	// Panggil fungsi repository untuk memperbarui profil
	err = u.usersRepo.UpdateProfile(ctx.Request.Context(), id, updatedProfile.Firstname, updatedProfile.Lastname, updatedProfile.PhoneNumber, updatedProfile.ProfileImage)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, models.Response{
				Status:  "gagal",
				Message: "user or profile not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "failed to update profile",
		})
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		Status:  "berhasil",
		Message: "successfully updated user profile",
	})
}
