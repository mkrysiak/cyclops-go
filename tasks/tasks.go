package tasks

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/mkrysiak/cyclops-go/models"
	log "github.com/sirupsen/logrus"
)

type Tasks struct {
	requestStorage             *models.RequestStorage
	sentryProjects             *models.SentryProjects
	defaultTaskDumpInterval    time.Duration
	projectsUpdaterTicker      *time.Ticker
	sendRequestsToSentryTicker *time.Ticker
	sentryClient               *http.Client
}

func StartBackgroundTaskRunners(sentryProjects *models.SentryProjects, requestStorage *models.RequestStorage) *Tasks {
	tasks := &Tasks{
		requestStorage:          requestStorage,
		sentryProjects:          sentryProjects,
		defaultTaskDumpInterval: 20 * time.Millisecond,
		sentryClient:            &http.Client{Timeout: 5 * time.Second},
	}
	tasks.projectsUpdaterTicker = tasks.startSendRequestsToSentryTask()
	tasks.sendRequestsToSentryTicker = tasks.startSentryProjectsUpdaterTask()
	return tasks
}

func (t *Tasks) startSendRequestsToSentryTask() *time.Ticker {
	ticker := time.NewTicker(t.defaultTaskDumpInterval)
	go func() {
		for {
			select {
			case <-ticker.C:
				t.dequeueRequestAndSend()
			}
		}
	}()
	return ticker
}

//TODO: Should use a pg listener rather than polling
func (t *Tasks) startSentryProjectsUpdaterTask() *time.Ticker {
	t.sentryProjects.UpdateProjects()
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				t.sentryProjects.UpdateProjects()
			}
		}
	}()
	return ticker
}

func (t *Tasks) dequeueRequestAndSend() {
	message, existsMessage := t.requestStorage.GetNextMessage()
	if !existsMessage {
		return
	}

	resp, err := t.sendRequest(message)
	if err != nil {
		log.Error(err)
		retryAfterDuration := time.Duration(30) * time.Second
		t.requestStorage.SetRetryAfter(message.ProjectId, retryAfterDuration)
		t.requestStorage.LPush(message)
		return
	}
	retryAfterDuration := t.handleResponse(resp)
	if retryAfterDuration > 0 {
		t.requestStorage.SetRetryAfter(message.ProjectId, retryAfterDuration)
		t.requestStorage.LPush(message)
	}
}

func (t *Tasks) Shutdown() {
	log.Info("Shutting down background tasks")
	t.projectsUpdaterTicker.Stop()
	t.sendRequestsToSentryTicker.Stop()
}

func (t *Tasks) sendRequest(message *models.Message) (*http.Response, error) {
	req, err := http.NewRequest(message.RequestMethod, message.OriginUrl, bytes.NewReader(message.RequestBody))
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Info("Unable to generate request to the origin Sentry server")
	}
	req.Header = stripHopByHopHeaders(message.Headers)
	resp, err := t.sentryClient.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Url:":  req.URL,
			// "Body":  string(message.RequestBody),
		}).Error("Failed to make a request to the backend Sentry server")
		return nil, err
	}
	return resp, nil
}

func (t *Tasks) handleResponse(resp *http.Response) time.Duration {
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.WithFields(log.Fields{
			"RspCode": resp.StatusCode,
			"Headers": resp.Header,
		}).Info("Bad response from Sentry:")
	}

	// Retry-After can be a time or a duration
	if resp.Header.Get("Retry-After") == "" {
		return 0
	}
	retryAfter, err := strconv.Atoi(resp.Header.Get("Retry-After"))
	// If we fail to convert to an int, it must be in Date format
	if err != nil {
		// According to the spec, http Date header must be sent in GMT
		retryDate, err := time.Parse(time.RFC1123, resp.Header.Get("Retry-After"))
		if err != nil {
			return 0
		}
		return (retryDate.Sub(time.Now().In(time.UTC)))

	}
	// Retry-After delay must be in seconds
	return time.Duration(retryAfter) * time.Second
}

func stripHopByHopHeaders(header http.Header) http.Header {

	connectionHeaders := header[http.CanonicalHeaderKey("Connection")]
	for _, v := range connectionHeaders {
		header.Del(v)
	}
	header.Del("Connection")
	header.Del("Keep-Alive")
	header.Del("Proxy-Authenticate")
	header.Del("Transfer-Encoding")
	header.Del("TE")
	header.Del("Trailer")
	header.Del("Upgrade")
	return header
}
