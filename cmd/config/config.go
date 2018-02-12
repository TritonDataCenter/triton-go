package config

import (
	"strings"

	triton "github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/joyent/triton-go/cmd/internal/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type TritonClientConfig struct {
	Config *triton.ClientConfig
}

func New() (*TritonClientConfig, error) {
	viper.AutomaticEnv()

	var signer authentication.Signer
	var err error

	signer, err = authentication.NewSSHAgentSigner(authentication.SSHAgentSignerInput{
		KeyID:       GetTritonKeyID(),
		AccountName: GetTritonAccount(),
	})
	if err != nil {
		log.Fatal().Str("func", "initConfig").Msg("Error Creating SSH Agent Signer")
		return nil, err
	}

	config := &triton.ClientConfig{
		TritonURL:   GetTritonUrl(),
		AccountName: GetTritonAccount(),
		Signers:     []authentication.Signer{signer},
	}

	return &TritonClientConfig{
		Config: config,
	}, nil
}

var envPrefixes = []string{"TRITON", "SDC"}

func getEnvVar(name string) string {
	for _, prefix := range envPrefixes {
		if val := viper.GetString(prefix + "_" + name); val != "" {
			return val
		}
	}

	return ""
}

func GetTritonUrl() string {
	url := viper.GetString(config.KeyUrl)
	if url == "" {
		url = getEnvVar("URL")
	}

	return url
}

func GetTritonAccount() string {
	account := viper.GetString(config.KeyAccount)
	if account == "" {
		account = getEnvVar("ACCOUNT")
	}

	return account
}

func GetTritonKeyID() string {
	keyID := viper.GetString(config.KeySshKeyID)
	if keyID == "" {
		keyID = getEnvVar("KEY_ID")
	}

	return keyID
}

func GetPkgID() string {
	return viper.GetString(config.KeyPackageId)
}

func GetPkgName() string {
	return viper.GetString(config.KeyPackageName)
}

func GetImgID() string {
	return viper.GetString(config.KeyImageId)
}

func GetImgName() string {
	return viper.GetString(config.KeyImageName)
}

func GetMachineID() string {
	return viper.GetString(config.KeyInstanceID)
}

func GetMachineName() string {
	return viper.GetString(config.KeyInstanceName)
}

func GetMachineState() string {
	return viper.GetString(config.KeyInstanceState)
}

func GetMachineBrand() string {
	return viper.GetString(config.KeyInstanceBrand)
}

func GetMachineFirewall() bool {
	return viper.GetBool(config.KeyInstanceFirewall)
}

func GetMachineNetworks() []string {
	if viper.IsSet(config.KeyInstanceNetwork) {
		var networks []string
		cfg := viper.GetStringSlice(config.KeyInstanceNetwork)
		for _, i := range cfg {
			networks = append(networks, i)
		}

		return networks
	}
	return nil
}

func GetMachineAffinityRules() []string {
	if viper.IsSet(config.KeyInstanceAffinityRule) {
		var rules []string
		cfg := viper.GetStringSlice(config.KeyInstanceAffinityRule)
		for _, i := range cfg {
			rules = append(rules, i)
		}

		return rules
	}
	return nil
}

func GetMachineTags() map[string]string {
	if viper.IsSet(config.KeyInstanceTag) {
		tags := make(map[string]string, 0)
		cfg := viper.GetStringSlice(config.KeyInstanceTag)
		for _, i := range cfg {
			m := strings.Split(i, "=")
			tags[m[0]] = m[1]
		}

		return tags
	}

	return nil
}

func GetSearchTags() map[string]interface{} {
	if viper.IsSet(config.KeyInstanceTag) {
		tags := make(map[string]interface{}, 0)
		cfg := viper.GetStringSlice(config.KeyInstanceTag)
		for _, i := range cfg {
			m := strings.Split(i, "=")
			tags[m[0]] = m[1]
		}

		return tags
	}

	return nil
}

func GetMachineMetadata() map[string]string {
	if viper.IsSet(config.KeyInstanceMetadata) {
		metadata := make(map[string]string, 0)
		cfg := viper.GetStringSlice(config.KeyInstanceMetadata)
		for _, i := range cfg {
			m := strings.Split(i, "=")
			metadata[m[0]] = m[1]
		}

		return metadata
	}

	return nil
}

func GetMachineUserdata() string {
	return viper.GetString(config.KeyInstanceUserdata)
}

func IsBlockingAction() bool {
	return viper.GetBool(config.KeyInstanceWait)
}
