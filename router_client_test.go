package gibson

import (
	"time"

	"github.com/apcera/nats"
	"github.com/cloudfoundry/yagnats/fakeyagnats"
	. "launchpad.net/gocheck"
)

type RCSuite struct{}

func init() {
	Suite(&RCSuite{})
}

func (s *RCSuite) TestRouterClientRegistering(c *C) {
	mbus := fakeyagnats.Connect()

	routerClient := NewCFRouterClient("1.2.3.4", mbus)

	routerClient.Register(123, "abc")

	registrations := mbus.PublishedMessages("router.register")

	c.Assert(len(registrations), Not(Equals), 0)
	c.Assert(string(registrations[0].Data), Equals,
		`{"uris":["abc"],"host":"1.2.3.4","port":123,"private_instance_id":"`+routerClient.PrivateInstanceId+`"}`)
}

func (s *RCSuite) TestRouterClientUnregistering(c *C) {
	mbus := fakeyagnats.Connect()

	routerClient := NewCFRouterClient("1.2.3.4", mbus)

	routerClient.Unregister(123, "abc")

	unregistrations := mbus.PublishedMessages("router.unregister")

	c.Assert(len(unregistrations), Not(Equals), 0)
	c.Assert(string(unregistrations[0].Data), Equals,
		`{"uris":["abc"],"host":"1.2.3.4","port":123,"private_instance_id":"`+routerClient.PrivateInstanceId+`"}`)
}

func (s *RCSuite) TestRouterClientRouterStartHandling(c *C) {
	mbus := fakeyagnats.Connect()

	mbus.WhenSubscribing("router.start", func(handler nats.MsgHandler) error {
		handler(&nats.Msg{
			Data: []byte(`{"minimumRegisterIntervalInSeconds":1}`),
		})

		return nil
	})

	routerClient := NewCFRouterClient("1.2.3.4", mbus)

	err := routerClient.Greet()
	c.Assert(err, IsNil)

	routerClient.Register(123, "abc")

	c.Assert(len(mbus.PublishedMessages("router.register")), Equals, 1)

	time.Sleep(600 * time.Millisecond)

	c.Assert(len(mbus.PublishedMessages("router.register")), Equals, 1)

	time.Sleep(600 * time.Millisecond)

	c.Assert(len(mbus.PublishedMessages("router.register")), Equals, 2)

	time.Sleep(600 * time.Millisecond)

	c.Assert(len(mbus.PublishedMessages("router.register")), Equals, 2)

	time.Sleep(600 * time.Millisecond)

	c.Assert(len(mbus.PublishedMessages("router.register")), Equals, 3)
}
