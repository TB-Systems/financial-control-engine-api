-- name: CreateAnnualTransaction :one
INSERT INTO annual_transactions (
    user_id,
    name,
    value,
    day,
    month,
    category_id,
    credit_card_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetShortAnnualTransactionByID :one
SELECT *
FROM annual_transactions
WHERE id = $1;

-- name: GetAnnualTransactionByID :one
SELECT 
    at.id,
    at.user_id, 
    at.name, 
    at.value,
    at.day, 
    at.month,
    at.created_at, 
    at.updated_at,

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
FROM annual_transactions at
LEFT JOIN categories c ON at.category_id = c.id
LEFT JOIN credit_cards cc ON at.credit_card_id = cc.id
WHERE at.id = $1;

-- name: ListAnnualTransactionsByUserIDPaginated :many
SELECT 
    at.id,
    at.user_id, 
    at.name, 
    at.value,
    at.day, 
    at.month,
    at.created_at, 
    at.updated_at,

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
FROM annual_transactions at
LEFT JOIN categories c ON at.category_id = c.id
LEFT JOIN credit_cards cc ON at.credit_card_id = cc.id
WHERE at.user_id = $1
ORDER BY at.month ASC, at.day ASC
LIMIT $2 OFFSET $3;

-- name: ListAnnualTransactionsIDs :many
SELECT at.id
FROM annual_transactions at
LEFT JOIN categories c ON at.category_id = c.id
LEFT JOIN credit_cards cc ON at.credit_card_id = cc.id
WHERE at.user_id = sqlc.arg(user_id)
  AND NOT EXISTS (
      SELECT 1
      FROM transactions t
      WHERE t.user_id = at.user_id
        AND t.annual_transactions_id = at.id
        AND (
            (
                c.transaction_type IN (0, 1)
                AND EXTRACT(YEAR FROM t.date) = sqlc.arg(year)::int
                AND EXTRACT(MONTH FROM t.date) = sqlc.arg(month)::int
            )
            OR
            (
                c.transaction_type = 2
                AND t.credit_card_id IS NOT NULL
                AND t.date > (
                    make_date(sqlc.arg(year)::int, sqlc.arg(month)::int, 1) - INTERVAL '2 months' + (cc.close_day || ' days')::INTERVAL
                )
                AND t.date <= (
                    make_date(sqlc.arg(year)::int, sqlc.arg(month)::int, 1) - INTERVAL '1 month' + (cc.close_day || ' days')::INTERVAL
                )
            )
        )
  )
ORDER BY at.month, at.day, at.id;

-- name: UpdateAnnualTransaction :one
UPDATE annual_transactions
SET
    name = $2,
    value = $3,
    day = $4,
    month = $5,
    category_id = $6,
    credit_card_id = $7,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteAnnualTransaction :exec
DELETE FROM annual_transactions
WHERE id = $1;