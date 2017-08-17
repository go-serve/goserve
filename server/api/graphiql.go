package api

import (
	"bytes"
	"net/http"
	"text/template"
)

// GraphiqlHTML is the html template for serving the graphiql page
const GraphiqlHTML = `<!DOCTYPE html>
<head>
  <style>body {height: 100vh; margin: 0; width: 100%; overflow: hidden;}</style>
	<title>GraphiQL for goserve</title>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.2/graphiql.min.css" />
	<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.6.1/react.min.js"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/react/15.6.1/react-dom.min.js"></script>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/graphiql/0.11.2/graphiql.min.js"></script>
  <script>
    (function () {
      document.addEventListener('DOMContentLoaded', function () {
        endpoint = window.location.origin + '{{ .Endpoint }}';
        function fetcher(params) {
          var options = {
            method: 'post',
            headers: {'Accept': 'application/json', 'Content-Type': 'application/json'},
            body: JSON.stringify(params),
            credentials: 'include',
          };
          return fetch(endpoint, options)
            .then(function (res) { return res.json() });
        }
        var body = React.createElement(GraphiQL, {fetcher: fetcher, query: '', variables: ''});
        ReactDOM.render(body, document.body);
      });
    }());
  </script>
</head>
<body>
</body>
`

// GraphiQLHandler generates a handler function for HTTP servers
// given the GraphQL endpoint
func GraphiQLHandler(endpoint string) http.HandlerFunc {

	var err error
	tpl := template.New("graphiql.html")
	if tpl, err = tpl.Parse(GraphiqlHTML); err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer([]byte{})
	err = tpl.Execute(buf, map[string]interface{}{
		"Endpoint": endpoint,
	})
	if err != nil {
		panic(err)
	}
	contents := []byte(buf.String())

	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")
		w.Write(contents)
	}
}
