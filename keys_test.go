package triton

import (
	"github.com/abdullin/seq"
	"testing"
)

func TestAccKey_Create(t *testing.T) {
	keyName := RandPrefixString("TestAccCreateKey", 32)

	AccTest(t, TestCase{
		Steps: []Step{
			&StepAPICall{
				StateBagKey: "key",
				CallFunc: func(client *Client) (interface{}, error) {
					return client.Keys().CreateKey(&CreateKeyInput{
						Name: keyName,
						Key:  testAccCreateKeyMaterial,
					})
				},
				CleanupFunc: func(client *Client, callState interface{}) {
					client.Keys().DeleteKey(&DeleteKeyInput{KeyName: keyName})
				},
			},
			&StepAssert{
				StateBagKey: "key",
				Assertions: seq.Map{
					"name":        keyName,
					"key":         testAccCreateKeyMaterial,
					"fingerprint": testAccCreateKeyFingerprint,
				},
			},
		},
	})
}

func TestAccKey_Get(t *testing.T) {
	keyName := RandPrefixString("TestAccGetKey", 32)

	AccTest(t, TestCase{
		Steps: []Step{
			&StepAPICall{
				StateBagKey: "key",
				CallFunc: func(client *Client) (interface{}, error) {
					return client.Keys().CreateKey(&CreateKeyInput{
						Name: keyName,
						Key:  testAccCreateKeyMaterial,
					})
				},
				CleanupFunc: func(client *Client, callState interface{}) {
					client.Keys().DeleteKey(&DeleteKeyInput{KeyName: keyName})
				},
			},
			&StepAPICall{
				StateBagKey: "getKey",
				CallFunc: func(client *Client) (interface{}, error) {
					return client.Keys().GetKey(&GetKeyInput{KeyName: keyName})
				},
			},
			&StepAssert{
				StateBagKey: "getKey",
				Assertions: seq.Map{
					"name":        keyName,
					"key":         testAccCreateKeyMaterial,
					"fingerprint": testAccCreateKeyFingerprint,
				},
			},
		},
	})
}

func TestAccKey_Delete(t *testing.T) {
	keyName := RandPrefixString("TestAccGetKey", 32)

	AccTest(t, TestCase{
		Steps: []Step{
			&StepAPICall{
				StateBagKey: "key",
				CallFunc: func(client *Client) (interface{}, error) {
					return client.Keys().CreateKey(&CreateKeyInput{
						Name: keyName,
						Key:  testAccCreateKeyMaterial,
					})
				},
				CleanupFunc: func(client *Client, callState interface{}) {
					client.Keys().DeleteKey(&DeleteKeyInput{KeyName: keyName})
				},
			},
			&StepAPICall{
				StateBagKey: "noop",
				CallFunc: func(client *Client) (interface{}, error) {
					return nil, client.Keys().DeleteKey(&DeleteKeyInput{KeyName: keyName})
				},
			},
			&StepAPICall{
				ErrorKey: "getKeyError",
				CallFunc: func(client *Client) (interface{}, error) {
					return client.Keys().GetKey(&GetKeyInput{KeyName: keyName})
				},
			},
			&StepAssertTritonError{
				ErrorKey: "getKeyError",
				Code:     "ResourceNotFound",
			},
		},
	})
}

const testAccCreateKeyMaterial = `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDBOJ5z6jTdY3SYK2Nc+MQLSQstAOzxFqDN00MJ9SMhJea8ZQbZFlhCAZBFE4TUBDI3zXBxFjygh84lb1QlNu1dmZeoQ10MThuowZllBAfg9Eb5RkXqLvDdYh9+rLdEdUL4+aiYZ8JYtQ+K5ZnogZoxdzNQ3WnVhMGJIrj1zcRveUSvQ6tMhaEQDxDWrAMDLxnLI/6SNmkhdF1ZKE8iQ+BnazYp0vg5jAzkHzEYJY9kFUOubupOxio93B9OTkpQ0jZD+J9iR1t8Me3JdhHy85inaAFc0fkjznDYluV8aqfIprD/WE9grQ/GfEYfsvQdQr1ljLBJZdad7DvnKqU0M4vJ James@jn-mpb15`
const testAccCreateKeyFingerprint = `ab:f4:8f:bc:26:e1:cf:1d:06:a3:9d:40:39:7c:5a:78`
