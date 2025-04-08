import React, { useState, useEffect } from "react";
import { useWalletSelector } from "@near-wallet-selector/react-hook";
import { utils } from "near-api-js";

const CONTRACT_ID = "neargocli.testnet";

const SmartContractOperations = () => {
  const { wallet, signedAccountId, signAndSendTransactions } = useWalletSelector();
  const [result, setResult] = useState(null);
  const [loggedIn, setLoggedIn] = useState(false);

  useEffect(() => {
    setLoggedIn(!!signedAccountId);
  }, [signedAccountId]);

  const executeFunction = async (methodName, args = {}, deposit = "0", gas = "300000000000000") => {
    if (!signedAccountId) {
      alert("Please sign in first");
      return;
    }

    try {
      const transaction = {
        signerId: signedAccountId,
        receiverId: CONTRACT_ID,
        actions: [
          {
            type: "FunctionCall",
            params: {
              methodName,
              args,
              gas,
              deposit: utils.format.parseNearAmount(deposit),
            },
          },
        ],
      };
      
      const result = await wallet.signAndSendTransaction(transaction);

      setResult(JSON.stringify(result.transaction.hash, null, 2));
    } catch (error) {
      console.error("Transaction Error:", error);
      setResult(`Error: ${error.message}`);
    }
  };

  return (
    <div>
      <h2>Smart Contract Operations</h2>
      {loggedIn ? (
        <>
          <button style={{ padding: "10px", margin: "5px" }} onClick={() => executeFunction("InitContract")}>Init Contract</button>
          <button style={{ padding: "10px", margin: "5px" }} onClick={() => executeFunction("WriteData", { key: "testKey", data: "lalalla" }, "1")}>Write Data</button>
          <button style={{ padding: "10px", margin: "5px" }} onClick={() => executeFunction("ReadData", { key: "testKey" })}>Read Data</button>
          <button style={{ padding: "10px", margin: "5px" }} onClick={() => executeFunction("AcceptPayment", {}, "1")}>Accept Payment</button>
          <button style={{ padding: "10px", margin: "5px" }} onClick={() => executeFunction("ReadIncommingTxData")}>Read Incoming Tx Data</button>
          <button style={{ padding: "10px", margin: "5px" }} onClick={() => executeFunction("ReadBlockchainData")}>Read Blockchain Data</button>
          {result && (
          <div style={{ marginTop: "10px", padding: "10px", border: "1px solid #ccc", borderRadius: "5px" }}>
            <p>Transaction Result:</p>
            {result.startsWith("Error") ? (
              <pre style={{ color: "red" }}>{result}</pre>
            ) : (
              <p>
                <a href={`https://testnet.nearblocks.io/txns/${result.slice(1, -1)}`} target="_blank" rel="noopener noreferrer">
                  View Transaction on NEAR Explorer
                </a>
              </p>
            )}
          </div>
        )}        
        </>
      ) : (
        <p>Please sign in to interact with the contract.</p>
      )}
    </div>
  );
};

export default SmartContractOperations;
