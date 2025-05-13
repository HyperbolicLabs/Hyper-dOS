package sh

import (
	"fmt"

	"epitome.hyperbolic.xyz/config"
	"github.com/chzyer/readline"
	"go.uber.org/zap"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type session struct {
	cfg           *config.Config
	logger        *zap.Logger
	clientset     kubernetes.Interface
	dynamicClient *dynamic.DynamicClient
	namespace     *string
	rl            *readline.Instance
	completions   readline.DynamicPrefixCompleterInterface
	cdCompletions []string
}

func (s *session) write(msg string) {
	s.rl.Write([]byte(msg))
}

func (s *session) writeln(msg string) {
	s.rl.Write([]byte(msg + "\n"))
}

func (s *session) writeErr(msg string) {
	// TODO some formatting
	s.rl.Write([]byte(msg))
}

func (s *session) initReadline() error {

	// Initialize with static commands first
	s.completions = readline.NewPrefixCompleter(
		readline.PcItem("cd",
			readline.PcItemDynamic(s.getCdCompletions),
		),
		readline.PcItem("ls"),
		readline.PcItem("exit"),
		readline.PcItem("help"),
		readline.PcItem("clear"),
	)

	readlineInstance, err := readline.NewEx(&readline.Config{
		Prompt:          s.getPrompt(),
		AutoComplete:    s.completions, // Use our dynamic completer
		HistoryFile:     "/tmp/epitome_readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return err
	}
	s.rl = readlineInstance

	// first and foremost - set vim mode
	s.rl.SetVimMode(s.cfg.Shell.VIM_MODE)

	return nil
}

func (s *session) getPrompt() string {
	path := ""
	if s.namespace != nil {
		path = "/" + *s.namespace
	}
	return fmt.Sprintf("%sepitomesh%s%s )%s ",
		config.ShellPromptColor,
		path,
		config.ShellResetColor,
		config.ShellResetColor)
}
