use near_gas::NearGas;
use near_workspaces::types::NearToken;
use serde_json::json;

async fn deploy_contract(
    worker: &near_workspaces::Worker<near_workspaces::network::Sandbox>,
) -> anyhow::Result<near_workspaces::Contract> {
    const NFT_WASM_FILEPATH: &str = "../main.wasm";
    let wasm = std::fs::read(NFT_WASM_FILEPATH)?;
    let contract = worker.dev_deploy(&wasm).await?;
    Ok(contract)
}

async fn call_integration_test_function(
    contract: &near_workspaces::Contract,
    function_name: &str,
    args: serde_json::Value,
    deposit: NearToken,
    gas: NearGas,
) -> anyhow::Result<()> {
    let outcome = contract
        .call(function_name)
        .args_json(args)
        .deposit(deposit)
        .gas(gas)
        .transact()
        .await?;

    let result = unsafe { outcome.clone().json::<i8>().unwrap_unchecked() };

    if result == 1 {
        println!("{} result: Test succeeded", function_name);
        println!("result: Functions Logs: {:#?}", outcome.logs());
    } else {
        println!("{} result: Test failed", function_name);
    }

    Ok(())
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let worker = near_workspaces::sandbox().await?;
    let contract = deploy_contract(&worker).await?;
    let standard_deposit = NearToken::from_near(3);
    let standard_gas = NearGas::from_tgas(300);

    call_integration_test_function(
        &contract,
        "InitContract",
        json!({ "testInputKey": "testInputValue" }),
        standard_deposit,
        standard_gas,
    )
    .await?;

    // Registers API
    call_integration_test_function(
        &contract,
        "TestWriteReadRegisterSafe",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;
    // Registers API

    // Storage API

    call_integration_test_function(
        &contract,
        "TestStorageWrite",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestStorageRead",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestStorageHasKey",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestStorageRemove",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestStateWrite",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestStateRead",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestStateExists",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    // // Not Working, throws unreachable
    // call_integration_test_function(
    //     &contract,
    //     "TestStorageGetEvicted",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // ).await?;

    // Storage API

    // Context API
    call_integration_test_function(
        &contract,
        "TestGetCurrentAccountId",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestGetSignerAccountID",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestGetSignerAccountPK",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestGetPredecessorAccountID",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestGetCurrentBlockHeight",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestGetBlockTimeMs",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestGetEpochHeight",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestGetStorageUsage",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestContractInputRawBytes",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestContractInputJSON",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;
    // Context API

    // Economics API
    call_integration_test_function(
        &contract,
        "TestGetAccountBalance",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestGetAccountLockedBalance",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestGetAttachedDeposit",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestGetPrepaidGas",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestGetUsedGas",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;
    // Economics API

    // Math API
    call_integration_test_function(
        &contract,
        "TestGetRandomSeed",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestSha256Hash",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestKeccak256Hash",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestKeccak512Hash",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestRipemd160Hash",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    // call_integration_test_function(
    //     &contract,
    //     "TestEcrecoverPubKey",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // ).await?;

    // call_integration_test_function(
    //     &contract,
    //     "TestEd25519VerifySig",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // ).await?;

    // call_integration_test_function(
    //     &contract,
    //     "TestAltBn128G1MultiExp",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // ).await?;

    // call_integration_test_function(
    //     &contract,
    //     "TestAltBn128G1Sum",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // ).await?;

    // call_integration_test_function(
    //     &contract,
    //     "TestAltBn128PairingCheck",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // ).await?;

    // Math API

    // Validator API
    call_integration_test_function(
        &contract,
        "TestValidatorStakeAmount",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestValidatorTotalStakeAmount",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;
    // Validator API

    // Miscellaneous API

    call_integration_test_function(
        &contract,
        "TestContractValueReturn",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    // It's unnesessary to test, but if u want unncoment and tests it, it will panic in tests
    // call_integration_test_function(
    //     &contract,
    //     "TestPanicStr",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // ).await?;

    // It's unnesessary to test, but if u want unncoment and tests it, it will panic in tests
    // call_integration_test_function(
    //     &contract,
    //     "TestAbortExecution",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // ).await?;

    call_integration_test_function(
        &contract,
        "TestLogString",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestLogStringUtf8",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestLogStringUtf16",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    // Miscellaneous API

    // Promises API
    call_integration_test_function(
        &contract,
        "TestPromiseCreate",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestPromiseThen",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestPromiseAnd",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestPromiseBatchCreate",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestPromiseBatchThen",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    // Promises API

    // Promises API Action
    call_integration_test_function(
        &contract,
        "TestPromiseBatchActionCreateAccount",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestPromiseBatchActionDeployContract",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestPromiseBatchActionFunctionCall",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestPromiseBatchActionFunctionCallWeight",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestPromiseBatchActionTransfer",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestPromiseBatchActionStake",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    // call_integration_test_function(
    //     &contract,
    //     "TestPromiseBatchActionAddKeyWithFullAccess",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // )
    // .await?;

    // call_integration_test_function(
    //     &contract,
    //     "TestPromiseBatchActionAddKeyWithFunctionCall",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // )
    // .await?;

    // call_integration_test_function(
    //     &contract,
    //     "TestPromiseBatchActionDeleteKey",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // )
    // .await?;

    call_integration_test_function(
        &contract,
        "TestPromiseBatchActionDeleteAccount",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    // call_integration_test_function(
    //     &contract,
    //     "TestPromiseYieldCreate",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // )
    // .await?;

    // call_integration_test_function(
    //     &contract,
    //     "TestPromiseYieldResume",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // )
    // .await?;
    // Promises API Action

    // Promise API Results
    call_integration_test_function(
        &contract,
        "TestPromiseResultsCount",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestPromiseResult",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;

    call_integration_test_function(
        &contract,
        "TestPromiseReturn",
        json!({}),
        standard_deposit,
        standard_gas,
    )
    .await?;
    // Promise API Results

    Ok(())
}
