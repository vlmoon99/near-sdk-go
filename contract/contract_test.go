package contract

import (
	"testing"

	"github.com/vlmoon99/near-sdk-go/borsh"
	"github.com/vlmoon99/near-sdk-go/env"
	"github.com/vlmoon99/near-sdk-go/system"
	"github.com/vlmoon99/near-sdk-go/types"
)

func init() {
	systemMock := system.NewMockSystem()

	expectedAccountID := "test.account"
	expectedSignerID := "test.signer"
	expectedPredecessorID := "test.predecessor"
	expectedGas := uint64(1000)
	expectedDeposit, _ := types.U128FromString("10000000000000000000000000")

	systemMock.CurrentAccountIdSys = expectedAccountID
	systemMock.SignerAccountIdSys = expectedSignerID
	systemMock.PredecessorAccountIdSys = expectedPredecessorID
	systemMock.AttachedDepositSys = expectedDeposit
	systemMock.PrepaidGasSys = expectedGas

	env.SetEnv(systemMock)
}

type TestStruct struct {
	Value string
}

func (t TestStruct) Serialize() ([]byte, error) {
	return borsh.Serialize(t)
}

func TestReturnValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected []byte
		wantErr  bool
	}{
		{
			name:     "Return byte array",
			input:    []byte("test"),
			expected: []byte("test"),
			wantErr:  false,
		},
		{
			name:     "Return serializable struct",
			input:    TestStruct{Value: "test"},
			expected: []byte("test"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ReturnValue(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReturnValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetContext(t *testing.T) {
	expectedAccountID := "test.account"
	expectedSignerID := "test.signer"
	expectedPredecessorID := "test.predecessor"
	expectedGas := uint64(1000)
	expectedDeposit, _ := types.U128FromString("10000000000000000000000000")

	context := GetContext()

	if context.AccountID != expectedAccountID {
		t.Errorf("GetContext().AccountID = %v, want %v", context.AccountID, expectedAccountID)
	}
	if context.SignerID != expectedSignerID {
		t.Errorf("GetContext().SignerID = %v, want %v", context.SignerID, expectedSignerID)
	}
	if context.PredecessorID != expectedPredecessorID {
		t.Errorf("GetContext().PredecessorID = %v, want %v", context.PredecessorID, expectedPredecessorID)
	}
	if context.AttachedDeposit != expectedDeposit {
		t.Errorf("GetContext().AttachedDeposit = %v, want %v", context.AttachedDeposit, expectedDeposit)
	}
	if context.PrepaidGas != expectedGas {
		t.Errorf("GetContext().PrepaidGas = %v, want %v", context.PrepaidGas, expectedGas)
	}
}

func TestRequireDeposit(t *testing.T) {
	equalDeposit, _ := types.U128FromString("10000000000000000000000000")
	sufficientDeposit, _ := types.U128FromString("11000000000000000000000000")
	insufficientDeposit, _ := types.U128FromString("1000")

	tests := []struct {
		name            string
		minDeposit      types.Uint128
		attachedDeposit types.Uint128
		wantErr         bool
	}{
		{
			name:            "Sufficient deposit",
			minDeposit:      equalDeposit,
			attachedDeposit: sufficientDeposit,
			wantErr:         false,
		},
		{
			name:            "Insufficient deposit",
			minDeposit:      equalDeposit,
			attachedDeposit: insufficientDeposit,
			wantErr:         true,
		},
		{
			name:            "Equal deposit",
			minDeposit:      equalDeposit,
			attachedDeposit: equalDeposit,
			wantErr:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSys := env.NearBlockchainImports.(*system.MockSystem)
			mockSys.AttachedDepositSys = tt.attachedDeposit

			err := RequireDeposit(tt.minDeposit)
			if (err != nil) != tt.wantErr {
				t.Errorf("RequireDeposit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHandleClientJSONInput(t *testing.T) {
	mockSys := env.NearBlockchainImports.(*system.MockSystem)
	jsonData := []byte(`{"test": "value"}`)
	mockSys.ContractInput = jsonData

	success := false
	HandleClientJSONInput(func(input *ContractInput) error {
		success = true
		return nil
	})

	if !success {
		t.Error("HandleClientJSONInput() did not execute callback")
	}

}

func TestHandleClientRawBytesInput(t *testing.T) {
	mockSys := env.NearBlockchainImports.(*system.MockSystem)
	rawData := []byte("raw bytes")
	mockSys.ContractInput = rawData

	success := false
	HandleClientRawBytesInput(func(input *ContractInput) error {
		success = true
		return nil
	})

	if !success {
		t.Error("HandleClientRawBytesInput() did not execute callback")
	}

}
