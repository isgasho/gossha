package gossha

import (
	//	"fmt"
	"net"
	"testing"
	//	"time"
)

type connFakeAddr struct{}

func (c connFakeAddr) Network() string {
	return ""
}
func (c connFakeAddr) String() string {
	return "127.0.0.1:22"
}

type connMetadataFake struct{}

func (c connMetadataFake) User() string {
	return "testUser"
}
func (c connMetadataFake) SessionID() []byte {
	return []byte("lalalalala")
}
func (c connMetadataFake) ClientVersion() []byte {
	return []byte("1")
}
func (c connMetadataFake) ServerVersion() []byte {
	return []byte("1")
}
func (c connMetadataFake) LocalAddr() net.Addr {
	return connFakeAddr{}
}
func (c connMetadataFake) RemoteAddr() net.Addr {
	return connFakeAddr{}
}

func TestInitDatabase(t *testing.T) {
	err := InitDatabase("sqlite3", ":memory:")
	if err != nil {
		t.Error(err.Error())
	}
	err = CreateUser("testUser", "testPassword", false)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestLoginByUsernameAndPassword(t *testing.T) {
	h1 := New()
	m := connMetadataFake{}
	err := h1.LoginByUsernameAndPassword(m, "testPassword")
	if err != nil {
		t.Error(err.Error())
	}
	if h1.CurrentUser.Name != "testUser" {
		t.Error("We have authorized wrong user via Handler.LoginByUsernameAndPassword!")
	}
	if h1.SessionId != "lalalalala" {
		t.Error("Session Id is not set via Handler.LoginByUsernameAndPassword!")
	}
	if h1.Ip != "127.0.0.1" {
		t.Error("Remote IP  is not set via Handler.LoginByUsernameAndPassword!")
	}

	h2 := New()
	err = h2.LoginByUsernameAndPassword(m, "wrongTestPassword")
	if err != nil {
		if err.Error() != "Wrong password for user testUser!" {
			t.Error(err.Error())
		}
	} else {
		t.Error("No error for wrong password via Handler.LoginByUsernameAndPassword!")
	}
	if h2.CurrentUser.Name != "" {
		t.Error("We have authorized wrong user via Handler.LoginByUsernameAndPassword!")
	}
}
