package main

import (
	"os"
	"path/filepath"

	"github.com/marekaf/gcr-lifecycle-policy/pkg/worker"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	credsFile      string   // path of credentials json file
	repoFilter     []string // list of regions we want to check
	logLevel       string
	keepTags       int
	retentionDays  int
	kubeconfigPath string
	registryPrefix string

	// default values
	repoFilterDefault     = []string{}
	logLevelDefault       = "ERROR"
	keepTagsDefault       = 10
	retentionDaysDefault  = 365
	credsFileDefault      = "./creds/serviceaccount.json"
	registryPrefixDefault = "eu.gcr.io"
	kubeconfigPathDefault = filepath.Join(homeDir(), ".kube", "config")

	// commands
	rootCmd = &cobra.Command{
		Use:   "gcr",
		Short: "", // add some clever but short description
		Long:  "", // add even more clever description
	}

	listCmd = &cobra.Command{
		Use:   "list",
		Short: "",
		Long:  "",
		Run:   list,
	}

	listCatalogCmd = &cobra.Command{
		Use:   "list-catalog",
		Short: "",
		Long:  "",
		Run:   listCatalog,
	}

	listReposCmd = &cobra.Command{
		Use:   "list-repos",
		Short: "",
		Long:  "",
		Run:   listRepos,
	}

	listClusterCmd = &cobra.Command{
		Use:   "list-cluster",
		Short: "",
		Long:  "",
		Run:   listCluster,
	}

	cleanupCmd = &cobra.Command{
		Use:   "cleanup",
		Short: "",
		Long:  "",
		Run:   cleanup,
	}
)

func init() {

	// commands
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(listReposCmd)
	rootCmd.AddCommand(listCatalogCmd)
	rootCmd.AddCommand(listClusterCmd)
	rootCmd.AddCommand(cleanupCmd)

	// root command
	rootCmd.PersistentFlags().StringVar(&credsFile, "creds", credsFileDefault, "credential file")
	rootCmd.PersistentFlags().StringVar(&registryPrefix, "registry", registryPrefixDefault, "GCR url to use")
	rootCmd.PersistentFlags().StringArrayVar(&repoFilter, "repos", repoFilterDefault, "list of repos you want to work with")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", logLevelDefault, "log level")
	rootCmd.PersistentFlags().StringVar(&kubeconfigPath, "cluster", kubeconfigPathDefault, "kubeconfig path")

	// cleanup command
	cleanupCmd.PersistentFlags().IntVar(&keepTags, "keep-tags", keepTagsDefault, "number of tags to keep per image")
	cleanupCmd.PersistentFlags().IntVar(&retentionDays, "retention", retentionDays, "number of days of retention to keep images")

}

func setLogLevel() {
	switch logLevel {
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	case "DEBUG":
		log.SetLevel(log.DebugLevel)
	case "FATAL":
		log.SetLevel(log.FatalLevel)
	case "TRACE":
		log.SetLevel(log.TraceLevel)
	case "PANIC":
		log.SetLevel(log.PanicLevel)
	default:
		log.WithField("log-level", logLevel).Warning("Wrong log level set. Falling back to ERROR")
		log.SetLevel(log.ErrorLevel)
	}
}

func main() {
	rootCmd.Execute()
}

func cleanup(cmd *cobra.Command, args []string) {

	// set loglevel
	setLogLevel()

	config := worker.Config{
		CredsFile:      credsFile,
		RepoFilter:     repoFilter,
		KeepTags:       keepTags,
		RetentionDays:  retentionDays,
		RegistryURL:    registryPrefix,
		KubeconfigPath: kubeconfigPath,
	}

	worker.HandleCleanup(config)
}

func list(cmd *cobra.Command, args []string) {

	// set loglevel
	setLogLevel()

	config := worker.Config{
		CredsFile:   credsFile,
		RepoFilter:  repoFilter,
		RegistryURL: registryPrefix,
	}

	result := worker.HandleList(config)
	worker.PrintList(result)
}

func listCatalog(cmd *cobra.Command, args []string) {

	// set loglevel
	setLogLevel()

	config := worker.Config{
		CredsFile:   credsFile,
		RepoFilter:  repoFilter,
		RegistryURL: registryPrefix,
	}

	result := worker.HandleListCatalog(config)
	worker.PrintListRepos(result)
}

func listRepos(cmd *cobra.Command, args []string) {

	// set loglevel
	setLogLevel()

	config := worker.Config{
		CredsFile:   credsFile,
		RepoFilter:  repoFilter,
		RegistryURL: registryPrefix,
	}

	result := worker.HandleListRepos(config)
	worker.PrintListRepos(result)
}

func listCluster(cmd *cobra.Command, args []string) {

	// set loglevel
	setLogLevel()

	config := worker.Config{
		CredsFile:      credsFile,
		RepoFilter:     repoFilter,
		KubeconfigPath: kubeconfigPath,
		RegistryURL:    registryPrefix,
	}

	result := worker.HandleListCluster(config)
	worker.PrintListCluster(result)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
