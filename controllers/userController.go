package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"vetclinic-rest-api/structs"
	"vetclinic-rest-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Register new user
func RegisterUser(c *gin.Context, db *sql.DB) {
	var newUser structs.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Validate required fields
	if newUser.Name == "" || newUser.Email == "" || newUser.PasswordHash == "" || newUser.Role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name, Email, Password, and Role are required"})
		return
	}

	// Validate roles
	validRoles := []string{"Staff", "Doctor", "Admin"}
	isValid := false
	for _, r := range validRoles {
		if newUser.Role == r {
			isValid = true
			break
		}
	}
	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role can only be Staff, Doctor, or Admin"})
		return
	}

	// Hash password before saving
	hashedPassword, err := utils.HashPassword(newUser.PasswordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Default values
	newUser.Id = uuid.New()
	newUser.PasswordHash = hashedPassword
	newUser.ActiveStatus = 1
	newUser.CreatedAt = time.Now()
	newUser.CreatedBy = newUser.Id.String()
	newUser.ModifiedAt = newUser.CreatedAt
	newUser.ModifiedBy = newUser.Id.String()

	// Insert into Users table
	query := `INSERT INTO "Users" (id, name, email, phone, password_hash, role, active_status, created_at, created_by, modified_at, modified_by) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err = db.Exec(query,
		newUser.Id, newUser.Name, newUser.Email, newUser.Phone,
		newUser.PasswordHash, newUser.Role, newUser.ActiveStatus,
		newUser.CreatedAt, newUser.CreatedBy, newUser.ModifiedAt, newUser.ModifiedBy)

	if err != nil {
		log.Println("Error inserting new User:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Return user data (without password hash)
	c.JSON(http.StatusCreated, gin.H{
		"id":            newUser.Id,
		"name":          newUser.Name,
		"email":         newUser.Email,
		"phone":         newUser.Phone,
		"role":          newUser.Role,
		"active_status": newUser.ActiveStatus,
		"created_at": newUser.CreatedAt,
		"created_by": newUser.CreatedBy,
	})
}

func LoginUser(c *gin.Context, db *sql.DB) {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Get user from DB
	var user structs.User
	query := `SELECT id, name, email, phone, password_hash, role, active_status FROM "Users" WHERE email=$1 AND active_status=1`
	err := db.QueryRow(query, req.Email).Scan(
		&user.Id, &user.Name, &user.Email, &user.Phone,
		&user.PasswordHash, &user.Role, &user.ActiveStatus,
	)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Check password with bcrypt
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id":    user.Id.String(), // convert uuid.UUID to string
		"email": user.Email,
		"role":  user.Role,
		"exp":   time.Now().Add(time.Hour * 2).Unix(), // expires in 2 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(utils.JwtSecret) // jwtSecret loaded from .env
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Return token to client
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func FetchProfile(c *gin.Context, db *sql.DB) {
    userId := c.Param("id")

    var user structs.User
    query := `SELECT id, name, email, phone, role, active_status, created_at, created_by, modified_at, modified_by FROM "Users" 
			WHERE id=$1 AND active_status=1`
    err := db.QueryRow(query, userId).Scan(
        &user.Id, &user.Name, &user.Email, &user.Phone,
        &user.Role, &user.ActiveStatus,
        &user.CreatedAt, &user.CreatedBy, &user.ModifiedAt, &user.ModifiedBy,
    )
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":            user.Id,
        "name":          user.Name,
        "email":         user.Email,
        "phone":         user.Phone,
        "role":          user.Role,
        "active_status": user.ActiveStatus,
        "created_at":    user.CreatedAt,
        "created_by":    user.CreatedBy,
        "modified_at":   user.ModifiedAt,
        "modified_by":   user.ModifiedBy,
    })
}

func GetUserByRole(c *gin.Context, db *sql.DB) {
    role := c.Param("role")

    // Validate role
    validRoles := []string{"Staff", "Doctor", "Admin"}
    isValid := false
    for _, r := range validRoles {
        if role == r {
            isValid = true
            break
        }
    }
    if !isValid {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be Staff, Doctor, or Admin"})
        return
    }

    query := `SELECT id, name, email, phone, role, active_status, created_at, created_by, modified_at, modified_by
            FROM "Users"
            WHERE role=$1 AND active_status=1`

    rows, err := db.Query(query, role)
    if err != nil {
        log.Println("Error querying Users by role:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
        return
    }
    defer rows.Close()

    var users []structs.User
    for rows.Next() {
        var user structs.User
        if err := rows.Scan(
            &user.Id, &user.Name, &user.Email, &user.Phone,
            &user.Role, &user.ActiveStatus,
            &user.CreatedAt, &user.CreatedBy,
            &user.ModifiedAt, &user.ModifiedBy,
        ); err != nil {
            log.Println("Error scanning User row:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse users"})
            return
        }
        users = append(users, user)
    }

    if len(users) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No active users found for this role"})
        return
    }

    c.JSON(http.StatusOK, users)
}

func UpdateUser(c *gin.Context, db *sql.DB) {
    userId := c.Param("id")
    var req struct {
        Name  string `json:"name"`
        Email string `json:"email"`
        Phone string `json:"phone"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    // Get existing user
    var existing structs.User
    query := `SELECT id, name, email, phone, role, active_status, created_at, created_by FROM "Users" 
			WHERE id=$1 AND active_status=1`
    err := db.QueryRow(query, userId).Scan(
        &existing.Id, &existing.Name, &existing.Email, &existing.Phone,
        &existing.Role, &existing.ActiveStatus,
        &existing.CreatedAt, &existing.CreatedBy,
    )
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Update fields if provided
    if req.Name != "" {
        existing.Name = req.Name
    }
    if req.Email != "" {
        existing.Email = req.Email
    }
    if req.Phone != "" {
        existing.Phone = req.Phone
    }

    existing.ModifiedAt = time.Now()
    existing.ModifiedBy = userId // or from JWT claims

    // Update query
    updateQuery := `UPDATE "Users"
                    SET name=$1, email=$2, phone=$3, modified_at=$4, modified_by=$5
                    WHERE id=$6 AND active_status=1`
    _, err = db.Exec(updateQuery,
        existing.Name, existing.Email, existing.Phone,
        existing.ModifiedAt, existing.ModifiedBy, existing.Id,
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":            existing.Id,
        "name":          existing.Name,
        "email":         existing.Email,
        "phone":         existing.Phone,
        "role":          existing.Role,
        "active_status": existing.ActiveStatus,
        "created_at":    existing.CreatedAt,
        "created_by":    existing.CreatedBy,
        "modified_at":   existing.ModifiedAt,
        "modified_by":   existing.ModifiedBy,
    })
}

