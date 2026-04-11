package repository

import (
	"context"
	"errors"
	"log"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/entity"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidUser  = errors.New("invalid user")
)

type UserRepositoryImp struct {
	DB *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &UserRepositoryImp{DB: db}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 11)
	return string(bytes), err
}

// Insert implements [UserRepository].
func (u *UserRepositoryImp) Insert(ctx context.Context, user entity.UserParam) (entity.User, error) {
	query := `
	insert into users (username, password)
	values ($1, $2)
	RETURNING id
	`

	newUser := entity.User{}
	hashPass, err := HashPassword(user.Password)
	if err != nil {
		return newUser, err
	}
	// newUser.Username = user.Username
	// newUser.Password = string(pass)

	var id int
	row := u.DB.QueryRow(ctx, query, user.Username, hashPass)
	err = row.Scan(&id)
	if err != nil {
		return newUser, err
	}

	newUser.Id = uint16(id)
	newUser.Username = user.Username
	newUser.Password = user.Password

	log.Println("new user :", newUser)

	return newUser, nil

	// panic("unimplemented")
}

// Login implements [UserRepository].
func (u *UserRepositoryImp) Login(ctx context.Context, user entity.UserParam) (entity.User, error) {
	query := `
		SELECT id, username, password FROM users 
		WHERE username = $1
	`
	userRes := entity.User{}
	// var currPass string
	// hashPass :=

	//chek is username available
	row, err := u.DB.Query(ctx, query, user.Username)
	if err != nil {
		return userRes, err
	}
	defer row.Close()

	if row.Next() {
		err = row.Scan(&userRes.Id, &userRes.Username, &userRes.Password)
		if err != nil {
			return userRes, err
		}
	} else {
		return userRes, ErrUserNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(userRes.Password), []byte(user.Password))
	if err != nil {
		return userRes, err
	}

	// log.Println("pass true")
	userRes.Password = ""

	return userRes, nil
}

// GetAll implements [UserRepository].
func (u *UserRepositoryImp) GetAll(ctx context.Context) ([]entity.User, error) {
	panic("unimplemented")
}

// Get implements [UserRepository].
func (u *UserRepositoryImp) Get(ctx context.Context, id uint16) (entity.User, error) {
	query := `
	SELECT id, username, password FROM users WHERE id = $1 AND pas
	`
	user := entity.User{}
	row, err := u.DB.Query(ctx, query, id)
	if err != nil {
		return user, nil
	}
	defer row.Close()

	if row.Next() {
		err := row.Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			return user, err
		}
		return user, nil
	} else {
		return user, ErrUserNotFound
	}
}

// Update implements [UserRepository].
func (u *UserRepositoryImp) Update(ctx context.Context, newUser entity.UserParam, id uint16) (entity.User, error) {
	query := `
	UPDATE users
	SET username = $1,
		password = $2,
	WHERE id = $3
	RETURNING *
	`

	user := entity.User{}
	row, err := u.DB.Query(ctx, query, newUser.Username, newUser.Password, id)
	if err != nil {
		return user, err
	}
	err = row.Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		return user, err
	}
	return user, nil

}

// Delete implements [UserRepository].
func (u *UserRepositoryImp) Delete(ctx context.Context, id uint16) error {
	query := `
	DELETE FROM users
	WHERE id = $1
	`
	cmdTag, err := u.DB.Exec(ctx, query, id)
	if err != nil {
		// fmt.Println("exec error")
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		// fmt.Println("task not fon")
		return ErrUserNotFound
	}

	return nil
}
