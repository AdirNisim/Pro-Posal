-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users"(
    "id" UUID NOT NULL,
    "first_name" TEXT NOT NULL,
    "last_name" TEXT NOT NULL,
    "phone" TEXT NOT NULL,
    "email" TEXT NOT NULL,
    "email_hash" TEXT NOT NULL,
    "password_hash" TEXT NOT NULL,
    "invited_by" UUID NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "users" ADD CONSTRAINT "users_email_hash_unique" UNIQUE("email_hash");
ALTER TABLE
    "users" ADD PRIMARY KEY("id");
CREATE TABLE "offers"(
    "id" UUID NOT NULL,
    "created_by" UUID NOT NULL,
    "customer_id" UUID NOT NULL,
    "company_id" UUID NOT NULL,
    "contract_template_id" UUID NOT NULL,
    "arguments" jsonb NOT NULL,
    "finalized_offer" TEXT NULL,
    "finalized_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "sent_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "opened_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "accepted_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "rejected_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "rejection_reason" TEXT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX "offers_customer_id_index" ON
    "offers"("customer_id");
ALTER TABLE
    "offers" ADD PRIMARY KEY("id");
CREATE TABLE "companies"(
    "id" UUID NOT NULL,
    "name" TEXT NOT NULL,
    "contact_id" UUID NOT NULL,
    "address" TEXT NOT NULL,
    "logo_base64" TEXT NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "companies" ADD PRIMARY KEY("id");
CREATE TABLE "permissions"(
    "id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "company_id" UUID NOT NULL,
    "role" TEXT NOT NULL,
    "contract_id" TEXT NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "permissions" ADD PRIMARY KEY("id");
CREATE TABLE "contract_templates"(
    "id" UUID NOT NULL,
    "name" TEXT NOT NULL,
    "company_id" UUID NOT NULL,
    "template" TEXT NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX "contract_templates_company_id_index" ON
    "contract_templates"("company_id");
ALTER TABLE
    "contract_templates" ADD PRIMARY KEY("id");
CREATE TABLE "categories"(
    "id" UUID NOT NULL,
    "company_id" UUID NOT NULL,
    "category_id" UUID NULL,
    "description" TEXT NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX "categories_category_id_index" ON
    "categories"("category_id");
ALTER TABLE
    "categories" ADD PRIMARY KEY("id");
ALTER TABLE
    "categories" ADD CONSTRAINT "categories_company_id_foreign" FOREIGN KEY("company_id") REFERENCES "companies"("id");
ALTER TABLE
    "offers" ADD CONSTRAINT "offers_contract_template_id_foreign" FOREIGN KEY("contract_template_id") REFERENCES "contract_templates"("id");
ALTER TABLE
    "companies" ADD CONSTRAINT "companies_contact_id_foreign" FOREIGN KEY("contact_id") REFERENCES "users"("id");
ALTER TABLE
    "offers" ADD CONSTRAINT "offers_created_by_foreign" FOREIGN KEY("created_by") REFERENCES "users"("id");
ALTER TABLE
    "permissions" ADD CONSTRAINT "permissions_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users"("id");
ALTER TABLE
    "offers" ADD CONSTRAINT "offers_customer_id_foreign" FOREIGN KEY("customer_id") REFERENCES "users"("id");
ALTER TABLE
    "offers" ADD CONSTRAINT "offers_company_id_foreign" FOREIGN KEY("company_id") REFERENCES "companies"("id");
ALTER TABLE
    "contract_templates" ADD CONSTRAINT "contract_templates_company_id_foreign" FOREIGN KEY("company_id") REFERENCES "companies"("id");
ALTER TABLE
    "permissions" ADD CONSTRAINT "permissions_company_id_foreign" FOREIGN KEY("company_id") REFERENCES "companies"("id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "categories";
DROP TABLE "contract_templates";
DROP TABLE "permissions";
DROP TABLE "companies";
DROP INDEX "offers_customer_id_index";
DROP TABLE "offers";
DROP TABLE "users";
-- +goose StatementEnd
