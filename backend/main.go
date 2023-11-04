package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("./serviceAccount.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	fireauth, err := app.Auth(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	firestore, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer firestore.Close()

	ah := &AuthHandler{
		DB:   firestore,
		Auth: fireauth,
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.CleanPath)
	r.Use(middleware.Timeout(time.Second * 30)) // set timeout of res to 30s
	r.Use(httprate.LimitByIP(100, 1*time.Minute))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("healthy")) // todo: should perform basic platform checks like db connection etc.
	})

	r.Post("/login", ah.Login)
	r.Post("/register", ah.Register)

	http.ListenAndServe(":8000", r)
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

	log.Printf("retrieved user: %v, %v, %v", user.Email, user.ID, user.Password)

	d := doc.Data()
	log.Printf("retrieved user: %v", d)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		render.Render(w, r, Forbidden(err))
		return
	}

	token, err := ah.Auth.CustomToken(context.Background(), user.ID)
	if err != nil {
		log.Printf("failed to generate custom token: %v", err)
		render.Render(w, r, InternalServerError(err))
		return
	}

	log.Printf("created token: %v for user with email: %v", token, data.Email)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]interface{}{
		"token": token,
	})
}

type User struct {
	ID       string `json:"id" firestore:"id"`
	Email    string `json:"email" firestore:"email"`
	Password string `json:"password" firestore:"password"`
}

type AuthHandler struct {
	DB   *firestore.Client
	Auth *auth.Client
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

	// checks if user already exists in db or not
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
		"password": hashedPassword,
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

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func InvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

func InternalServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 500,
		StatusText:     "Internal server error.",
		ErrorText:      err.Error(),
	}
}

func Forbidden(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 403,
		StatusText:     "Forbidden.",
		ErrorText:      err.Error(),
	}
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}
