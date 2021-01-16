package synclog

import (
	"fmt"
	"github.com/hjertnes/roam/commands/help"
	"github.com/hjertnes/roam/state"
	"github.com/hjertnes/roam/utils"
	"github.com/rotisserie/eris"
	"os"
	"strings"
)

func Run(path string) error{
	s, err := state.New(path)
	if err != nil {
		return eris.Wrap(err, "Failed to create state")
	}

	if len(os.Args) == 2{
		help.Run()
		return nil
	}

	if os.Args[2] == "clear"{
		err = s.Dal.ClearLog()
		if err != nil{
			return eris.Wrap(err, "failed to clear log")
		}
		return nil
	}

	if os.Args[2] == "list"{
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

	help.Run()
	return nil
}
