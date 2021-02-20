package main

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/marekaf/gcr-lifecycle-policy/internal/worker"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	credsFile       string   // path of credentials json file
	repoFilter      []string // list of regions we want to check
	logLevel        string
	keepTags        int
	retentionDays   int
	kubeconfigPath  string
	registryPrefix  string
	sortby          string
	protectTagRegex string
	dryRun          bool

	// default values
	repoFilterDefault     = []string{}
	dryRunDefault         = true
	logLevelDefault       = "ERROR"
	keepTagsDefault       = 10
	retentionDaysDefault  = 365
	credsFileDefault      = "./creds/serviceaccount.json"
	registryPrefixDefault = "eu.gcr.io"
	sortbyDefault         = "timeCreatedMs"
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
	rootCmd.AddCommand(listClusterCmd)
	rootCmd.AddCommand(cleanupCmd)

	// root command
	rootCmd.PersistentFlags().StringVar(&credsFile, "creds", credsFileDefault, "credential file")
	rootCmd.PersistentFlags().StringVar(&registryPrefix, "registry", registryPrefixDefault, "GCR url to use")
	rootCmd.PersistentFlags().StringVar(&sortby, "sort-by", sortbyDefault, "field to sort images by")
	rootCmd.PersistentFlags().StringArrayVar(&repoFilter, "repos", repoFilterDefault, "list of repos you want to work with")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", logLevelDefault, "log level")
	rootCmd.PersistentFlags().StringVar(&kubeconfigPath, "cluster", kubeconfigPathDefault, "kubeconfig path")

	// cleanup command
	cleanupCmd.PersistentFlags().IntVar(&keepTags, "keep-tags", keepTagsDefault, "number of tags to keep per image")
	cleanupCmd.PersistentFlags().IntVar(&retentionDays, "retention", retentionDaysDefault, "number of days of retention to keep images")
	cleanupCmd.PersistentFlags().StringVar(&protectTagRegex, "protect-tag-regex", "", "regex to protect matched tag")
	cleanupCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", dryRunDefault, "dry-run for images cleaning")

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
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("root cmd execute failed")
	}
}

func cleanup(_ *cobra.Command, _ []string) {

	// set loglevel
	setLogLevel()

	if protectTagRegex != "" {
		_, err := regexp.Compile(protectTagRegex)
		if err != nil {
			log.Fatalf("Bad regex string: %s", protectTagRegex)
		}
	}

	config := worker.Config{
		CredsFile:       credsFile,
		RepoFilter:      repoFilter,
		KeepTags:        keepTags,
		RetentionDays:   retentionDays,
		RegistryURL:     registryPrefix,
		SortBy:          sortby,
		KubeconfigPath:  kubeconfigPath,
		ProtectTagRegex: protectTagRegex,
		DryRun:          dryRun,
	}

	worker.HandleCleanup(config)
}

func list(_ *cobra.Command, _ []string) {

	// set loglevel
	setLogLevel()

	config := worker.Config{
		CredsFile:   credsFile,
		RepoFilter:  repoFilter,
		RegistryURL: registryPrefix,
		SortBy:      sortby,
	}

	result := worker.HandleList(config)
	worker.PrintList(result)
}

func listRepos(_ *cobra.Command, _ []string) {

	// set loglevel
	setLogLevel()

	config := worker.Config{
		CredsFile:   credsFile,
		RepoFilter:  repoFilter,
		RegistryURL: registryPrefix,
		SortBy:      sortby,
	}

	result := worker.HandleListRepos(config)
	worker.PrintListRepos(result)
}

func listCluster(_ *cobra.Command, _ []string) {

	// set loglevel
	setLogLevel()

	config := worker.Config{
		CredsFile:      credsFile,
		RepoFilter:     repoFilter,
		KubeconfigPath: kubeconfigPath,
		RegistryURL:    registryPrefix,
		SortBy:         sortby,
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
