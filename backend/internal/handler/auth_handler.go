package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"log-management-backend/internal/auth"
	"log-management-backend/internal/models"
	"log-management-backend/internal/repository"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
}

func NewAuthHandler(userRepo *repository.UserRepository) *AuthHandler {
	return &AuthHandler{userRepo: userRepo}
}

// LoginRequest - request body สำหรับ login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse - response หลัง login สำเร็จ
type LoginResponse struct {
	Token    string `json:"token"`
	Username string `json:"username"`
	Role     string `json:"role"`
	TenantID string `json:"tenant_id"`
}

// RegisterRequest - request body สำหรับ register
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`      // "admin" หรือ "viewer"
	TenantID string `json:"tenant_id"` // tenant ที่จะผูกกับ user
}

// LoginHandler - POST /login
func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Parse request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. ตรวจสอบ username ว่ามี user ไหม
	user, err := h.userRepo.GetUserByUsername(r.Context(), req.Username)
	if err != nil {
		log.Printf("Login failed - user not found: %s", req.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 3. ตรวจสอบ password (bcrypt compare)
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Printf("Login failed - wrong password: %s", req.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// 4. สร้าง JWT token
	token, err := auth.GenerateToken(user.ID.String(), user.Username, user.Role, user.TenantID) // เพิ่ม .String()
	if err != nil {
		log.Printf("Failed to generate token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 5. ส่ง token กลับ
	response := LoginResponse{
		Token:    token,
		Username: user.Username,
		Role:     user.Role,
		TenantID: user.TenantID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(response)

	log.Printf("✅ Login success: %s (role: %s, tenant: %s)", user.Username, user.Role, user.TenantID)
}

// RegisterHandler - POST /register
func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Parse request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Validate input
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	if req.Role != "admin" && req.Role != "viewer" {
		http.Error(w, "Role must be 'admin' or 'viewer'", http.StatusBadRequest)
		return
	}

	if req.TenantID == "" {
		http.Error(w, "Tenant ID is required", http.StatusBadRequest)
		return
	}

	// 3. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// 4. สร้าง user ใหม่
	newUser := models.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		TenantID:     req.TenantID,
	}

	if err := h.userRepo.CreateUser(r.Context(), newUser); err != nil {
		log.Printf("Failed to create user: %v", err)
		http.Error(w, "Username already exists or database error", http.StatusConflict)
		return
	}

	// 5. ส่ง response สำเร็จ
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User created successfully",
	})

	log.Printf("✅ User registered: %s (role: %s, tenant: %s)", req.Username, req.Role, req.TenantID)
}
