CREATE TABLE public.provinces (
	id SERIAL PRIMARY KEY,
	code VARCHAR(10) UNIQUE NOT NULL,
	"name" VARCHAR(100) UNIQUE NOT NULL,
	lat VARCHAR(100) NULL,
	long VARCHAR(100) NULL,
	"status" BOOLEAN DEFAULT false NOT NULL,
	created_date TIMESTAMPTZ(0) NOT NULL,
    created_by VARCHAR(10) NOT NULL,
	updated_date TIMESTAMPTZ(0) NULL,
    updated_by VARCHAR(10) NULL,
	deleted_date TIMESTAMPTZ(0) NULL,
    deleted_by VARCHAR(10) NULL,
    "version" INTEGER NOT NULL
);