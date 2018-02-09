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
		KeyID:       viper.GetString("SDC_KEY_ID"),
		AccountName: viper.GetString("SDC_ACCOUNT"),
	})
	if err != nil {
		log.Fatal().Str("func", "initConfig").Msg("Error Creating SSH Agent Signer")
		return nil, err
	}

	config := &triton.ClientConfig{
		TritonURL:   viper.GetString("SDC_URL"),
		AccountName: viper.GetString("SDC_ACCOUNT"),
		Signers:     []authentication.Signer{signer},
	}

	return &TritonClientConfig{
		Config: config,
	}, nil
}

func GetPkgID() string {
	if viper.IsSet(config.KeyPackageId) {
		return viper.GetString(config.KeyPackageId)
	}
	return ""
}

func GetPkgName() string {
	if viper.IsSet(config.KeyPackageName) {
		return viper.GetString(config.KeyPackageName)
	}
	return ""
}

func GetImgID() string {
	if viper.IsSet(config.KeyImageId) {
		return viper.GetString(config.KeyImageId)
	}
	return ""
}

func GetImgName() string {
	if viper.IsSet(config.KeyImageName) {
		return viper.GetString(config.KeyImageName)
	}
	return ""
}

func GetMachineID() string {
	if viper.IsSet(config.KeyInstanceID) {
		return viper.GetString(config.KeyInstanceID)
	}
	return ""
}

func GetMachineName() string {
	if viper.IsSet(config.KeyInstanceName) {
		return viper.GetString(config.KeyInstanceName)
	}
	return ""
}

func GetMachineNamePrefix() string {
	return viper.GetString(config.KeyInstanceNamePrefix)
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
		metadata := make(map[string]string, 0)
		cfg := viper.GetStringSlice(config.KeyInstanceTag)
		for _, i := range cfg {
			m := strings.Split(i, "=")
			metadata[m[0]] = m[1]
		}

		return metadata
	}

	return nil
}

func GetSearchTags() map[string]interface{} {
	if viper.IsSet(config.KeyInstanceSearchTag) {
		tags := make(map[string]interface{}, 0)
		cfg := viper.GetStringSlice(config.KeyInstanceSearchTag)
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
	if viper.IsSet(config.KeyInstanceUserdata) {
		return viper.GetString(config.KeyInstanceUserdata)
	}

	return ""
}

func IsBlockingAction() bool {
	return viper.GetBool(config.KeyInstanceWait)
}
