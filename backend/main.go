package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

// 服务地址, casdoor提供
const (
	casdoorAuthBase    = "http://localhost:8000/login/oauth/authorize"
	casdoorTokenURL    = "http://localhost:8000/api/login/oauth/access_token"
	casdoorUserInfoURL = "http://localhost:8000/api/userinfo"
)

// 应用账号, casdoor提供
const (
	clientID     = "45f4c9041d1a76681629"
	clientSecret = "4dd4cbc7eadd36d4c93a14281a5a8edb6bdd2a4c"
)

// 前端应用提供
const callbackURL = "http://localhost:8080/#/callback" // 前端的一个中间页, 夹在login和home之间

func main() {
	http.HandleFunc("/login", corsMiddleware(login)) // optional 可以在前端做
	http.HandleFunc("/token", corsMiddleware(getToken))
	http.HandleFunc("/userinfo", corsMiddleware(getUserInfo)) // optional 可以在前端做

	log.Println("Starting server on :9000")
	if err := http.ListenAndServe(":9000", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// https://casdoor.org/docs/basic/core-concepts/#login-urls
func login(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(casdoorAuthBase)
	if err != nil {
		http.Error(w, "Failed to parse URL", http.StatusInternalServerError)
		log.Printf("Failed to parse URL: %v", err)
		return
	}

	q := url.Values{}
	q.Add("client_id", clientID)
	q.Add("response_type", "code")
	q.Add("redirect_uri", callbackURL)
	q.Add("scope", "profile") // https://casdoor.org/docs/how-to-connect/oauth/#available-scopes
	q.Add("state", "STATE")

	u.RawQuery = q.Encode()

	casdoorAuthURL := u.String()

	http.Redirect(w, r, casdoorAuthURL, http.StatusFound)
}

// https://door.casdoor.com/swagger/#/Token%20API/ApiController.GetOAuthToken
type TokenRequest struct {
	Code string `json:"code"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

func getToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var tokenReq TokenRequest
	if err := json.NewDecoder(r.Body).Decode(&tokenReq); err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		log.Printf("Failed to decode request body: %v", err)
		return
	}

	// 用code换token
	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", tokenReq.Code)
	data.Set("redirect_uri", callbackURL)

	resp, err := http.PostForm(casdoorTokenURL, data)
	if err != nil {
		http.Error(w, "Failed to exchange code for token", http.StatusInternalServerError)
		log.Printf("Failed to exchange code for token: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		log.Printf("Failed to read response body: %v", err)
		return
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		http.Error(w, "Failed to unmarshal token response", http.StatusInternalServerError)
		log.Printf("Failed to unmarshal token response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokenResponse); err != nil {
		http.Error(w, "Failed to encode token response", http.StatusInternalServerError)
		log.Printf("Failed to encode token response: %v", err)
		return
	}
}

// https://door.casdoor.com/swagger/#/Account%20API/ApiController.UserInfo

type UserInfoResponse struct {
	Name string `json:"name"`
	// omit the others
}

func getUserInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	req, err := http.NewRequest("GET", casdoorUserInfoURL, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		log.Printf("Failed to create request: %v", err)
		return
	}
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to get user info", http.StatusInternalServerError)
		log.Printf("Failed to get user info: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		log.Printf("Failed to read response body: %v", err)
		return
	}

	var userInfo UserInfoResponse
	if err := json.Unmarshal(body, &userInfo); err != nil {
		http.Error(w, "Failed to unmarshal user info response", http.StatusInternalServerError)
		log.Printf("Failed to unmarshal user info response: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userInfo); err != nil {
		http.Error(w, "Failed to encode user info response", http.StatusInternalServerError)
		log.Printf("Failed to encode user info response: %v", err)
		return
	}
}

// 跨域设置
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
