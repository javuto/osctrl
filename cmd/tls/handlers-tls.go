package main

import (
	"compress/gzip"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/javuto/osctrl/pkg/carves"
	"github.com/javuto/osctrl/pkg/environments"
	"github.com/javuto/osctrl/pkg/nodes"
	"github.com/javuto/osctrl/pkg/queries"
	"github.com/javuto/osctrl/pkg/settings"
)

const (
	metricEnrollReq = "enroll-req"
	metricEnrollErr = "enroll-err"
	metricEnrollOK  = "enroll-ok"
	metricLogReq    = "log-req"
	metricLogErr    = "log-err"
	metricLogOK     = "log-ok"
	metricConfigReq = "config-req"
	metricConfigErr = "config-err"
	metricConfigOK  = "config-ok"
	metricReadReq   = "read-req"
	metricReadErr   = "read-err"
	metricReadOK    = "read-ok"
	metricWriteReq  = "write-req"
	metricWriteErr  = "write-err"
	metricWriteOK   = "write-ok"
	metricInitReq   = "init-req"
	metricInitErr   = "init-err"
	metricInitOK    = "init-ok"
	metricBlockReq  = "block-req"
	metricBlockErr  = "block-err"
	metricBlockOK   = "block-ok"
)

// JSONApplication for Content-Type headers
const JSONApplication string = "application/json"

// TextPlain for Content-Type headers
const TextPlain string = "text/plain"

// JSONApplicationUTF8 for Content-Type headers, UTF charset
const JSONApplicationUTF8 string = JSONApplication + "; charset=UTF-8"

// TextPlainUTF8 for Content-Type headers, UTF charset
const TextPlainUTF8 string = TextPlain + "; charset=UTF-8"

// Handler to be used as health check
func okHTTPHandler(w http.ResponseWriter, r *http.Request) {
	// Send response
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("💥"))
}

// Handle testing requests
func testingHTTPHandler(w http.ResponseWriter, r *http.Request) {
	// Send response
	w.Header().Set("Content-Type", JSONApplicationUTF8)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("test"))
}

// Handle error requests
func errorHTTPHandler(w http.ResponseWriter, r *http.Request) {
	// Send response
	w.Header().Set("Content-Type", JSONApplicationUTF8)
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("oh no..."))
}

// Function to handle the enroll requests from osquery nodes
func enrollHandler(w http.ResponseWriter, r *http.Request) {
	incMetric(metricEnrollReq)
	var response []byte
	// Retrieve environment variable
	vars := mux.Vars(r)
	env, ok := vars["environment"]
	if !ok {
		incMetric(metricEnrollErr)
		log.Println(" Environment is missing")
		return
	}
	// Check if environment is valid
	if !envs.Exists(env) {
		incMetric(metricEnrollErr)
		log.Printf("error unknown environment (%s)", env)
		return
	}
	// Debug HTTP for environment
	debugHTTPDump(r, envsmap[env].DebugHTTP, true)
	// Decode read POST body
	var t EnrollRequest
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		incMetric(metricEnrollErr)
		log.Printf("error parsing POST body %v", err)
		return
	}
	// Check if received secret is valid
	var nodeKey string
	var newNode nodes.OsqueryNode
	nodeInvalid := true
	if checkValidSecret(t.EnrollSecret, env) {
		// Generate node_key using UUID as entropy
		nodeKey = generateNodeKey(t.HostIdentifier)
		newNode = nodeFromEnroll(t, env, r.Header.Get("X-Real-IP"), nodeKey)
		// Check if UUID exists already, if so archive node and enroll new node
		if nodesmgr.CheckByUUID(t.HostIdentifier) {
			err := nodesmgr.Archive(t.HostIdentifier, "exists")
			if err != nil {
				incMetric(metricEnrollErr)
				log.Printf("error archiving node %v", err)
			}
			// Update existing with new enroll data
			err = nodesmgr.UpdateByUUID(newNode, t.HostIdentifier)
			if err != nil {
				incMetric(metricEnrollErr)
				log.Printf("error updating existing node %v", err)
			} else {
				nodeInvalid = false
			}
		} else { // New node, persist it
			err := nodesmgr.Create(newNode)
			if err != nil {
				incMetric(metricEnrollErr)
				log.Printf("error creating node %v", err)
			} else {
				nodeInvalid = false
			}
		}
	} else {
		incMetric(metricEnrollErr)
		log.Printf("error invalid enrolling secret %s", t.EnrollSecret)
	}
	// Prepare response
	response, err = json.Marshal(EnrollResponse{NodeKey: nodeKey, NodeInvalid: nodeInvalid})
	if err != nil {
		log.Printf("error formating response %v", err)
		return
	}
	// Debug HTTP
	if envsmap[env].DebugHTTP {
		log.Printf("Response: %s", string(response))
	}
	// Send response
	w.Header().Set("Content-Type", JSONApplicationUTF8)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
	incMetric(metricEnrollOK)
}

