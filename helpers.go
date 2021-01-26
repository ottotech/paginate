package paginate

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode"
)

func getParameters(colNames []string, u url.URL, c chan parameters) {
	decodedURL, _ := url.PathUnescape(u.String())
	list := make(parameters, 0)
	i := strings.Index(decodedURL, "?")
	if i == -1 {
		c <- list
		return
	}

	getP := func(key, val, char string) (bool, parameter) {
		p := parameter{}
		if strings.Contains(val, char) {
			if val[:len(char)] == char && len(val) > len(char) {
				p.name = key
				p.sign = char
				p.value = val[len(char):]
				return true, p
			}
		}
		return false, p
	}
	params := strings.Split(decodedURL[i+1:], "&")
	for _, n := range colNames {
		for _, p := range params {
			if len(p) <= len(n) {
				continue
			}
			key, value := p[:len(n)], p[len(n):]
			if key != n {
				continue
			}
			// order matters
			if ok, newP := getP(key, value, gte); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getP(key, value, lte); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getP(key, value, ne); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getP(key, value, gt); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getP(key, value, lt); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getP(key, value, eq); ok {
				list = append(list, newP)
				continue
			}
		}
	}
	// as an special case we need to also get our custom sort parameter
	sort := "sort"
	for _, p := range params {
		if len(p) <= len(sort) {
			continue
		}
		key, value := p[:len(sort)], p[len(sort):]
		if key != sort {
			continue
		}
		if ok, newP := getP(key, value, eq); ok {
			list = append(list, newP)
			continue
		}
	}
	c <- list
}

func getRequestData(v url.Values) paginationRequest {
	p := paginationRequest{}
	if page := v.Get("page"); page != "" {
		page, err := strconv.Atoi(page)
		if err != nil {
			page = 1
		}
		if page <= 0 {
			page = 1
		}
		p.pageNumber = page
	} else {
		p.pageNumber = 1
	}

	if pageSize := v.Get("page_size"); pageSize != "" {
		pageSize, err := strconv.Atoi(pageSize)
		if err != nil {
			pageSize = PageSize
		}
		if pageSize <= 0 {
			pageSize = PageSize
		}
		p.pageSize = pageSize
	} else {
		p.pageSize = PageSize
	}
	return p
}

func createWhereClause(colNames []string, params parameters, c chan whereClause) {
	w := whereClause{}
	var WHERE = " WHERE "
	var AND = " AND "
	var separator string
	var clauses []string
	var values []interface{}

	// map all db column names with the url parameters
	for _, name := range colNames {
		for _, p := range params {
			if p.name == name {
				values = append(values, p.value)
				clauses = append(clauses, p.name+" "+p.sign+" $%v")
			}
		}
	}
	// use appropriate `separator` to join the clauses
	if len(clauses) == 1 {
		separator = ""
	} else {
		separator = AND
	}
	// let's map the clause and args to the whereClause struct, and specify if there were some where clauses at all
	w.clause = WHERE + strings.Join(clauses, separator)
	w.args = values
	w.exists = len(clauses) > 0
	c <- w
}

func createPaginationClause(pageNumber int, pageSize int, c chan string) {
	var clause string
	var offset int

	if pageSize > PageSize {
		pageSize = PageSize
	} else if pageSize < 0 {
		pageSize = PageSize
	}

	clause += fmt.Sprintf(" LIMIT %v ", pageSize)

	if pageNumber < 0 || pageNumber == 0 || pageNumber == 1 {
		offset = 0
	} else {
		offset = pageSize * (pageNumber - 1)
	}

	clause += fmt.Sprintf("OFFSET %v", offset)
	c <- clause
}

func createOrderByClause(params parameters, colNames []string, c chan string) {
	var ASC = "ASC"
	var DESC = "DESC"
	clauses := make([]string, 0)
	sort, exists := params.getParameter("sort")
	if !exists {
		c <- " ORDER BY id"
		return
	}
	fields := strings.Split(sort.value, ",")
	for _, v := range fields {
		orderBy := string(v[0])
		field := v[1:]
		for _, f := range colNames {
			if f == "id" {
				// we will always order the records by ID (see below). In order
				// to keep the same order between pages or results
				continue
			}
			if field == f {
				if orderBy == "+" {
					clauses = append(clauses, field+" "+ASC)
				}
				if orderBy == "-" {
					clauses = append(clauses, field+" "+DESC)
				}
			}
		}
	}
	clauses = append(clauses, "id")
	clauseSTR := strings.Join(clauses, ",")
	c <- " ORDER BY " + clauseSTR
}

// parseCamelCaseToSnake parses a camelcase string to a snake case.
// So for example, if we use as input for this function the following
// string "myCamelCaseVar" the output would be "my_camel_case_var".
func parseCamelCaseToSnake(camelCase string) string {
	var s []string
	for i := len(camelCase) - 1; i >= 0; i-- {
		if unicode.IsUpper(rune(camelCase[i])) {
			s = append(s, camelCase[i:])
			camelCase = camelCase[:i]
		}
	}
	var orderedSlice []string
	for i := len(s) - 1; i >= 0; i-- {
		orderedSlice = append(orderedSlice, s[i])
	}
	return strings.ToLower(strings.Join(orderedSlice, "_"))
}
