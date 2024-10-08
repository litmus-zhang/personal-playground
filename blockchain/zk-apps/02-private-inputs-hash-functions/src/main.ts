import {IncrementSecrete} from './IncrementSecrete.js'
import {Field, Mina, PrivateKey, AccountUpdate} from 'o1js'

// This is the main entry point for the zk-app.
const setup = async () => {
    const useProof = false;

const Local = await Mina.LocalBlockchain({proofsEnabled: useProof});
Mina.setActiveInstance(Local);
    const deployerAccount = Local.testAccounts[0];
    const deployerKey = deployerAccount.key;
    const senderAccount = Local.testAccounts[1];
    const senderKey = senderAccount.key;
    
    const zkAppPrivateKey = PrivateKey.random();
    const zkAppAddress = zkAppPrivateKey.toPublicKey();
    const zkAppInstance = new IncrementSecrete(zkAppAddress);
    return {deployerAccount, deployerKey, senderAccount, senderKey, zkAppInstance, zkAppPrivateKey};
}

const {deployerAccount, deployerKey, senderAccount, senderKey, zkAppInstance, zkAppPrivateKey} = await setup()

const salt = Field.random()

const deployTxn = await Mina.transaction(deployerAccount, async()=>{
    AccountUpdate.fundNewAccount(deployerAccount);
    await zkAppInstance.deploy();
    await zkAppInstance.initState(salt, Field(750))
})
await deployTxn.prove()
await deployTxn.sign([deployerKey, zkAppPrivateKey]).send();

const num0 = zkAppInstance.x.get();
console.log("State after init:", num0.toString());


const txn1 = await Mina.transaction(senderAccount, async()=>{ 
    await zkAppInstance.incrementSecret(salt, Field(750))
 })
 await txn1.prove()
 await txn1.sign([senderKey]).send();

 const num1 = zkAppInstance.x.get();
 console.log("State after increment:", num1.toString());