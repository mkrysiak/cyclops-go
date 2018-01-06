package tasks

import (
	"bytes"
	"net/http"
	"time"

	"github.com/mkrysiak/cyclops-go/models"
	log "github.com/sirupsen/logrus"
)

type Tasks struct {
	requestStorage             *models.RequestStorage
	sentryProjects             *models.SentryProjects
	projectsUpdaterTicker      *time.Ticker
	sendRequestsToSentryTicker *time.Ticker
}

func StartBackgroundTaskRunners(sentryProjects *models.SentryProjects, requestStorage *models.RequestStorage) *Tasks {
	tasks := &Tasks{
		requestStorage: requestStorage,
		sentryProjects: sentryProjects,
	}
	tasks.projectsUpdaterTicker = tasks.startSendRequestsToSentryTask()
	tasks.sendRequestsToSentryTicker = tasks.startSentryProjectsUpdaterTask()
	return tasks
}

func (t *Tasks) startSendRequestsToSentryTask() *time.Ticker {
	ticker := time.NewTicker(20 * time.Millisecond)
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

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest(message.RequestMethod, message.OriginUrl, bytes.NewReader(message.RequestBody))
	if err != nil {
		log.WithFields(log.Fields{"Error": err}).Info("Unable to generate request to the origin Sentry server")
	}
	req.Header = stripHopByHopHeaders(message.Headers)
	resp, err := client.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Url:":  req.URL,
			"Body":  string(message.RequestBody),
		}).Error("Failed to make a request to the backend Sentry server")

		// TODO: Handle a down Sentry instance, and respond to the client.
		// Currently the request is lost (dequeued, but not sent)
		return
	}

	handleResponse(resp)
}

func (t *Tasks) Shutdown() {
	log.Info("Shutting down background tasks")
	t.projectsUpdaterTicker.Stop()
	t.sendRequestsToSentryTicker.Stop()
}

//TODO: Handle back-offs when Sentry is overloaded by honoring Sentry's
// retry-after (429) header.  https://docs.sentry.io/clientdev/overview/
func handleResponse(resp *http.Response) {
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		log.WithFields(log.Fields{
			"RspCode": resp.StatusCode,
			"Headers": resp.Header,
		}).Info("Bad response from Sentry:")
	}
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
