package batch

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// ReplayCollection is used to unmarshal the JSON results of a list replays API call
type ReplayCollection struct {
	Name string `db:"collection_name" json:"name"`
}

// ReplayDestination is used to unmarshal the JSON results of a list replays API call
type ReplayDestination struct {
	Name string `db:"destination_name" json:"name"`
}

// Replay is used to unmarshal the JSON results of a list replays API call
type Replay struct {
	ID                 string `header:"Replay ID" json:"id"`
	Name               string `header:"Name" json:"name"`
	Type               string `header:"Type" json:"type"`
	Query              string `header:"Query" json:"query"`
	Paused             bool   `header:"Is Paused" json:"paused"`
	*ReplayDestination `json:"destination"`
	*ReplayCollection  `json:"collection"`
}

// ReplayOutput is used for displaying replays as a table
type ReplayOutput struct {
	Name        string `header:"Name" json:"name"`
	ID          string `header:"Replay ID" json:"id"`
	Type        string `header:"Type" json:"type"`
	Query       string `header:"Query" json:"query"`
	Collection  string `header:"Collection Name"`
	Destination string `header:"Destination Name"`
	Paused      bool   `header:"Is Paused" json:"paused"`
}

var (
	errReplayListFailed   = errors.New("unable to get list of replays")
	errNoReplays          = errors.New("you have no replays")
	errCreateReplayFailed = errors.New("failed to create new replay")
)

// ListReplays lists all of an account's replays
func (b *Batch) ListReplays() error {
	output, err := b.listReplays()
	if err != nil {
		return err
	}

	b.Printer(output)

	return nil
}

func (b *Batch) listReplays() ([]ReplayOutput, error) {
	res, _, err := b.Get("/v1/replay", nil)
	if err != nil {
		return nil, errReplayListFailed
	}

	replays := make([]*Replay, 0)

	err = json.Unmarshal(res, &replays)
	if err != nil {
		return nil, errReplayListFailed
	}

	if len(replays) == 0 {
		return nil, errNoReplays
	}

	output := make([]ReplayOutput, 0)
	for _, r := range replays {
		output = append(output, ReplayOutput{
			ID:          r.ID,
			Name:        r.Name,
			Type:        r.Type,
			Query:       r.Query,
			Collection:  r.ReplayCollection.Name,
			Destination: r.ReplayDestination.Name,
			Paused:      r.Paused,
		})
	}

	return output, nil
}

func (b *Batch) pauseReplay() error {
	// TODO

	return nil
}

func (b *Batch) resumeReplay() error {
	// TODO

	return nil
}

func (b *Batch) createReplay(query string) (*Replay, error) {
	p := map[string]interface{}{
		"name":           b.Opts.Batch.ReplayName,
		"type":           b.Opts.Batch.ReplayType,
		"query":          query,
		"collection_id":  b.Opts.Batch.CollectionID,
		"destination_id": b.Opts.Batch.DestinationID,
	}

	res, code, err := b.Post("/v1/replay", p)
	if code > 299 {
		errResponse := &BlunderErrorResponse{}
		if err := json.Unmarshal(res, errResponse); err != nil {
			return nil, errCreateReplayFailed
		}

		for _, e := range errResponse.Errors {
			err := fmt.Errorf("%s: '%s' %s", errCreateReplayFailed, e.Field, e.Message)
			b.Log.Error(err)
		}

		return nil, err
	}

	replay := &Replay{}

	if err := json.Unmarshal(res, replay); err != nil {
		return nil, errors.New("failed to unmarshal response")
	}

	return replay, nil
}

func (b *Batch) CreateReplay() error {

	query, err := b.generateReplayQuery()
	if err != nil {
		return err
	}

	replay, err := b.createReplay(query)
	if err != nil {
		return err
	}

	// TODO: Watch replay events and errors using generated ID

	b.Log.Infof("Replay started with id '%s'", replay.ID)

	return nil
}

func (b *Batch) generateReplayQuery() (string, error) {
	from, err := time.Parse("2006-01-02T15:04:05Z", b.Opts.Batch.ReplayFrom)
	if err != nil {
		return "", fmt.Errorf("--from-timestamp '%s' is not a valid RFC3339 date", b.Opts.Batch.ReplayFrom)
	}

	to, err := time.Parse("2006-01-02T15:04:05Z", b.Opts.Batch.ReplayTo)
	if err != nil {
		return "", fmt.Errorf("--to-timestamp '%s' is not a valid RFC3339 date", b.Opts.Batch.ReplayTo)
	}

	if b.Opts.Batch.Query == "*" {
		return fmt.Sprintf("batch.info.date_human: [%s TO %s]", from.Format(time.RFC3339), to.Format(time.RFC3339)), nil
	}

	return fmt.Sprintf("%s AND batch.info.date_human: [%s TO %s]", b.Opts.Batch.Query, from.Format(time.RFC3339), to.Format(time.RFC3339)), nil
}