// Function to handle the configuration requests from osquery nodes
func configHandler(w http.ResponseWriter, r *http.Request) {
	incMetric(metricConfigReq)
	var response []byte
	// Retrieve environment variable
	vars := mux.Vars(r)
	env, ok := vars["environment"]
	if !ok {
		incMetric(metricConfigErr)
		log.Println(" Environment is missing")
		return
	}
	// Check if environment is valid
	if !envs.Exists(env) {
		incMetric(metricConfigErr)
		log.Printf("error unknown environment (%s)", env)
		return
	}
	// Debug HTTP for environment
	debugHTTPDump(r, envsmap[env].DebugHTTP, true)
	// Get environment
	e, err := envs.Get(env)
	if err != nil {
		incMetric(metricConfigErr)
		log.Printf("error getting environment %v", err)
		return
	}
	// Decode read POST body
	var t ConfigRequest
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		incMetric(metricConfigErr)
		log.Printf("error parsing POST body %v", err)
		return
	}
	// Check if provided node_key is valid and if so, update node
	if nodesmgr.CheckByKey(t.NodeKey) {
		err = nodesmgr.UpdateIPAddressByKey(r.Header.Get("X-Real-IP"), t.NodeKey)
		if err != nil {
			incMetric(metricConfigErr)
			log.Printf("error updating IP address %v", err)
		}
		// Refresh last config for node
		err = nodesmgr.RefreshLastConfig(t.NodeKey)
		if err != nil {
			incMetric(metricConfigErr)
			log.Printf("error refreshing last config %v", err)
		}
		response = []byte(e.Configuration)
	} else {
		response, err = json.Marshal(ConfigResponse{NodeInvalid: true})
		if err != nil {
			incMetric(metricConfigErr)
			log.Printf("error formating response %v", err)
			return
		}
	}
	// Debug HTTP
	if envsmap[env].DebugHTTP {
		log.Printf("Configuration: %s", string(response))
	}
	// Send response
	w.Header().Set("Content-Type", JSONApplicationUTF8)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
	incMetric(metricConfigOK)
}

// Function to handle the log requests from osquery nodes, both status and results
func logHandler(w http.ResponseWriter, r *http.Request) {
	incMetric(metricLogReq)
	var response []byte
	// Retrieve environment variable
	vars := mux.Vars(r)
	env, ok := vars["environment"]
	if !ok {
		incMetric(metricLogErr)
		log.Println(" Environment is missing")
		return
	}
	// Check if environment is valid
	if !envs.Exists(env) {
		incMetric(metricLogErr)
		log.Printf("error unknown environment (%s)", env)
		return
	}
	// Check if body is compressed, if so, uncompress
	var err error
	if r.Header.Get("Content-Encoding") == "gzip" {
		r.Body, err = gzip.NewReader(r.Body)
		if err != nil {
			incMetric(metricLogErr)
			log.Printf("error decoding gzip body %v", err)
		}
		//defer r.Body.Close()
		defer func() {
			err := r.Body.Close()
			if err != nil {
				incMetric(metricLogErr)
				log.Printf("Failed to close body %v", err)
			}
		}()
	}
	// Debug HTTP here so the body will be uncompressed
	debugHTTPDump(r, envsmap[env].DebugHTTP, true)
	// Extract POST body and decode JSON
	var t LogRequest
	err = json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		incMetric(metricLogErr)
		log.Printf("error parsing POST body %v", err)
		return
	}
	//defer r.Body.Close()
	defer func() {
		err := r.Body.Close()
		if err != nil {
			incMetric(metricLogErr)
			log.Printf("Failed to close body %v", err)
		}
	}()
	var nodeInvalid bool
	// Check if provided node_key is valid and if so, update node
	if nodesmgr.CheckByKey(t.NodeKey) {
		nodeInvalid = false
		// Process logs and update metadata
		processLogs(t.Data, t.LogType, env, r.Header.Get("X-Real-IP"))
	} else {
		nodeInvalid = true
	}
	// Prepare response
	response, err = json.Marshal(LogResponse{NodeInvalid: nodeInvalid})
	if err != nil {
		incMetric(metricLogErr)
		log.Printf("error preparing response %v", err)
		response = []byte("")
	}
	// Debug
	if envsmap[env].DebugHTTP {
		log.Printf("Response: %s", string(response))
	}
	// Send response
	w.Header().Set("Content-Type", JSONApplicationUTF8)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
	incMetric(metricLogOK)
}

