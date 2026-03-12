package constants

// API
const (
	EnvDBUser      = "FINANCIAL_CONTROL_DATABASE_USER"
	EnvDBPassword  = "FINANCIAL_CONTROL_DATABASE_PASSWORD"
	EnvDBHost      = "FINANCIAL_CONTROL_DATABASE_HOST"
	EnvDBPort      = "FINANCIAL_CONTROL_DATABASE_PORT"
	EnvDBName      = "FINANCIAL_CONTROL_DATABASE_NAME"
	EnvAppPort     = "FINANCIAL_CONTROL_APP_PORT"
	DefaultAppPort = "3080"
)

// COMMONS
const (
	EmptyString        = ""
	Success            = "success"
	ID                 = "id"
	LimitText          = "limit"
	LimitDefaultString = "10"
	LimitDefault       = 10
	PageText           = "page"
	PageDefaultString  = "1"
	PageDefault        = 1
	StartDateText      = "start_date"
	EndDateText        = "end_date"
	MonthText          = "month"
	YearText           = "year"
)

// STORE ERRORS
const (
	StoreErrorNoRowsMsg = "no rows in result set"
)

// ERRORS_USER_MESSAGES
const (
	UserUnauthorized                    = "USER_UNAUTHORIZED"
	InvalidStartDate                    = "INVALID_START_DATE"
	InvalidEndDate                      = "INVALID_END_DATE"
	EndDateBeforeStartDate              = "END_DATE_BEFORE_START_DATE"
	TransactionTypeMsg                  = "TRANSACTION_TYPE_INVALID"
	TransactionTypeEmptyMsg             = "TRANSACTION_TYPE_EMPTY"
	NameEmptyMsg                        = "NAME_EMPTY"
	IconEmptyMsg                        = "ICON_EMPTY"
	NameInvalidCharsCountMsg            = "NAME_INVALID_CHARS_COUNT"
	IconInvalidCharsCountMsg            = "ICON_INVALID_CHARS_COUNT"
	LimitReachedMsg                     = "LIMIT_REACHED"
	CannotBeDeletedMsg                  = "CANNOT_BE_DELETED_BECAUSE_IT_HAS_ASSOCIATED_TRANSACTIONS"
	ValueInvalidMsg                     = "VALUE_INVALID"
	DateEmptyMsg                        = "DATE_EMPTY_OR_INVALID"
	DateInvalidMsg                      = "DATE_INVALID"
	CreditcardLimitExceededMsg          = "CREDITCARD_LIMIT_EXCEEDED"
	InvalidData                         = "INVALID_DATA"
	InvalidID                           = "INVALID_ID"
	CategoryNotFoundMsg                 = "CATEGORY_NOT_FOUND"
	CreditcardNotFoundMsg               = "CREDITCARD_NOT_FOUND"
	TransactionNotFoundMsg              = "TRANSACTION_NOT_FOUND"
	MonthlyTransactionNotFoundMsg       = "MONTHLY_TRANSACTION_NOT_FOUND"
	AnnualTransactionNotFoundMsg        = "ANNUAL_TRANSACTION_NOT_FOUND"
	InstallmentTransactionNotFoundMsg   = "INSTALLMENT_TRANSACTION_NOT_FOUND"
	DayInvalidMsg                       = "DAY_INVALID"
	MonthInvalidMsg                     = "MONTH_INVALID"
	YearInvalidMsg                      = "YEAR_INVALID"
	YearMustBe1970OrLaterMsg            = "YEAR_MUST_BE_1970_OR_LATER"
	InitialDateEmptyMsg                 = "INITIAL_DATE_EMPTY"
	FinalDateEmptyMsg                   = "FINAL_DATE_EMPTY"
	FinalDateBeforeInitialDateMsg       = "FINAL_DATE_BEFORE_INITIAL_DATE"
	InitialDateEqualsFinalDateMsg       = "INITIAL_DATE_EQUALS_FINAL_DATE"
	AnnualAndInstallmentTransactionMsg  = "ANNUAL_AND_INSTALLMENT_TRANSACTION"
	AnnualAndMonthlyTransactionMsg      = "ANNUAL_AND_MONTHLY_TRANSACTION"
	InstallmentAndMonthlyTransactionMsg = "INSTALLMENT_AND_MONTHLY_TRANSACTION"
)