func UpdateRole(c *gin.Context, db *sql.DB) {
    targetUserId := c.Param("id") // the user whose role is being updated
    var req struct {
        Role string `json:"role"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    // Validate role
    validRoles := []string{"Staff", "Doctor", "Admin"}
    isValid := false
    for _, r := range validRoles {
        if req.Role == r {
            isValid = true
            break
        }
    }
    if !isValid {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Role must be Staff, Doctor, or Admin"})
        return
    }

    // Get the user_id from JWT context (the one performing the update)
    userIdVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
        return
    }

    // Type assertion to string
    createdBy, ok := userIdVal.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
        return
    }

    // Update only role + modified fields
    updateQuery := `UPDATE "Users"
                    SET role=$1, modified_at=$2, modified_by=$3
                    WHERE id=$4 AND active_status=1`
    _, err := db.Exec(updateQuery, req.Role, time.Now(), createdBy, targetUserId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update role"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":      targetUserId,
        "role":    req.Role,
        "updated_by": createdBy,
        "message": "Role updated successfully",
    })
}

func ChangePassword(c *gin.Context, db *sql.DB) {
    userId := c.Param("id")
    var req struct {
        OldPassword     string `json:"old_password"`
        NewPassword     string `json:"new_password"`
        ConfirmPassword string `json:"confirm_password"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    // Get existing user
    var existing structs.User
    query := `SELECT id, password_hash FROM "Users" WHERE id=$1 AND active_status=1`
    err := db.QueryRow(query, userId).Scan(&existing.Id, &existing.PasswordHash)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // Verify old password
    if !utils.CheckPasswordHash(req.OldPassword, existing.PasswordHash) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
        return
    }

    // Verify new + confirm match
    if req.NewPassword != req.ConfirmPassword {
        c.JSON(http.StatusBadRequest, gin.H{"error": "New password and confirmation do not match"})
        return
    }

    // Hash new password
    hashedPassword, err := utils.HashPassword(req.NewPassword)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
        return
    }

    // Update DB
    updateQuery := `UPDATE "Users" SET password_hash=$1, modified_at=$2, modified_by=$3 WHERE id=$4`
    _, err = db.Exec(updateQuery, hashedPassword, time.Now(), userId, userId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func UpdateUserActiveStatus(c *gin.Context, db *sql.DB) {
    targetUserId := c.Param("id") // the user being deactivated

    // Get the user_id from JWT (the one performing the action)
    userIdVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
        return
    }

    createdBy, ok := userIdVal.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
        return
    }

    query := `UPDATE "Users"
            SET active_status=0, modified_at=$1, modified_by=$2
            WHERE id=$3 AND active_status=1`

    _, err := db.Exec(query, time.Now(), createdBy, targetUserId)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate user"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":          targetUserId,
        "deactivated_by": createdBy,
        "message":     "User deactivated successfully",
    })
}