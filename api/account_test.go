package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	mockdb "github.com/emon46/bank-application/db/mock"
	db "github.com/emon46/bank-application/db/sqlc"
	"github.com/emon46/bank-application/token"
	"github.com/emon46/bank-application/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetAccountApi(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		accountId     int64
		setupAuth     func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "ok",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, r, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "not found",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, r, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "internal error",
			accountId: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, r, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "invalid id",
			accountId: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, r, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			testStore := mockdb.NewMockStore(ctrl)
			// stubs
			testCase.buildStubs(testStore)
			// start test server and send request
			server := newTestServer(t, testStore)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts/%d", testCase.accountId)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)
			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})
	}

}

func TestCreateAccountAPI(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)

	testCases := []struct {
		name          string
		body          gin.H
		setupAuth     func(t *testing.T, r *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				// just call the add authorization which we add the token bearer in the request body
				addAuthorization(t, r, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateAccountParams{
					Owner:    account.Owner,
					Balance:  0,
					Currency: account.Currency,
				}
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "Connection check",
			body: gin.H{
				"currency": account.Currency,
			},
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				// just call the add authorization which we add the token bearer in the request body
				addAuthorization(t, r, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "Invalid currency",
			body: gin.H{
				"currency": "invalid",
			},
			setupAuth: func(t *testing.T, r *http.Request, tokenMaker token.Maker) {
				// just call the add authorization which we add the token bearer in the request body
				addAuthorization(t, r, tokenMaker, authorizationTypeBearer, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			testStore := mockdb.NewMockStore(ctrl)
			testCase.buildStubs(testStore)

			//create a test server with test store
			server := newTestServer(t, testStore)
			recorder := httptest.NewRecorder()

			// need tp create a recorder to record response
			data, err := json.Marshal(testCase.body)
			require.NoError(t, err)

			url := fmt.Sprintf("/accounts")
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			testCase.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			testCase.checkResponse(t, recorder)
		})

	}
}

//func TestListAccountsAPI(t *testing.T) {
//	user := randomUser()
//
//	n := 5
//	accounts := make([]db.Account, n)
//	for i := 0; i < n; i++ {
//		accounts[i] = randomAccount(user)
//	}
//	type Query struct {
//		pageID   int
//		pageSize int
//	}
//
//	testCases := []struct {
//		name          string
//		query         Query
//		buildStubs    func(store *mockdb.MockStore)
//		checkResponse func(recoder *httptest.ResponseRecorder)
//	}{
//		{
//			name: "OK",
//			query: Query{
//				pageID:   1,
//				pageSize: n,
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				arg := db.ListAccountParams{
//					Limit:  int32(n),
//					Offset: 0,
//				}
//
//				store.EXPECT().
//					ListAccount(gomock.Any(), gomock.Eq(arg)).
//					Times(1).
//					Return(accounts, nil)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusOK, recorder.Code)
//				requireBodyMatchAccounts(t, recorder.Body, accounts)
//			},
//		},
//		{
//			name: "InternalError",
//			query: Query{
//				pageID:   1,
//				pageSize: n,
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					ListAccount(gomock.Any(), gomock.Any()).
//					Times(1).
//					Return([]db.Account{}, sql.ErrConnDone)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusInternalServerError, recorder.Code)
//			},
//		},
//		{
//			name: "InvalidPageID",
//			query: Query{
//				pageID:   -1,
//				pageSize: n,
//			},
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					ListAccount(gomock.Any(), gomock.Any()).
//					Times(0)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusBadRequest, recorder.Code)
//			},
//		},
//		{
//			name: "InvalidPageSize",
//			query: Query{
//				pageID:   1,
//				pageSize: 100000,
//			},
//
//			buildStubs: func(store *mockdb.MockStore) {
//				store.EXPECT().
//					ListAccount(gomock.Any(), gomock.Any()).
//					Times(0)
//			},
//			checkResponse: func(recorder *httptest.ResponseRecorder) {
//				require.Equal(t, http.StatusBadRequest, recorder.Code)
//			},
//		},
//	}
//
//	for i := range testCases {
//		tc := testCases[i]
//
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//
//			store := mockdb.NewMockStore(ctrl)
//			tc.buildStubs(store)
//
//			server := NewServer(store)
//			recorder := httptest.NewRecorder()
//
//			url := "/accounts"
//			request, err := http.NewRequest(http.MethodGet, url, nil)
//			require.NoError(t, err)
//
//			// Add query parameters to request URL
//			q := request.URL.Query()
//			q.Add("page_id", fmt.Sprintf("%d", tc.query.pageID))
//			q.Add("page_size", fmt.Sprintf("%d", tc.query.pageSize))
//			request.URL.RawQuery = q.Encode()
//
//			server.router.ServeHTTP(recorder, request)
//			tc.checkResponse(recorder)
//		})
//	}
//}

func randomAccount(user string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    user,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, accounts, gotAccounts)
}
func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, acc db.Account) {
	data, err := ioutil.ReadAll(body)

	var resultAcc db.Account
	err = json.Unmarshal(data, &resultAcc)
	require.NoError(t, err)
	require.Equal(t, acc, resultAcc)
}

func randomUser(t *testing.T) (user db.User, password string) {
	password = util.RandomString(6)
	hashedPassword, err := util.GetHashPassword(password)
	require.NoError(t, err)

	user = db.User{
		Username:       util.RandomName(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomString(8),
		Email:          util.RandomEmail(),
	}
	return
}
