package promise

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
	"github.com/vlmoon99/near-sdk-go/types"
)

func init() {
	systemMock := system.NewMockSystem()
	env.SetEnv(systemMock)
}

func TestNewPromise(t *testing.T) {
	promiseID := uint64(123)
	promise := NewPromise(promiseID)

	if promise.promiseID != promiseID {
		t.Errorf("NewPromise() promiseID = %v, want %v", promise.promiseID, promiseID)
	}
	if promise.gas != DefaultGas {
		t.Errorf("NewPromise() gas = %v, want %v", promise.gas, DefaultGas)
	}
	if promise.deposit != (types.Uint128{Hi: 0, Lo: 0}) {
		t.Errorf("NewPromise() deposit = %v, want %v", promise.deposit, types.Uint128{Hi: 0, Lo: 0})
	}
}

func TestPromise_Gas(t *testing.T) {
	promise := NewPromise(123)
	gas := uint64(1000)
	newPromise := promise.Gas(gas)

	if newPromise.gas != gas {
		t.Errorf("Gas() = %v, want %v", newPromise.gas, gas)
	}
	if newPromise.promiseID != promise.promiseID {
		t.Errorf("Gas() promiseID = %v, want %v", newPromise.promiseID, promise.promiseID)
	}
}

func TestPromise_Deposit(t *testing.T) {
	promise := NewPromise(123)
	deposit := types.Uint128{Hi: 1, Lo: 2}
	newPromise := promise.Deposit(deposit)

	if newPromise.deposit != deposit {
		t.Errorf("Deposit() = %v, want %v", newPromise.deposit, deposit)
	}
	if newPromise.promiseID != promise.promiseID {
		t.Errorf("Deposit() promiseID = %v, want %v", newPromise.promiseID, promise.promiseID)
	}
}

func TestPromise_DepositYocto(t *testing.T) {
	promise := NewPromise(123)
	amount := uint64(1000)
	newPromise := promise.DepositYocto(amount)

	expected := types.U64ToUint128(amount)
	if newPromise.deposit != expected {
		t.Errorf("DepositYocto() = %v, want %v", newPromise.deposit, expected)
	}
}

func TestNewPromiseBatch(t *testing.T) {
	promiseID := uint64(123)
	batch := NewPromiseBatch(promiseID)

	if batch.promiseID != promiseID {
		t.Errorf("NewPromiseBatch() promiseID = %v, want %v", batch.promiseID, promiseID)
	}
	if batch.gas != DefaultGas {
		t.Errorf("NewPromiseBatch() gas = %v, want %v", batch.gas, DefaultGas)
	}
}

func TestPromiseBatch_Gas(t *testing.T) {
	batch := NewPromiseBatch(123)
	gas := uint64(1000)
	newBatch := batch.Gas(gas)

	if newBatch.gas != gas {
		t.Errorf("Gas() = %v, want %v", newBatch.gas, gas)
	}
	if newBatch.promiseID != batch.promiseID {
		t.Errorf("Gas() promiseID = %v, want %v", newBatch.promiseID, batch.promiseID)
	}
}

func TestNewCrossContract(t *testing.T) {
	accountID := "test.account"
	contract := NewCrossContract(accountID)

	if contract.accountID != accountID {
		t.Errorf("NewCrossContract() accountID = %v, want %v", contract.accountID, accountID)
	}
	if contract.gas != DefaultGas {
		t.Errorf("NewCrossContract() gas = %v, want %v", contract.gas, DefaultGas)
	}
	if contract.deposit != (types.Uint128{Hi: 0, Lo: 0}) {
		t.Errorf("NewCrossContract() deposit = %v, want %v", contract.deposit, types.Uint128{Hi: 0, Lo: 0})
	}
}

func TestCrossContract_Gas(t *testing.T) {
	contract := NewCrossContract("test.account")
	gas := uint64(1000)
	newContract := contract.Gas(gas)

	if newContract.gas != gas {
		t.Errorf("Gas() = %v, want %v", newContract.gas, gas)
	}
	if newContract.accountID != contract.accountID {
		t.Errorf("Gas() accountID = %v, want %v", newContract.accountID, contract.accountID)
	}
}

