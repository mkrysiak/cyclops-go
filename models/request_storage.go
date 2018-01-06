package models

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
)

type Message struct {
	ProjectId     int
	RequestMethod string
	Headers       http.Header
	OriginUrl     string
	RequestBody   []byte
}

type RequestStorage struct {
	Client      *redis.Client
	projectsKey string
	random      *rand.Rand
}

func NewRequestStorage(client *redis.Client) *RequestStorage {

	s := rand.NewSource(time.Now().Unix())

	return &RequestStorage{
		Client:      client,
		projectsKey: "cyclops:projects",
		random:      rand.New(s),
	}
}

func (s *RequestStorage) Put(projectId int, message *Message) {
	err := s.Client.SAdd(s.projectsKey, projectId).Err()
	if err != nil {
		log.Error(err)
	}

	b, err := msgpack.Marshal(message)
	if err != nil {
		log.Error(err)
	}

	err = s.Client.RPush(getQueueKey(strconv.Itoa(projectId)), b).Err()
	if err != nil {
		log.Error(err)
	}
}

// func (s *RequestStorage) getProjects() []string {
// 	projects, err := s.Client.SMembers(s.projectsKey).Result()
// 	if err != nil {
// 		log.Error(err)
// 		return nil
// 	}
// 	if len(projects) < 1 {
// 		return nil
// 	}
// 	return projects

// }

// func (s *RequestStorage) getSize(projectId int) int64 {
// 	length, err := s.Client.LLen(strconv.Itoa(projectId)).Result()
// 	if err != nil {
// 		log.Error(err)
// 	}
// 	return length
// }

func (s *RequestStorage) GetNextMessage() (Message, bool) {
	// TODO: This could waste time picking from projects with empty queues.
	// Filter out empty queues
	var message Message
	projects := s.Client.SMembers(s.projectsKey).Val()
	if len(projects) == 0 {
		return message, false
	}

	projectId := s.random.Intn(len(projects))
	msg, _ := s.Client.RPop(getQueueKey(projects[projectId])).Bytes()
	if len(msg) == 0 {
		return message, false
	}
	var m Message
	err := msgpack.Unmarshal(msg, &m)
	if err != nil {
		log.Error(err)
		return message, false
	}
	return m, true
}

func getQueueKey(projectId string) string {
	return "cyclops:queue:" + projectId
}
