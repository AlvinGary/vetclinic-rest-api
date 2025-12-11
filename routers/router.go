package routers

import (
	"database/sql"
	"vetclinic-rest-api/controllers"
	"vetclinic-rest-api/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, db *sql.DB) {
	usersGroup := router.Group("api/users")
	{
		// register user
		usersGroup.POST("/register", func(c *gin.Context) {
			controllers.RegisterUser(c, db)
		})
		// login user
		usersGroup.POST("/login", func(c *gin.Context) {
			controllers.LoginUser(c, db)
		})
		// Get user profile (all roles)
		usersGroup.GET("/:id/profile", middleware.JWTAuth("Staff","Doctor","Admin"), func(c *gin.Context) {
			controllers.FetchProfile(c, db)
		})
		// Get users by role (all roles)
		usersGroup.GET("/role/:role", middleware.JWTAuth("Staff", "Doctor", "Admin"), func(c *gin.Context) {
			controllers.GetUserByRole(c, db)
		})
		// Update user (all roles)
		usersGroup.PUT("/:id/update", middleware.JWTAuth("Staff","Doctor","Admin"), func(c *gin.Context) {
			controllers.UpdateUser(c, db)
		})
		// Change password (all roles)
		usersGroup.PUT("/:id/change-password", middleware.JWTAuth("Staff","Doctor","Admin"), func(c *gin.Context) {
			controllers.ChangePassword(c, db)
		})
		// Update role (Admin only)
		usersGroup.PUT("/:id/role", middleware.JWTAuth("Admin"), func(c *gin.Context) {
			controllers.UpdateRole(c, db)
		})
		// Users soft delete (Admin only)
		usersGroup.PUT("/:id/active-status", middleware.JWTAuth("Admin"), func(c *gin.Context) {
			controllers.UpdateUserActiveStatus(c, db)
		})
	}
	petsGroup := router.Group("api/pets")
	{
		// Get pets data (all roles)
		petsGroup.GET("/:id/profile", middleware.JWTAuth("Staff","Doctor", "Admin"), func(c *gin.Context) {
			controllers.FetchPetProfile(c, db)
		})
		// Get pets data based on owner name and phone (all roles)
		petsGroup.GET("/by-owner/:owner_name/:owner_phone", middleware.JWTAuth("Staff","Doctor", "Admin"), func(c *gin.Context) {
			controllers.FetchPetsByOwner(c, db)
		})
		// Create new pets data (Staff and Admin)
		petsGroup.POST("", middleware.JWTAuth("Staff", "Admin"), func(c *gin.Context) {
			controllers.CreatePet(c, db)
		})
		// Update pets data (Staff and Admin)
		petsGroup.PUT("/:id", middleware.JWTAuth("Staff", "Admin"), func(c *gin.Context) {
			controllers.UpdatePet(c, db)
		})
		// Pets soft delete (Staff and Admin)
		petsGroup.PUT("/:id/active-status", middleware.JWTAuth("Staff","Admin"), func(c *gin.Context) {
			controllers.UpdatePetActiveStatus(c, db)
		})
	}
	appointmentsGroup := router.Group("api/appointments")
	{
		// Create appointment (Staff and Admin)
		appointmentsGroup.POST("", middleware.JWTAuth("Staff", "Admin"), func(c *gin.Context) {
			controllers.CreateAppointment(c, db)
		})
		// Fetch appointment by ID (all roles)
		appointmentsGroup.GET("/:id", middleware.JWTAuth("Staff","Doctor", "Admin"), func(c *gin.Context) {
			controllers.FetchAppointment(c, db)
		})
		// Update appointment details (Staff and Admin)
		appointmentsGroup.PUT("/:id", middleware.JWTAuth("Staff", "Admin"), func(c *gin.Context) {
			controllers.UpdateAppointment(c, db)
		})
		// Update appointment status (all roles)
		appointmentsGroup.PUT("/:id/status", middleware.JWTAuth("Staff","Doctor", "Admin"), func(c *gin.Context) {
			controllers.UpdateAppointmentStatus(c, db)
		})
		// Soft delete appointment (Staff and Admin)
		appointmentsGroup.PUT("/:id/active-status", middleware.JWTAuth("Staff","Admin"), func(c *gin.Context) {
			controllers.UpdateAppointmentActiveStatus(c, db)
		})
		// Get Appointments by pet id (all roles)
		appointmentsGroup.GET("/pet/:pet_id", middleware.JWTAuth("Staff","Doctor", "Admin"), func(c *gin.Context) {
			controllers.GetAppointmentsByPetId(c, db)
		})
		// Get Appointments by doctor (all roles)
		appointmentsGroup.GET("/doctor/:doctor_id", middleware.JWTAuth("Staff","Doctor", "Admin"), func(c *gin.Context) {
			controllers.GetAppointmentsByDoctorId(c, db)
		})
		// Get Appointments by apointment date (all roles)
		appointmentsGroup.GET("/date/:date", middleware.JWTAuth("Staff","Doctor", "Admin"), func(c *gin.Context) {
			controllers.GetAppointmentsByAppointmentDate(c, db)
		})
		// Get Appointment Full Detail by id (all roles)
		appointmentsGroup.GET("/:id/full", middleware.JWTAuth("Staff","Doctor", "Admin"), func(c *gin.Context) {
			controllers.GetFullAppointmentDetail(c, db)
		})
	}
	medicalGroup := router.Group("api/medical-records")
	{
		// Create new Medical Records (Doctor and Admin)
		medicalGroup.POST("", middleware.JWTAuth("Doctor", "Admin"), func(c *gin.Context) {
			controllers.CreateMedicalRecord(c, db)
		})
		// Get Medical Records based on appointment id (all roles)
		medicalGroup.GET("/appointment/:appointment_id", middleware.JWTAuth("Doctor","Staff", "Admin"), func(c *gin.Context) {
			controllers.GetMedicalRecordByAppointmentId(c, db)
		})
		// Update Medical Records (Doctor and Admin)
		medicalGroup.PUT("/:id", middleware.JWTAuth("Doctor", "Admin"), func(c *gin.Context) {
			controllers.UpdateMedicalRecord(c, db)
		})
		// Soft delete Medical Records (Admin)
		medicalGroup.PUT("/:id/active-status", middleware.JWTAuth("Admin"), func(c *gin.Context) {
			controllers.UpdateMedicalRecordActiveStatus(c, db)
		})
	}
	treatmentGroup := router.Group("api/treatments")
	{
		// Create new treatment (Doctor and Admin)
		treatmentGroup.POST("", middleware.JWTAuth("Doctor", "Admin"), func(c *gin.Context) {
			controllers.CreateTreatment(c, db)
		})
		// Get Treatement based on medical records id (all roles)
		treatmentGroup.GET("/medicalrecord/:medicalrecord_id", middleware.JWTAuth("Doctor","Staff", "Admin"), func(c *gin.Context) {
			controllers.GetTreatmentsByMedicalRecordId(c, db)
		})
		// Update treatment (Doctor and Admin)
		treatmentGroup.PUT("/:id", middleware.JWTAuth("Doctor", "Admin"), func(c *gin.Context) {
			controllers.UpdateTreatment(c, db)
		})
		// Soft delete treatment (Doctor and Admin)
		treatmentGroup.PUT("/:id/active-status", middleware.JWTAuth("Doctor", "Admin"), func(c *gin.Context) {
			controllers.UpdateTreatmentActiveStatus(c, db)
		})
	}
}