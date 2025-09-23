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
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/habibmrizki/back-end-tickitz/internal/configs"
	"github.com/habibmrizki/back-end-tickitz/internal/models"
	"github.com/habibmrizki/back-end-tickitz/internal/repositories"
	"github.com/habibmrizki/back-end-tickitz/pkg"
	"github.com/redis/go-redis/v9"
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
// @Param                   user body models.UsersStruct true "Input email and password"2
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
	if err := ctx.ShouldBind(&newDataUser); err != nil {
		log.Println("[ERROR] :", err.Error())

		// error jika format email salah
		if strings.Contains(err.Error(), "Error:Field validation for 'Email'") {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  "salah",
				Message: "incorrect email format",
			})
			return
		}

		// error jika panjang karakter password kurang dari 8 karakter
		if strings.Contains(err.Error(), "Error:Field validation for 'Password'") {
			ctx.JSON(http.StatusBadRequest, models.Response{
				Status:  "salah",
				Message: "password length must be at least 8 characters",
			})
			return
		}

		// 		// error yang lainnya
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
	ctx.JSON(http.StatusCreated, models.ResponseProfile{
		Status:  "berhasil",
		Message: "successfully create an account",
		Data:    gin.H{"email": newDataUser.Email},
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

	if err := ctx.ShouldBind(&userData); err != nil {
		log.Println("[ERROR] :", err.Error())
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid data sent",
		})
		return
	}
	log.Printf("[DEBUG] Email: %s, Password: %s\n", userData.Email, userData.Password)

	// Ambil data user dari database
	user, err := u.usersRepo.GetUserByEmail(ctx.Request.Context(), userData.Email)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid email or password",
		})
		return
	}

	// Verifikasi password
	hash := pkg.NewHashConfig()
	// Ubah nama fungsi di sini
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
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "failed to generate token"})
		return
	}

	// Jika berhasil, kirim respons sukses
	ctx.JSON(http.StatusOK, models.ResponseProfile{
		Status:  "berhasil",
		Message: "login successful",
		Data:    gin.H{"id": user.Id, "email": user.Email, "token": token}, // Add the token here
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
	claims, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Status:  "gagal",
			Message: "unauthorized",
		})
		return
	}

	userClaims, ok := claims.(pkg.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Status:  "gagal",
			Message: "invalid claims type",
		})
		return
	}

	profile, err := u.usersRepo.GetProfileById(ctx.Request.Context(), userClaims.UserId)
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

