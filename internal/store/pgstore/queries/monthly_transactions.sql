-- name: CreateMonthlyTransaction :one
INSERT INTO monthly_transactions (
    user_id,
    name,
    value,
    day,
    category_id,
    credit_card_id
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetShortMonthlyTransactionByID :one
SELECT *
FROM monthly_transactions
WHERE id = $1;

-- name: GetMonthlyTransactionByID :one
SELECT 
    mt.id,
    mt.user_id, 
    mt.name, 
    mt.value,
    mt.day, 
    mt.created_at, 
    mt.updated_at,

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
FROM monthly_transactions mt
LEFT JOIN categories c ON mt.category_id = c.id
LEFT JOIN credit_cards cc ON mt.credit_card_id = cc.id
WHERE mt.id = $1;

-- name: ListMonthlyTransactionsByUserIDPaginated :many
SELECT 
    mt.id,
    mt.user_id, 
    mt.name, 
    mt.value,
    mt.day, 
    mt.created_at, 
    mt.updated_at,

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
FROM monthly_transactions mt
LEFT JOIN categories c ON mt.category_id = c.id
LEFT JOIN credit_cards cc ON mt.credit_card_id = cc.id
WHERE mt.user_id = $1
ORDER BY mt.day ASC
LIMIT $2 OFFSET $3;

-- name: ListMonthlyTransactionsIDs :many
SELECT DISTINCT t.monthly_transactions_id::uuid
FROM transactions t
LEFT JOIN categories c ON t.category_id = c.id
LEFT JOIN credit_cards cc ON t.credit_card_id = cc.id
WHERE t.user_id = sqlc.arg(user_id)
  AND t.monthly_transactions_id IS NOT NULL
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
ORDER BY t.monthly_transactions_id;

-- name: UpdateMonthlyTransaction :one
UPDATE monthly_transactions
SET
    name = $2,
    value = $3,
    day = $4,
    category_id = $5,
    credit_card_id = $6,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteMonthlyTransaction :exec
DELETE FROM monthly_transactions
WHERE id = $1;