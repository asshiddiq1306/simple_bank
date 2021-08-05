package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/asshiddiq1306/simple_bank/db/mock"
	db "github.com/asshiddiq1306/simple_bank/db/sql"
	"github.com/asshiddiq1306/simple_bank/token"
	"github.com/asshiddiq1306/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateNewAccount(t *testing.T) {
	user, _ := createRandomUser(t)
	account := createRandomAccount(user.Username)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, time.Minute)
			},
			body: gin.H{
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateNewAccountArgs{
					Owner:    account.Owner,
					Balance:  account.Balance,
					Currency: account.Currency,
				}

				store.EXPECT().CreateNewAccount(gomock.Any(), gomock.Eq(arg)).Times(1).Return(account, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			body: gin.H{
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateNewAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredAuth",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, -time.Minute)
			},
			body: gin.H{
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateNewAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, time.Minute)
			},
			body: gin.H{
				"currency": "unsupported",
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateNewAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, time.Minute)
			},
			body: gin.H{
				"currency": account.Currency,
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().CreateNewAccount(gomock.Any(), gomock.Any()).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newServerTest(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/account"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func TestGetAccountByIDAPI(t *testing.T) {
	user, _ := createRandomUser(t)
	account1 := createRandomAccount(user.Username)

	testCases := []struct {
		name          string
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, time.Minute)
			},
			accountID: account1.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByID(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(account1, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account1)
			},
		},
		{
			name: "NoAuthorization",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			accountID: account1.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByID(gomock.Any(), gomock.Eq(account1.ID)).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredToken",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, -time.Minute)
			},
			accountID: account1.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByID(gomock.Any(), gomock.Eq(account1.ID)).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, time.Minute)
			},
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByID(gomock.Any(), gomock.Eq(account1.ID)).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, time.Minute)
			},
			accountID: account1.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByID(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "AccountNotFound",
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, time.Minute)
			},
			accountID: account1.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountByID(gomock.Any(), gomock.Eq(account1.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newServerTest(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/account/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)

		})
	}
}

func TestGetListAccounts(t *testing.T) {
	n := 5
	accounts := make([]db.Account, n)
	user, _ := createRandomUser(t)
	for i := 0; i < n; i++ {
		accounts[i] = createRandomAccount(user.Username)
	}

	type Query struct {
		PageID   int
		PageSize int
	}

	testCases := []struct {
		name          string
		query         Query
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			query: Query{
				PageID:   1,
				PageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.GetListAccountsArgs{
					Owner:  user.Username,
					Limit:  int32(n),
					Offset: 0,
				}

				store.EXPECT().GetListAccounts(gomock.Any(), gomock.Eq(arg)).Times(1).Return(accounts, nil)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		},
		{
			name: "NoAuthorization",
			query: Query{
				PageID:   1,
				PageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "ExpiredAuthorization",
			query: Query{
				PageID:   1,
				PageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, -time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name: "InternalError",
			query: Query{
				PageID:   1,
				PageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetListAccounts(gomock.Any(), gomock.Any()).Times(1).Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "BadRequest",
			query: Query{
				PageID:   0,
				PageSize: n,
			},
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuth(t, request, tokenMaker, authorizationBearerTypeKey, user.Username, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newServerTest(t, store)
			recorder := httptest.NewRecorder()

			url := "/accounts"
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			q := request.URL.Query()
			q.Add("page_id", fmt.Sprintf("%d", tc.query.PageID))
			q.Add("page_size", fmt.Sprintf("%d", tc.query.PageSize))
			request.URL.RawQuery = q.Encode()

			tc.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(recorder)
		})
	}
}

func createRandomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    owner,
		Balance:  0,
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var getAccount db.Account
	err = json.Unmarshal(data, &getAccount)
	require.NoError(t, err)
	require.Equal(t, account.Owner, getAccount.Owner)
	require.Equal(t, account.Balance, getAccount.Balance)
	require.Equal(t, account.Currency, getAccount.Currency)
}
