package main

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/secrets"
)

func process(command *Command, secrets secrets.Secrets) (interface{}, error) {
	switch command.Command {
	case PingCommand:
		return "pong", nil
	case StatusCommand:
		return secrets.Status(context.Background())
	case LockCommand:
		return nil, secrets.Lock(context.Background())
	case UnlockCommand:
		var unlockArgs UnlockArgs
		if err := json.Unmarshal(command.Args, &unlockArgs); err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal unlockArgs")
		}
		return nil, secrets.Unlock(context.Background(), unlockArgs.Name, unlockArgs.Email, unlockArgs.Passphrase)
	case IdentitiesCommand:
		return secrets.Identities(context.Background())
	case ListCommand:
		return secrets.List(context.Background())
	case AddCommand:
		var addArgs AddArgs
		if err := json.Unmarshal(command.Args, &addArgs); err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal addArgs")
		}
		return nil, secrets.Add(context.Background(), addArgs.ID, addArgs.Type, addArgs.Version)
	case GetCommand:
		var getArgs GetArgs
		if err := json.Unmarshal(command.Args, &getArgs); err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal getArgs")
		}
		return secrets.Get(context.Background(), getArgs.ID)
	case EstimateCommand:
		var estimateArgs EstimateArgs
		if err := json.Unmarshal(command.Args, &estimateArgs); err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal estimateArgs")
		}
		return secrets.EstimateStrength(context.Background(), estimateArgs.Password, estimateArgs.Inputs)
	}
	return nil, nil
}
