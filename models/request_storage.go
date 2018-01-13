package models

import (
	"bytes"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack"
)

type Message struct {
	ProjectId     string
	RequestMethod string
	Headers       http.Header
	OriginUrl     string
	RequestBody   []byte
}

// Sentry has global and per-project rate limiting. This will store the per-project retry time.
// There is no way to distinguish between a global or per-project rate limit from the response,
// so we'll ignore it.  It's not ideal, but not terrible either.
type RequestStorage struct {
	Client          *redis.Client
	projectsKey     string
	requestQueueKey string
	retryAfterKey   string
	scriptSha       string
}

var random *rand.Rand

func init() {
	s := rand.NewSource(time.Now().Unix())
	random = rand.New(s)
}

func NewRequestStorage(client *redis.Client) *RequestStorage {

	r := &RequestStorage{
		Client:          client,
		projectsKey:     "cyclops:projects",
		requestQueueKey: "c:q:",
		retryAfterKey:   "c:r:",
	}
	r.scriptSha = r.loadScript()
	return r
}

func (s *RequestStorage) Put(message *Message) {
	b, err := msgpack.Marshal(message)
	if err != nil {
		log.Error(err)
	}

	queueLength, err := s.Client.RPush(s.GetQueueKey(message.ProjectId), b).Result()
	if err != nil {
		log.Error(err)
	}

	err = s.Client.SAdd(s.projectsKey, message.ProjectId).Err()
	if err != nil {
		log.Error(err)
	}

	log.WithFields(log.Fields{
		"Queue Key":   s.GetQueueKey(message.ProjectId),
		"QueueLength": queueLength,
	}).Debug("RequestStorage::Put")
}

func (s *RequestStorage) LPush(message *Message) {
	b, err := msgpack.Marshal(message)

	queueLength, err := s.Client.LPush(s.GetQueueKey(message.ProjectId), b).Result()
	if err != nil {
		log.Error(err)
	}
	log.WithFields(log.Fields{
		"Queue Key":   s.GetQueueKey(message.ProjectId),
		"QueueLength": queueLength,
	}).Debug("RequestStoage::LPush")
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

func (s *RequestStorage) GetNextMessage() (*Message, bool) {
	var message Message
	msg, err := s.selectRandomRequest()
	if msg == nil || err != nil {
		return &message, false
	}

	var m Message
	err = msgpack.Unmarshal(msg, &m)
	if err != nil {
		log.Error(err)
		return &message, false
	}

	return &m, true
}

func (s *RequestStorage) loadScript() string {

	scriptSha, err := s.Client.Get("cyclops:scriptsha").Result()
	if err == nil {
		return scriptSha
	}

	script := `redis.replicate_commands()
	local randProjId=redis.call('srandmember', KEYS[1])
	if randProjId == false then
		return ""
	end
	local retryAfter=redis.call('get', KEYS[2] .. randProjId)
	if retryAfter ~= false then
		return ""
	end
	local projectQueue=redis.call('rpop', KEYS[3] .. randProjId)
	if projectQueue == false then
		redis.call('srem', KEYS[1], randProjId)
		return ""
	end
	return projectQueue`
	scriptSha, err = s.Client.ScriptLoad(script).Result()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Unable to load Lua script to Redis")
	}

	err = s.Client.Set("cyclops:scriptsha", scriptSha, 0).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Unable to set cyclops:scriptsha")
	}
	return scriptSha
}

func (s *RequestStorage) selectRandomRequest() ([]byte, error) {
	projectQueue, err := s.Client.EvalSha(s.scriptSha, []string{s.projectsKey, s.retryAfterKey, s.requestQueueKey}).Result()
	if err != nil {
		log.Error(err)
	}
	if projectQueue == nil || projectQueue.(string) == "" {
		return nil, err
	}
	return []byte(projectQueue.(string)), err
}

func (s *RequestStorage) SetRetryAfter(projectId string, d time.Duration) {
	err := s.Client.Set(s.GetRetryAfterKey(projectId), 1, d).Err()
	if err != nil {
		log.Error(err)
	}
}

func (s *RequestStorage) GetQueueKey(projectId string) string {
	var buffer bytes.Buffer
	buffer.WriteString(s.requestQueueKey)
	buffer.WriteString(projectId)
	return buffer.String()
}

func (s *RequestStorage) GetRetryAfterKey(projectId string) string {
	var buffer bytes.Buffer
	buffer.WriteString(s.retryAfterKey)
	buffer.WriteString(projectId)
	return buffer.String()
}
