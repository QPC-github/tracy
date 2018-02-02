package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tracy/configure"
	"tracy/log"
	"tracy/tracer/types"
)

/*AddTracers takes multiple tracer structs and sends them to the tracer API. This client request can return multiple errors,
 * up to one per tracer sent. */
func AddTracers(request types.Request) error {
	log.Trace.Printf("Adding the following tracers: %+v", request.Tracers)
	var requestJSON string
	var err error
	if requestJSON, err = json.Marshal(request); err == nil {
		log.Trace.Printf("Decoded the tracer into the following JSON: %s", requestJSON)

		/* Send the request off to the API. We don't need the response.*/
		var tracerServer interface{}
		if tracerServer, err = configure.ReadConfig("tracer-server"); err == nil {
			url := fmt.Sprintf("http://%s/tracers", tracerServer.(string))
			contentType := "application/json; charset=UTF-8"
			log.Trace.Printf("Sending POST request to %s %s", url, contentType)
			_, err = http.Post(url, contentType, bytes.NewBuffer(requestJSON))
		}
	}

	if err != nil {
		log.Warning.Printf(err.Error())
	}

	return err
}

/*GetTracers gets a list of the current tracers in the database. */
func GetTracers() ([]types.Tracer, error) {
	log.Trace.Printf("Getting all the tracers")
	ret := make([]types.Tracer, 0)
	var err error

	var tracerServer interface{}
	if tracerServer, err = configure.ReadConfig("tracer-server"); err == nil {
		url := fmt.Sprintf("http://%s/tracers", tracerServer.(string))
		log.Trace.Printf("Sending GET request to %s", url)
		var tracers *http.Response
		if tracers, err = http.Get(url); err == nil {
			log.Trace.Printf("Request submitted successfully")
			var tracersBody []byte
			if tracersBody, err = ioutil.ReadAll(tracers.Body); err == nil {
				log.Trace.Printf("Read the following from the request response: %s", tracersBody)
				/* Last success case. Unmarshal the tracers and check for parsing errors. */
				err = json.Unmarshal(tracersBody, &ret)
			}
			defer tracers.Body.Close()
		}
	}

	if err != nil {
		log.Warning.Printf(err.Error())
	}

	return ret, err
}

/*AddTracerEvents takes multiple tracer event structs and adds to them to a tracer using the tracer API. This client
 * request can return multiple errors, up to one per tracer sent. */
func AddTracerEvents(tracerEvents map[int]types.TracerEvent) []error {
	log.Trace.Printf("Adding the following tracer events: %+v", tracerEvents)
	ret := make([]error, 0)

	for tracerID, tracerEvent := range tracerEvents {
		/* Using the tracer ID associated with the event, add it to the API. */
		if err := AddTracerEvent(tracerEvent, tracerID); err != nil {
			/* If there is an error, record it and continue. */
			log.Warning.Printf(err.Error())
			ret = append(ret, err)
		}

	}

	return ret
}

/*AddTracerEvent adds a single tracer event struct to a tracer using the tracer API. */
func AddTracerEvent(tracerEvent types.TracerEvent, tracerID int) error {
	log.Trace.Printf("Adding the following tracer event: %+v, tracer ID: %s", tracerEvent, tracerID)
	var err error
	tracerEvent.TracerID = tracerID

	var eventJSON string
	if eventJSON, err = json.Marshal(tracerEvent); err == nil {
		var tracerServer interface{}
		if tracerServer, err = configure.ReadConfig("tracer-server"); err == nil {
			url := fmt.Sprintf("http://%s/tracers/%d/events", tracerServer.(string), tracerID)
			contentType := "application/json; charset=UTF-8"
			log.Trace.Printf("Sending POST request with %s to %s %s", eventJSON, url, contentType)
			_, err = http.Post(url, contentType, bytes.NewBuffer(eventJSON))
		}
	}

	if err != nil {
		log.Warning.Println(err)
	}

	return err
}

/*AddLabel adds a single label to the tracer API. */
func AddLabel(label types.Label) error {
	log.Trace.Printf("Adding the following label: %+v", label)
	var ret error

	labelJSON, err := json.Marshal(label)
	if err == nil {
		var tracerServer interface{}
		tracerServer, err = configure.ReadConfig("tracer-server")
		if err == nil {
			url := fmt.Sprintf("http://%s/labels", tracerServer.(string))
			contentType := "application/json; charset=UTF-8"
			log.Trace.Printf("Sending POST request with %s to %s %s", labelJSON, url, contentType)
			_, err = http.Post(url, contentType, bytes.NewBuffer(labelJSON))
		}
	}

	if err != nil {
		log.Warning.Println(err)
	}

	return err
}

/*GetLabels gets a list of the all the labels in the database. */
func GetLabels() ([]types.Label, error) {
	log.Trace.Printf("Getting all the labels")
	ret := make([]types.Label, 0)
	var err error

	var tracerServer interface{}
	if tracerServer, err = configure.ReadConfig("tracer-server"); err == nil {
		url := fmt.Sprintf("http://%s/labels", tracerServer.(string))
		log.Trace.Printf("Sending GET request to %s", url)

		var labels *http.Response
		if labels, err = http.Get(url); err == nil {
			log.Trace.Printf("Request submitted successfully")
			var tracersBody []byte
			if tracersBody, err = ioutil.ReadAll(labels.Body); err == nil {
				log.Trace.Printf("Read the following from the request response: %s", tracersBody)
				/* Last success case. Unmarshal the tracers and check for parsing errors. */
				err = json.Unmarshal(tracersBody, &ret)
			}
			defer labels.Body.Close()
		}
	}

	if err != nil {
		log.Warning.Println(err)
	}

	return ret, err
}

/*GetLabel gets the label with the ID in the API. */
func GetLabel(ID int) (types.Label, error) {
	log.Trace.Printf("Getting the label %d", ID)
	ret := types.Label{}
	var err error

	var tracerServer interface{}
	if tracerServer, err = configure.ReadConfig("tracer-server"); err == nil {
		url := fmt.Sprintf("http://%s/tracers/%d", tracerServer.(string), ID)
		log.Trace.Printf("Sending GET request to %s", url)
		var label *http.Response
		if label, err = http.Get(url); err == nil {
			log.Trace.Printf("Request submitted successfully")
			var labelBody []byte
			if labelBody, err = ioutil.ReadAll(label.Body); err == nil {
				log.Trace.Printf("Read the following from the request response: %s", labelBody)
				/* Last success case. Unmarshal the label and check for parsing errors. */
				err = json.Unmarshal(labelBody, &ret)
			}
			defer label.Body.Close()
		}
	}

	if err != nil {
		log.Warning.Printf(err.Error())
	}

	return ret, err
}
