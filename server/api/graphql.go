package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/graphql-go/graphql"
	linq "gopkg.in/ahmetb/go-linq.v3"
)

func graphStatFile(ctx context.Context, filepath string) (resp *FileInfo, err error) {

	fs := getFilesystem(ctx)

	// replace os.Stat with FileSystem read
	fsEntry, err := fs.Open(filepath)
	if os.IsNotExist(err) {
		err = NewStatError(http.StatusNotFound, filepath)
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
	mimeType := ""
	statName := stat.Name()
	if stat.Mode().IsRegular() {
		statType = "file"
		mimeType = mime.TypeByExtension(strings.ToLower(path.Ext(stat.Name())))
	} else if stat.Mode().IsDir() {
		statType = "directory"
	}

	if filepath == "" {
		statName = "/"
	}

	resp = &FileInfo{
		Name:  statName,
		Type:  statType,
		Mime:  mimeType,
		Path:  "/" + filepath,
		Size:  stat.Size(),
		MTime: stat.ModTime(),
	}

	return
}

func hasIndex(fs http.FileSystem, filepath string) bool {
	fileIndex := path.Join(filepath, "index.html")
	fi, err := fs.Open(fileIndex)
	if err != nil {
		return false
	}
	defer fi.Close()
	return true
}

func graphListFiles(ctx context.Context, filepath string) (list []*FileInfo, err error) {
	fs := getFilesystem(ctx)
	graphCtx := getGraphContext(ctx)
	args := graphCtx.Args

	// replace os.Stat with FileSystem read
	fsEntry, err := fs.Open(filepath)
	if os.IsNotExist(err) {
		err = NewStatError(http.StatusNotFound, filepath)
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
			err = NewStatError(http.StatusForbidden, filepath)
		}
		return
	}

	// for directories
	if stat.Mode().IsDir() {

		var d http.File
		files := make([]os.FileInfo, 0, 40)
		if d, err = fs.Open(filepath); err != nil {
			log.Printf("Error listing filepath %#v:%s", filepath, err)
			err = NewStatError(http.StatusInternalServerError, filepath)
			return
		}
		defer d.Close()

		files, err = d.Readdir(0)
		if err != nil {
			log.Printf("Error listing filepath %#v:%s", filepath, err)
			return
		}

		listLen := len(files)
		list = make([]*FileInfo, listLen)
		for i := 0; i < listLen; i++ {
			item := files[i]

			// parse item URL
			itemPath := filepath + "/" + item.Name()
			if filepath == "." || filepath == "" {
				itemPath = item.Name()
			}

			itemType := "other"
			mimeType := ""
			if item.Mode().IsRegular() {
				itemType = "file"
				mimeType = mime.TypeByExtension(strings.ToLower(path.Ext(item.Name())))
			} else if item.IsDir() {
				itemType = "directory"
			}

			list[i] = &FileInfo{
				Name:     item.Name(),
				Type:     itemType,
				Mime:     mimeType,
				HasIndex: hasIndex(fs, itemPath),
				Path:     "/" + itemPath,
				Size:     item.Size(),
				MTime:    item.ModTime(),
			}
		}

		s := "-mtime"
		nameLike := ""
		nameLikeMe := false

		if args != nil {
			// sort list according to query
			if argSort, ok := args["sort"]; ok && argSort != "" {
				argSortString, ok := argSort.(string)
				if ok {
					s = argSortString
				}
			}
			// filters: nameLike
			if argNameLike, ok := args["nameLike"]; ok && argNameLike != "" {
				argNameLikeStr := argNameLike.(string)
				nameLike = argNameLikeStr
			}
			// filters: nameLikeMe
			if argNameLikeMe, ok := args["nameLikeMe"]; ok && argNameLikeMe != "" {
				nameLikeMe, _ = argNameLikeMe.(bool)
			}
		}

		ss := strings.Split(s, ",")

		op := linq.From(list)
		var orderedOp linq.OrderedQuery

		if nameLikeMe && graphCtx.Source != nil {
			src := graphCtx.Source.(*FileInfo)
			nameBase := path.Base(src.Name)
			if src.Type != "directory" {
				nameBase = nameBase[:len(src.Name)-len(path.Ext(src.Name))]
			}

			// add where to linq
			op = op.WhereT(func(fi *FileInfo) bool {
				return strings.HasPrefix(fi.Name, nameBase)
			})
		}

		if nameLike != "" {
			var nameLikeRE *regexp.Regexp
			nameLike = strings.Replace(nameLike, "*", ".*", -1)
			nameLikeRE, err = regexp.Compile(nameLike)
			if err != nil {
				err = newError(http.StatusBadRequest, err)
				return
			}

			// add where to linq
			op = op.WhereT(func(fi *FileInfo) bool {
				return nameLikeRE.MatchString(fi.Name)
			})
		}

		for i, sort := range ss {

			// determine order direction
			isAsc := true
			if sort[0] == '-' {
				isAsc = false
				sort = sort[1:]
			}

			// determine field
			var selectorFn interface{}
			switch sort {
			case "mtime":
				selectorFn = func(fi *FileInfo) int64 {
					return fi.MTime.Unix()
				}
			case "name":
				selectorFn = func(fi *FileInfo) string {
					return strings.ToLower(fi.Name)
				}
			case "type":
				selectorFn = func(fi *FileInfo) string {
					return fi.Type
				}
			}

			// do first order sorting
			if i == 0 {
				if isAsc {
					orderedOp = op.OrderByT(selectorFn)
					continue
				}
				orderedOp = op.OrderByDescendingT(selectorFn)
				continue
			}

			// do later order sorting
			if isAsc {
				orderedOp = orderedOp.ThenByT(selectorFn)
				continue
			}
			orderedOp = orderedOp.ThenByDescendingT(selectorFn)
		}

		orderedOp.ToSlice(&list)
		return
	}
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
			"mime": &graphql.Field{
				Type: graphql.String,
			},
			"hasIndex": &graphql.Field{
				Type: graphql.Boolean,
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
		Args: graphql.FieldConfigArgument{
			"nameLike": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "string, with wildcard *, to match file name",
			},
			"nameLikeMe": &graphql.ArgumentConfig{
				Type:        graphql.Boolean,
				Description: "bool, if the children has a filename (without extension) prefixed with name of self",
			},
		},
		Resolve: func(p graphql.ResolveParams) (resp interface{}, err error) {
			switch src := p.Source.(type) {
			case *FileInfo:
				resp, err = graphListFiles(withGraphContext(p.Context, &graphContext{
					Source: p.Source,
					Args:   p.Args,
				}), src.Path)
			}
			return
		},
	})
	fileInfoType.AddFieldConfig("siblings", &graphql.Field{
		Type: fileInfosType,
		Args: graphql.FieldConfigArgument{
			"nameLike": &graphql.ArgumentConfig{
				Type:        graphql.String,
				Description: "string, with wildcard *, to match file name",
			},
			"nameLikeMe": &graphql.ArgumentConfig{
				Type:        graphql.Boolean,
				Description: "bool, if the sibling has a filename (without extension) prefixed with name of self",
			},
		},
		Resolve: func(p graphql.ResolveParams) (resp interface{}, err error) {
			switch src := p.Source.(type) {
			case *FileInfo:
				dir := strings.TrimLeft(path.Dir(src.Path), "/")
				resp, err = graphListFiles(withGraphContext(p.Context, &graphContext{
					Source: p.Source,
					Args:   p.Args,
				}), dir)
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
					"nameLike": &graphql.ArgumentConfig{
						Type:        graphql.String,
						Description: "string, with wildcard *, to match file name",
					},
					"path": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
					"sort": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (resp interface{}, err error) {
					path := p.Args["path"].(string)
					path = strings.TrimLeft(path, "/")
					resp, err = graphListFiles(withGraphContext(p.Context, &graphContext{
						Args: p.Args,
					}), path)
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
					path = strings.TrimLeft(path, "/")
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
