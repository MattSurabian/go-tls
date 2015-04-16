/********
 * Config
 *
 * This file provides all the implementation necessary to read and write INI style
 * configurationManager files and to intelligently accept all configurationManager options as command
 * line flags. Getters are provided to retrieve configurationManager data.
 *
 * Reading configurationManager Options:
 *  When configurationManager data is read in from a file, any file paths present in the
 *  configurationManager values will be expanded relative to the location of said configurationManager
 *  file. When configurationManager data containing paths is passed in via the command line
 *  paths will be expanded relative to the current working directory. Command
 *  line flags take precedence and override any values loaded from a configurationManager file.
 *
 * Writing configurationManager Options:
 *  For convenience a configurationManager wizard is implemented which will allow the user to
 *  create a configurationManager file interactively. Any values containing a file path will
 *  expand that path relative to the current working directory. As a result the generated
 *  configurationManager file will contain only absolute paths.
 *
 * Searching For a Config File:
 *  When this module is loaded it attempts to load a configurationManager file, if the -config
 *  flag was passed in and the file exists it will be loaded, otherwise the program will
 *  begin searching up through the current directory hierarchy until it finds one. If it still does
 *  not find a configurationManager file a warning will be output suggesting the user generates one.
 *  It is possible to pass all configurationManager options with command line flags and avoid
 *  using a configurationManager file.
 */

package cliUtils

