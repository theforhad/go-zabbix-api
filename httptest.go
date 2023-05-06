package zabbix

// httptest represent Zabbix httptest type returned from Zabbix API
// https://www.zabbix.com/documentation/3.2/en/manual/api/reference/httptest/object
type HttpTest struct {
	HttpTestId      string       `json:"httptestid,omitempty"`
	HostId          string       `json:"hostid,omitempty"`
	TemplatesClear  TemplateIDs  `json:"templates_clear,omitempty"`
	Description     string       `json:"description,omitempty"`
	Name            string       `json:"name,omitempty"`
	Steps						Steps			 	 `json:"steps"`
}

type Step struct {
	Name  					string `json:"name"`
	Url							string `json:"url"`
	StatusCodes			string `json:"status_codes"`
	FollowRedirects string `json:"follow_redirects,omitempty"`
	No 							string `json:"no,omitempty"`
	Headers 				Headers `json:"headers,omitempty"`
}

type Header struct{
	Name string  `json:"name"`
	Value string  `json:"value"`
}

type Headers []Header

type Steps []Step

// Templates is an Array of Template structs.
type HttpTests []HttpTest

// HttpTestGet Wrapper for httptest.get
// https://www.zabbix.com/documentation/current/en/manual/api/reference/httptest/object
func (api *API) HttpTestGet(params Params) (res HttpTests, err error) {
	if _, present := params["output"]; !present {
		params["output"] = "extend"
	}
	err = api.CallWithErrorParse("httptest.get", params, &res)
	return
}

// HttpTestGetByID Gets httptest(webscenario) by Id only if there is exactly 1 matching web scenario.
func (api *API) HttpTestGetByID(id string) (httptest *HttpTest, err error) {
	httptests, err := api.HttpTestGet(Params{"httptestid": id})
	if err != nil {
		return
	}

	if len(httptests) == 1 {
		httptest = &httptests[0]
	} else {
		e := ExpectedOneResult(len(httptests))
		err = &e
	}
	return
}

// HttpTestCreate Wrapper for httptest.create
// https://www.zabbix.com/documentation/3.2/en/manual/api/reference/httptest/create
func (api *API) HttpTestCreate(httptests HttpTests) (err error) {
	response, err := api.CallWithError("httptest.create", httptests)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	httptestids := result["httptestids"].([]interface{})
	
	for i, id := range httptestids {
		httptests[i].HttpTestId = id.(string)
	}
	return
}

// HttpTestUpdate Wrapper for httptest.update
// https://www.zabbix.com/documentation/3.2/manual/api/reference/httptest/update
func (api *API) HttpTestUpdate(httptests HttpTests) (err error) {
	_, err = api.CallWithError("httptest.update", httptests)
	return
}

// HttpTestDelete Wrapper for httptest.delete
// https://www.zabbix.com/documentation/3.2/manual/api/reference/httptest/delete
func (api *API) HttpTestDelete(httptests HttpTests) (err error) {
	httptestids := make([]string, len(httptests))
	for i, httptest := range httptests {
		httptestids[i] = httptest.HttpTestId
	}

	err = api.HttpTestDeleteByIds(httptestids)
	if err == nil {
		for i := range httptestids {
			httptests[i].HttpTestId = ""
		}
	}
	return
}

// HttpTestDeleteByIds Wrapper for httptest.delete
// Use httptest's id to delete the httptest
// https://www.zabbix.com/documentation/3.2/manual/api/reference/httptest/delete
func (api *API) HttpTestDeleteByIds(ids []string) (err error) {
	response, err := api.CallWithError("httptest.delete", ids)
	if err != nil {
		return
	}

	result := response.Result.(map[string]interface{})
	httptestids := result["httptestids"].([]interface{})
	if len(ids) != len(httptestids) {
		err = &ExpectedMore{len(ids), len(httptestids)}
	}
	return
}
