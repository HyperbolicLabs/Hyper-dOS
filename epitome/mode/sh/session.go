package sh

import (
	"fmt"

	"epitome.hyperbolic.xyz/config"
	"github.com/chzyer/readline"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

type session struct {
	cfg           *config.Config
	logger        *zap.Logger
	clientset     kubernetes.Interface
	namespace     *string
	rl            *readline.Instance
	completions   readline.DynamicPrefixCompleterInterface
	cdCompletions []string
}

func NewSession(
	cfg *config.Config,
	logger *zap.Logger,
	clientset kubernetes.Interface,
) (*session, CloseFunc, error) {
	s := &session{
		cfg:       cfg,
		logger:    logger,
		clientset: clientset,
	}

	if err := s.initReadline(); err != nil {
		return nil, nil, fmt.Errorf("failed to create session: %v", err)
	}

	closeF := func() {
		s.rl.Close()
	}
	return s, closeF, nil
}

type CloseFunc func()

func (s *session) write(msg string) {
	s.rl.Write([]byte(msg))
}

func (s *session) writeln(msg string) {
	s.rl.Write([]byte(msg + "\n"))
}

func (s *session) writeErr(err error) {
	msg := fmt.Sprintf("Error: %v\n", err)
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