import (
	"flag"
	"github.com/mitchellh/cli"
	"github.com/rakyll/globalconf"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// In the event the config flag isn't passed this is the filename that will be searched for
const defaultConfigFileName = ".config"

var currentWorkingDirectory string
var userHomeDir string
var flagBasePath string
var configFilePath string

// configurationManager Values
var host string
var port string
var rootCert string
var serverTLSCert string
var serverTLSKey string
var clientTLSCert string
var clientTLSKey string
var rootName string

// Globalconf is used to intelligently merge flags and INI config values as well as
// persist changes to disk
var configurationManager *globalconf.GlobalConf
var ui = &cli.BasicUi{
	Writer: os.Stdout,
	Reader: os.Stdin,
}

//////////////// PUBLIC GETTERS /////////////////////

func GetConfigFilePath() string {
	return configFilePath
}
func GetRootCert() string {
	return rootCert
}

func GetHostAndPort() string {
	return host + ":" + port
}

func GetServerTLSCertPath() string {
	return serverTLSCert
}

func GetServerTLSKeyPath() string {
	return serverTLSKey
}

func GetClientTLSCertPath() string {
	return clientTLSCert
}

func GetClientTLSKeyPath() string {
	return clientTLSKey
}

func GetRootName() string {
	return rootName
}

/**
 * init
 * Initialize flags and set helper variables like currentWorkingDirectory and userHomeDir.
 * Search for a configurationManager file if necessary and correctly parse all values found.
 */
func init() {
	flag.StringVar(&configFilePath, "config", "", "What is the path to the configurationManager file?")
	flag.StringVar(&rootCert, "root-cert", "", "What is the path to the root CA certificate for TLS?")
	flag.StringVar(&rootName, "root-name", "", "What is the name on the CA cert?")
	flag.StringVar(&host, "host", "", "What is the domain name or ip address of the server?")
	flag.StringVar(&port, "port", "", "What port should the server be listening on?")
	flag.StringVar(&serverTLSCert, "server-tls-cert", "", "What is the path to the server's TLS certificate?")
	flag.StringVar(&serverTLSKey, "server-tls-key", "", "What is the path to the server's TLS key?")
	flag.StringVar(&clientTLSCert, "client-tls-cert", "", "What is the path to the TLS client certificate?")
	flag.StringVar(&clientTLSKey, "client-tls-key", "", "What is the path to the TLS client key?")

	var err error
	currentWorkingDirectory, err = os.Getwd()
	if err != nil {
		log.Println(err)
	}

	usr, _ := user.Current()
	userHomeDir = usr.HomeDir

	// If any flags containing a file path were passed in
	// on the command line we want to resolve them to absolute paths
	// relative to the current working directory. We'll make this call
	// again when we've parsed the config file, but those will be resolved
	// relative to the directory containing the config file.
	flag.Parse()
	flagBasePath = currentWorkingDirectory + "/"
	flag.Visit(resolveAbsoluteFlagPaths)

	// If a config file wasn't passed in on the command line we go looking for one
	// starting at the current working directory and traveling up the hierarchy
	_, err = os.Stat(configFilePath)
	if configFilePath == "" || os.IsNotExist(err) {
		configFilePath = findConfigFile(currentWorkingDirectory)
	}

	if configFilePath != "" {
		loadConfFile()
	} else {
		ui.Info("WARNING: No config file found")
	}
}

/**
 * loadGlobalConf
 * Helper method to load a configurationManager file and parse values, used by init and during
 * interactive configurationManager file generation as the globalconf package is able to handle
 * persisting flag values to disk.
 */
func loadConfFile() *globalconf.GlobalConf {
	var err error
	configurationManager, err = globalconf.NewWithOptions(&globalconf.Options{
		Filename: configFilePath,
	})

	if err != nil {
		log.Println(err)
	}

	absConfigFilePath, _ := filepath.Abs(configFilePath)
	configFileBasePath, _ := filepath.Split(absConfigFilePath)
	// Reads configurationManager data as provided in the config file
	// Path data provided will be expanded relative to the config file
	// any flags passed in via CLI will already be absolute and unaffected by
	// this repeated call
	configurationManager.ParseAll()
	flagBasePath = configFileBasePath
	flag.Visit(resolveAbsoluteFlagPaths)
	return configurationManager
}

/**
 * findConfigFile
 * Recursive method to walk backwards through a path looking for a configurationManager file
 */
func findConfigFile(directory string) string {
	filePath := checkDirForConfigFile(directory)
	if filePath == "" {
		climbIndex := strings.LastIndex(directory, "/")
		if climbIndex != -1 {
			return findConfigFile(directory[0:climbIndex])
		}
	} else {
		return filePath
	}

	return ""
}

/**
 * checkDirForConfigFile
 * Helper which either returns the full path of the configurationManager file if one is found
 * or returns the empty string if one is not.
 */
func checkDirForConfigFile(directory string) string {
	if len(directory) > 0 && directory[:len(directory)-1] != "/" {
		directory += "/"
	}
	filePath := directory + defaultConfigFileName
	if _, err := os.Stat(filePath); err == nil {
		return filePath
	} else {
		return ""
	}
}

/**
 * getAbsPath
 * Returns the absolute path given a string representation of a file path.In contrast to the
 * built in filepath.Abs method which always evaluates a path relative to the current
 * working directory, this method allows you to set a starting point for file path resolution.
 * This method also expands ~ to the home directory, which Abs does not support.
 *
 * To avoid confusion this method is always used to expand paths even when evaluating
 * paths relative to the current working directory.
 */
func getAbsPath(path string, base string) string {
	if base != currentWorkingDirectory {
		err := os.Chdir(base)
		// change working directory only while necessary
		defer os.Chdir(currentWorkingDirectory)
		if err != nil {
			panic(err)
		}
	}

	if "~" == path[:1] {
		path = userHomeDir + path[1:]
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return absPath
}

/**
 * flagStoresPathString
 * Helper method that is used to filter out flags which don't store
 * file path data.
 */
func flagStoresPathString(flagName string) bool {
	switch flagName {
	case "host", "port", "root-name":
		return false
	default:
		return true
	}
}

/**
 * resolveAbsoluteFlagPaths
 * Helper method which is passed to flag.Visit and ensures the file path strings
 * are properly expanded.
 */
func resolveAbsoluteFlagPaths(f *flag.Flag) {
	if flagStoresPathString(f.Name) {
		f.Value.Set(getAbsPath(f.Value.String(), flagBasePath))
	}
}

/**
 * GenerateServerConfig
 * Exported method that determines whether a configuration file needs to be created or whether
 * an existing file should be updated, then iterates over all server configuration options to
 * allow a user to set them.
 */
func GenerateServerConfig() {
	generateConfigFile()
	flag.VisitAll(iterateOverServerConfigFlags)
	ui.Info("All provided configuration information has been persisted to disk.\n")
}

/**
 * GenerateClientConfig
 * Exported method that determines whether a configuration file needs to be created or whether
 * an existing file should be updated, then iterates over all client configuration options to
 * allow a user to set them.
 */
func GenerateClientConfig() {
	generateConfigFile()
	flag.VisitAll(iterateOverClientConfigFlags)
	ui.Info("All provided configuration information has been persisted to disk.\n")
}

/**
 * generateConfigFile
 * Helper method that guides the user through interactive prompts and determines whether the
 * intention is to create a new configuration file or update the existing one.
 */
func generateConfigFile() {
	if configFilePath == "" {
		ok := getOrCreateConfigFile()
		if !ok {
			return
		}
	} else {
		ui.Info("A config file is already loaded from: " + configFilePath)
		resp, err := ui.Ask("U to Update the existing file, C to Create a new file somewhere else [U/C]:")
		if err != nil {
			log.Fatal(err)
		}
		if len(resp) > 1 {
			resp = resp[:1]
		}
		resp = strings.ToLower(resp)
		switch {
		default:
			ui.Info("No input detected, exiting...")
			return
		case "c" == resp:
			// this method will create an empty file
			ok := getOrCreateConfigFile()
			if !ok {
				return
			}
		case "u" == resp:
			ui.Info("Updating existing file...\n")
		}
	}

	ui.Info("Config file will be written to: " + configFilePath + "\n")
	loadConfFile()

	ui.Info("All paths entered on these prompts are relative to the current working directory")
	ui.Info("Any of the following options can be skipped by hitting return.")
	ui.Info("Skipped responses do not overwrite existing settings.\n")
}

/**
 * promptForAndSetConfigFilePath
 * Helper method to handle the logic of creating a new configurationManager file.
 */
func getOrCreateConfigFile() bool {
	cp, err := ui.Ask("Where should we write a new config file?")
	if err != nil {
		log.Fatal(err)
	}
	if cp == "" {
		ui.Warn("No input detected, exiting...")
		return false
	}
	fullPath := getAbsPath(cp, currentWorkingDirectory)
	filePath := checkDirForConfigFile(fullPath)
	if filePath == "" {
		filePath = fullPath + "/" + defaultConfigFileName
		_, err := os.Create(filePath)
		if err != nil {
			ui.Warn("Path invalid, path directories must already exist, exiting...")
			return false
		}
	} else {
		ui.Info("Config file " + filePath + " exists, updating in place...\n")
	}
	configFilePath = filePath
	return true
}

/**
 * iterateOverClientConfigFlags
 * Method which is meant to be called by flag.Visit or flag.VisitAll and offer an interactive
 * prompt to change and persist the value of flags which are used by the client.
 */
func iterateOverClientConfigFlags(f *flag.Flag) {
	if isClientConfigFlag(f.Name) {
		promptForAndPersistFlagValue(f)
	}
}

/**
 * isClientConfigFlag
 * Helper that returns true if a provided flagName is needed by the client binary.
 */
func isClientConfigFlag(flagName string) bool {
	switch flagName {
	case "host", "port", "root-cert", "root-name", "client-tls-cert", "client-tls-key":
		return true
	default:
		return false
	}
}

/**
 * iterateOverServerConfigFlags
 * Method which is meant to be called by flag.Visit or flag.VisitAll and offer an interactive
 * prompt to change and persist the value of flags which are used by the server.
 */
func iterateOverServerConfigFlags(f *flag.Flag) {
	if isServerConfigFlag(f.Name) {
		promptForAndPersistFlagValue(f)
	}
}

/**
 * isServerConfigFlag
 * Helper that returns true if a provided flagName is needed by the server binary.
 */
func isServerConfigFlag(flagName string) bool {
	switch flagName {
	case "host", "port", "root-cert", "server-tls-cert", "server-tls-key":
		return true
	default:
		return false
	}
}

/**
 * promptForAndPersistFlagValues
 * Helper method which uses the flag's usage string to prompt the user to enter a value
 * for the flag. That value is then assigned to the flag and persisted to disk using
 * globalconf's Set method.
 */
func promptForAndPersistFlagValue(f *flag.Flag) {
	response, err := ui.Ask(f.Usage)
	if err != nil {
		log.Fatal(err)
	}
	if response != "" {
		if flagStoresPathString(f.Name) {
			response = getAbsPath(response, currentWorkingDirectory)
		}
		f.Value.Set(response)
		configurationManager.Set("", f)
	}
}
