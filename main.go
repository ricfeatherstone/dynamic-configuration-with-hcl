package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
)

const configFile = "config.hcl"

func main() {
	var err error

	_ = os.Setenv("message", "Hello World!")

	var evalContext *hcl.EvalContext
	if evalContext, err = awsStsEvalContext(context.Background()); err != nil {
		slog.Error("failed to load aws context", "error", err)
		os.Exit(1)
	}

	var data Config
	if err = hclsimple.DecodeFile(configFile, evalContext, &data); err != nil {
		slog.Error("failed to load configuration", "error", err)
	}

	slog.Info("finished", "data", data)
}

type Config struct {
	AccountId string `hcl:"account_id"`
	Region    string `hcl:"region"`
	Principal string `hcl:"principal"`
	Email     string `hcl:"email"`
	Message   string `hcl:"message"`
}

func awsStsEvalContext(ctx context.Context) (evalCtx *hcl.EvalContext, err error) {
	var cfg aws.Config
	if cfg, err = config.LoadDefaultConfig(ctx); err != nil {
		slog.Error("failed to load aws default config", "error", err)
		return evalCtx, err
	}

	client := sts.NewFromConfig(cfg)

	var result *sts.GetCallerIdentityOutput
	if result, err = client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{}); err != nil {
		slog.Error("failed to get sts caller identity", "error", err)
		return evalCtx, err
	}

	evalCtx = &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"aws_account_id": cty.StringVal(*result.Account),
			"aws_principal":  cty.StringVal(*result.UserId),
			"aws_region":     cty.StringVal(cfg.Region),
		},
		Functions: map[string]function.Function{
			"split":    stdlib.SplitFunc,
			"upper":    stdlib.UpperFunc,
			"from_env": FromEnv,
		},
	}

	return evalCtx, err
}

var FromEnv = function.New(&function.Spec{
	Description: "Return a message.",
	Params: []function.Parameter{
		{
			Name:             "key",
			Type:             cty.String,
			AllowDynamicType: true,
		},
	},
	Type: function.StaticReturnType(cty.String),
	Impl: func(args []cty.Value, retType cty.Type) (cty.Value, error) {
		message := os.Getenv(args[0].AsString())
		return cty.StringVal(message), nil
	},
})
