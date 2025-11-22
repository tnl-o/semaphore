package api

import (
	"github.com/semaphoreui/semaphore/api/helpers"
	"github.com/semaphoreui/semaphore/db"
	"net/http"
)

// nolint: gocyclo
// getEvents retrieves events from the database based on user permissions and project context
// If a project is specified in the context, it returns events for that project
// Otherwise, it returns events for the current user
func getEvents(w http.ResponseWriter, r *http.Request, limit int) {
	user := helpers.GetFromContext(r, "user").(*db.User)
	projectObj, exists := helpers.GetOkFromContext(r, "project")

	var err error
	var events []db.Event

	if exists {
		project := projectObj.(db.Project)

		if !user.Admin { // check permissions to view events
			_, err = helpers.Store(r).GetProjectUser(project.ID, user.ID)
		}

		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		events, err = helpers.Store(r).GetEvents(project.ID, db.RetrieveQueryParams{Count: limit})
	} else {
		events, err = helpers.Store(r).GetUserEvents(user.ID, db.RetrieveQueryParams{Count: limit})
	}

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, events)
}

// getLastEvents retrieves the last 200 events for the current user or project
func getLastEvents(w http.ResponseWriter, r *http.Request) {
	getEvents(w, r, 200)
}

// getAllEvents retrieves all events for the current user or project without limit
func getAllEvents(w http.ResponseWriter, r *http.Request) {
	getEvents(w, r, 0)
}
