import React from "react";
import { useWalletSelector } from "@near-wallet-selector/react-hook";

const BlockchainDataInfo = () => {
  const { signedAccountId, signIn, signOut } = useWalletSelector();

  console.log(`signedAccountId: ${signedAccountId}`);

  return (
    <div>
      <h2>NEAR Blockchain Data</h2>
      {signedAccountId ? (
        <>
          <p>Signed in as: {signedAccountId}</p>
          <button onClick={signOut}>Sign Out</button>
        </>
      ) : (
        <>
          <p>You are not signed in.</p>
          <button onClick={signIn}>Sign In</button>
        </>
      )}
    </div>
  );
};

export default BlockchainDataInfo;
