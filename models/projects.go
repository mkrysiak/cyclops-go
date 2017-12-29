package models

import (
	"sync"

	"github.com/cyclops-go/conf"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type SentryProject struct {
	ProjectId int    `db:"project_id"`
	PublicKey string `db:"public_key"`
	SecretKey string `db:"secret_key"`
}

type SentryProjects struct {
	SentryProjects []SentryProject
	Cfg            *conf.Config
	Db             *sqlx.DB
}

var projectsById sync.Map

func NewSentryProjects(cfg *conf.Config) *SentryProjects {

	// scheme, url, _ := getDbUrl(cfg, env)
	scheme, url, _ := cfg.GetDatabaseSchemeAndUrl()

	db, err := sqlx.Connect(scheme, url)
	if err != nil {
		log.Fatal(err)
	}
	//TODO:  Log the full URL, after stripping the password
	log.WithFields(log.Fields{
		"DriverName": db.DriverName(),
	}).Info("Connected to database:")

	if scheme == "sqlite3" {
		db.SetMaxOpenConns(1)
	}

	return &SentryProjects{Cfg: cfg, Db: db}
}

func (p *SentryProjects) IsValidProjectAndPublicKey(projectId int, publicKey string) bool {
	if p, ok := projectsById.Load(projectId); ok {
		if p.(SentryProject).PublicKey == publicKey {
			return true
		}
	}
	return false
}

func (p *SentryProjects) UpdateProjects() {
	log.Info("Updating projects map")
	projects := []SentryProject{}
	err := p.Db.Select(&projects, "SELECT project_id, public_key, secret_key FROM sentry_projectkey")
	if err != nil {
		log.Error(err)
	} else {
		//TODO: Handle deletion of projects once moved to a pg listener
		for _, v := range projects {
			if _, ok := projectsById.LoadOrStore(v.ProjectId, v); !ok {
				log.WithFields(log.Fields{"ProjectId": v.ProjectId}).Info("Loaded new project: ")
			}

		}
	}

}

func (p *SentryProjects) Shutdown() {
	err := p.Db.Close()
	if err != nil {
		log.Error("Unable to cleanly close the database connection")
	}
	log.Info("Closed database connection.")
}
