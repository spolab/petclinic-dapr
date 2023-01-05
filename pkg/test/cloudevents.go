package test

import (
	"bytes"
	"strings"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/rs/zerolog/log"
)

// A gomock matcher that matches test events without timestamp and uuid verification
// It can compare more stuff, for now that´s sufficient.
type CloudEventArrayMatcher struct {
	Expected []*cloudevents.Event
}

func (m CloudEventArrayMatcher) Matches(arg any) bool {
	actual, ok := arg.([]*cloudevents.Event)
	if !ok {
		log.Error().Msg("source is not an array of cloudevents")
		return false
	}
	if len(actual) != len(m.Expected) {
		log.Error().Msg("arrays are of different size")
		return false
	}
	for i := range actual {
		if !stringsEqual(m.Expected[i].SpecVersion(), actual[i].SpecVersion(), "specversion", i) {
			return false
		}
		if !stringsEqual(m.Expected[i].Type(), actual[i].Type(), "type", i) {
			return false
		}
		if !stringsEqual(m.Expected[i].Source(), actual[i].Source(), "source", i) {
			return false
		}
		if !stringsEqual(m.Expected[i].DataContentType(), actual[i].DataContentType(), "dataContentType", i) {
			log.Error().Int("index", i).Msg("dataContentType")
			return false
		}
		if !bytesEqual(m.Expected[i].Data(), actual[i].Data(), "data", i) {
			return false
		}
	}
	return true
}

func stringsEqual(expected string, actual string, field string, index int) bool {
	if strings.Compare(expected, actual) != 0 {
		log.Error().Int("index", index).Str("expected", expected).Str("actual", actual).Msg(field)
		return false
	}
	return true
}

func bytesEqual(expected []byte, actual []byte, field string, index int) bool {
	if !bytes.Equal(expected, actual) {
		log.Error().Int("index", index).Bytes("expected", expected).Bytes("actual", actual).Msg(field)
		return false
	}
	return true
}

func (m CloudEventArrayMatcher) String() string {
	return "not available"
}
