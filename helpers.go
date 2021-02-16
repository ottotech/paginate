package paginate

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode"
)

func getParameters(colNames, filters []string, mappers mappers, u url.URL) parameters {
	list := make(parameters, 0)
	decodedURL, _ := url.PathUnescape(u.String())

	getParameter := func(key, val, char string) (bool, parameter) {
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

	i := strings.Index(decodedURL, "?")
	if i == -1 {
		return list
	}

	params := strings.Split(decodedURL[i+1:], "&")

	for _, colName := range colNames {

		// If colName is not in filters we will not try to build
		// a where clause condition
		if !isStringIn(colName, filters) {
			continue
		}

		columnIsMapped, customParameterName := mappers.isColumnMapped(colName)

		for _, requestParameterAndValue := range params {
			// ``key`` and ``value`` will be the variables used
			// to build sql where clauses.
			var key, value string

			if columnIsMapped {
				if len(requestParameterAndValue) <= len(customParameterName) {
					continue
				}
				key, value = requestParameterAndValue[:len(customParameterName)], requestParameterAndValue[len(customParameterName):]
				if key != customParameterName {
					continue
				}
				// We do this to map the column name with the parameter
				// from the request. A ``key`` should be always a column name.
				// But if the user specified its customParameter we know that the
				// above key is not a column name, so that's why we override
				// the value here.
				key = colName
			} else {
				if len(requestParameterAndValue) <= len(colName) {
					continue
				}
				key, value = requestParameterAndValue[:len(colName)], requestParameterAndValue[len(colName):]
				if key != colName {
					continue
				}
			}

			// order matters
			if ok, newP := getParameter(key, value, gte); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getParameter(key, value, lte); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getParameter(key, value, ne); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getParameter(key, value, gt); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getParameter(key, value, lt); ok {
				list = append(list, newP)
				continue
			}
			if ok, newP := getParameter(key, value, eq); ok {
				list = append(list, newP)
				continue
			}
		}
	}

	// As an special case if there are repeated parameters with the ``eq`` sign
	// we will group them together under the ``_in`` sign, in order to allow
	// the creation of an sql IN clause.
	listWithoutEqualDuplicates := make(parameters, 0)
	eqDuplicates := make(map[string][]string, 0)
	for i, x := range list {
		if x.sign != eq {
			listWithoutEqualDuplicates = append(listWithoutEqualDuplicates, x)
			continue
		}

		isDuplicated := false

		for j, y := range list {
			if i == j {
				continue
			}
			if x.sign == y.sign && x.name == y.name {
				if _, exists := eqDuplicates[x.name]; exists {
					if !isStringIn(x.value, eqDuplicates[x.name]) {
						eqDuplicates[x.name] = append(eqDuplicates[x.name], x.value)
					}
				} else {
					eqDuplicates[x.name] = append(eqDuplicates[x.name], x.value)
				}
				isDuplicated = true
			}
		}
		if !isDuplicated {
			listWithoutEqualDuplicates = append(listWithoutEqualDuplicates, x)
		}
	}

	list = listWithoutEqualDuplicates

	for k, v := range eqDuplicates {
		inVal := strings.Join(v, ",")
		p := parameter{
			name:  k,
			sign:  _in,
			value: inVal,
		}
		list = append(list, p)
	}

	// As an special case if there are repeated parameters with the ``ne`` sign
	// we will group them together under the ``_notin`` sign, in order to allow
	// the creation of an sql NOT IN clause.
	listWithoutNotEqualDuplicates := make(parameters, 0)
	notEqDuplicates := make(map[string][]string, 0)
	for i, x := range list {
		if x.sign != ne {
			listWithoutNotEqualDuplicates = append(listWithoutNotEqualDuplicates, x)
			continue
		}

		isDuplicated := false

		for j, y := range list {
			if i == j {
				continue
			}
			if x.sign == y.sign && x.name == y.name {
				if _, exists := notEqDuplicates[x.name]; exists {
					if !isStringIn(x.value, notEqDuplicates[x.name]) {
						notEqDuplicates[x.name] = append(notEqDuplicates[x.name], x.value)
					}
				} else {
					notEqDuplicates[x.name] = append(notEqDuplicates[x.name], x.value)
				}
				isDuplicated = true
			}
		}
		if !isDuplicated {
			listWithoutNotEqualDuplicates = append(listWithoutNotEqualDuplicates, x)
		}
	}

	list = listWithoutNotEqualDuplicates

	for k, v := range notEqDuplicates {
		inVal := strings.Join(v, ",")
		p := parameter{
			name:  k,
			sign:  _notin,
			value: inVal,
		}
		list = append(list, p)
	}

	// As an special case we need to also get our custom sort parameter.
	sort := "sort"
	for _, p := range params {
		if len(p) <= len(sort) {
			continue
		}
		key, value := p[:len(sort)], p[len(sort):]
		if key != sort {
			continue
		}
		if ok, newP := getParameter(key, value, eq); ok {
			list = append(list, newP)
			continue
		}
	}

	return list
}

func getRequestData(v url.Values) paginationRequest {
	p := paginationRequest{}
	if page := v.Get("page"); page != "" {
		page, err := strconv.Atoi(page)
		if err != nil {
			page = defaultPageNumber
		}
		if page <= 0 {
			page = defaultPageNumber
		}
		p.pageNumber = page
	} else {
		p.pageNumber = defaultPageNumber
	}

	if pageSize := v.Get("page_size"); pageSize != "" {
		pageSize, err := strconv.Atoi(pageSize)
		if err != nil {
			pageSize = defaultPageSize
		}
		if pageSize <= 0 {
			pageSize = defaultPageSize
		}
		p.pageSize = pageSize
	} else {
		p.pageSize = defaultPageSize
	}

	return p
}

func createWhereClause(dialect string, colNames []string, params parameters, extraWhereClauses []RawWhereClause, c chan whereClause) {
	w := whereClause{}
	var WHERE = " WHERE "
	var AND = " AND "
	var separator string
	var clauses []string
	var values []interface{}

	for _, name := range colNames {
		for _, p := range params {
			if p.name == name {
				switch p.sign {
				case _in, _notin:
					vals := strings.Split(p.value, ",")
					for _, v := range vals {
						values = append(values, v)
					}
					placeholder := dialectPlaceholder.GetPlaceHolder(dialect)
					str := ""
					for i := 0; i < len(vals); i++ {
						if i == len(vals)-1 {
							str += placeholder
						} else {
							str += placeholder + ","
						}
					}
					clauses = append(clauses, p.name+" "+p.sign+fmt.Sprintf("(%s)", str))
				default:
					values = append(values, p.value)
					clauses = append(
						clauses,
						fmt.Sprintf("%s %s %s", p.name, p.sign, dialectPlaceholder.GetPlaceHolder(dialect)),
					)
				}
			}
		}
	}

	// If there are extra custom where clauses we append them here.
	for _, predicate := range extraWhereClauses {
		clauses = append(clauses, predicate.String())
		values = append(values, predicate.args...)
	}

	// Let's use an appropriate `separator` to join the clauses.
	if len(clauses) == 1 {
		separator = ""
	} else {
		separator = AND
	}

	w.clause = WHERE + strings.Join(clauses, separator)
	w.args = values
	w.exists = len(clauses) > 0
	c <- w
}

func createPaginationClause(pageNumber int, pageSize int, c chan string) {
	var clause string
	var offset int

	clause += fmt.Sprintf(" LIMIT %v ", pageSize)

	if pageNumber < 0 || pageNumber == 0 || pageNumber == 1 {
		offset = 0
	} else {
		offset = pageSize * (pageNumber - 1)
	}

	clause += fmt.Sprintf("OFFSET %v", offset)

	c <- clause
}

func createOrderByClause(params parameters, colNames []string, id string, c chan string) {
	var ASC = "ASC"
	var DESC = "DESC"

	clauses := make([]string, 0)

	sort, exists := params.getParameter("sort")
	if !exists {
		c <- fmt.Sprintf(" ORDER BY %s", id)
		return
	}

	fields := strings.Split(sort.value, ",")
	for _, v := range fields {
		orderBy := string(v[0])
		field := v[1:]
		for _, f := range colNames {
			if f == id {
				// we will always order the records by ID (see below). In order
				// to keep the same order between pages or results deterministic.
				// See: https://use-the-index-luke.com/sql/partial-results/fetch-next-page
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

	clauses = append(clauses, id)
	clauseSTR := strings.Join(clauses, ",")
	c <- " ORDER BY " + clauseSTR
}

// parseCamelCaseToSnakeLowerCase parses a camelcase string to a snake case
// lower cased. So for example, if we use as input for this function the following
// string "myCamelCaseVar" the output would be "my_camel_case_var".
func parseCamelCaseToSnakeLowerCase(camelCase string) string {
	var s []string
	var lastIndexChecked int
	var isSneakCase = true
	var allAreUpperCase = true

	for j := 0; j < len(camelCase); j++ {
		if !unicode.IsUpper(rune(camelCase[j])) {
			allAreUpperCase = false
		}
	}

	if allAreUpperCase {
		return strings.ToLower(camelCase)
	}

	for i := len(camelCase) - 1; i >= 0; i-- {
		if unicode.IsUpper(rune(camelCase[i])) {
			s = append(s, camelCase[i:])
			camelCase = camelCase[:i]
			lastIndexChecked = i
			isSneakCase = false
		}

		if i == 0 && lastIndexChecked > 0 {
			s = append(s, camelCase[:lastIndexChecked])
		}
	}

	if isSneakCase {
		return strings.ToLower(camelCase)
	}

	var orderedSlice []string
	for i := len(s) - 1; i >= 0; i-- {
		orderedSlice = append(orderedSlice, s[i])
	}

	return strings.ToLower(strings.Join(orderedSlice, "_"))
}

// isStringIn checks whether the given string ``s`` is in the given slice ``in``.
func isStringIn(s string, in []string) bool {
	for _, elem := range in {
		if s == elem {
			return true
		}
	}
	return false
}
