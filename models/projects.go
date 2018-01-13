package models

import (
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

type SentryProject struct {
	ProjectId int    `db:"id"`
	Name      string `db:"name"`
	PublicKey string `db:"public_key"`
	SecretKey string `db:"secret_key"`
}

type SentryProjects struct {
	SentryProjects []SentryProject
	Db             *sqlx.DB
}

var projectsById sync.Map

func NewSentryProjects(scheme string, url string) *SentryProjects {

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

	return &SentryProjects{
		Db: db,
	}
}

func (p *SentryProjects) GetProjectName(projectId int) string {
	if p, ok := projectsById.Load(projectId); ok {
		return p.(SentryProject).Name
	}
	return ""
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
	log.Debug("Updating projects map from the Database")

	projects := []SentryProject{}
	err := p.Db.Select(&projects, "SELECT sentry_project.id, name, public_key, secret_key FROM sentry_project INNER JOIN sentry_projectkey ON sentry_projectkey.project_id = sentry_project.id")
	if err != nil {
		log.Error(err)
		return
	}

	//TODO: Handle deletion of projects once moved to a pg listener
	for _, v := range projects {
		if _, ok := projectsById.LoadOrStore(v.ProjectId, v); !ok {
			log.WithFields(log.Fields{"ProjectId": v.ProjectId}).Info("Loaded new project: ")
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
