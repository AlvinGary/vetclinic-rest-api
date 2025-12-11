package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"time"
	"vetclinic-rest-api/structs"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreatePet(c *gin.Context, db *sql.DB) {
    var newPet structs.Pet
    if err := c.ShouldBindJSON(&newPet); err != nil {
        log.Println("Error binding JSON for new Pet:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    if newPet.Name == "" || newPet.Species == "" || newPet.Gender == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Name, Species, and Gender are required"})
        return
    }

    // Get user_id from JWT
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

    newPet.Id = uuid.New()
    newPet.ActiveStatus = 1
    newPet.CreatedAt = time.Now()
    newPet.CreatedBy = createdBy
    newPet.ModifiedAt = newPet.CreatedAt
    newPet.ModifiedBy = createdBy

    query := `INSERT INTO "Pets"
        (id, name, species, breed, gender, birth_date, owner_name, owner_phone,
		active_status, created_at, created_by, modified_at, modified_by)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`

    _, err := db.Exec(query,
        newPet.Id, newPet.Name, newPet.Species, newPet.Breed, newPet.Gender,
        newPet.BirthDate, newPet.OwnerName, newPet.OwnerPhone,
        newPet.ActiveStatus, newPet.CreatedAt, newPet.CreatedBy,
        newPet.ModifiedAt, newPet.ModifiedBy,
    )
    if err != nil {
        log.Println("Error inserting new Pet:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pet"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "id":          newPet.Id,
        "name":        newPet.Name,
        "species":     newPet.Species,
        "breed":       newPet.Breed,
        "gender":      newPet.Gender,
        "birth_date":  newPet.BirthDate,
        "owner_name":  newPet.OwnerName,
        "owner_phone": newPet.OwnerPhone,
    })
}

func FetchPetProfile(c *gin.Context, db *sql.DB) {
    petId := c.Param("id")
    var pet structs.Pet

    query := `SELECT id, name, species, breed, gender, birth_date, owner_name, owner_phone, active_status, created_at, created_by, modified_at, modified_by 
			FROM "Pets" WHERE id=$1 AND active_status=1`
    err := db.QueryRow(query, petId).Scan(
        &pet.Id, &pet.Name, &pet.Species, &pet.Breed, &pet.Gender,
        &pet.BirthDate, &pet.OwnerName, &pet.OwnerPhone,
        &pet.ActiveStatus, &pet.CreatedAt, &pet.CreatedBy,
        &pet.ModifiedAt, &pet.ModifiedBy,
    )
    if err != nil {
        log.Println("Error fetching Pet profile:", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
        return
    }

    c.JSON(http.StatusOK, pet)
}

func UpdatePet(c *gin.Context, db *sql.DB) {
    petId := c.Param("id")

    // 1. Fetch existing pet
    var existing structs.Pet
    fetchQuery := `SELECT id, name, species, breed, gender, birth_date,
                        owner_name, owner_phone, active_status,
                        created_at, created_by, modified_at, modified_by
                    FROM "Pets"
                    WHERE id=$1 AND active_status=1`

    err := db.QueryRow(fetchQuery, petId).Scan(
        &existing.Id, &existing.Name, &existing.Species, &existing.Breed,
        &existing.Gender, &existing.BirthDate, &existing.OwnerName,
        &existing.OwnerPhone, &existing.ActiveStatus,
        &existing.CreatedAt, &existing.CreatedBy,
        &existing.ModifiedAt, &existing.ModifiedBy,
    )
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
        return
    }

    // 2. Bind incoming JSON
    var req structs.Pet
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Println("Error binding JSON for UpdatePet:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    // 3. Merge fields (only overwrite if provided)
    if req.Name != "" {
        existing.Name = req.Name
    }
    if req.Species != "" {
        existing.Species = req.Species
    }
    if req.Breed != "" {
        existing.Breed = req.Breed
    }
    if req.Gender != "" {
        existing.Gender = req.Gender
    }
    if req.BirthDate != "" {
        existing.BirthDate = req.BirthDate
    }
    if req.OwnerName != "" {
        existing.OwnerName = req.OwnerName
    }
    if req.OwnerPhone != "" {
        existing.OwnerPhone = req.OwnerPhone
    }

    // 4. Get user_id from JWT
    userIdVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
        return
    }
    modifiedBy := userIdVal.(string)

    // 5. Update query
    updateQuery := `UPDATE "Pets"
                    SET name=$1, species=$2, breed=$3, gender=$4, birth_date=$5,
                        owner_name=$6, owner_phone=$7,
                        modified_at=$8, modified_by=$9
                    WHERE id=$10 AND active_status=1`

    _, err = db.Exec(updateQuery,
        existing.Name, existing.Species, existing.Breed, existing.Gender,
        existing.BirthDate, existing.OwnerName, existing.OwnerPhone,
        time.Now(), modifiedBy, petId,
    )
    if err != nil {
        log.Println("Error updating Pet:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pet"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Pet updated successfully"})
}

func UpdatePetActiveStatus(c *gin.Context, db *sql.DB) {
    petId := c.Param("id")

    // Get user_id from JWT
    userIdVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
        return
    }
    modifiedBy, ok := userIdVal.(string)
    if !ok {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
        return
    }

    query := `UPDATE "Pets"
            SET active_status=0, modified_at=$1, modified_by=$2
            WHERE id=$3`

    _, err := db.Exec(query, time.Now(), modifiedBy, petId)
    if err != nil {
        log.Println("Error soft deleting Pet:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate pet"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":            petId,
        "deactivated_by": modifiedBy,
        "message":       "Pet deactivated successfully",
    })
}

func FetchPetsByOwner(c *gin.Context, db *sql.DB) {
    ownerName := c.Param("owner_name")
    ownerPhone := c.Param("owner_phone")

    query := `SELECT id, name, species, breed, gender, birth_date, owner_name, owner_phone, active_status, created_at, created_by, modified_at, modified_by
            FROM "Pets"
            WHERE owner_name=$1 AND owner_phone=$2 AND active_status=1`

    rows, err := db.Query(query, ownerName, ownerPhone)
    if err != nil {
        log.Println("Error querying Pets by owner:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pets"})
        return
    }
    defer rows.Close()

    var pets []structs.Pet
    for rows.Next() {
        var pet structs.Pet
        if err := rows.Scan(
            &pet.Id, &pet.Name, &pet.Species, &pet.Breed, &pet.Gender,
            &pet.BirthDate, &pet.OwnerName, &pet.OwnerPhone,
            &pet.ActiveStatus, &pet.CreatedAt, &pet.CreatedBy,
            &pet.ModifiedAt, &pet.ModifiedBy,
        ); err != nil {
            log.Println("Error scanning Pet row:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse pets"})
            return
        }
        pets = append(pets, pet)
    }

    if len(pets) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "No active pets found for this owner"})
        return
    }

    c.JSON(http.StatusOK, pets)
}