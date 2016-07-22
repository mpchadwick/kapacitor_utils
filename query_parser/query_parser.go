/**
 * Takes in a stream from logparser and tracks usage of query parameters in the
 * URL
 *
 * query_param,key=key value=value
 */
package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/influxdata/kapacitor/udf"
	"github.com/influxdata/kapacitor/udf/agent"
)

type queryParserHandler struct {
	agent *agent.Agent
}

func newQueryParserHandler(a *agent.Agent) *queryParserHandler {
	return &queryParserHandler{agent: a}
}

func (*queryParserHandler) Info() (*udf.InfoResponse, error) {
	info := &udf.InfoResponse{
		Wants:    udf.EdgeType_STREAM,
		Provides: udf.EdgeType_STREAM,
		Options:  map[string]*udf.OptionInfo{},
	}
	return info, nil
}

func (q *queryParserHandler) Init(r *udf.InitRequest) (*udf.InitResponse, error) {
	init := &udf.InitResponse{
		Success: true,
		Error:   "",
	}
	return init, nil
}

func (*queryParserHandler) Snaphost() (*udf.SnapshotResponse, error) {
	return &udf.SnapshotResponse{}, nil
}

func (*queryParserHandler) Restore(req *udf.RestoreRequest) (*udf.RestoreResponse, error) {
	return &udf.RestoreResponse{
		Success: true,
	}, nil
}

func (*queryParserHandler) BeginBatch(begin *udf.BeginBatch) error {
	return errors.New("batching not supported")
}

func (q *queryParserHandler) Point(p *udf.Point) error {
	value := p.FieldsString["request"]
	pos := strings.Index(value, "?")
	if pos == -1 {
		// No query string? See ya!
		return nil
	}

	query := value[pos+1 : len(value)]
	params := strings.Split(query, "&")

	for i := 0; i < len(params); i++ {
		parts := strings.Split(params[i], "=")
		value := ""
		if len(parts) == 2 {
			value = parts[1]
		}

		newPoint := new(udf.Point)
		newPoint.Time = p.Time
		newPoint.Tags = map[string]string{
			"k": parts[0],
		}
		newPoint.FieldsString = map[string]string{
			"v": value,
		}

		q.agent.Responses <- &udf.Response{
			Message: &udf.Response_Point{
				Point: newPoint,
			},
		}
	}

	return nil
}

func (*queryParserHandler) EndBatch(end *udf.EndBatch) error {
	return errors.New("batching not supported")
}

func (q *queryParserHandler) Stop() {
	close(q.agent.Responses)
}

func main() {
	a := agent.New(os.Stdin, os.Stdout)
	h := newQueryParserHandler(a)
	a.Handler = h

	log.Println("Starting agent")
	a.Start()
	err := a.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
