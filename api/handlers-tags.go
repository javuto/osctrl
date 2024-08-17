package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmpsec/osctrl/settings"
	"github.com/jmpsec/osctrl/users"
	"github.com/jmpsec/osctrl/utils"
)

const (
	metricAPITagsReq = "tags-req"
	metricAPITagsErr = "tags-err"
	metricAPITagsOK  = "tags-ok"
)

// GET Handler for multiple JSON tags
func apiTagsHandler(w http.ResponseWriter, r *http.Request) {
	incMetric(metricAPITagsReq)
	utils.DebugHTTPDump(r, settingsmgr.DebugHTTP(settings.ServiceAPI, settings.NoEnvironmentID), false)
	// Get context data and check access
	ctx := r.Context().Value(contextKey(contextAPI)).(contextValue)
	if !apiUsers.CheckPermissions(ctx[ctxUser], users.AdminLevel, users.NoEnvironment) {
		apiErrorResponse(w, "no access", http.StatusForbidden, fmt.Errorf("attempt to use API by user %s", ctx[ctxUser]))
		incMetric(metricAPITagsErr)
		return
	}
	// Get tags
	tags, err := tagsmgr.All()
	if err != nil {
		apiErrorResponse(w, "error getting tags", http.StatusInternalServerError, err)
		incMetric(metricAPITagsErr)
		return
	}
	// Serialize and serve JSON
	if settingsmgr.DebugService(settings.ServiceAPI) {
		log.Println("DebugService: Returned tags")
	}
	utils.HTTPResponse(w, utils.JSONApplicationUTF8, http.StatusOK, tags)
	incMetric(metricAPITagsOK)
}

// GET Handler to return tags for one environment as JSON
func apiTagsEnvHandler(w http.ResponseWriter, r *http.Request) {
	incMetric(metricAPITagsReq)
	utils.DebugHTTPDump(r, settingsmgr.DebugHTTP(settings.ServiceAPI, settings.NoEnvironmentID), false)
	vars := mux.Vars(r)
	// Extract environment
	envVar, ok := vars["env"]
	if !ok {
		apiErrorResponse(w, "error getting environment", http.StatusInternalServerError, nil)
		incMetric(metricAPIEnvsErr)
		return
	}
	// Get environment by name
	env, err := envs.Get(envVar)
	if err != nil {
		if err.Error() == "record not found" {
			apiErrorResponse(w, "environment not found", http.StatusNotFound, err)
		} else {
			apiErrorResponse(w, "error getting environment", http.StatusInternalServerError, err)
		}
		incMetric(metricAPIEnvsErr)
		return
	}
	// Get context data and check access
	ctx := r.Context().Value(contextKey(contextAPI)).(contextValue)
	if !apiUsers.CheckPermissions(ctx[ctxUser], users.AdminLevel, users.NoEnvironment) {
		apiErrorResponse(w, "no access", http.StatusForbidden, fmt.Errorf("attempt to use API by user %s", ctx[ctxUser]))
		incMetric(metricAPITagsErr)
		return
	}
	// Get tags
	tags, err := tagsmgr.GetByEnv(env.ID)
	if err != nil {
		apiErrorResponse(w, "error getting tags", http.StatusInternalServerError, err)
		incMetric(metricAPITagsErr)
		return
	}
	// Serialize and serve JSON
	if settingsmgr.DebugService(settings.ServiceAPI) {
		log.Println("DebugService: Returned tags")
	}
	utils.HTTPResponse(w, utils.JSONApplicationUTF8, http.StatusOK, tags)
	incMetric(metricAPITagsOK)
}
