package structs

import (
	"time"

	"github.com/google/uuid"
)

// USERS
type User struct {
    Id           uuid.UUID `json:"id"`
    Name         string    `json:"name"`
    Email        string    `json:"email"`
    Phone        string    `json:"phone"`
    PasswordHash string    `json:"password_hash"`
    Role         string    `json:"role"`          // staff, doctor
    ActiveStatus int       `json:"active_status"` // 1 = aktif
    CreatedAt    time.Time `json:"created_at"`
    CreatedBy    string    `json:"created_by"`
    ModifiedAt   time.Time `json:"modified_at"`
    ModifiedBy   string    `json:"modified_by"`
}

// PETS
type Pet struct {
    Id          	uuid.UUID 	`json:"id"`
    Name        	string    	`json:"name"`
    Species     	string    	`json:"species"`
    Breed       	string    	`json:"breed"`
    Gender      	string    	`json:"gender"`
    BirthDate   	string 		`json:"birth_date"`
    OwnerName   	string    	`json:"owner_name"`
    OwnerPhone  	string    	`json:"owner_phone"`
    ActiveStatus 	int      	`json:"active_status"`
    CreatedAt   	time.Time 	`json:"created_at"`
    CreatedBy   	string    	`json:"created_by"`
    ModifiedAt  	time.Time 	`json:"modified_at"`
    ModifiedBy  	string    	`json:"modified_by"`
}

// APPOINTMENTS
type Appointment struct {
    Id                  uuid.UUID `json:"id"`
    PetId               uuid.UUID `json:"pet_id"`
    DoctorId            uuid.UUID `json:"doctor_id"`
    Status              string    `json:"status"` // pending, cancelled, completed
    AppointmentDatetime time.Time `json:"appointment_datetime"`
    Notes               string    `json:"notes"`
    ActiveStatus        int       `json:"active_status"`
    CreatedAt           time.Time `json:"created_at"`
    CreatedBy           string    `json:"created_by"`
    ModifiedAt          time.Time `json:"modified_at"`
    ModifiedBy          string    `json:"modified_by"`
}

// MEDICAL RECORDS
type MedicalRecord struct {
    Id            uuid.UUID `json:"id"`
    AppointmentId uuid.UUID `json:"appointment_id"`
    PetId         uuid.UUID `json:"pet_id"`
    Diagnosis     string    `json:"diagnosis"`
    Notes         string    `json:"notes"`
    ActiveStatus  int       `json:"active_status"`
    CreatedAt     time.Time `json:"created_at"`
    CreatedBy     string    `json:"created_by"`
    ModifiedAt    time.Time `json:"modified_at"`
    ModifiedBy    string    `json:"modified_by"`
}

// TREATMENTS
type Treatment struct {
    Id             	uuid.UUID `json:"id"`
    MedicalRecordId uuid.UUID `json:"medicalrecord_id"`
    DoctorId       	uuid.UUID `json:"doctor_id"`
    Description    	string    `json:"description"`
    Cost           	int       `json:"cost"`
    ActiveStatus   	int       `json:"active_status"`
    CreatedAt      	time.Time `json:"created_at"`
    CreatedBy      	string    `json:"created_by"`
    ModifiedAt     	time.Time `json:"modified_at"`
    ModifiedBy     	string    `json:"modified_by"`
}