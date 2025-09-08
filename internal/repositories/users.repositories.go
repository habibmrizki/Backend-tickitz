package repositories

import (
	"context"
	"errors"
	"log"
	"time" // Tambahkan import ini

	"github.com/habibmrizki/back-end-tickitz/internal/models"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	// db ini digunakan untuk menyimpan koneksi database
	// *pgxpool.Pool digunakan untuk mengelola banyak koneksi ke database Postgresql secara bersamaan
	db *pgxpool.Pool
}

// function cunstructor memmbuat dan menganalisa sebuah instance baru dari UserRepository
func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// Repository add user
// pgconn.CommandTage digunakan untuk memberi status dari opearasi database yang dilakukan biasan jumlha bars baru byang ditampilkan
func (u *UserRepository) UserRegister(ctx context.Context, email string, password string, role string) (pgconn.CommandTag, error) {
	// Begin digunakann untuk memulai sebuah transaksi database
	tx, err := u.db.Begin(ctx)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return pgconn.CommandTag{}, err
	}

	// Rollback digunakan untuk membatalkan transaksi jika ada error di tengah jalan
	// jika transaksi berhasil dan dicommit , makan perintah ini tidak berpengaruh
	defer tx.Rollback(ctx)

	// ambil waktu saat ini
	currentTime := time.Now()

	// menambahkan user baru dengan mengembalikan id user baru
	queryUser := "INSERT INTO users (email, password, role, created_at, update_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	var userID int
	// queryRow digunakan menjalanka query dan mengembalikan satu baris hasil
	if err := tx.QueryRow(ctx, queryUser, email, password, role, currentTime, currentTime).Scan(&userID); err != nil {
		return pgconn.CommandTag{}, err
	}

	// default value
	first_name := ""
	last_name := ""
	phone_number := ""
	photo_path := ""
	point := 0

	// menambahkan baris baru untuk data profile user baru namun dengan default value
	queryProfile := "INSERT INTO profile (users_id, first_name, last_name, phone_number, profile_image, point) VALUES ($1, $2, $3, $4, $5, $6)"
	// Exec menjalankan query yang tidak mengembalikan baris data dan mengembalikan informasi tentang hasil eksekusi
	cmd, err := tx.Exec(ctx, queryProfile, userID, first_name, last_name, phone_number, photo_path, point)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return pgconn.CommandTag{}, err
	}
	// Commit menyimpan semua perubahan secara permanen ke database
	return cmd, tx.Commit(ctx)
}

func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.UsersStruct, error) {
	sql := "SELECT id, email, password, role FROM users WHERE email = $1"

	var user models.UsersStruct
	// Scan digunakan untuk membaca hasil dari QueryRow yang hasilnya akan di masukkan ke dalam variable
	if err := u.db.QueryRow(ctx, sql, email).Scan(&user.Id, &user.Email, &user.Password, &user.Role); err != nil {
		// ErrNoRows error spesifik ketika tidak ditemukan baris yang cocok
		if err == pgx.ErrNoRows {
			return models.UsersStruct{}, errors.New("user not found")
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		return models.UsersStruct{}, err
	}
	return user, nil
}

// func (u *UserRepository) GetProfileById(ctx context.Context, idInt int) ([]models.Profile, error) {

// 	// jika tidak terdapat data di redis maka jalankan query GET profile berikut ini
// 	query := "SELECT p.user_id, p.first_name, p.last_name, p.phone_number, p.photo_path, p.title, p.point, u.email FROM profiles p join users u on u.id = p.user_id WHERE user_id = $1"
// 	rows, err := u.db.Query(ctx, query, idInt)
// 	if err != nil {
// 		return nil, err
// 	}

// 	defer rows.Close()
// 	var result []models.Profile

// 	for rows.Next() {
// 		var profile models.Profile
// 		if err := rows.Scan(&profile.User_Id, &profile.Firstname, &profile.Lastname, &profile.Email, &profile.PhoneNumber, &profile.ProfileImage, &profile.Point); err != nil {
// 			return []models.Profile{}, err
// 		}
// 		result = append(result, profile)
// 	}

// 	return result, nil
// }

func (u *UserRepository) GetProfileById(ctx context.Context, idInt int) (models.Profile, error) {

	// Ubah p.user_id menjadi p.users_id di kueri SQL
	query := "SELECT p.users_id, p.first_name, p.last_name, p.phone_number, p.profile_image, p.point FROM profile p JOIN users u ON u.id = p.users_id WHERE p.users_id = $1"

	var profile models.Profile
	// mengambil data dari hasil kueri database dan memasukannya ke dalam variabel Go.
	err := u.db.QueryRow(ctx, query, idInt).Scan(
		&profile.User_Id,
		&profile.Firstname,
		&profile.Lastname,
		&profile.PhoneNumber,
		&profile.ProfileImage,
		&profile.Point,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models.Profile{}, errors.New("profile not found")
		}
		log.Println("[ERROR] : ", err.Error())
		return models.Profile{}, err
	}

	return profile, nil
}

func (u *UserRepository) UpdateProfile(ctx context.Context, id int, firstName string, lastName string, phoneNumber string, profileImage string) error {
	tx, err := u.db.Begin(ctx)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return err
	}

	defer tx.Rollback(ctx)

	// Periksa apakah pengguna dengan ID tersebut ada di tabel 'users'
	var usersID int
	err = tx.QueryRow(ctx, "SELECT id FROM users WHERE id = $1", id).Scan(&usersID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("user not found")
		}
		log.Println("[ERROR] : ", err.Error())
		return err
	}

	// Lakukan pembaruan di tabel 'profile'
	query := `UPDATE profile SET first_name = $1, last_name = $2, phone_number = $3, profile_image = $4 WHERE users_id = $5`

	// exec menjalanka perintah sebuah sql yang dimana memodifikasi data dan berjalan secara atomik (semua gagal atau semua berhasil)
	_, err = tx.Exec(ctx, query, firstName, lastName, phoneNumber, profileImage, id)
	if err != nil {
		log.Println("[ERROR] : ", err.Error())
		return err
	}

	return tx.Commit(ctx)
}
