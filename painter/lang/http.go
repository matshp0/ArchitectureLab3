package lang

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/matshp0/ArchitectureLab3/painter"
)

// HttpHandler конструює обробник HTTP запитів, який дані з запиту віддає у Parser, а потім відправляє отриманий список
// операцій у painter.Loop.
func HttpHandler(loop *painter.Loop, p *Parser) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		var in io.Reader = r.Body
		if r.Method == http.MethodGet {
			in = strings.NewReader(r.URL.Query().Get("cmd"))
		}

		cmds, err := p.Parse(in)
		if err != nil {
			log.Printf("Bad script: %s", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, op := range cmds {
			fmt.Println(op)
			loop.Post(op)
		}
		rw.WriteHeader(http.StatusOK)
	})
}
