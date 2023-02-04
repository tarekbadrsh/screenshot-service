-- Table: public."screenshot"

-- DROP TABLE public."screenshot";

CREATE TABLE public.screenshot
(
    "id" SERIAL PRIMARY KEY,
    "url" TEXT NOT NULL,
    "url_hash" VARCHAR (32)  NOT NULL,
    "is_image_created" BOOLEAN,
    "image_path" TEXT,
    "created_at" BIGINT
);