// ERRORS
const (
	UserIDNotFound       = "USER_ID_NOT_FOUND"
	UserIDInvalid        = "USER_ID_INVALID"
	InternalServerError  = "INTERNAL_SERVER_ERROR"
	NilValueError        = "NIL_VALUE"
	UnsupportedTypeError = "UNSUPPORTED_TYPE"
	DecodeJsonError      = "DECODE_JSON_ERROR"
	EncodeJsonError      = "ENCODE_JSON_ERROR"
	InvalidFieldError    = "INVALID_FIELD_ERROR"
	LimitError           = "LIMIT_ERROR"
	NotFoundError        = "NOT_FOUND_ERROR"
	StoreError           = "STORE_ERROR"
	UnauthorizedError    = "UNAUTHORIZED_ERROR"
	InvalidPageParam     = "INVALID_PAGE_PARAM"
	CustomError          = "CUSTOM_ERROR"
	BadRequestError      = "BAD_REQUEST_ERROR"
)

// USER
const (
	UserID = "user_id"
)

// CREDIT CARDS
const (
	FirstFourNumbersInvalidMsg          = "FIRST_FOUR_NUMBERS_INVALID"
	LimitInvalidMsg                     = "LIMIT_INVALID"
	ClosingDayInvalidMsg                = "CLOSING_DAY_INVALID"
	ExpireDayInvalidMsg                 = "EXPIRE_DAY_INVALID"
	BackgroundColorEmptyMsg             = "BACKGROUND_COLOR_EMPTY"
	BackgroundColorInvalidCharsCountMsg = "BACKGROUND_COLOR_INVALID_CHARS_COUNT"
	TextColorEmptyMsg                   = "TEXT_COLOR_EMPTY"
	TextColorInvalidCharsCountMsg       = "TEXT_COLOR_INVALID_CHARS_COUNT"
)

// TRANSACTIONS
const (
	CreditWithoutCreditcardMsg                       = "TRANSACTION_CREDIT_WITHOUT_CREDITCARD"
	DebitOrIncomeWithCreditcardMsg                   = "TRANSACTION_DEBIT_OR_INCOME_WITH_CREDITCARD"
	AnnualTransactionWithDifferentCategoryMsg        = "ANNUAL_TRANSACTION_WITH_DIFFERENT_CATEGORY"
	AnnualTransactionWithCreditcardMsg               = "ANNUAL_TRANSACTION_WITH_CREDITCARD"
	AnnualTransactionWithoutCreditcardMsg            = "ANNUAL_TRANSACTION_WITHOUT_CREDITCARD"
	AnnualTransactionWithDifferentCreditcardMsg      = "ANNUAL_TRANSACTION_WITH_DIFFERENT_CREDITCARD"
	InstallmentTransactionWithDifferentCategoryMsg   = "INSTALLMENT_TRANSACTION_WITH_DIFFERENT_CATEGORY"
	InstallmentTransactionWithCreditcardMsg          = "INSTALLMENT_TRANSACTION_WITH_CREDITCARD"
	InstallmentTransactionWithoutCreditcardMsg       = "INSTALLMENT_TRANSACTION_WITHOUT_CREDITCARD"
	InstallmentTransactionWithDifferentCreditcardMsg = "INSTALLMENT_TRANSACTION_WITH_DIFFERENT_CREDITCARD"
	MonthlyTransactionWithDifferentCategoryMsg       = "MONTHLY_TRANSACTION_WITH_DIFFERENT_CATEGORY"
	MonthlyTransactionWithCreditcardMsg              = "MONTHLY_TRANSACTION_WITH_CREDITCARD"
	MonthlyTransactionWithoutCreditcardMsg           = "MONTHLY_TRANSACTION_WITHOUT_CREDITCARD"
	MonthlyTransactionWithDifferentCreditcardMsg     = "MONTHLY_TRANSACTION_WITH_DIFFERENT_CREDITCARD"
)