// UpdateProfile godoc
// @Summary      Update user profile
// @Description  Update user profile data including optional profile image and password
// @Tags         Users
// @Accept       multipart/form-data
// @Produce      json
// @Param        id            path      int     true   "User ID"
// @Param        fullname      formData  string  false  "Full name"
// @Param        email         formData  string  false  "Email"
// @Param        password      formData  string  false  "New password"
// @Param        old_password  formData  string  false  "Old password (required if changing password)"
// @Param        profile_image formData  file    false  "Profile image file"
// @Success      200 {object} models.Response
// @Failure      400 {object} models.Response
// @Failure      404 {object} models.Response
// @Failure      500 {object} models.Response
// @Router       /users/{id} [patch]
// @Security     ApiKeyAuth
func (u *UsersHandler) UpdateProfile(ctx *gin.Context) {
	// Ambil claims
	claims, isExist := ctx.Get("claims")
	if !isExist {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Status:  "gagal",
			Message: "unauthorized",
		})
		return
	}

	user, ok := claims.(pkg.Claims)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, models.Response{
			Status:  "gagal",
			Message: "invalid claims",
		})
		return
	}

	id := user.UserId
	log.Println("JWT Claims User ID:", id)

	req := models.UpdateProfileRequest{}
	if err := ctx.ShouldBindWith(&req, binding.FormMultipart); err != nil {
		log.Println("[ERROR] :", err.Error())
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "invalid form data: " + err.Error(),
		})
		return
	}

	if req.Email != nil {
		ctx.JSON(http.StatusBadRequest, models.Response{
			Status:  "gagal",
			Message: "email cannot be updated",
		})
		return
	}

	var profileImagePath string
	if req.ProfileImage != nil {
		ext := filepath.Ext(req.ProfileImage.Filename)
		filename := fmt.Sprintf("profile_%d_%d%s", id, time.Now().UnixNano(), ext)
		location := filepath.Join("public", "images", filename)

		if err := os.MkdirAll(filepath.Dir(location), 0755); err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Status:  "gagal",
				Message: "failed to create directory",
			})
			return
		}

		if err := ctx.SaveUploadedFile(req.ProfileImage, location); err != nil {
			ctx.JSON(http.StatusInternalServerError, models.Response{
				Status:  "gagal",
				Message: "failed to upload profile image",
			})
			return
		}
		profileImagePath = "/img/" + filename
	}

	// PERBAIKAN: Panggil repo dan tangkap data yang dikembalikan
	updatedUser, err := u.usersRepo.UpdateProfile(ctx.Request.Context(), id, &req, profileImagePath)
	if err != nil {
		log.Println(err.Error())
		if err.Error() == "user not found" {
			ctx.JSON(http.StatusNotFound, models.Response{
				Status:  "gagal",
				Message: err.Error(),
			})
			return
		}
		if err.Error() == "invalid old password" {
			ctx.JSON(http.StatusUnauthorized, models.Response{
				Status:  "gagal",
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusInternalServerError, models.Response{
			Status:  "gagal",
			Message: "failed to update profile",
		})
		return
	}

	// PERBAIKAN: Kembalikan data profil yang diperbarui di dalam respons
	ctx.JSON(http.StatusOK, models.ResponseUpdateProfile{
		Status:  "berhasil",
		Message: "successfully update profile",
		Data:    updatedUser,
	})
}

// func (u *UsersHandler) UpdateProfile(ctx *gin.Context) {
// 	// Ambil user id dari JWT (middleware sudah set "claims" di context)
// 	claims, exists := ctx.Get("claims")
// 	if !exists {
// 		ctx.JSON(http.StatusUnauthorized, models.Response{
// 			Status:  "gagal",
// 			Message: "unauthorized",
// 		})
// 		return
// 	}

// 	userClaims, ok := claims.(pkg.Claims)
// 	if !ok {
// 		ctx.JSON(http.StatusUnauthorized, models.Response{
// 			Status:  "gagal",
// 			Message: "invalid claims type",
// 		})
// 		return
// 	}

// 	id := userClaims.UserId

// 	// Debug: cek semua form value
// 	//TABMAHAN
// 	if form, err := ctx.MultipartForm(); err == nil {
// 		fmt.Printf("Form Values: %+v\n", form.Value)
// 		fmt.Printf("Form Files: %+v\n", form.File)
// 	}

// 	// Bind request body
// 	var req models.UpdateProfileRequest
// 	if err := ctx.ShouldBind(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "invalid form data: " + err.Error(),
// 		})
// 		return
// 	}

// 	// Handle file upload
// 	var profileImagePath string
// 	if req.ProfileImage != nil {
// 		ext := filepath.Ext(req.ProfileImage.Filename)
// 		filename := fmt.Sprintf("profile_%d_%d%s", id, time.Now().UnixNano(), ext)
// 		location := filepath.Join("public", "images", filename)

// 		if err := os.MkdirAll(filepath.Dir(location), 0755); err != nil {
// 			ctx.JSON(http.StatusInternalServerError, models.Response{
// 				Status:  "gagal",
// 				Message: "failed to create directory",
// 			})
// 			return
// 		}

// 		if err := ctx.SaveUploadedFile(req.ProfileImage, location); err != nil {
// 			ctx.JSON(http.StatusInternalServerError, models.Response{
// 				Status:  "gagal",
// 				Message: "failed to upload profile image",
// 			})
// 			return
// 		}

// 		profileImagePath = "/public/images/" + filename
// 	}

// 	// Panggil repository update
// 	err := u.usersRepo.UpdateProfile(ctx.Request.Context(), id, &req, profileImagePath)
// 	if err != nil {
// 		if err.Error() == "user not found" {
// 			ctx.JSON(http.StatusNotFound, models.Response{
// 				Status:  "gagal",
// 				Message: "user not found",
// 			})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: err.Error(),
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, models.Response{
// 		Status:  "berhasil",
// 		Message: "profile updated successfully",
// 	})
// }

// Logout godoc
// @Summary      Logout user
// @Description  Logout user by blacklisting the current JWT token in Redis
// @Tags         Users
// @Produce      json
// @Success      200 {object} map[string]interface{} "Logout successful"
// @Failure      401 {object} map[string]interface{} "Authorization header missing or malformed"
// @Failure      500 {object} map[string]interface{} "Logout failed due to server error"
// @Router       /logout [post]
// @Security     BearerAuth
func LogoutHandler(rdb *redis.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bearerToken := ctx.GetHeader("Authorization")
		if bearerToken == "" || !strings.HasPrefix(bearerToken, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "Authorization header is missing or malformed",
			})
			return
		}
		token := strings.Split(bearerToken, " ")[1]
		err := configs.BlacklistToken(rdb, token, time.Hour)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Logout gagal",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Logout berhasil",
		})
	}
}

// func (u *UsersHandler) UpdateProfile(ctx *gin.Context) {
// 	paramId := ctx.Param("id")
// 	id, err := strconv.Atoi(paramId)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "invalid user ID",
// 		})
// 		return
// 	}

// 	var req models.UpdateProfileRequest
// 	if err := ctx.ShouldBind(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, models.Response{
// 			Status:  "gagal",
// 			Message: "invalid form data: " + err.Error(),
// 		})
// 		return
// 	}

// 	// handle file upload
// 	var profileImagePath string
// 	if req.ProfileImage != nil {
// 		ext := filepath.Ext(req.ProfileImage.Filename)
// 		filename := fmt.Sprintf("profile_%d_%d%s", id, time.Now().UnixNano(), ext)
// 		location := filepath.Join("public", "images", filename)

// 		if err := os.MkdirAll(filepath.Dir(location), 0755); err != nil {
// 			ctx.JSON(http.StatusInternalServerError, models.Response{
// 				Status:  "gagal",
// 				Message: "failed to create directory",
// 			})
// 			return
// 		}

// 		if err := ctx.SaveUploadedFile(req.ProfileImage, location); err != nil {
// 			ctx.JSON(http.StatusInternalServerError, models.Response{
// 				Status:  "gagal",
// 				Message: "failed to upload profile image",
// 			})
// 			return
// 		}
// 		profileImagePath = "/public/images/" + filename
// 	}

// 	// Panggil repo pakai pointer
// 	err = u.usersRepo.UpdateProfile(ctx.Request.Context(), id, &req, profileImagePath)
// 	if err != nil {
// 		if err.Error() == "user not found" {
// 			ctx.JSON(http.StatusNotFound, models.Response{
// 				Status:  "gagal",
// 				Message: "user not found",
// 			})
// 			return
// 		}
// 		ctx.JSON(http.StatusInternalServerError, models.Response{
// 			Status:  "gagal",
// 			Message: err.Error(),
// 		})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, models.Response{
// 		Status:  "berhasil",
// 		Message: "profile updated successfully",
// 	})
// }
