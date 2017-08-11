package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/graphql-go/graphql"
)

func graphStatFile(ctx context.Context, path string) (resp *FileInfo, err error) {

	fs := getFilesystem(ctx)

	// replace os.Stat with FileSystem read
	fsEntry, err := fs.Open(path)
	if os.IsNotExist(err) {
		err = NewStatError(http.StatusNotFound, path)
		return
	} else if err != nil {
		err = newError(http.StatusBadRequest, err)
		return
	}
	defer fsEntry.Close()

	stat, err := fsEntry.Stat()
	if err != nil {
		err = newError(http.StatusBadRequest, err)
		return
	}

	statType := "other"
	if stat.Mode().IsRegular() {
		statType = "file"
	} else if stat.Mode().IsDir() {
		statType = "directory"
	}

	resp = &FileInfo{
		Name:  stat.Name(),
		Type:  statType,
		Path:  path,
		Size:  stat.Size(),
		MTime: stat.ModTime(),
	}

	return
}

func graphListFiles(ctx context.Context, path string) (list []*FileInfo, err error) {
	fs := getFilesystem(ctx)

	// replace os.Stat with FileSystem read
	fsEntry, err := fs.Open(path)
	if os.IsNotExist(err) {
		err = NewStatError(http.StatusNotFound, path)
		return
	} else if err != nil {
		err = newError(http.StatusBadRequest, err)
		return
	}
	defer fsEntry.Close()

	stat, err := fsEntry.Stat()
	if err != nil {
		err = newError(http.StatusBadRequest, err)
		return
	}

	// permission problem
	if err != nil {
		perr, _ := err.(*os.PathError)
		if perr.Err.Error() == os.ErrPermission.Error() {
			err = NewStatError(http.StatusForbidden, path)
		}
		return
	}

	// for directories
	if stat.Mode().IsDir() {

		var d http.File
		files := make([]os.FileInfo, 0, 40)
		// TODO: use FileSystem for file access
		if d, err = fs.Open(path); err != nil {
			log.Printf("Error listing path %#v:%s", path, err)
			err = NewStatError(http.StatusInternalServerError, path)
			return
		}
		defer d.Close()

		files, err = d.Readdir(0)
		if err != nil {
			log.Printf("Error listing path %#v:%s", path, err)
			return
		}

		// sort according to query
		epCtx := getEndpointContext(ctx)
		s := epCtx.Sort
		if s == "" {
			s = "-mtime"
		}
		// TODO: rewrite with go-linq
		//QuerySort(s, files) // TODO: add error reporting here

		listLen := len(files)
		list = make([]*FileInfo, listLen)
		for i := 0; i < listLen; i++ {
			item := files[i]

			// parse item URL
			itemPath := path + "/" + item.Name()
			if path == "." {
				itemPath = item.Name()
			}

			itemType := "other"
			if item.Mode().IsRegular() {
				itemType = "file"
			} else if item.IsDir() {
				itemType = "directory"
			}

			list[i] = &FileInfo{
				Name:  item.Name(),
				Type:  itemType,
				Path:  itemPath,
				Size:  item.Size(),
				MTime: item.ModTime(),
			}
		}
		return
	}
	err = NewStatError(http.StatusBadRequest, path)
	return
}

type endpointError struct {
	code int
	err  error
}

func (err endpointError) Error() string {
	return err.err.Error()
}

func newError(code int, err error) *endpointError {
	return &endpointError{
		code: code,
		err:  err,
	}
}

func newInputError(err error) *endpointError {
	return &endpointError{
		code: http.StatusBadRequest,
		err:  err,
	}
}

func parseCode(err error) int {
	switch tErr := err.(type) {
	case endpointError:
		return tErr.code
	case *endpointError:
		return tErr.code
	default:
		return http.StatusInternalServerError
	}
}

