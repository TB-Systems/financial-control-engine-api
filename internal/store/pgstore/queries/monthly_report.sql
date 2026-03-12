-- name: GetMonthlyBalance :one
-- Parâmetros: $1 = user_id, $2 = ano, $3 = mês (1-12)
WITH credit_card_transactions AS (
    -- Transações de crédito: considera o ciclo de faturamento do cartão
    -- Se mês = 9 e close_day = 25, pega transações de 26/07 até 25/08
    SELECT 
        t.value,
        c.transaction_type
    FROM transactions t
    INNER JOIN categories c ON t.category_id = c.id
    INNER JOIN credit_cards cc ON t.credit_card_id = cc.id
    WHERE t.user_id = $1
      AND c.transaction_type = 2 -- Credit
      AND t.date > (
          -- Data inicial: dia após o fechamento, 2 meses antes
          make_date($2::int, $3::int, 1) - INTERVAL '2 months' + (cc.close_day || ' days')::INTERVAL
      )
      AND t.date <= (
          -- Data final: dia de fechamento, 1 mês antes
          make_date($2::int, $3::int, 1) - INTERVAL '1 month' + (cc.close_day || ' days')::INTERVAL
      )
),
debit_income_transactions AS (
    -- Transações de débito e income: considera o mês selecionado diretamente
    SELECT 
        t.value,
        c.transaction_type
    FROM transactions t
    INNER JOIN categories c ON t.category_id = c.id
    WHERE t.user_id = $1
      AND t.credit_card_id IS NULL
      AND c.transaction_type IN (0, 1) -- Income (0) ou Debit (1)
      AND EXTRACT(YEAR FROM t.date) = $2
      AND EXTRACT(MONTH FROM t.date) = $3
),
all_transactions AS (
    SELECT value, transaction_type FROM credit_card_transactions
    UNION ALL
    SELECT value, transaction_type FROM debit_income_transactions
)
SELECT 
    COALESCE(SUM(CASE WHEN transaction_type = 0 THEN value ELSE 0 END), 0)::NUMERIC(15,2) as total_income,
    COALESCE(SUM(CASE WHEN transaction_type = 1 THEN value ELSE 0 END), 0)::NUMERIC(15,2) as total_debit,
    COALESCE(SUM(CASE WHEN transaction_type = 2 THEN value ELSE 0 END), 0)::NUMERIC(15,2) as total_credit,
    (
        COALESCE(SUM(CASE WHEN transaction_type = 0 THEN value ELSE 0 END), 0) -
        COALESCE(SUM(CASE WHEN transaction_type = 1 THEN value ELSE 0 END), 0) -
        COALESCE(SUM(CASE WHEN transaction_type = 2 THEN value ELSE 0 END), 0)
    )::NUMERIC(15,2) as balance
FROM all_transactions;

-- name: GetCategoriesSpending :many
WITH credit_card_transactions AS (
    SELECT 
        t.value,
        t.category_id
    FROM transactions t
    INNER JOIN categories c ON t.category_id = c.id
    INNER JOIN credit_cards cc ON t.credit_card_id = cc.id
    WHERE t.user_id = $1
      AND c.transaction_type = 2 -- Credit
      AND t.date > (
          make_date($2::int, $3::int, 1) - INTERVAL '2 months' + (cc.close_day || ' days')::INTERVAL
      )
      AND t.date <= (
          make_date($2::int, $3::int, 1) - INTERVAL '1 month' + (cc.close_day || ' days')::INTERVAL
      )
),
debit_income_transactions AS (
    SELECT 
        t.value,
        t.category_id
    FROM transactions t
    INNER JOIN categories c ON t.category_id = c.id
    WHERE t.user_id = $1
      AND t.credit_card_id IS NULL
      AND c.transaction_type IN (0, 1) -- Income ou Debit
      AND EXTRACT(YEAR FROM t.date) = $2
      AND EXTRACT(MONTH FROM t.date) = $3
),
all_transactions_cat AS (
    SELECT value, category_id FROM credit_card_transactions
    UNION ALL
    SELECT value, category_id FROM debit_income_transactions
)
SELECT 
    c.id as category_id,
    c.name as category_name,
    c.icon as category_icon,
    c.transaction_type as category_transaction_type,
    COUNT(a.value)::int as transactions_count,
    COALESCE(SUM(a.value), 0)::NUMERIC(15,2) as total_spent
FROM categories c
LEFT JOIN all_transactions_cat a ON a.category_id = c.id
WHERE c.user_id = $1
GROUP BY c.id, c.name, c.icon, c.transaction_type
ORDER BY total_spent DESC;

-- name: GetCreditCardsSpending :many
WITH credit_card_transactions AS (
    SELECT 
        t.value,
        t.credit_card_id
    FROM transactions t
    INNER JOIN credit_cards cc ON t.credit_card_id = cc.id
    WHERE t.user_id = $1
      AND t.date > (
          make_date($2::int, $3::int, 1) - INTERVAL '2 months' + (cc.close_day || ' days')::INTERVAL
      )
      AND t.date <= (
          make_date($2::int, $3::int, 1) - INTERVAL '1 month' + (cc.close_day || ' days')::INTERVAL
      )
)
SELECT 
    cc.id as credit_card_id,
    cc.name as credit_card_name,
    cc.first_four_numbers as credit_card_first_four_numbers,
    cc.credit_limit as credit_card_limit,
    cc.close_day as credit_card_close_day,
    cc.expire_day as credit_card_expire_day,
    cc.background_color as credit_card_background_color,
    cc.text_color as credit_card_text_color,
    COUNT(t.value)::int as transactions_count,
    COALESCE(SUM(t.value), 0)::NUMERIC(15,2) as total_spent
FROM credit_cards cc
LEFT JOIN credit_card_transactions t ON t.credit_card_id = cc.id
WHERE cc.user_id = $1
GROUP BY cc.id, cc.name, cc.first_four_numbers, cc.credit_limit, cc.close_day, cc.expire_day, cc.background_color, cc.text_color
ORDER BY total_spent DESC;
