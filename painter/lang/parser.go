package lang

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/matshp0/ArchitectureLab3/painter"
	"io"
	"strconv"
	"strings"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation
	scanner := bufio.NewScanner(in)

	hasInput := false
	for scanner.Scan() {
		hasInput = true
		commandLine := scanner.Text()
		cmd, err := parse(commandLine)
		if err != nil {
			return nil, err
		}
		res = append(res, cmd)
	}

	if !hasInput {
		return nil, errors.New("input is empty")
	}

	return res, nil
}

func parse(str string) (painter.Operation, error) {
	fields := strings.Fields(str)
	if len(fields) == 0 {
		return nil, errors.New("empty command")
	}

	cmd := fields[0]
	options := make(map[string]float32)
	var f painter.OperationFunc

	switch cmd {
	case "white":
		f = painter.WhiteFill

	case "green":
		f = painter.GreenFill

	case "update":
		return painter.UpdateOp, nil

	case "bgrect":
		if len(fields) != 5 {
			return nil, errors.New("bgrect requires 4 coordinates")
		}
		x1, err := parseFloat(fields[1])
		if err != nil {
			return nil, err
		}
		y1, err := parseFloat(fields[2])
		if err != nil {
			return nil, err
		}
		x2, err := parseFloat(fields[3])
		if err != nil {
			return nil, err
		}
		y2, err := parseFloat(fields[4])
		if err != nil {
			return nil, err
		}
		options["x1"] = x1
		options["y1"] = y1
		options["x2"] = x2
		options["y2"] = y2
		f = painter.BGRect

	case "figure":
		if len(fields) != 3 {
			return nil, errors.New("figure requires 2 coordinates")
		}
		x, err := parseFloat(fields[1])
		if err != nil {
			return nil, err
		}
		y, err := parseFloat(fields[2])
		if err != nil {
			return nil, err
		}
		options["x"] = x
		options["y"] = y
		f = painter.Figure1

	case "move":
		if len(fields) != 3 {
			return nil, errors.New("move requires 2 coordinates")
		}
		x, err := parseFloat(fields[1])
		if err != nil {
			return nil, err
		}
		y, err := parseFloat(fields[2])
		if err != nil {
			return nil, err
		}
		options["x"] = x
		options["y"] = y
		f = painter.Move

	case "reset":
		f = painter.Reset

	default:
		return nil, errors.New("unknown command")
	}

	return painter.Command{
		F:       f,
		Options: options,
	}, nil
}

func parseFloat(s string) (float32, error) {
	val, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid float value: %s", s)
	}
	return float32(val), nil
}
