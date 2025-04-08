# 🚀 Full Stack NEAR Project with Go

This tutorial walks you through setting up a full-stack NEAR project using Go, React, Node.js, and NEAR smart contracts.

---

## 1. Create Project

Run the following command to generate a full-stack project template:

```bash
near-go create -p "full_stack_template" \
  -m "github.com/vlmoon99/near-sdk-go/full_stack_template" \
  -t "full-stack-react-nodejs"
```

Navigate into the project folder:

```bash
cd full_stack_template
```

Your project structure should look like this:

```bash
(base) user@pc1:~/dev/near-sdk-go/examples/full_stack_template$ ls
backend  client  contract  contract_listener
```

### Folder Overview:

- **`backend/`**  
  Backend server written in Node.js. Use this to:
  - Execute blockchain operations (transfers, balance checks, etc.)
  - Interact with smart contracts
  - Store data in your database

- **`client/`**  
  React + Vite frontend to call your smart contract functions directly from the browser.

- **`contract/`**  
  Smart contract written in Rust (or Go). You can build and deploy this to NEAR testnet or mainnet.

- **`contract_listener/`**  
  A listener for smart contract events using NEAR Lake framework.  
  ⚠️ Requires an AWS account and EC2 access to run properly.

---

## 2. Setup Accounts for Deployment

### 🧪 Testnet Account

All testnet accounts must end with `.testnet`.  
You can check availability on [NEAR Blocks Testnet](https://testnet.nearblocks.io/address/mytestnetaccount779.testnet).

Create a testnet account using the CLI:

```bash
near-go account create -n "testnet" -a "mytestnetaccount779.testnet"
```

![Check testnet account on Near Blocks](./docs_images/tutorial_account_creation_1.jpeg)

---

### 🌐 Mainnet Account

1. Create a mainnet account using any NEAR-based wallet (In this tutorial I will use Meteor Wallet) :
   - [NEAR Wallets](https://wallet.near.org/)
   - [Meteor Wallet](https://wallet.meteorwallet.app/add_wallet/create_new)

![Choose mainnet and Choose your nickname and his availability](./docs_images/tutorial_account_creation_2.jpeg)

![Fund your account](./docs_images/tutorial_account_creation_3.jpeg)

2. Import your mainnet wallet into the CLI.

 - Call account import function
```bash
near-go account import
```
 - Choose type of account import and save mainnet account into your system

```bash
All necessary programs are installed.
? How would you like to import the account?  
> using-web-wallet          - Import existing account using NEAR Wallet (a.k.a. "sign in")
  using-seed-phrase         - Import existing account using a seed phrase
  using-private-key         - Import existing account using a private key
  back
[↑↓ to move, enter to select, type to filter]
```




3.write code and tests ? (think anout how can I show Near Blockchain to the go user and connect contract->client, contract-> backend , contract -> contract_listner)



4.deploy on the testnet,prod

5.Manage in the production

Tell abouts things (what if I delete all keys , how much storage costs , what the limits of the Near Blockchain)