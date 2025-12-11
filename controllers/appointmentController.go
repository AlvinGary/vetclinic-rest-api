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

func CreateAppointment(c *gin.Context, db *sql.DB) {
    var newAppointment structs.Appointment
    if err := c.ShouldBindJSON(&newAppointment); err != nil {
        log.Println("Error binding JSON for new Appointment:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    // Required fields
    if newAppointment.PetId == uuid.Nil || newAppointment.AppointmentDatetime.IsZero() {
        c.JSON(http.StatusBadRequest, gin.H{"error": "PetId and AppointmentDatetime are required"})
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

    newAppointment.Id = uuid.New()
    newAppointment.Status = "Pending"
    newAppointment.ActiveStatus = 1
    newAppointment.CreatedAt = time.Now()
    newAppointment.CreatedBy = createdBy
    newAppointment.ModifiedAt = newAppointment.CreatedAt
    newAppointment.ModifiedBy = createdBy

    query := `INSERT INTO "Appointments"
        (id, pet_id, doctor_id, status, appointment_datetime, notes,
        active_status, created_at, created_by, modified_at, modified_by)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

    _, err := db.Exec(query,
        newAppointment.Id, newAppointment.PetId, newAppointment.DoctorId, newAppointment.Status,
        newAppointment.AppointmentDatetime, newAppointment.Notes, newAppointment.ActiveStatus,
        newAppointment.CreatedAt, newAppointment.CreatedBy, newAppointment.ModifiedAt, newAppointment.ModifiedBy,
    )
    if err != nil {
        log.Println("Error inserting new Appointment:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create appointment"})
        return
    }

    c.JSON(http.StatusCreated, newAppointment)
}

func FetchAppointment(c *gin.Context, db *sql.DB) {
    appointmentId := c.Param("id")
    var appt structs.Appointment

    query := `SELECT id, pet_id, doctor_id, status, appointment_datetime, notes, active_status, created_at, created_by, modified_at, modified_by
            FROM "Appointments"
            WHERE id=$1 AND active_status=1`
    err := db.QueryRow(query, appointmentId).Scan(
        &appt.Id, &appt.PetId, &appt.DoctorId, &appt.Status,
        &appt.AppointmentDatetime, &appt.Notes, &appt.ActiveStatus,
        &appt.CreatedAt, &appt.CreatedBy, &appt.ModifiedAt, &appt.ModifiedBy,
    )
    if err != nil {
        log.Println("Error fetching Appointment:", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
        return
    }

    c.JSON(http.StatusOK, appt)
}

func UpdateAppointment(c *gin.Context, db *sql.DB) {
    appointmentId := c.Param("id")

    // 1. Fetch existing appointment
    var existing structs.Appointment
    fetchQuery := `SELECT id, pet_id, doctor_id, status, appointment_datetime,
                          notes, active_status, created_at, created_by,
                          modified_at, modified_by
                   FROM "Appointments"
                   WHERE id=$1 AND active_status=1`

    err := db.QueryRow(fetchQuery, appointmentId).Scan(
        &existing.Id, &existing.PetId, &existing.DoctorId, &existing.Status,
        &existing.AppointmentDatetime, &existing.Notes, &existing.ActiveStatus,
        &existing.CreatedAt, &existing.CreatedBy,
        &existing.ModifiedAt, &existing.ModifiedBy,
    )
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
        return
    }

    // 2. Bind incoming JSON
    var req structs.Appointment
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Println("Error binding JSON for UpdateAppointment:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    // 3. Merge fields
    if req.PetId != uuid.Nil {
        existing.PetId = req.PetId
    }
    if req.DoctorId != uuid.Nil {
        existing.DoctorId = req.DoctorId
    }
    if !req.AppointmentDatetime.IsZero() {
        existing.AppointmentDatetime = req.AppointmentDatetime
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
    updateQuery := `UPDATE "Appointments"
                    SET pet_id=$1, doctor_id=$2, appointment_datetime=$3,
                        notes=$4, modified_at=$5, modified_by=$6
                    WHERE id=$7 AND active_status=1`

    _, err = db.Exec(updateQuery,
        existing.PetId, existing.DoctorId, existing.AppointmentDatetime,
        existing.Notes, time.Now(), modifiedBy, appointmentId,
    )
    if err != nil {
        log.Println("Error updating Appointment:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update appointment"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Appointment updated successfully"})
}

func UpdateAppointmentStatus(c *gin.Context, db *sql.DB) {
    appointmentId := c.Param("id")
    var req struct {
        Status string `json:"status"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        log.Println("Error binding JSON for UpdateAppointmentStatus:", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
        return
    }

    validStatuses := []string{"Pending", "Cancelled", "Completed"}
    isValid := false
    for _, s := range validStatuses {
        if req.Status == s {
            isValid = true
            break
        }
    }
    if !isValid {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status"})
        return
    }

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

    query := `UPDATE "Appointments"
            SET status=$1, modified_at=$2, modified_by=$3
            WHERE id=$4 AND active_status=1`

    _, err := db.Exec(query, req.Status, time.Now(), modifiedBy, appointmentId)
    if err != nil {
        log.Println("Error updating Appointment status:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Appointment status updated successfully"})
}

func UpdateAppointmentActiveStatus(c *gin.Context, db *sql.DB) {
    appointmentId := c.Param("id")

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

    query := `UPDATE "Appointments"
            SET active_status=0, modified_at=$1, modified_by=$2
            WHERE id=$3`

    _, err := db.Exec(query, time.Now(), modifiedBy, appointmentId)
    if err != nil {
        log.Println("Error soft deleting Appointment:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to deactivate appointment"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":             appointmentId,
        "deactivated_by": modifiedBy,
        "message":        "Appointment deactivated successfully",
    })
}

func GetAppointmentsByPetId(c *gin.Context, db *sql.DB) {
    petId := c.Param("pet_id")

    query := `SELECT id, pet_id, doctor_id, status, appointment_datetime, notes, active_status, created_at, created_by, modified_at, modified_by
            FROM "Appointments"
            WHERE pet_id=$1 AND active_status=1
            ORDER BY appointment_datetime DESC`// sort from newest to oldest

    rows, err := db.Query(query, petId)
    if err != nil {
        log.Println("Error fetching appointments by pet_id:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
        return
    }
    defer rows.Close()

    var appointments []structs.Appointment
    for rows.Next() {
        var appt structs.Appointment
        if err := rows.Scan(
            &appt.Id, &appt.PetId, &appt.DoctorId, &appt.Status,
            &appt.AppointmentDatetime, &appt.Notes, &appt.ActiveStatus,
            &appt.CreatedAt, &appt.CreatedBy, &appt.ModifiedAt, &appt.ModifiedBy,
        ); err != nil {
            log.Println("Error scanning appointment row:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse appointments"})
            return
        }
        appointments = append(appointments, appt)
    }

    c.JSON(http.StatusOK, appointments)
}

func GetAppointmentsByDoctorId(c *gin.Context, db *sql.DB) {
    doctorId := c.Param("doctor_id")

    query := `SELECT id, pet_id, doctor_id, status, appointment_datetime, notes, active_status, created_at, created_by, modified_at, modified_by
            FROM "Appointments"
            WHERE doctor_id=$1 AND active_status=1
            ORDER BY appointment_datetime DESC`

    rows, err := db.Query(query, doctorId)
    if err != nil {
        log.Println("Error fetching appointments by doctor_id:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
        return
    }
    defer rows.Close()

    var appointments []structs.Appointment
    for rows.Next() {
        var appt structs.Appointment
        if err := rows.Scan(
            &appt.Id, &appt.PetId, &appt.DoctorId, &appt.Status,
            &appt.AppointmentDatetime, &appt.Notes, &appt.ActiveStatus,
            &appt.CreatedAt, &appt.CreatedBy, &appt.ModifiedAt, &appt.ModifiedBy,
        ); err != nil {
            log.Println("Error scanning appointment row:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse appointments"})
            return
        }
        appointments = append(appointments, appt)
    }

    c.JSON(http.StatusOK, appointments)
}

func GetAppointmentsByAppointmentDate(c *gin.Context, db *sql.DB) {
    dateStr := c.Param("date") // YYYY-MM-DD

    query := `SELECT id, pet_id, doctor_id, status, appointment_datetime, notes, active_status, created_at, created_by, modified_at, modified_by
            FROM "Appointments"
            WHERE DATE(appointment_datetime) = $1
            AND active_status=1
            ORDER BY appointment_datetime DESC`

    rows, err := db.Query(query, dateStr)
    if err != nil {
        log.Println("Error fetching appointments by date:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch appointments"})
        return
    }
    defer rows.Close()

    var appointments []structs.Appointment
    for rows.Next() {
        var appt structs.Appointment
        if err := rows.Scan(
            &appt.Id, &appt.PetId, &appt.DoctorId, &appt.Status,
            &appt.AppointmentDatetime, &appt.Notes, &appt.ActiveStatus,
            &appt.CreatedAt, &appt.CreatedBy, &appt.ModifiedAt, &appt.ModifiedBy,
        ); err != nil {
            log.Println("Error scanning appointment row:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse appointments"})
            return
        }
        appointments = append(appointments, appt)
    }

    c.JSON(http.StatusOK, appointments)
}

func GetFullAppointmentDetail(c *gin.Context, db *sql.DB) {
    appointmentId := c.Param("id")

    // -----------------------------
    // 1. Fetch Appointment
    // -----------------------------
    var appointment structs.Appointment

    apptQuery := `SELECT id, pet_id, doctor_id, status, appointment_datetime, notes,
                	active_status, created_at, created_by, modified_at, modified_by
                  	FROM "Appointments"
                    WHERE id=$1 AND active_status=1`

    err := db.QueryRow(apptQuery, appointmentId).Scan(
        &appointment.Id, &appointment.PetId, &appointment.DoctorId, &appointment.Status,
        &appointment.AppointmentDatetime, &appointment.Notes, &appointment.ActiveStatus,
        &appointment.CreatedAt, &appointment.CreatedBy, &appointment.ModifiedAt, &appointment.ModifiedBy,
    )
    if err != nil {
        log.Println("Error fetching appointment:", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Appointment not found"})
        return
    }

    // -----------------------------
    // 2. Fetch Medical Record (1:1)
    // -----------------------------
    var medicalRecord structs.MedicalRecord

    mrQuery := `SELECT id, appointment_id, pet_id, diagnosis, notes,
                active_status, created_at, created_by, modified_at, modified_by
                FROM "MedicalRecords"
                WHERE appointment_id=$1 AND active_status=1`

    err = db.QueryRow(mrQuery, appointmentId).Scan(
        &medicalRecord.Id, &medicalRecord.AppointmentId, &medicalRecord.PetId,
        &medicalRecord.Diagnosis, &medicalRecord.Notes, &medicalRecord.ActiveStatus,
        &medicalRecord.CreatedAt, &medicalRecord.CreatedBy, &medicalRecord.ModifiedAt, &medicalRecord.ModifiedBy,
    )

    hasMedicalRecord := (err == nil)

    // -----------------------------
    // 3. Fetch Treatments (N)
    // -----------------------------
    var treatments []structs.Treatment
    totalCost := 0

    if hasMedicalRecord {
        tQuery := `SELECT id, medicalrecord_id, doctor_id, description, cost,
                	active_status, created_at, created_by, modified_at, modified_by
                    FROM "Treatments"
                    WHERE medicalrecord_id=$1 AND active_status=1
                    ORDER BY created_at DESC`

        rows, err := db.Query(tQuery, medicalRecord.Id)
        if err != nil {
            log.Println("Error fetching treatments:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch treatments"})
            return
        }
        defer rows.Close()

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
            totalCost += t.Cost
        }
    }

    // -----------------------------
    // 4. Build final response
    // -----------------------------
    response := gin.H{
        "appointment": appointment,
        "medical_record": func() interface{} {
            if hasMedicalRecord {
                return medicalRecord
            }
            return nil
        }(),
        "treatments":  treatments,
        "total_cost":  totalCost,
    }

    c.JSON(http.StatusOK, response)
}