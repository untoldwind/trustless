package main

import (
	"context"
	"encoding/json"

	"github.com/leanovate/microtools/logging"
	"github.com/pkg/errors"
	"github.com/untoldwind/trustless/api"
	"github.com/untoldwind/trustless/secrets"
)

func process(command *Command, secrets secrets.Secrets, logger logging.Logger) (interface{}, error) {
	switch command.Command {
	case PingCommand:
		return "pong", nil
	case StatusCommand:
		return secrets.Status(context.Background())
	case LockCommand:
		if err := secrets.Lock(context.Background()); err != nil {
			return nil, err
		}
		return secrets.Status(context.Background())
	case UnlockCommand:
		var unlockArgs UnlockArgs
		if err := json.Unmarshal(command.Args, &unlockArgs); err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal unlockArgs")
		}
		if err := secrets.Unlock(context.Background(), unlockArgs.Name, unlockArgs.Email, unlockArgs.Passphrase); err != nil {
			return nil, err
		}
		return secrets.Status(context.Background())
	case IdentitiesCommand:
		return secrets.Identities(context.Background())
	case ListCommand:
		var listFilter api.SecretListFilter
		if err := json.Unmarshal(command.Args, &listFilter); err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal listFilter")
		}
		return secrets.List(context.Background(), listFilter)
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
		var estimate api.PasswordEstimate
		if err := json.Unmarshal(command.Args, &estimate); err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal estimateArgs")
		}
		return secrets.EstimateStrength(context.Background(), estimate)
	case GenerateCommand:
		var parameter api.GenerateParameter
		if err := json.Unmarshal(command.Args, &parameter); err != nil {
			return nil, errors.Wrap(err, "Failed to unmarshal parameter")
		}
		return secrets.GeneratePassword(context.Background(), parameter)
	}
	return nil, nil
}
