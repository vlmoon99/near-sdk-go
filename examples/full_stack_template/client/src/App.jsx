import '@near-wallet-selector/modal-ui/styles.css';

import React from 'react';
import './App.css'

import { setupMeteorWallet } from "@near-wallet-selector/meteor-wallet";
import { setupLedger } from "@near-wallet-selector/ledger";
import { WalletSelectorProvider } from "@near-wallet-selector/react-hook";
import BlockchainDataInfo from "./BlockchainDataInfo";
import SmartContractOperations from "./SmartContractOperations";

const walletSelectorConfig = {
  network: "testnet", 
  createAccessKeyFor: "neargocli.testnet",
  modules: [
    setupMeteorWallet(),
    setupLedger()
  ],
}

export default function App({ Component }) {

  return (
    <WalletSelectorProvider config={walletSelectorConfig}>
      <SmartContractOperations/>
      <BlockchainDataInfo /> 
    </WalletSelectorProvider>
  );
}