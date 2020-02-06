package msg

// Mailbox is a global message manager.
//
// Only use this if the object being worked on is agnostic of messages
// of differing dispatch points.
var Mailbox = &MessageManager{}
