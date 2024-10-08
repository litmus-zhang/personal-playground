import { Square } from './Square.js'
import { Field, Mina, PrivateKey, AccountUpdate } from 'o1js'

// This is the main entry point for the zk-app.
const setup = async () => {
    const useProof = false;

    const Local = await Mina.LocalBlockchain({ proofsEnabled: useProof });
    Mina.setActiveInstance(Local);
    const deployerAccount = Local.testAccounts[0];
    const deployerKey = deployerAccount.key;
    const senderAccount = Local.testAccounts[1];
    const senderKey = senderAccount.key;

    const zkAppPrivateKey = PrivateKey.random();
    const zkAppAddress = zkAppPrivateKey.toPublicKey();

    const zkAppInstance = new Square(zkAppAddress);
    return { deployerAccount, deployerKey, senderAccount, senderKey, zkAppInstance, zkAppPrivateKey }
}

const { deployerAccount, deployerKey, senderAccount, senderKey, zkAppInstance, zkAppPrivateKey } = await setup()
// Deploy the zk-app
const deployTxn = await Mina.transaction(deployerAccount, async () => {
    AccountUpdate.fundNewAccount(deployerAccount);
    await zkAppInstance.deploy();
});
await deployTxn.sign([deployerKey, zkAppPrivateKey]).send();
const num0 = zkAppInstance.num.get();
console.log("Initial value of num: ", num0.toString());

// Update the zk-app
const txn1 = await Mina.transaction(senderAccount, async () => {
    await zkAppInstance.update(Field(9));
});
await txn1.prove();
await txn1.sign([senderKey]).send();

const num1 = zkAppInstance.num.get();
console.log("Updated value of num: ", num1.toString());

try {
    const txn2 = await Mina.transaction(senderAccount, async () => {
        await zkAppInstance.update(Field(81));
    });
    await txn2.prove();
    await txn2.sign([senderKey]).send();
} catch (error: any) {
    console.error(error.message);
}

const num2 = zkAppInstance.num.get();
console.log("Value of num after failed update: ", num2.toString());