// Helper to process logs
func processLogs(data json.RawMessage, logType, environment, ipaddress string) {
	// Parse log to extract metadata
	var logs []LogGenericData
	err := json.Unmarshal(data, &logs)
	if err != nil {
		// FIXME metrics for this
		log.Printf("error parsing log %s %v", string(data), err)
	}
	// Iterate through received messages to extract metadata
	var uuids, hosts, names, users, osqueryusers, hashes, osqueryversions []string
	for _, l := range logs {
		uuids = append(uuids, l.HostIdentifier)
		hosts = append(hosts, l.Decorations.Hostname)
		names = append(names, l.Decorations.LocalHostname)
		users = append(users, l.Decorations.Username)
		osqueryusers = append(osqueryusers, l.Decorations.OsqueryUser)
		hashes = append(hashes, l.Decorations.ConfigHash)
		osqueryversions = append(osqueryversions, l.Version)
	}
	// FIXME it only uses the first element from the []string that uniq returns
	uuid := uniq(uuids)[0]
	user := uniq(users)[0]
	osqueryuser := uniq(osqueryusers)[0]
	host := uniq(hosts)[0]
	name := uniq(names)[0]
	hash := uniq(hashes)[0]
	osqueryversion := uniq(osqueryversions)[0]
	// Dispatch logs and update metadata
	dispatchLogs(data, uuid, ipaddress, user, osqueryuser, host, name, hash, osqueryversion, logType, environment)
}

// Helper to dispatch logs
func dispatchLogs(data []byte, uuid, ipaddress, user, osqueryuser, hostname, localname, hash, osqueryversion, logType, environment string) {
	// Send data to storage
	// FIXME allow multiple types of logging
	switch tlsConfig.Logging {
	case settings.LoggingGraylog:
		go graylogSend(data, environment, logType, uuid, tlsConfig.LoggingCfg)
	case settings.LoggingSplunk:
		go splunkSend(data, environment, logType, uuid, tlsConfig.LoggingCfg)
	case settings.LoggingDB:
		go postgresLog(data, environment, logType, uuid)
	case settings.LoggingStdout:
		log.Printf("LOG: %s from environment %s : %s", logType, environment, string(data))
	}
	// Use metadata to update record
	err := nodesmgr.UpdateMetadataByUUID(user, osqueryuser, hostname, localname, ipaddress, hash, osqueryversion, uuid)
	if err != nil {
		log.Printf("error updating metadata %s", err)
	}
	// Refresh last logging request
	if logType == statusLog {
		err := nodesmgr.RefreshLastStatus(uuid)
		if err != nil {
			log.Printf("error refreshing last status %v", err)
		}
	}
	if logType == resultLog {
		err := nodesmgr.RefreshLastResult(uuid)
		if err != nil {
			log.Printf("error refreshing last result %v", err)
		}
	}
}

// Helper to dispatch queries
func dispatchQueries(queryData QueryWriteData, node nodes.OsqueryNode) {
	// Prepare data to send
	data, err := json.Marshal(queryData)
	if err != nil {
		log.Printf("error preparing data %v", err)
	}
	// Send data to storage
	// FIXME allow multiple types of logging
	switch tlsConfig.Logging {
	case settings.LoggingGraylog:
		go graylogSend(data, node.Environment, queryLog, node.UUID, tlsConfig.LoggingCfg)
	case settings.LoggingSplunk:
		go splunkSend(data, node.Environment, queryLog, node.UUID, tlsConfig.LoggingCfg)
	case settings.LoggingDB:
		go postgresQuery(data, queryData.Name, node, queryData.Status)
	case settings.LoggingStdout:
		log.Printf("QUERY: %s from environment %s : %s", "query", node.Environment, string(data))
	}
	// Refresh last query write request
	err = nodesmgr.RefreshLastQueryWrite(node.UUID)
	if err != nil {
		log.Printf("error refreshing last query write %v", err)
	}
}