func TestCrossContract_Deposit(t *testing.T) {
	contract := NewCrossContract("test.account")
	deposit := types.Uint128{Hi: 1, Lo: 2}
	newContract := contract.Deposit(deposit)

	if newContract.deposit != deposit {
		t.Errorf("Deposit() = %v, want %v", newContract.deposit, deposit)
	}
	if newContract.accountID != contract.accountID {
		t.Errorf("Deposit() accountID = %v, want %v", newContract.accountID, contract.accountID)
	}
}

func TestCrossContract_DepositYocto(t *testing.T) {
	contract := NewCrossContract("test.account")
	amount := uint64(1000)
	newContract := contract.DepositYocto(amount)

	expected := types.U64ToUint128(amount)
	if newContract.deposit != expected {
		t.Errorf("DepositYocto() = %v, want %v", newContract.deposit, expected)
	}
}

func TestPromiseResult_Unwrap(t *testing.T) {
	tests := []struct {
		name        string
		result      PromiseResult
		wantData    []byte
		wantErr     bool
		errContains string
	}{
		{
			name:        "Successful result",
			result:      NewPromiseResult(1, []byte("success")),
			wantData:    []byte("success"),
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "Failed result",
			result:      NewPromiseResult(2, []byte("error")),
			wantData:    nil,
			wantErr:     true,
			errContains: "promise failed with status: Failed",
		},
		{
			name:        "Not ready result",
			result:      NewPromiseResult(0, nil),
			wantData:    nil,
			wantErr:     true,
			errContains: "promise failed with status: NotReady",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.result.Unwrap()
			if (err != nil) != tt.wantErr {
				t.Errorf("PromiseResult.Unwrap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errContains {
				t.Errorf("PromiseResult.Unwrap() error = %v, want %v", err, tt.errContains)
			}
			if !tt.wantErr && string(data) != string(tt.wantData) {
				t.Errorf("PromiseResult.Unwrap() = %v, want %v", data, tt.wantData)
			}
		})
	}
}

func TestPromiseResult_UnwrapOr(t *testing.T) {
	tests := []struct {
		name           string
		result         PromiseResult
		defaultValue   []byte
		expectedResult []byte
	}{
		{
			name:           "Successful result",
			result:         NewPromiseResult(1, []byte("success")),
			defaultValue:   []byte("default"),
			expectedResult: []byte("success"),
		},
		{
			name:           "Failed result",
			result:         NewPromiseResult(2, []byte("error")),
			defaultValue:   []byte("default"),
			expectedResult: []byte("default"),
		},
		{
			name:           "Not ready result",
			result:         NewPromiseResult(0, nil),
			defaultValue:   []byte("default"),
			expectedResult: []byte("default"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.result.UnwrapOr(tt.defaultValue)
			if string(result) != string(tt.expectedResult) {
				t.Errorf("PromiseResult.UnwrapOr() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestPromiseResult_UnwrapToParser(t *testing.T) {
	tests := []struct {
		name        string
		result      PromiseResult
		wantErr     bool
		errContains string
	}{
		{
			name:        "Successful result",
			result:      NewPromiseResult(1, []byte(`{"key":"value"}`)),
			wantErr:     false,
			errContains: "",
		},
		{
			name:        "Failed result",
			result:      NewPromiseResult(2, []byte("error")),
			wantErr:     true,
			errContains: "promise failed with status: Failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := tt.result.Unwrap()
			if (err != nil) != tt.wantErr {
				t.Errorf("PromiseResult.UnwrapToParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errContains {
				t.Errorf("PromiseResult.UnwrapToParser() error = %v, want %v", err, tt.errContains)
			}
			if !tt.wantErr && data == nil {
				t.Error("PromiseResult.UnwrapToParser() data is nil")
			}
		})
	}
}
