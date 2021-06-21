package core

type CtxKeys string

const(
	AppCtxKey    CtxKeys = "app_ctx_key"
	PsqlCtxKey   CtxKeys = "psql_ctx_Key"
	MongoCtxKey   CtxKeys = "mongo_ctx_Key"
	ModelsCtxKey CtxKeys = "models_ctx_key"
)
