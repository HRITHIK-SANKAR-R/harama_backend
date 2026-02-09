package middleware

import (
	"encoding/json"
	"fmt"
	"harama/internal/auth"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// SupabaseAuthMiddleware validates Supabase JWT tokens by calling Supabase Auth API.
// This supports both HS256 and ES256 tokens and ensures the token is not revoked.
func SupabaseAuthMiddleware(supabaseURL, anonKey, jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Try Authorization: Bearer <token> first
			authHeader := r.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				
				// Validate with Supabase API
				user, err := validateTokenWithSupabase(token, supabaseURL, anonKey)
				if err != nil {
					http.Error(w, "invalid or expired token: "+err.Error(), http.StatusUnauthorized)
					return
				}

				userID, err := uuid.Parse(user.ID)
				if err != nil {
					http.Error(w, "invalid user id in token", http.StatusUnauthorized)
					return
				}

				// Use user_id as tenant_id (single-tenant per user model)
				ctx := auth.WithTenantID(r.Context(), userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Fallback: X-Tenant-ID header (for dev/testing)
			tenantIDStr := r.Header.Get("X-Tenant-ID")
			if tenantIDStr != "" {
				tenantID, err := uuid.Parse(tenantIDStr)
				if err != nil {
					http.Error(w, "invalid X-Tenant-ID", http.StatusUnauthorized)
					return
				}
				ctx := auth.WithTenantID(r.Context(), tenantID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			http.Error(w, "authorization required: provide Bearer token or X-Tenant-ID", http.StatusUnauthorized)
		})
	}
}

type SupabaseUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func validateTokenWithSupabase(token, supabaseURL, anonKey string) (*SupabaseUser, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("GET", supabaseURL+"/auth/v1/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", anonKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase rejected token: %s", resp.Status)
	}

	var user SupabaseUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}