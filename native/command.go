package main

import (
	"encoding/binary"
	"encoding/json"
	"io"
	"runtime"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
)

type CommandName string

const (
	PingCommand       CommandName = "ping"
	StatusCommand     CommandName = "status"
	LockCommand       CommandName = "lock"
	UnlockCommand     CommandName = "unlock"
	IdentitiesCommand CommandName = "identities"
	ListCommand       CommandName = "list"
	AddCommand        CommandName = "add"
	GetCommand        CommandName = "get"
	EstimateCommand   CommandName = "estimate"
)

type Command struct {
	Command CommandName     `json:"command"`
	Args    json.RawMessage `json:"args"`
}

type CommandReply struct {
	Command CommandName     `json:"command"`
	Reply   json.RawMessage `json:"response"`
}

type UnlockArgs struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Passphrase string `json:"passphrase"`
}

type AddArgs struct {
	ID      string            `json:"id"`
	Type    api.SecretType    `json:"type"`
	Version api.SecretVersion `json:"version"`
}

type GetArgs struct {
	ID string `json:"id"`
}

type EstimateArgs struct {
	Password string   `json:"password"`
	Inputs   []string `json:"inputs"`
}

func readCommand(reader io.Reader) (*Command, error) {
	var size uint32
	if err := binary.Read(reader, nativeByteOrder(), &size); err == io.EOF {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "Failed to read message size")
	}
	message := make([]byte, size)
	if _, err := io.ReadFull(reader, message); err != nil {
		return nil, errors.Wrap(err, "Failed to read message")
	}
	var command Command
	if err := json.Unmarshal(message, &command); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal message")
	}
	return &command, nil
}

func writeReply(writer io.Writer, command CommandName, reply interface{}) error {
	replyRaw, err := json.Marshal(reply)
	if err != nil {
		return errors.Wrap(err, "Failed to encode reply")
	}
	commandReply := &CommandReply{
		Command: command,
		Reply:   json.RawMessage(replyRaw),
	}
	message, err := json.Marshal(commandReply)
	if err != nil {
		return errors.Wrap(err, "Failed to encode commandReply")
	}
	size := uint32(len(message))
	if err := binary.Write(writer, nativeByteOrder(), &size); err != nil {
		return errors.Wrap(err, "Failed to write message size")
	}
	if _, err := writer.Write(message); err != nil {
		return errors.Wrap(err, "Failed to write emssage")
	}
	return nil
}

func nativeByteOrder() binary.ByteOrder {
	switch runtime.GOARCH {
	case "arm":
		return binary.BigEndian
	default:
		return binary.LittleEndian
	}
}
