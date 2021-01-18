package synclog

import (
	"fmt"
	"github.com/hjertnes/roam/commands/help"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	"github.com/rotisserie/eris"
	"strings"
)

func Run(path string, args []string) error{
	s, err := state.New(path, args)
	if err != nil {
		return eris.Wrap(err, "Failed to create state")
	}

	if len(s.Arguments) == 2{
		help.Run([]string{})
		return nil
	}

	if s.Arguments[2] == "clear"{
		err = s.Dal.ClearLog()
		if err != nil{
			return eris.Wrap(err, "failed to clear log")
		}
		return nil
	}

	if s.Arguments[2] == "list"{
		logs, err := s.Dal.GetLog()
		if err != nil{
			return eris.Wrap(err, "failed to get log")
		}

		res := make([]string, 0)

		for l := range logs{
			status := "Success"
			if logs[l].Failure{
				status = "Failure"
			}
			res = append(res, fmt.Sprintf("- %s - %s", logs[l].Timestamp.Format(s.Conf.DateTimeFormat), status))
		}

		err = utils.RenderMarkdown(strings.Join(res, "\n"))
		if err != nil{
			return eris.Wrap(err, "failed to render markdown")
		}
		return nil
	}

	help.Run([]string{})
	return nil
}
