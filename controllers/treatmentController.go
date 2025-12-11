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

func CreateTreatment(c *gin.Context, db *sql.DB) {
    var treatment structs.Treatment
    if err := c.ShouldBindJSON(&treatment); err != nil {
        log.Println("Error binding JSON for new Treatment:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    if treatment.MedicalRecordId == uuid.Nil || treatment.Description == "" || treatment.Cost <= 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "MedicalRecordId, Description, and Cost are required"})
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

    treatment.Id = uuid.New()
    treatment.ActiveStatus = 1
    treatment.CreatedAt = time.Now()
    treatment.CreatedBy = createdBy
    treatment.ModifiedAt = treatment.CreatedAt
    treatment.ModifiedBy = createdBy

    query := `INSERT INTO "Treatments"
        (id, medicalrecord_id, doctor_id, description, cost,
        active_status, created_at, created_by, modified_at, modified_by)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

    _, err := db.Exec(query,
        treatment.Id, treatment.MedicalRecordId, treatment.DoctorId,
        treatment.Description, treatment.Cost,
        treatment.ActiveStatus, treatment.CreatedAt, treatment.CreatedBy,
        treatment.ModifiedAt, treatment.ModifiedBy,
    )
    if err != nil {
        log.Println("Error inserting Treatment:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create treatment"})
        return
    }

    c.JSON(http.StatusCreated, treatment)
}

func GetTreatmentsByMedicalRecordId(c *gin.Context, db *sql.DB) {
    recordId := c.Param("medicalrecord_id")

    query := `SELECT id, medicalrecord_id, doctor_id, description, cost,
            active_status, created_at, created_by, modified_at, modified_by
            FROM "Treatments"
            WHERE medicalrecord_id=$1 AND active_status=1
            ORDER BY created_at DESC`

    rows, err := db.Query(query, recordId)
    if err != nil {
        log.Println("Error fetching treatments:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch treatments"})
        return
    }
    defer rows.Close()

    var treatments []structs.Treatment
    for rows.Next() {
        var t structs.Treatment
        if err := rows.Scan(
            &t.Id, &t.MedicalRecordId, &t.DoctorId, &t.Description, &t.Cost,
            &t.ActiveStatus, &t.CreatedAt, &t.CreatedBy, &t.ModifiedAt, &t.ModifiedBy,
        ); err != nil {
            log.Println("Error scanning treatment row:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse treatments"})
            return
        }
        treatments = append(treatments, t)
    }

    c.JSON(http.StatusOK, treatments)
}

func UpdateTreatment(c *gin.Context, db *sql.DB) {
    treatmentId := c.Param("id")

    // 1. Fetch existing treatment
    var existing structs.Treatment
    fetchQuery := `SELECT id, medicalrecord_id, doctor_id, description, cost,
                    active_status, created_at, created_by, modified_at, modified_by
                    FROM "Treatments"
                    WHERE id=$1 AND active_status=1`

    err := db.QueryRow(fetchQuery, treatmentId).Scan(
        &existing.Id, &existing.MedicalRecordId, &existing.DoctorId,
        &existing.Description, &existing.Cost, &existing.ActiveStatus,
        &existing.CreatedAt, &existing.CreatedBy,
        &existing.ModifiedAt, &existing.ModifiedBy,
    )
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Treatment not found"})
        return
    }

    // 2. Bind incoming JSON
    var req structs.Treatment
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Println("Error binding JSON for UpdateTreatment:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    // 3. Merge fields
    if req.Description != "" {
        existing.Description = req.Description
    }
    if req.Cost != 0 {
        existing.Cost = req.Cost
    }

    // 4. Get user_id from JWT
    userIdVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
        return
    }
    modifiedBy := userIdVal.(string)

    // 5. Update query
    updateQuery := `UPDATE "Treatments"
                    SET description=$1, cost=$2, modified_at=$3, modified_by=$4
                    WHERE id=$5 AND active_status=1`

    _, err = db.Exec(updateQuery,
        existing.Description, existing.Cost, time.Now(), modifiedBy, treatmentId,
    )
    if err != nil {
        log.Println("Error updating Treatment:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update treatment"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Treatment updated successfully"})
}

func UpdateTreatmentActiveStatus(c *gin.Context, db *sql.DB) {
    treatmentId := c.Param("id")

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

    query := `UPDATE "Treatments"
            SET active_status=0, modified_at=$1, modified_by=$2
            WHERE id=$3`

    _, err := db.Exec(query, time.Now(), modifiedBy, treatmentId)
    if err != nil {
        log.Println("Error soft deleting Treatment:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate treatment"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":             treatmentId,
        "deactivated_by": modifiedBy,
        "message":        "Treatment deactivated successfully",
    })
}