// Function to handle on-demand queries to osquery nodes
func queryReadHandler(w http.ResponseWriter, r *http.Request) {
	incMetric(metricReadReq)
	// Retrieve environment variable
	vars := mux.Vars(r)
	env, ok := vars["environment"]
	if !ok {
		incMetric(metricReadErr)
		log.Println(" Environment is missing")
		return
	}
	// Check if environment is valid
	if !envs.Exists(env) {
		incMetric(metricReadErr)
		log.Printf("error unknown environment (%s)", env)
		return
	}
	// Debug HTTP
	debugHTTPDump(r, envsmap[env].DebugHTTP, true)
	// Decode read POST body
	var response []byte
	var t QueryReadRequest
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		incMetric(metricReadErr)
		log.Printf("error parsing POST body %v", err)
		return
	}
	var nodeInvalid bool
	qs := make(queries.QueryReadQueries)
	// Lookup node by node_key
	node, err := nodesmgr.GetByKey(t.NodeKey)
	if err == nil {
		err = nodesmgr.UpdateIPAddress(r.Header.Get("X-Real-IP"), node)
		if err != nil {
			incMetric(metricReadErr)
			log.Printf("error updating IP Address %v", err)
		}
		nodeInvalid = false
		qs, err = queriesmgr.NodeQueries(node)
		if err != nil {
			incMetric(metricReadErr)
			log.Printf("error getting queries from db %v", err)
		}
		// Refresh last query read request
		err = nodesmgr.RefreshLastQueryRead(t.NodeKey)
		if err != nil {
			incMetric(metricReadErr)
			log.Printf("error refreshing last query read %v", err)
		}
	} else {
		nodeInvalid = true
	}
	// Prepare response for invalid key
	response, err = json.Marshal(QueryReadResponse{Queries: qs, NodeInvalid: nodeInvalid})
	if err != nil {
		incMetric(metricReadErr)
		log.Printf("error formating response %v", err)
		response = []byte("")
	}
	// Debug HTTP
	if envsmap[env].DebugHTTP {
		log.Printf("Response: %s", string(response))
	}
	// Send response
	w.Header().Set("Content-Type", JSONApplicationUTF8)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
	incMetric(metricReadOK)
}

// Function to handle distributed query results from osquery nodes
func queryWriteHandler(w http.ResponseWriter, r *http.Request) {
	incMetric(metricWriteReq)
	// Retrieve environment variable
	vars := mux.Vars(r)
	env, ok := vars["environment"]
	if !ok {
		incMetric(metricWriteErr)
		log.Println(" Environment is missing")
		return
	}
	// Check if environment is valid
	if !envs.Exists(env) {
		incMetric(metricWriteErr)
		log.Printf("error unknown environment (%s)", env)
		return
	}
	// Debug HTTP
	debugHTTPDump(r, envsmap[env].DebugHTTP, true)
	// Decode read POST body
	var response []byte
	var t QueryWriteRequest
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		incMetric(metricWriteErr)
		log.Printf("error parsing POST body %v", err)
		return
	}
	var nodeInvalid bool
	// Check if provided node_key is valid and if so, update node
	if nodesmgr.CheckByKey(t.NodeKey) {
		err = nodesmgr.UpdateIPAddressByKey(r.Header.Get("X-Real-IP"), t.NodeKey)
		if err != nil {
			incMetric(metricWriteErr)
			log.Printf("error updating IP Address %v", err)
		}
		nodeInvalid = false
		// Process submitted results
		go processLogQueryResult(t.Queries, t.Statuses, t.NodeKey, env)
	} else {
		nodeInvalid = true
	}
	// Prepare response
	response, err = json.Marshal(QueryWriteResponse{NodeInvalid: nodeInvalid})
	if err != nil {
		incMetric(metricWriteErr)
		log.Printf("error formating response %v", err)
		response = []byte("")
	}
	// Debug HTTP
	if envsmap[env].DebugHTTP {
		log.Printf("Response: %s", string(response))
	}
	// Send response
	w.Header().Set("Content-Type", JSONApplicationUTF8)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
	incMetric(metricWriteOK)
}

