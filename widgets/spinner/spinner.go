package spinner

import (
	"fmt"
	"time"

	"github.com/rotisserie/eris"
	"github.com/theckman/yacspin"
)

func defaultMessage(message string) string {
	if message == "" {
		return "Loading"
	}

	return message
}

func Run(message string) (*yacspin.Spinner, error) {
	cfg := yacspin.Config{
		Frequency:       100 * time.Millisecond,
		CharSet:         yacspin.CharSets[69],
		Suffix:          fmt.Sprintf(" %s", defaultMessage(message)),
		SuffixAutoColon: true,
		Message:         "",
		StopCharacter:   "✓",
		StopColors:      []string{"fgGreen"},
	}

	spinner, err := yacspin.New(cfg)
	if err != nil {
		return nil, eris.Wrap(err, "failed to create spinner")
	}

	return spinner, nil
}
