package types

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var SessionMedia = MediaType("session", func() {
	Description("A session for a user, associated with a specific browser")
	Attributes(func() {
		Attribute("id", String, "Unique unchanging session ID", func() {
			Metadata("struct:field:type", "string")
		})
		Attribute("userId", String, "ID of the user this session is for", func() {
			Metadata("struct:field:type", "string")
		})
		Attribute("lastUsed", DateTime, "Time that this session was last used")
		Attribute("browser", String, "The browser and browser version connected with this session")
		Attribute("os", String, "The OS of the system where this session was used")
		Attribute("ip", String, "The last IP address where this session was used")
		Attribute("location", String, "A humanReadable string describing the last known location of the session")
		Attribute("coordinates", String, "The latitude and longitude of the last known location of the session", func() {
			Example("513452x54123")
		})
		Attribute("isMobile", Boolean, "Whether the session was from a mobile device")
		Attribute("mapUrl", String, "The URL of the Google map to show the location, suitable for using in an img tag")
		Required("id", "userId", "lastUsed", "browser", "os", "ip", "location", "coordinates", "isMobile", "mapUrl")
	})

	View("default", func() {
		Attribute("id")
		Attribute("userId")
		Attribute("lastUsed")
		Attribute("browser")
		Attribute("os")
		Attribute("ip")
		Attribute("location")
		Attribute("coordinates")
		Attribute("isMobile")
		Attribute("mapUrl")
	})
})

var AllSessionsMedia = MediaType("all-sessions", func() {
	Description("All of the sessions associated with a user")
	Attributes(func() {
		Attribute("currentSession", SessionMedia)
		Attribute("otherSessions", CollectionOf(SessionMedia))
	})

	View("default", func() {
		Attribute("currentSession")
		Attribute("otherSessions")
	})
})
