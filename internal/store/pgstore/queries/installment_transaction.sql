-- name: CreateInstallmentTransaction :one
INSERT INTO installment_transactions (
    user_id,
    name,
    value,
    initial_date,
    final_date,
    category_id,
    credit_card_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetShortInstallmentTransactionByID :one
SELECT *
FROM installment_transactions
WHERE id = $1;

-- name: GetInstallmentTransactionByID :one
SELECT 
    it.id,
    it.user_id, 
    it.name, 
    it.value,
    it.initial_date, 
    it.final_date,
    it.created_at, 
    it.updated_at,

    c.id as category_id, 
    c.transaction_type as category_transaction_type, 
    c.name as category_name, 
    c.icon as category_icon,
    c.created_at as category_created_at,
    c.updated_at as category_updated_at,

    cc.id as creditcard_id, 
    cc.name as creditcard_name, 
    cc.first_four_numbers as creditcard_first_four_numbers, 
    cc.credit_limit as creditcard_credit_limit, 
    cc.close_day as creditcard_close_day, 
    cc.expire_day as creditcard_expire_day, 
    cc.background_color as creditcard_background_color, 
    cc.text_color as creditcard_text_color,
    cc.created_at as creditcard_created_at,
    cc.updated_at as creditcard_updated_at
FROM installment_transactions it
LEFT JOIN categories c ON it.category_id = c.id
LEFT JOIN credit_cards cc ON it.credit_card_id = cc.id
WHERE it.id = $1;

-- name: ListInstallmentTransactionsByUserIDPaginated :many
SELECT 
    it.id,
    it.user_id, 
    it.name, 
    it.value,
    it.initial_date, 
    it.final_date,
    it.created_at, 
    it.updated_at,

    c.id as category_id, 
    c.transaction_type as category_transaction_type, 
    c.name as category_name, 
    c.icon as category_icon,
    c.created_at as category_created_at,
    c.updated_at as category_updated_at,

    cc.id as creditcard_id, 
    cc.name as creditcard_name, 
    cc.first_four_numbers as creditcard_first_four_numbers, 
    cc.credit_limit as creditcard_credit_limit, 
    cc.close_day as creditcard_close_day, 
    cc.expire_day as creditcard_expire_day, 
    cc.background_color as creditcard_background_color, 
    cc.text_color as creditcard_text_color,
    cc.created_at as creditcard_created_at,
    cc.updated_at as creditcard_updated_at,

    COUNT(*) OVER() as total_count
FROM installment_transactions it
LEFT JOIN categories c ON it.category_id = c.id
LEFT JOIN credit_cards cc ON it.credit_card_id = cc.id
WHERE it.user_id = $1
ORDER BY it.initial_date ASC, it.final_date ASC
LIMIT $2 OFFSET $3;

-- name: UpdateInstallmentTransaction :one
UPDATE installment_transactions
SET
    name = $2,
    value = $3,
    initial_date = $4,
    final_date = $5,
    category_id = $6,
    credit_card_id = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteInstallmentTransaction :exec
DELETE FROM installment_transactions
WHERE id = $1;