package api

import (
	"context"
	"errors"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	ID       string `json:"id" firestore:"id"`
	Email    string `json:"email" firestore:"email"`
	Password string `json:"password" firestore:"password"`
	Username string `json:"username" firestore:"username"`
}

type AuthHandler struct {
	DB   *firestore.Client
	Auth *auth.Client
}

type LoginReq struct {
	*User
	ID string `json:"-"`
}

func (lr *LoginReq) Bind(r *http.Request) error {
	if lr.Email == "" || lr.Password == "" {
		return errors.New("missing required fields")
	}
	return nil
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	data := &LoginReq{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, InvalidRequest(err))
		return
	}

	var user User
	ctx := context.Background()

	log.Printf("received login request from user with email: %v", data.Email)

	doc, err := ah.DB.Collection("users").Where("email", "==", data.Email).Documents(ctx).Next()
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Fatalf("no user with email found: %v", err)
		} else {
			log.Fatalf("failed to retrieve user with email: %v from database: %v", data.Email, err)
		}
		render.Render(w, r, InternalServerError(err))
		return
	}
	err = doc.DataTo(&user)
	if err != nil {
		log.Fatalf("no user with email found: %v", err)
	}

	log.Printf("retrieved user data with email: %v", user.Email)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		render.Render(w, r, Forbidden(err))
		return
	}

	customToken, err := ah.Auth.CustomToken(context.Background(), user.ID)
	if err != nil {
		log.Printf("failed to generate custom token: %v", err)
		render.Render(w, r, InternalServerError(err))
		return
	}

	log.Printf("created firebase token for user with email: %v", data.Email)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"token": customToken,
	})
}

type RegisterReq struct {
	*User
	ID string `json:"-"`
}

func (rr *RegisterReq) Bind(r *http.Request) error {
	if rr.Email == "" || rr.Password == "" {
		return errors.New("missing required fields")
	}
	return nil
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	data := &RegisterReq{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, InvalidRequest(err))
		return
	}

	ctx := context.Background()

	iter := ah.DB.Collection("users").Where("email", "==", data.Email).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("failed to retrieve users with email: %v with error: %v", data.Email, err)
			render.Render(w, r, InternalServerError(err))
			return
		}
		if existing_email, _ := doc.DataAt("email"); existing_email == data.Email {
			// todo: what's the proper server-side response for `account already exists` scenarios?
			log.Fatalf("user with email: %v already exists", data.Email)
			render.Render(w, r, InternalServerError(errors.New("account already exists")))
			return
		}
	}

	// todo: check if email and password are valid email and password values
	uid := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("failed to hash password: %v", err)
		render.Render(w, r, InternalServerError(err))
		return
	}

	_, _, err = ah.DB.Collection("users").Add(ctx, map[string]interface{}{
		"id":       uid,
		"email":    data.Email,
		"password": string(hashedPassword),
	})
	if err != nil {
		log.Printf("failed to insert user into database: %v", err)
		render.Render(w, r, InternalServerError(err))
		return
	}

	customToken, err := ah.Auth.CustomToken(context.Background(), uid)
	if err != nil {
		log.Printf("failed to create custom token for user: %v", err)
		render.Render(w, r, InternalServerError(err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"token": customToken,
	})
}
