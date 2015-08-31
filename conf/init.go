package conf

type App struct {
	DS Datastore
}

type Mongo struct {
	db    string
	hosts []string
}

type Datastore struct {
	Mongo Mongo
}

func (m Mongo) Hosts() []string {
	return m.hosts
}

func (m Mongo) DB() string {
	return m.db
}

func (m *Mongo) SetHosts(hosts ...string) {
	m.hosts = append(m.hosts, hosts...)
}

func (m *Mongo) SetDB(name string) {
	m.db = name
}
