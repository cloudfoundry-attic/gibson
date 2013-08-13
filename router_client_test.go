package gibson

import (
	"github.com/cloudfoundry/go_cfmessagebus/mock_cfmessagebus"
	. "launchpad.net/gocheck"
	"time"
)

type RCSuite struct{}

func init() {
	Suite(&RCSuite{})
}

func (s *RCSuite) TestRouterClientRegistering(c *C) {
	mbus := mock_cfmessagebus.NewMockMessageBus()

	routerClient := NewCFRouterClient("1.2.3.4", mbus)

	registered := make(chan []byte)

	mbus.Subscribe("router.register", func(msg []byte) {
		registered <- msg
	})

	routerClient.Register(123, "abc")

	select {
	case msg := <-registered:
		c.Assert(string(msg), Equals, `{"uris":["abc"],"host":"1.2.3.4","port":123}`)
	case <-time.After(500 * time.Millisecond):
		c.Error("did not receive a router.register!")
	}
}

func (s *RCSuite) TestRouterClientUnregistering(c *C) {
	mbus := mock_cfmessagebus.NewMockMessageBus()

	routerClient := NewCFRouterClient("1.2.3.4", mbus)

	registered := make(chan []byte)

	mbus.Subscribe("router.unregister", func(msg []byte) {
		registered <- msg
	})

	routerClient.Unregister(123, "abc")

	select {
	case msg := <-registered:
		c.Assert(string(msg), Equals, `{"uris":["abc"],"host":"1.2.3.4","port":123}`)
	case <-time.After(500 * time.Millisecond):
		c.Error("did not receive a router.unregister!")
	}
}

func (s *RCSuite) TestRouterClientRouterStartHandling(c *C) {
	mbus := mock_cfmessagebus.NewMockMessageBus()

	routerClient := NewCFRouterClient("1.2.3.4", mbus)

	times := make(chan time.Time)

	mbus.Subscribe("router.register", func(msg []byte) {
		times <- time.Now()
	})

	err := routerClient.Greet()
	c.Assert(err, IsNil)

	mbus.Publish("router.start", []byte(`{"minimumRegisterIntervalInSeconds":1}`))

	routerClient.Register(123, "abc")

	initialRegister := timedReceive(times, 1*time.Second)
	c.Assert(initialRegister, NotNil)

	time1 := timedReceive(times, 2*time.Second)
	c.Assert(time1, NotNil)

	time2 := timedReceive(times, 2*time.Second)
	c.Assert(time2, NotNil)

	c.Assert((*time2).Sub(*time1) >= 1*time.Second, Equals, true)
}

func (s *RCSuite) TestRouterClientGreeting(c *C) {
	mbus := mock_cfmessagebus.NewMockMessageBus()

	routerClient := NewCFRouterClient("1.2.3.4", mbus)

	times := make(chan time.Time)

	mbus.Subscribe("router.register", func(msg []byte) {
		times <- time.Now()
	})

	routerClient.Register(123, "abc")

	initialRegister := timedReceive(times, 1*time.Second)
	c.Assert(initialRegister, NotNil)

	mbus.RespondToChannel("router.greet", func([]byte) []byte {
		return []byte(`{"minimumRegisterIntervalInSeconds":1}`)
	})

	err := routerClient.Greet()
	c.Assert(err, IsNil)

	time1 := timedReceive(times, 2*time.Second)
	c.Assert(time1, NotNil)

	time2 := timedReceive(times, 2*time.Second)
	c.Assert(time2, NotNil)

	c.Assert((*time2).Sub(*time1) >= 1*time.Second, Equals, true)
}

func timedReceive(from chan time.Time, giveup time.Duration) *time.Time {
	select {
	case val := <-from:
		return &val
	case <-time.After(giveup):
		return nil
	}
}
