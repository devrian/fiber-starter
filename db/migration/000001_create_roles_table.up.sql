CREATE TABLE public.roles (
	id SERIAL PRIMARY KEY,
	code VARCHAR(10) UNIQUE NOT NULL,
	"name" VARCHAR(50) UNIQUE NOT NULL,
	slug VARCHAR(50) UNIQUE NOT NULL,
	"status" BOOLEAN DEFAULT false NOT NULL,
	created_date TIMESTAMPTZ(0) NOT NULL,
    created_by VARCHAR(10) NOT NULL,
	updated_date TIMESTAMPTZ(0) NULL,
    updated_by VARCHAR(10) NULL,
	deleted_date TIMESTAMPTZ(0) NULL,
    deleted_by VARCHAR(10) NULL,
    "version" INTEGER NOT NULL
);