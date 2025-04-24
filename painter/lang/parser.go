package lang

import (
	"bufio"
	"errors"
	"github.com/matshp0/ArchitectureLab3/painter"
	"io"
	"strings"
)

var (
	InvalidParameters = "received invalid amount of parameters"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		commandLine := scanner.Text()
		cmd, err := parse(commandLine)
		if err != nil {
			return nil, err
		}
		res = append(res, cmd)
	}
	return res, nil
}

func parse(str string) (painter.Operation, error) {
	fields := strings.Fields(str)
	if len(fields) == 0 {
		return nil, errors.New("empty command")
	}
	cmd := fields[0]
	var options map[string]float32
	var f painter.OperationFunc
	switch cmd {
	case "white":
		f = painter.WhiteFill
		break
	case "green":
		f = painter.GreenFill
		break
	case "update":
		return painter.UpdateOp, nil

	default:
		return nil, errors.New("unknown command")
	}
	return painter.Command{
		F:       f,
		Options: options,
	}, nil
}
