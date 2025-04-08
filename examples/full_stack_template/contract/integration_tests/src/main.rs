use near_gas::NearGas;
use near_workspaces::types::NearToken;
use serde_json::json;

async fn deploy_contract(
    worker: &near_workspaces::Worker<near_workspaces::network::Sandbox>,
) -> anyhow::Result<near_workspaces::Contract> {
    const WASM_FILEPATH: &str = "../main.wasm";
    let wasm = std::fs::read(WASM_FILEPATH)?;
    let contract = worker
    .dev_deploy(&wasm)
    .await?;

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
        .await;

    match outcome {
        Ok(result) => {
            println!("✅ Function '{}' executed successfully.", function_name);
            println!("Success: {:#?}", result.clone().is_success());
            println!("Logs: {:#?}", result.clone().logs());
            Ok(())
        }
        Err(err) => {
            println!("❌ Test failed for '{}' with error: {:#?}", function_name, err);
            Err(err.into())
        }
    }
}



#[tokio::main]
async fn main() -> anyhow::Result<()>  {
    let worker = near_workspaces::sandbox().await?;
    let contract = deploy_contract(&worker).await?;
    let standard_deposit = NearToken::from_near(3);
    let standard_gas = NearGas::from_tgas(300);
    println!("Dev Account ID: {}", contract.id());

    let success_results = vec![
        call_integration_test_function(
            &contract,
            "InitContract",
            json!({}),
            standard_deposit,
            standard_gas,
        ).await,
        call_integration_test_function(
            &contract,
            "WriteData",
            json!({ "key": "testKey", "data": "testData" }),
            standard_deposit,
            standard_gas,
        ).await,
        call_integration_test_function(
            &contract,
            "ReadData",
            json!({ "key": "testKey" }),
            standard_deposit,
            standard_gas,
        ).await,
        call_integration_test_function(
            &contract,
            "AcceptPayment",
            json!({}),
            standard_deposit,
            standard_gas,
        ).await,
        call_integration_test_function(
            &contract,
            "ReadIncommingTxData",
            json!({}),
            standard_deposit,
            standard_gas,
        ).await,
        call_integration_test_function(
            &contract,
            "ReadBlockchainData",
            json!({}),
            standard_deposit,
            standard_gas,
        ).await,
    ];

    for result in success_results {
        if let Err(e) = result {
            eprintln!("Error: {:?}", e);
        }
    }

    Ok(())
}
