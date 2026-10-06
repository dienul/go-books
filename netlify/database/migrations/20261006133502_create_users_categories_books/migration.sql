CREATE TABLE "books" (
	"id" bigserial PRIMARY KEY,
	"title" varchar(255) NOT NULL,
	"description" varchar(1000) NOT NULL,
	"image_url" varchar(1000) NOT NULL,
	"release_year" integer NOT NULL,
	"price" integer NOT NULL,
	"total_page" integer NOT NULL,
	"thickness" varchar(10) NOT NULL,
	"category_id" bigint NOT NULL,
	"created_at" timestamp with time zone DEFAULT now() NOT NULL,
	"created_by" varchar(100) NOT NULL,
	"modified_at" timestamp with time zone DEFAULT now() NOT NULL,
	"modified_by" varchar(100) NOT NULL,
	CONSTRAINT "books_release_year_check" CHECK ("release_year" BETWEEN 1980 AND 2024),
	CONSTRAINT "books_price_check" CHECK ("price" >= 0),
	CONSTRAINT "books_total_page_check" CHECK ("total_page" > 0),
	CONSTRAINT "books_thickness_check" CHECK ("thickness" IN ('tipis', 'tebal'))
);
--> statement-breakpoint
CREATE TABLE "categories" (
	"id" bigserial PRIMARY KEY,
	"name" varchar(100) NOT NULL UNIQUE,
	"created_at" timestamp with time zone DEFAULT now() NOT NULL,
	"created_by" varchar(100) NOT NULL,
	"modified_at" timestamp with time zone DEFAULT now() NOT NULL,
	"modified_by" varchar(100) NOT NULL
);
--> statement-breakpoint
CREATE TABLE "users" (
	"id" bigserial PRIMARY KEY,
	"username" varchar(100) NOT NULL UNIQUE,
	"password" varchar(255) NOT NULL,
	"created_at" timestamp with time zone DEFAULT now() NOT NULL,
	"created_by" varchar(100) NOT NULL,
	"modified_at" timestamp with time zone DEFAULT now() NOT NULL,
	"modified_by" varchar(100) NOT NULL
);
--> statement-breakpoint
CREATE INDEX "idx_books_category_id" ON "books" ("category_id");--> statement-breakpoint
ALTER TABLE "books" ADD CONSTRAINT "books_category_id_categories_id_fkey" FOREIGN KEY ("category_id") REFERENCES "categories"("id") ON DELETE RESTRICT;