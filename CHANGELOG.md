## Unreleased

- Introduce CloudAPI's ListRulesMachines under networking

## 0.2.1 (November 8)

- Fixing a bug where CreateUser and UpdateUser didn't return the UserID

## 0.2.0 (November 7)

- Introduce CloudAPI's Ping under compute
- Introduce CloudAPI's RebootMachine under compute instances
- Introduce CloudAPI's ListUsers, GetUser, CreateUser, UpdateUser and DeleteUser under identity package
- Introduce CloudAPI's ListMachineSnapshots, GetMachineSnapshot, CreateSnapshot, DeleteMachineSnapshot and StartMachineFromSnapshot under compute package
- tools: Introduce unit testing and scripts for linting, etc.
- bug: Fix the `compute.ListMachineRules` endpoint

## 0.1.0 (November 2)

- Initial release of a versioned SDK