// Helper to process on-demand query result logs
func processLogQueryResult(queries QueryWriteQueries, statuses QueryWriteStatuses, nodeKey string, environment string) {
	// Retrieve node
	node, err := nodesmgr.GetByKey(nodeKey)
	if err != nil {
		log.Printf("error retrieving node %s", err)
	}
	// Tap into results so we can update internal metrics
	for q, r := range queries {
		// Dispatch query name, result and status
		d := QueryWriteData{
			Name:   q,
			Result: r,
			Status: statuses[q],
		}
		go dispatchQueries(d, node)
		// Update internal metrics per query
		var err error
		if statuses[q] != 0 {
			err = queriesmgr.IncError(q)
		} else {
			err = queriesmgr.IncExecution(q)
		}
		if err != nil {
			log.Printf("error updating query %s", err)
		}
		// Add a record for this query
		err = queriesmgr.TrackExecution(q, node.UUID, statuses[q])
		if err != nil {
			log.Printf("error adding query execution %s", err)
		}
	}
}

// Function to handle the endpoint for quick enrollment script distribution
func quickEnrollHandler(w http.ResponseWriter, r *http.Request) {
	// FIXME metrics
	// Retrieve environment variable
	vars := mux.Vars(r)
	env, ok := vars["environment"]
	if !ok {
		log.Println(" Environment is missing")
		return
	}
	// Check if environment is valid
	if !envs.Exists(env) {
		log.Printf("error unknown environment (%s)", env)
		return
	}
	// Debug HTTP
	debugHTTPDump(r, envsmap[env].DebugHTTP, true)
	e, err := envs.Get(env)
	if err != nil {
		log.Printf("error getting environment %v", err)
		return
	}
	// Retrieve type of script
	script, ok := vars["script"]
	if !ok {
		log.Println("Script is missing")
		return
	}
	// Retrieve SecretPath variable
	secretPath, ok := vars["secretpath"]
	if !ok {
		log.Println("Path is missing")
		return
	}
	// Check if provided SecretPath is valid and is not expired
	if strings.HasPrefix(script, "enroll") {
		if !checkValidEnrollSecretPath(env, secretPath) {
			log.Println("Invalid Path")
			return
		}
	} else if strings.HasPrefix(script, "remove") {
		if !checkValidRemoveSecretPath(env, secretPath) {
			log.Println("Invalid Path")
			return
		}
	}
	// Prepare response with the script
	quickScript, err := environments.QuickAddScript(projectName, script, e)
	if err != nil {
		log.Printf("error getting script %v", err)
		return
	}
	// Send response
	w.Header().Set("Content-Type", TextPlainUTF8)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(quickScript))
}

// Function to initialize a file carve from a node
func processCarveInit(req CarveInitRequest, sessionid, environment string) {
	// Retrieve node
	node, err := nodesmgr.GetByKey(req.NodeKey)
	if err != nil {
		incMetric(metricInitErr)
		log.Printf("error retrieving node %s", err)
	}
	// Prepare carve to initialize
	carve := carves.CarvedFile{
		CarveID:         req.CarveID,
		RequestID:       req.RequestID,
		SessionID:       sessionid,
		UUID:            node.UUID,
		 Environment:         environment,
		CarveSize:       req.CarveSize,
		BlockSize:       req.BlockSize,
		TotalBlocks:     req.BlockCount,
		CompletedBlocks: 0,
		CarvedPath:      "",
		DestPath:        "",
		Status:          carves.StatusInitialized,
	}
	// Create File Carve
	err = filecarves.CreateCarve(carve)
	if err != nil {
		incMetric(metricInitErr)
		log.Printf("error creating  CarvedFile %v", err)
	}
}

