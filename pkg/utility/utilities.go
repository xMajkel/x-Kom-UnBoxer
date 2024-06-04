package utility

import (
	"compress/gzip"
	"io"
	"log"
	"strconv"
	"strings"
)

func ReadHttpResponseBody(contentEncoding string, body io.ReadCloser) (r io.ReadCloser, err error) {
	switch contentEncoding {
	case "gzip":
		r, err = gzip.NewReader(body)
		if err != nil {
			return r, err
		}
		return r, nil
	default:
		return body, nil
	}
}

func ParsePreferredRollTime(t string) (h int, m int) {
	var err error

	timeSplit := strings.Split(t, ":")

	if len(timeSplit) != 2 {
		log.Print("[ERROR] \"preferred_roll_time\" should be in such format: hh:mm; using 00:00\n")
		h = 0
		m = 0
	} else {
		h, err = strconv.Atoi(timeSplit[0])
		if err != nil {
			panic(err)
		}
		m, err = strconv.Atoi(timeSplit[1])
		if err != nil {
			panic(err)
		}
	}
	return h, m
}
