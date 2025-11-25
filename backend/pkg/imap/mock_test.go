package imap

import (
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

type MockBackend struct {
}

func (b *MockBackend) Login(conn *imap.ConnInfo, username, password string) (backend.User, error) {
	return &MockUser{username: username}, nil
}

type MockUser struct {
	username string
}

func (u *MockUser) Username() string {
	return u.username
}

func (u *MockUser) ListMailboxes(subscribed bool) ([]backend.Mailbox, error) {
	return []backend.Mailbox{
		&MockMailbox{name: "INBOX"},
	}, nil
}

func (u *MockUser) GetMailbox(name string) (backend.Mailbox, error) {
	return &MockMailbox{name: name}, nil
}

func (u *MockUser) CreateMailbox(name string) error                  { return nil }
func (u *MockUser) DeleteMailbox(name string) error                  { return nil }
func (u *MockUser) RenameMailbox(existingName, newName string) error { return nil }
func (u *MockUser) Logout() error                                    { return nil }

type MockMailbox struct {
	name string
}

func (m *MockMailbox) Name() string { return m.name }
func (m *MockMailbox) Info() (*imap.MailboxInfo, error) {
	return &imap.MailboxInfo{Name: m.name}, nil
}
func (m *MockMailbox) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	return &imap.MailboxStatus{
		Name:        m.name,
		Messages:    1,
		Unseen:      1,
		Recent:      1,
		UidNext:     2,
		UidValidity: 1,
	}, nil
}
func (m *MockMailbox) SetSubscribed(subscribed bool) error { return nil }
func (m *MockMailbox) Check() error                        { return nil }
func (m *MockMailbox) ListMessages(uid bool, seqset *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	defer close(ch)
	// Emit a dummy message
	msg := imap.NewMessage(1, items)
	envelope := imap.Envelope{
		Subject: "Test Subject",
		From:    []*imap.Address{{PersonalName: "Sender", MailboxName: "sender", HostName: "example.com"}},
		Date:    time.Now(),
	}
	msg.Envelope = &envelope

	// Add Body
	section, _ := imap.ParseBodySectionName("BODY[]")
	bodyString := "Date: Mon, 7 Feb 1994 21:52:25 -0800 (PST)\r\n" +
		"From: Fred Foobar <foobar@example.com>\r\n" +
		"Subject: afternoon meeting\r\n" +
		"To: mooch@example.com\r\n" +
		"Message-Id: <B85893d97@example.com>\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=utf-8\r\n" +
		"\r\n" +
		"This is a test body."
	msg.Body[section] = strings.NewReader(bodyString)

	// Also set Uid if requested
	msg.Uid = 1
	ch <- msg
	return nil
}

func (m *MockMailbox) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	return []uint32{1}, nil
}
func (m *MockMailbox) Expunge() error { return nil }
func (m *MockMailbox) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, operation imap.FlagsOp, flags []string) error {
	return nil
}
func (m *MockMailbox) CopyMessages(uid bool, seqset *imap.SeqSet, dest string) error { return nil }
func (m *MockMailbox) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	return nil
}
