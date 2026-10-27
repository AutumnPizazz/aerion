package imap

import (
	"github.com/emersion/go-imap/v2"
)

// sendID identifies this client to the server with the IMAP ID command
// (RFC 2971) when the server advertises the ID capability.
//
// Some providers, most notably 163.com (Coremail), refuse to open a mailbox
// for clients that never identified themselves. Every SELECT fails with
//
//	NO SELECT Unsafe Login. Please contact kefu@188.com for help
//
// which left the account unusable: folders never synced and no message was
// ever stored. Sending ID right after login makes SELECT succeed.
//
// Every field in an ID response is optional, so only the client name and
// vendor are reported: a version string here would duplicate app.Version and
// silently drift away from it on every release.
func (c *Client) sendID() {
	if c.client == nil || !c.caps.Has(imap.CapID) {
		return
	}

	idData := &imap.IDData{
		Name:   "Aerion",
		Vendor: "Aerion",
	}

	if _, err := c.client.ID(idData).Wait(); err != nil {
		// A failed ID is not fatal: the server still accepts the session and
		// most servers do not require it at all.
		c.log.Debug().Err(err).Msg("IMAP ID command failed")
		return
	}

	c.log.Debug().Msg("Sent IMAP ID (RFC 2971)")
}
