package imap

import (
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type EmailData struct {
	Subject   string
	Sender    string
	Date      time.Time
	MessageID string
	BodyText  string
	BodyHTML  string
}

// FetchEmails fetches the latest N messages' data (including body) from the specified mailbox.
func FetchEmails(c *client.Client, mailbox string, limit int) ([]EmailData, error) {
	// Select Mailbox
	_, err := c.Select(mailbox, false)
	if err != nil {
		return nil, err
	}

	// Explicitly fetch status
	mbox, err := c.Status(mailbox, []imap.StatusItem{imap.StatusMessages})
	if err != nil {
		return nil, err
	}

	// Get the last N messages
	if mbox.Messages == 0 {
		return []EmailData{}, nil
	}

	from := uint32(1)
	if mbox.Messages > uint32(limit) {
		from = mbox.Messages - uint32(limit) + 1
	}
	to := mbox.Messages

	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	// Fetch Envelope and Body
	section := &imap.BodySectionName{} // Empty section name means the whole message body (RFC 822 style)
	items := []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}
	
	messages := make(chan *imap.Message, limit)
	done := make(chan error, 1)
	
	go func() {
		done <- c.Fetch(seqset, items, messages)
	}()

	var results []EmailData
	for msg := range messages {
		if msg.Envelope == nil {
			continue
		}
        
        sender := "Unknown"
        if len(msg.Envelope.From) > 0 {
            sender = msg.Envelope.From[0].Address()
        }

		// Extract Body
		var bodyText, bodyHTML string
		
		// We requested only one body section, so we can just take the first one found.
		// This avoids potential issues with BodySectionName pointer equality in tests/mocks.
		var r imap.Literal
		for _, literal := range msg.Body {
			r = literal
			break
		}
		
		if r != nil {
			bodyText, bodyHTML, _ = ExtractBody(r)
		}

		results = append(results, EmailData{
			Subject:   msg.Envelope.Subject,
			Sender:    sender,
			Date:      msg.Envelope.Date,
			MessageID: msg.Envelope.MessageId,
			BodyText:  bodyText,
			BodyHTML:  bodyHTML,
		})
	}

	if err := <-done; err != nil {
		return nil, err
	}

	return results, nil
}
