package errstring

const (
	// ItemNotExist indicate the record not exist in the database
	ItemNotExist = "aurora_database.record_not_exist.domain.app_error"

	// ItemAlreadyExist indicate the record not exist in the database
	ItemAlreadyExist = "aurora_database.record_already_exist.domain.app_error"

	// DBSetUpFail error string
	DBSetUpFail = "aurora_database.database.db.setup_connection_failed.domain.app_erro"

	// DBSetUpMaxRetry error string
	DBSetUpMaxRetry = "aurora_database.database.db.setup_connection_with_max_attempts_or_canceled_failed.domain.app_error"

	// ESSetUpConnFail error string
	ESSetUpConnFail = "aurora_database.searching.elastic.setup_connection_client_failed.domain.app_error"

	// ESSetUpMaxRetry error string
	ESSetUpMaxRetry = "aurora_database.searching.elastic.setup_connection_with_max_attempts_or_canceled_failed.domain.app_error"

	// ESSetUpPingFailed error string
	ESSetUpPingFailed = "aurora_database.searching.elastic.setup_connection_ping_failed.domain.app_error"

	// ESCreateFailed error string
	ESCreateFailed = "aurora_database.searching.elastic.create_item_failed.domain.app_error"

	// ESCheckIndexExistFailed error string
	ESCheckIndexExistFailed = "aurora_database.searching.elastic.check_index_failed.domain.app_error"

	// ESCreateIndexFailed error string
	ESCreateIndexFailed = "aurora_database.searching.elastic.create_index_failed.domain.app_error"
)