// Function to process one block from a file carve
// FIXME it can be more efficient on db access
func processCarveBlock(req CarveBlockRequest, environment string) {
	// Prepare carve block
	block := carves.CarvedBlock{
		RequestID: req.RequestID,
		SessionID: req.SessionID,
		 Environment:   environment,
		BlockID:   req.BlockID,
		Data:      req.Data,
	}
	// Create Block
	err := filecarves.CreateBlock(block)
	if err != nil {
		incMetric(metricBlockErr)
		log.Printf("error creating CarvedBlock %v", err)
	}
	// Bump block completion
	err = filecarves.CompleteBlock(req.SessionID)
	if err != nil {
		incMetric(metricBlockErr)
		log.Printf("error completing block %v", err)
	}
	// If it is completed, set status
	if filecarves.Completed(req.SessionID) {
		err = filecarves.ChangeStatus(carves.StatusCompleted, req.SessionID)
		if err != nil {
			incMetric(metricBlockErr)
			log.Printf("error completing status %v", err)
		}
	} else {
		err = filecarves.ChangeStatus(carves.StatusInProgress, req.SessionID)
		if err != nil {
			incMetric(metricBlockErr)
			log.Printf("error progressing status %v", err)
		}
	}
}

// Function to handle the initialization of the file carver
func carveInitHandler(w http.ResponseWriter, r *http.Request) {
	incMetric(metricInitReq)
	// Retrieve environment variable
	vars := mux.Vars(r)
	env, ok := vars["environment"]
	if !ok {
		incMetric(metricInitErr)
		log.Println(" Environment is missing")
		return
	}
	// Check if environment is valid
	if !envs.Exists(env) {
		incMetric(metricInitErr)
		log.Printf("error unknown environment (%s)", env)
		return
	}
	// Debug HTTP
	debugHTTPDump(r, envsmap[env].DebugHTTP, true)
	// Decode read POST body
	var response []byte
	var t CarveInitRequest
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		incMetric(metricInitErr)
		log.Printf("error parsing POST body %v", err)
		return
	}
	var initCarve bool
	var carveSessionID string
	// Check if provided node_key is valid and if so, update node
	if nodesmgr.CheckByKey(t.NodeKey) {
		err = nodesmgr.UpdateIPAddressByKey(r.Header.Get("X-Real-IP"), t.NodeKey)
		if err != nil {
			incMetric(metricInitErr)
			log.Printf("error updating IP Address %v", err)
		}
		initCarve = true
		carveSessionID = generateCarveSessionID()
		// Process carve init
		go processCarveInit(t, carveSessionID, env)
	} else {
		initCarve = false
	}
	// Prepare response
	response, err = json.Marshal(CarveInitResponse{Success: initCarve, SessionID: carveSessionID})
	if err != nil {
		incMetric(metricInitErr)
		log.Printf("error formating response %v", err)
		response = []byte("")
	}
	// Debug HTTP
	if envsmap[env].DebugHTTP {
		log.Printf("Response: %s", string(response))
	}
	// Send response
	w.Header().Set("Content-Type", JSONApplicationUTF8)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
	incMetric(metricInitOK)
}

// Function to handle the blocks of the file carver
func carveBlockHandler(w http.ResponseWriter, r *http.Request) {
	incMetric(metricBlockReq)
	// Retrieve environment variable
	vars := mux.Vars(r)
	env, ok := vars["environment"]
	if !ok {
		incMetric(metricBlockErr)
		log.Println(" Environment is missing")
		return
	}
	// Check if environment is valid
	if !envs.Exists(env) {
		incMetric(metricBlockErr)
		log.Printf("error unknown environment (%s)", env)
		return
	}
	// Debug HTTP
	debugHTTPDump(r, envsmap[env].DebugHTTP, true)
	// Decode read POST body
	var response []byte
	var t CarveBlockRequest
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		incMetric(metricBlockErr)
		log.Printf("error parsing POST body %v", err)
		return
	}
	var blockCarve bool
	// Check if provided node_key is valid and if so, update node
	if filecarves.CheckCarve(t.SessionID, t.RequestID) {
		blockCarve = true
		// Process received block
		go processCarveBlock(t, env)
	} else {
		blockCarve = false
	}
	// Prepare response
	response, err = json.Marshal(CarveBlockResponse{Success: blockCarve})
	if err != nil {
		incMetric(metricBlockErr)
		log.Printf("error formating response %v", err)
		response = []byte("")
	}
	// Debug HTTP
	if envsmap[env].DebugHTTP {
		log.Printf("Response: %s", string(response))
	}
	// Send response
	w.Header().Set("Content-Type", JSONApplicationUTF8)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response)
	incMetric(metricBlockOK)
}