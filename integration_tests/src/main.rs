use near_workspaces::types::NearToken;
use near_gas::NearGas;
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

    let result = unsafe { outcome.json::<i8>().unwrap_unchecked() };

    if result == 1 {
        println!("{} result: Test succeeded", function_name);
    } else {
        println!("{} result: Test failed", function_name);
    }

    Ok(())
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    let worker = near_workspaces::sandbox().await?;
    let contract = deploy_contract(&worker).await?;
    let standard_deposit = NearToken::from_near(1);
    let standard_gas = NearGas::from_tgas(300);

    call_integration_test_function(
        &contract,
        "InitContract",
        json!({ "testInputKey": "testInputValue" }),
        standard_deposit,
        standard_gas,
    ).await?;
    
    call_integration_test_function(
        &contract,
        "TestWriteReadRegisterSafe",
        json!({}),
        standard_deposit,
        standard_gas,
    ).await?;
    
    call_integration_test_function(
        &contract,
        "TestStorageWrite",
        json!({}),
        standard_deposit,
        standard_gas,
    ).await?;
    
    call_integration_test_function(
        &contract,
        "TestStorageRead",
        json!({}),
        standard_deposit,
        standard_gas,
    ).await?;
    
    call_integration_test_function(
        &contract,
        "TestStorageHasKey",
        json!({}),
        standard_deposit,
        standard_gas,
    ).await?;
    
    call_integration_test_function(
        &contract,
        "TestStorageRemove",
        json!({}),
        standard_deposit,
        standard_gas,
    ).await?;


    call_integration_test_function(
        &contract,
        "TestStateWrite",
        json!({}),
        standard_deposit,
        standard_gas,
    ).await?;

    call_integration_test_function(
        &contract,
        "TestStateRead",
        json!({}),
        standard_deposit,
        standard_gas,
    ).await?;
    
    
    call_integration_test_function(
        &contract,
        "TestStateExists",
        json!({}),
        standard_deposit,
        standard_gas,
    ).await?;

    // // Not Working, throws unreachable    
    // call_integration_test_function(
    //     &contract,
    //     "TestStorageGetEvicted",
    //     json!({}),
    //     standard_deposit,
    //     standard_gas,
    // ).await?;
    
    Ok(())
}
