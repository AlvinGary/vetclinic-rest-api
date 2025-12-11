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

func CreateMedicalRecord(c *gin.Context, db *sql.DB) {
    var record structs.MedicalRecord
    if err := c.ShouldBindJSON(&record); err != nil {
        log.Println("Error binding JSON for new MedicalRecord:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    if record.AppointmentId == uuid.Nil || record.PetId == uuid.Nil || record.Diagnosis == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "AppointmentId, PetId, and Diagnosis are required"})
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

    record.Id = uuid.New()
    record.ActiveStatus = 1
    record.CreatedAt = time.Now()
    record.CreatedBy = createdBy
    record.ModifiedAt = record.CreatedAt
    record.ModifiedBy = createdBy

    query := `INSERT INTO "MedicalRecords"
        (id, appointment_id, pet_id, diagnosis, notes,
        active_status, created_at, created_by, modified_at, modified_by)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

    _, err := db.Exec(query,
        record.Id, record.AppointmentId, record.PetId, record.Diagnosis, record.Notes,
        record.ActiveStatus, record.CreatedAt, record.CreatedBy, record.ModifiedAt, record.ModifiedBy,
    )
    if err != nil {
        log.Println("Error inserting MedicalRecord:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medical record"})
        return
    }

    c.JSON(http.StatusCreated, record)
}

func GetMedicalRecordByAppointmentId(c *gin.Context, db *sql.DB) {
    appointmentId := c.Param("appointment_id")

    query := `SELECT id, appointment_id, pet_id, diagnosis, notes,
            active_status, created_at, created_by, modified_at, modified_by
            FROM "MedicalRecords"
            WHERE appointment_id=$1 AND active_status=1`

    var record structs.MedicalRecord

    err := db.QueryRow(query, appointmentId).Scan(
        &record.Id, &record.AppointmentId, &record.PetId, &record.Diagnosis,
        &record.Notes, &record.ActiveStatus, &record.CreatedAt, &record.CreatedBy,
        &record.ModifiedAt, &record.ModifiedBy,
    )
    if err != nil {
        log.Println("Error fetching MedicalRecord:", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
        return
    }

    c.JSON(http.StatusOK, record)
}

func UpdateMedicalRecord(c *gin.Context, db *sql.DB) {
    recordId := c.Param("id")

    // 1. Fetch existing record
    var existing structs.MedicalRecord
    fetchQuery := `SELECT id, appointment_id, pet_id, diagnosis, notes,
                          active_status, created_at, created_by,
                          modified_at, modified_by
                   FROM "MedicalRecords"
                   WHERE id=$1 AND active_status=1`

    err := db.QueryRow(fetchQuery, recordId).Scan(
        &existing.Id, &existing.AppointmentId, &existing.PetId,
        &existing.Diagnosis, &existing.Notes, &existing.ActiveStatus,
        &existing.CreatedAt, &existing.CreatedBy,
        &existing.ModifiedAt, &existing.ModifiedBy,
    )
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Medical record not found"})
        return
    }

    // 2. Bind incoming JSON
    var req structs.MedicalRecord
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Println("Error binding JSON for UpdateMedicalRecord:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    // 3. Merge fields
    if req.Diagnosis != "" {
        existing.Diagnosis = req.Diagnosis
    }
    if req.Notes != "" {
        existing.Notes = req.Notes
    }

    // 4. Get user_id from JWT
    userIdVal, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
        return
    }
    modifiedBy := userIdVal.(string)

    // 5. Update query
    updateQuery := `UPDATE "MedicalRecords"
                    SET diagnosis=$1, notes=$2, modified_at=$3, modified_by=$4
                    WHERE id=$5 AND active_status=1`

    _, err = db.Exec(updateQuery,
        existing.Diagnosis, existing.Notes, time.Now(), modifiedBy, recordId,
    )
    if err != nil {
        log.Println("Error updating MedicalRecord:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update medical record"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Medical record updated successfully"})
}

func UpdateMedicalRecordActiveStatus(c *gin.Context, db *sql.DB) {
    recordId := c.Param("id")

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

    query := `UPDATE "MedicalRecords"
            SET active_status=0, modified_at=$1, modified_by=$2
            WHERE id=$3`

    _, err := db.Exec(query, time.Now(), modifiedBy, recordId)
    if err != nil {
        log.Println("Error soft deleting MedicalRecord:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate medical record"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":             recordId,
        "deactivated_by": modifiedBy,
        "message":        "Medical record deactivated successfully",
    })
}