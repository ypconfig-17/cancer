package stomp

// STOMP frame commands. Used upper case naming
// convention to avoid clashing with STOMP header names.
const (
	//Client commands.
	CommandAbort       = "ABORT"
	CommandACK         = "ACK"
	CommandBegin       = "BEGIN"
	CommandCommit      = "COMMIT"
	CommandConnect     = "CONNECT"
	CommandDisconnect  = "DISCONNECT"
	CommandNack        = "NACK"
	CommandSend        = "SEND"
	CommandStomp       = "STOMP"
	CommandSubscribe   = "SUBSCRIBE"
	CommandUnsubscribe = "UNSUBSCRIBE"

	// Server commands.
	CommandConnected = "CONNECTED"
	CommandError     = "ERROR"
	CommandMessage   = "MESSAGE"
	CommandHeartbeat = "RECEIPT"
)

type command struct {
	// Client commands.
	Abort       string
	ACK         string
	Begin       string
	Commit      string
	Connect     string
	Disconnect  string
	Nack        string
	Send        string
	Stomp       string
	Subscribe   string
	Unsubscribe string

	// Server commands.
	Connected string
	Error     string
	Message   string
	Heartbeat string
}

var FrameC = &command{
	// Client commands.
	Abort:       CommandAbort,
	ACK:         Ack,
	Begin:       CommandBegin,
	Commit:      CommandCommit,
	Connect:     CommandConnect,
	Disconnect:  CommandDisconnect,
	Nack:        CommandNack,
	Send:        CommandSend,
	Stomp:       CommandStomp,
	Subscribe:   CommandSubscribe,
	Unsubscribe: CommandUnsubscribe,

	// Server commands.
	Connected: CommandConnected,
	Error:     CommandError,
	Message:   CommandMessage,
	Heartbeat: CommandHeartbeat,
}

// ValidateClientCommand  验证是否是客户端的命令
func ValidateClientCommand(command string) (bool, error) {
	switch command {
	case CommandAbort, CommandACK, CommandBegin, CommandCommit, CommandConnect, CommandDisconnect, CommandNack, CommandSend, CommandStomp, CommandSubscribe, CommandUnsubscribe:
		return true, nil
	default:
		return false, ErrInvalidCommand
	}
}

// ValidateServerCommand  是否是server的包
func ValidateServerCommand(command string) error {
	switch command {
	case CommandConnected, CommandError, CommandMessage, CommandHeartbeat:
		{
			return nil
		}
	default:
		{
			return ErrInvalidCommand
		}
	}
}