func getSchema() (graphql.Schema, error) {

	// Recipe type
	fileInfoType := graphql.NewObject(graphql.ObjectConfig{
		Name:        "FileStat",
		Description: "Information about a file or directory",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"path": &graphql.Field{
				Type: graphql.String,
			},
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"size": &graphql.Field{
				Type: graphql.Int,
			},
			"mtime": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	fileInfosType := graphql.NewList(fileInfoType)

	fileInfoType.AddFieldConfig("children", &graphql.Field{
		Type: fileInfosType,
		Resolve: func(p graphql.ResolveParams) (resp interface{}, err error) {
			switch src := p.Source.(type) {
			case *FileInfo:
				if src.Type == "directory" {
					resp, err = graphListFiles(p.Context, src.Path)
				} else {
					resp = []int{}
				}
				return
			}
			return
		},
	})

	// Root Query Schema
	rootQuery := graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"list": &graphql.Field{
				Type:        fileInfosType,
				Description: "List of files and directories within a directory of given path",
				Args: graphql.FieldConfigArgument{
					"path": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (resp interface{}, err error) {
					path := p.Args["path"].(string)
					resp, err = graphListFiles(p.Context, path)
					return
				},
			},
			"stat": &graphql.Field{
				Type:        fileInfoType,
				Description: "Information about a file or a directory",
				Args: graphql.FieldConfigArgument{
					"path": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (resp interface{}, err error) {
					path := p.Args["path"].(string)
					resp, err = graphStatFile(p.Context, path)
					return
				},
			},
		},
	}
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}

	return graphql.NewSchema(schemaConfig)
}

type graphPostRequest struct {
	Query         string                 `json:"query"`
	OperationName string                 `json:"operationName"`
	Variables     map[string]interface{} `json:"variables"`
}

func deccodeGraphRequest(ctx context.Context, r *http.Request) (req interface{}, err error) {
	vreq := &graphPostRequest{}
	switch r.Method {
	case "GET":

		// get query
		vreq.Query = r.URL.Query().Get("query")
		if vreq.Query == "" {
			err = newInputError(fmt.Errorf("requires argument 'query' in GET request"))
			return
		}

		// get operation name
		vreq.OperationName = r.URL.Query().Get("operationName")

		// get variables
		if variables := r.URL.Query().Get("variables"); variables != "" {
			err = json.Unmarshal([]byte(variables), vreq.Variables)
		}
		if err != nil {
			err = newInputError(fmt.Errorf("error decoding 'variables' in GET request"))
			return
		}

		// return request
		req = vreq
		return
	case "POST":
		contentType := r.Header.Get("Content-Type")
		switch contentType {
		case "application/json":
			fallthrough
		default:
			dec := json.NewDecoder(r.Body)
			err = dec.Decode(vreq)
			if err != nil {
				err = newInputError(err)
				return
			}

			// return request
			req = vreq
			return
		}
	}
	return
}

func encodeGraphErrorResponse(ctx context.Context, err error, w http.ResponseWriter) {
	errCode := parseCode(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errCode)
	enc := json.NewEncoder(w)
	enc.Encode(struct {
		Code    int    `json:"code"`
		Status  string `json:"status"`
		Error   string `json:"error"`
		Message string `json:"message"`
	}{
		Code:    errCode,
		Status:  "error",
		Error:   http.StatusText(errCode),
		Message: err.Error(),
	})
}

func encodeGraphResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	return enc.Encode(resp)
}

func graphEndpoint(ctx context.Context, req interface{}) (resp interface{}, err error) {

	schema, err := getSchema()
	if err != nil {
		return
	}

	// get request string
	vreq, ok := req.(*graphPostRequest)
	if !ok {
		err = newInputError(fmt.Errorf("graphEndpoint expect req to be a string, got %#v instead", req))
		return
	}
	params := graphql.Params{
		Schema:         schema,
		RequestString:  vreq.Query,
		Context:        ctx,
		OperationName:  vreq.OperationName,
		VariableValues: vreq.Variables,
	}
	graphResp := graphql.Do(params)
	if len(graphResp.Errors) > 0 {
		err = newInputError(graphResp.Errors[0])
		return
	}
	resp = graphResp
	return
}

// GraphQLHandler returns http.Handler for the
// graphql endpoint
func GraphQLHandler() http.Handler {
	return httptransport.NewServer(
		graphEndpoint,
		deccodeGraphRequest,
		encodeGraphResponse,
		httptransport.ServerErrorEncoder(encodeGraphErrorResponse),
	)
}
