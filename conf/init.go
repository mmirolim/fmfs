package conf

type App struct {
	DS Datastore
}

type Datastore struct {
	MongoHosts  []string
	MongoDbName string
}
