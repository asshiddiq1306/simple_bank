package api

import (
	"os"
	"testing"
	"time"

	db "github.com/asshiddiq1306/simple_bank/db/sql"
	"github.com/asshiddiq1306/simple_bank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newServerTest(t *testing.T, store db.Store) *Server {
	config := util.Config{
		TokenSymmetricKey:   util.RandomString(32),
		TokenAccessDuration: time.Minute,
	}
	server, err := NewServer(config, store)
	require.NoError(t, err)
	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
