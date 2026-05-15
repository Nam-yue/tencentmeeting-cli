package middleware

import (
	"tmeet/internal/cmdutil"

	"github.com/spf13/cobra"
)

// WithApiCmd resolves the ApiCmd for this invocation (possibly choosing one
// out of several candidates based on flags) and writes it to the command's
// annotation so that downstream middlewares (e.g. WithCompact) and helpers
// can read it via cmdutil.GetApiCmdAnnotation(cmd).
//
// It MUST be placed BEFORE any middleware that reads the ApiCmd annotation,
// e.g. WithCompact. Typical usage:
//
//	RunE: middleWare.Chain(
//	    opts.Run,
//	    middleWare.WithApiCmd(cmdutil.FlagSwitch(
//	        cmdutil.FlagCase{
//	            When:   cmdutil.WhenStringFlagSet("meeting-id"),
//	            ApiCmd: cmdutil.ApiCmdMeetingGetById,
//	        },
//	        cmdutil.FlagCase{
//	            When:   cmdutil.WhenStringFlagSet("meeting-code"),
//	            ApiCmd: cmdutil.ApiCmdMeetingGetByCode,
//	        },
//	    )),
//	    middleWare.WithCompact(tmeet),
//	),
//
// When resolver is nil or resolves to an empty string, no annotation is
// written and the chain proceeds transparently.
func WithApiCmd(resolver cmdutil.ApiCmdResolver) CmdMiddleware {
	return func(next RunEFunc) RunEFunc {
		return func(cmd *cobra.Command, args []string) error {
			if resolver != nil {
				if name := resolver.Resolve(cmd); name != "" {
					cmdutil.InjectApiCmdAnnotation(cmd, name)
				}
			}
			return next(cmd, args)
		}
	}
}
