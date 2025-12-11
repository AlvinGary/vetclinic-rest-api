-- +migrate Up

---------------------------------------------------------
-- USERS
---------------------------------------------------------
CREATE TABLE IF NOT EXISTS "Users"
(
    id uuid NOT NULL,
    name character varying(100) NOT NULL,
    email character varying(100) NOT NULL,
    phone character varying(20),
    password_hash text NOT NULL,
    role character varying(20) NOT NULL, -- staff, doctor
    active_status integer NOT NULL DEFAULT 1,
    created_at timestamp(0) without time zone NOT NULL,
    created_by character varying(50) NOT NULL,
    modified_at timestamp(0) without time zone,
    modified_by character varying(50),
    CONSTRAINT "Users_pkey" PRIMARY KEY (id)
);

---------------------------------------------------------
-- PETS
---------------------------------------------------------
CREATE TABLE IF NOT EXISTS "Pets"
(
    id uuid NOT NULL,
    name character varying(100) NOT NULL,
    species character varying(50),
    breed character varying(100),
    gender character varying(10),
    birth_date date,
    owner_name character varying(100),
    owner_phone character varying(20),
    active_status integer NOT NULL DEFAULT 1,
    created_at timestamp(0) without time zone NOT NULL,
    created_by character varying(50) NOT NULL,
    modified_at timestamp(0) without time zone,
    modified_by character varying(50),
    CONSTRAINT "Pets_pkey" PRIMARY KEY (id)
);

---------------------------------------------------------
-- APPOINTMENTS
---------------------------------------------------------
CREATE TABLE IF NOT EXISTS "Appointments"
(
    id uuid NOT NULL,
    pet_id uuid NOT NULL,
    doctor_id uuid,
    status character varying(20) NOT NULL, -- pending, cancelled, completed
    appointment_datetime timestamp(0) without time zone NOT NULL, -- YYYY-MM-DD HH:mm:ss
    notes text,
    active_status integer NOT NULL DEFAULT 1,
    created_at timestamp(0) without time zone NOT NULL,
    created_by character varying(50) NOT NULL,
    modified_at timestamp(0) without time zone,
    modified_by character varying(50),
    CONSTRAINT "Appointments_pkey" PRIMARY KEY (id),
    CONSTRAINT appointments_pet_id_to_pets_id FOREIGN KEY (pet_id)
        REFERENCES "Pets" (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT appointments_doctor_id_to_users_id FOREIGN KEY (doctor_id)
        REFERENCES "Users" (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

---------------------------------------------------------
-- MEDICAL RECORDS (1 per appointment)
---------------------------------------------------------
CREATE TABLE IF NOT EXISTS "MedicalRecords"
(
    id uuid NOT NULL,
    appointment_id uuid NOT NULL,
    pet_id uuid NOT NULL,
    diagnosis text NOT NULL,
    notes text,
    active_status integer NOT NULL DEFAULT 1,
    created_at timestamp(0) without time zone NOT NULL,
    created_by character varying(50) NOT NULL,
    modified_at timestamp(0) without time zone,
    modified_by character varying(50),
    CONSTRAINT "MedicalRecords_pkey" PRIMARY KEY (id),
    CONSTRAINT medicalrecords_appointment_id_to_appointments_id FOREIGN KEY (appointment_id)
        REFERENCES "Appointments" (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT medicalrecords_pet_id_to_pets_id FOREIGN KEY (pet_id)
        REFERENCES "Pets" (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);


---------------------------------------------------------
-- TREATMENTS (many per medical record)
---------------------------------------------------------
CREATE TABLE IF NOT EXISTS "Treatments"
(
    id uuid NOT NULL,
    medicalrecord_id uuid NOT NULL,
    doctor_id uuid NOT NULL,
    description text NOT NULL,
    cost integer NOT NULL,
    active_status integer NOT NULL DEFAULT 1,
    created_at timestamp(0) without time zone NOT NULL,
    created_by character varying(50) NOT NULL,
    modified_at timestamp(0) without time zone,
    modified_by character varying(50),
    CONSTRAINT "Treatments_pkey" PRIMARY KEY (id),
    CONSTRAINT treatments_medicalrecord_id_to_medicalrecords_id FOREIGN KEY (medicalrecord_id)
        REFERENCES "MedicalRecords" (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT treatments_doctor_id_to_users_id FOREIGN KEY (doctor_id)
        REFERENCES "Users" (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);
