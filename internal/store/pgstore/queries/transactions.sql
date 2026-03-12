-- name: CreateTransaction :one
INSERT INTO transactions (
    user_id,
    name,
    date,
    value,
    paid,
    category_id,
    credit_card_id,
    monthly_transactions_id,
    annual_transactions_id,
    installment_transactions_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING id, name, date, value, paid, created_at, updated_at;

-- name: GetTransactionByID :one
SELECT 
    t.id,
    t.user_id, 
    t.name, 
    t.date, 
    t.value,
    t.paid,
    t.created_at, 
    t.updated_at,

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

    mt.id as monthly_transactions_id, 
    mt.day as monthly_transactions_day,
    mt.name as monthly_transactions_name,
    mt.value as monthly_transactions_value,
    mt.created_at as monthly_transactions_created_at,
    mt.updated_at as monthly_transactions_updated_at,

    at.id as annual_transactions_id, 
    at.month as annual_transactions_month, 
    at.day as annual_transactions_day,
    at.name as annual_transactions_name,
    at.value as annual_transactions_value,
    at.created_at as annual_transactions_created_at,
    at.updated_at as annual_transactions_updated_at,

    it.id as installment_transactions_id, 
    it.name as installment_transactions_name,
    it.value as installment_transactions_value,
    it.initial_date as installment_transactions_initial_date,  
    it.final_date as installment_transactions_final_date,
    it.created_at as installment_transactions_created_at,
    it.updated_at as installment_transactions_updated_at
FROM transactions t
LEFT JOIN categories c ON t.category_id = c.id
LEFT JOIN credit_cards cc ON t.credit_card_id = cc.id
LEFT JOIN monthly_transactions mt ON t.monthly_transactions_id = mt.id
LEFT JOIN annual_transactions at ON t.annual_transactions_id = at.id
LEFT JOIN installment_transactions it ON t.installment_transactions_id = it.id
WHERE t.id = $1;

-- name: ListTransactionsByUserIDPaginated :many
SELECT 
    t.id,
    t.user_id, 
    t.name, 
    t.date, 
    t.value,
    t.paid, 
    t.created_at, 
    t.updated_at,

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

    mt.id as monthly_transactions_id, 
    mt.day as monthly_transactions_day,
    mt.name as monthly_transactions_name,
    mt.value as monthly_transactions_value,
    mt.created_at as monthly_transactions_created_at,
    mt.updated_at as monthly_transactions_updated_at,

    at.id as annual_transactions_id, 
    at.month as annual_transactions_month, 
    at.day as annual_transactions_day,
    at.name as annual_transactions_name,
    at.value as annual_transactions_value,
    at.created_at as annual_transactions_created_at,
    at.updated_at as annual_transactions_updated_at,

    it.id as installment_transactions_id, 
    it.name as installment_transactions_name,
    it.value as installment_transactions_value,
    it.initial_date as installment_transactions_initial_date,  
    it.final_date as installment_transactions_final_date,
    it.created_at as installment_transactions_created_at,
    it.updated_at as installment_transactions_updated_at,

    COUNT(*) OVER() as total_count
FROM transactions t
LEFT JOIN categories c ON t.category_id = c.id
LEFT JOIN credit_cards cc ON t.credit_card_id = cc.id
LEFT JOIN monthly_transactions mt ON t.monthly_transactions_id = mt.id
LEFT JOIN annual_transactions at ON t.annual_transactions_id = at.id
LEFT JOIN installment_transactions it ON t.installment_transactions_id = it.id
WHERE t.user_id = $1
ORDER BY t.date DESC
LIMIT $2 OFFSET $3;

-- name: ListTransactionsByUserAndDate :many
SELECT 
    t.id,
    t.user_id, 
    t.name, 
    t.date, 
    t.value,
    t.paid, 
    t.created_at, 
    t.updated_at,

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

    mt.id as monthly_transactions_id, 
    mt.day as monthly_transactions_day,
    mt.name as monthly_transactions_name,
    mt.value as monthly_transactions_value,
    mt.created_at as monthly_transactions_created_at,
    mt.updated_at as monthly_transactions_updated_at,

    at.id as annual_transactions_id, 
    at.month as annual_transactions_month, 
    at.day as annual_transactions_day,
    at.name as annual_transactions_name,
    at.value as annual_transactions_value,
    at.created_at as annual_transactions_created_at,
    at.updated_at as annual_transactions_updated_at,

    it.id as installment_transactions_id, 
    it.name as installment_transactions_name,
    it.value as installment_transactions_value,
    it.initial_date as installment_transactions_initial_date,  
    it.final_date as installment_transactions_final_date,
    it.created_at as installment_transactions_created_at,
    it.updated_at as installment_transactions_updated_at,

    COUNT(*) OVER() as total_count
FROM transactions t
LEFT JOIN categories c ON t.category_id = c.id
LEFT JOIN credit_cards cc ON t.credit_card_id = cc.id
LEFT JOIN monthly_transactions mt ON t.monthly_transactions_id = mt.id
LEFT JOIN annual_transactions at ON t.annual_transactions_id = at.id
LEFT JOIN installment_transactions it ON t.installment_transactions_id = it.id
WHERE t.user_id = $1
  AND t.date >= $2
  AND t.date <= $3
ORDER BY t.date DESC
LIMIT $4 OFFSET $5;

-- name: ListTransactionsByUserAndMonthYearPaginated :many
SELECT 
    t.id,
    t.user_id, 
    t.name, 
    t.date, 
    t.value,
    t.paid, 
    t.created_at, 
    t.updated_at,

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

    mt.id as monthly_transactions_id, 
    mt.day as monthly_transactions_day,
    mt.name as monthly_transactions_name,
    mt.value as monthly_transactions_value,
    mt.created_at as monthly_transactions_created_at,
    mt.updated_at as monthly_transactions_updated_at,

    at.id as annual_transactions_id, 
    at.month as annual_transactions_month, 
    at.day as annual_transactions_day,
    at.name as annual_transactions_name,
    at.value as annual_transactions_value,
    at.created_at as annual_transactions_created_at,
    at.updated_at as annual_transactions_updated_at,

    it.id as installment_transactions_id, 
    it.name as installment_transactions_name,
    it.value as installment_transactions_value,
    it.initial_date as installment_transactions_initial_date,  
    it.final_date as installment_transactions_final_date,
    it.created_at as installment_transactions_created_at,
    it.updated_at as installment_transactions_updated_at,

    COUNT(*) OVER() as total_count
FROM transactions t
LEFT JOIN categories c ON t.category_id = c.id
LEFT JOIN credit_cards cc ON t.credit_card_id = cc.id
LEFT JOIN monthly_transactions mt ON t.monthly_transactions_id = mt.id
LEFT JOIN annual_transactions at ON t.annual_transactions_id = at.id
LEFT JOIN installment_transactions it ON t.installment_transactions_id = it.id
WHERE t.user_id = sqlc.arg(user_id)
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
ORDER BY t.date DESC
LIMIT sqlc.arg(page_limit) OFFSET sqlc.arg(page_offset);

-- name: UpdateTransaction :one
UPDATE transactions
SET
    name = $2,
    date = $3,
    value = $4,
    paid = $5,
    category_id = $6,
    credit_card_id = $7,
    monthly_transactions_id = $8,
    annual_transactions_id = $9,
    installment_transactions_id = $10,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: PayTransaction :exec
UPDATE transactions
SET
    paid = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteTransaction :exec
DELETE FROM transactions
WHERE id = $1;

-- name: HasTransactionsByCategory :one
SELECT EXISTS(
    SELECT 1 FROM transactions
    WHERE category_id = $1
);

-- name: HasTransactionsByCreditCard :one
SELECT EXISTS(
    SELECT 1 FROM transactions
    WHERE credit_card_id = $1
);