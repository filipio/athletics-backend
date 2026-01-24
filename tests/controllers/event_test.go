package controllers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
)

func TestEventPublishing(t *testing.T) {
	t.Run("Create event defaults to draft status", testCaseEvent(func(t *testing.T) {
		eventPayload := utils.AnyMap{
			"name":        "Test Event",
			"description": "Test Description",
			"deadline":    time.Now().Add(24 * time.Hour),
		}

		response, eventResponse, err := Post[models.Event]("/api/v1/events", eventPayload)

		if err != nil {
			t.Errorf("Error executing request: %s", err.Error())
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", response.StatusCode)
		}

		if eventResponse.Status != "draft" {
			t.Errorf("Expected event status 'draft', got '%s'", eventResponse.Status)
		}
	}))

	t.Run("Organizer can publish draft event", testCaseEvent(func(t *testing.T) {
		desc := "Test Description"
		event := &models.Event{
			Name:        "Test Event",
			Description: &desc,
			Deadline:    time.Now().Add(24 * time.Hour),
			Status:      "draft",
		}
		dbInstance.Save(event)

		response, publishedEvent, err := Post[models.Event](fmt.Sprintf("/api/v1/events/%d/publish", event.ID), nil)

		if err != nil {
			t.Errorf("Error executing request: %s", err.Error())
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", response.StatusCode)
		}

		if publishedEvent.Status != "published" {
			t.Errorf("Expected event status 'published', got '%s'", publishedEvent.Status)
		}
	}))

	t.Run("Publishing already-published event is idempotent", testCaseEvent(func(t *testing.T) {
		desc := "Test Description"
		event := &models.Event{
			Name:        "Test Event",
			Description: &desc,
			Deadline:    time.Now().Add(24 * time.Hour),
			Status:      "published",
		}
		dbInstance.Save(event)

		response, publishedEvent, err := Post[models.Event](fmt.Sprintf("/api/v1/events/%d/publish", event.ID), nil)

		if err != nil {
			t.Errorf("Error executing request: %s", err.Error())
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", response.StatusCode)
		}

		if publishedEvent.Status != "published" {
			t.Errorf("Expected event status 'published', got '%s'", publishedEvent.Status)
		}
	}))

	t.Run("Organizer can unpublish published event", testCaseEvent(func(t *testing.T) {
		desc := "Test Description"
		event := &models.Event{
			Name:        "Test Event",
			Description: &desc,
			Deadline:    time.Now().Add(24 * time.Hour),
			Status:      "published",
		}
		dbInstance.Save(event)

		response, unpublishedEvent, err := Post[models.Event](fmt.Sprintf("/api/v1/events/%d/unpublish", event.ID), nil)

		if err != nil {
			t.Errorf("Error executing request: %s", err.Error())
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", response.StatusCode)
		}

		if unpublishedEvent.Status != "draft" {
			t.Errorf("Expected event status 'draft', got '%s'", unpublishedEvent.Status)
		}
	}))

	t.Run("Unpublishing draft event is idempotent", testCaseEvent(func(t *testing.T) {
		desc := "Test Description"
		event := &models.Event{
			Name:        "Test Event",
			Description: &desc,
			Deadline:    time.Now().Add(24 * time.Hour),
			Status:      "draft",
		}
		dbInstance.Save(event)

		response, unpublishedEvent, err := Post[models.Event](fmt.Sprintf("/api/v1/events/%d/unpublish", event.ID), nil)

		if err != nil {
			t.Errorf("Error executing request: %s", err.Error())
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", response.StatusCode)
		}

		if unpublishedEvent.Status != "draft" {
			t.Errorf("Expected event status 'draft', got '%s'", unpublishedEvent.Status)
		}
	}))

	t.Run("Organizer sees both draft and published events", testCaseEvent(func(t *testing.T) {
		desc1 := "Test Description"
		desc2 := "Test Description"
		draftEvent := &models.Event{
			Name:        "Draft Event",
			Description: &desc1,
			Deadline:    time.Now().Add(24 * time.Hour),
			Status:      "draft",
		}
		publishedEvent := &models.Event{
			Name:        "Published Event",
			Description: &desc2,
			Deadline:    time.Now().Add(24 * time.Hour),
			Status:      "published",
		}

		dbInstance.Save(draftEvent)
		dbInstance.Save(publishedEvent)

		response, paginatedResponse, err := Get[utils.PaginatedResponse]("/api/v1/events")

		if err != nil {
			t.Errorf("Error executing request: %s", err.Error())
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", response.StatusCode)
		}

		if len(paginatedResponse.Data.([]interface{})) != 2 {
			t.Errorf("Expected 2 events, got %d", len(paginatedResponse.Data.([]interface{})))
		}
	}))

	t.Run("Event not found returns 404", testCaseEvent(func(t *testing.T) {
		response, _, err := Post[utils.ErrorsResponse]("/api/v1/events/99999/publish", nil)

		if err != nil {
			t.Errorf("Error executing request: %s", err.Error())
		}

		if response.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status code 404, got %d", response.StatusCode)
		}
	}))
}

func beforeEachEvent() {
	dbInstance.Where("1 = 1").Delete(&models.Event{})
	dbInstance.Where("1 = 1").Delete(&models.Question{})
}

func afterEachEvent() {
}

func testCaseEvent(test func(t *testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		beforeEachEvent()
		defer afterEachEvent()
		test(t)
	}
}
