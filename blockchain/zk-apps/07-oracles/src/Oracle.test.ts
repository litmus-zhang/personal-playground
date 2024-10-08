import { Oracle } from './Oracle';
import {Field, Mina, PrivateKey, PublicKey, AccountUpdate, Signature} from 'o1js'

const ORACLE_PUBLIC_KEY = 'B62qoAE4rBRuTgC42vqvEyUqCGhaZsW58SKVW4Ht8aYqP9UTvxFWBgy';
let proofsEnabled = false;

describe('Oracle.js', () => {
  let deployerAccount: Mina.TestPublicKey,
  deployerKey: PrivateKey,
  senderAccount: Mina.TestPublicKey,
  senderKey: PrivateKey,
  zkAppAddress: PublicKey,
  zkAppPrivateKey: PrivateKey,
  zkApp: Oracle;


  beforeAll(async ()=>{
    if (proofsEnabled) await Oracle.compile();
  })

  beforeEach(async () => {
    const Local = await Mina.LocalBlockchain({
      proofsEnabled,
    })
    Mina.setActiveInstance(Local);
    deployerAccount = Local.testAccounts[0];
    deployerKey = deployerAccount.key;
    senderAccount = Local.testAccounts[1];
    senderKey = senderAccount.key;
    zkAppPrivateKey = PrivateKey.random();
    zkAppAddress = zkAppPrivateKey.toPublicKey();
    zkApp = new Oracle(zkAppAddress);
  });

  async function localDeploy() {
    const txn = await Mina.transaction(deployerAccount, async() => {
      AccountUpdate.fundNewAccount(deployerAccount);
      await zkApp.deploy();
    });
    await txn.prove();
    await txn.sign([deployerKey, zkAppPrivateKey]).send();
  }
  it('Generate and deploy the smartcontract', async () => {
    await localDeploy();

    const oraclePublicKey = zkApp.oraclePublicKey.get();
    expect(oraclePublicKey).toEqual(PublicKey.fromBase58(ORACLE_PUBLIC_KEY));
  });
  describe('Oracle()', () => {
    it('Emit an id evnet containing user id if credit score is greater than 700 with valid signature', async ()=>{
      await localDeploy();
      const id = Field(1);
      const creditScore = Field(800);
      const signature = Signature.fromBase58("7mXXnqMx6YodEkySD3yQ5WK7CCqRL1MBRTASNhrm48oR4EPmenD2NjJqWpFNZnityFTZX5mWuHS1WhRnbdxSTPzytuCgMGuL");
      const txn = await Mina.transaction(senderAccount, async ()=>{
        await zkApp.verify(id, creditScore, signature);
      })
      await txn.prove();
      await txn.sign([senderKey]).send();

      const events = await zkApp.fetchEvents();
      const verifiedEventsValue = events[0].event.data.toFields(null)[0];
      expect(verifiedEventsValue).toEqual(id);
    });
    it('Throws an error if credit score is lesser than 700 or invalid signature', async ()=>{
      await localDeploy();
      const id = Field(1);
      const creditScore = Field(600);
      const signature = Signature.fromBase58("7mXXnqMx6YodEkySD3yQ5WK7CCqRL1MBRTASNhrm48oR4EPmenD2NjJqWpFNZnityFTZX5mWuHS1WhRnbdxSTPzytuCgMGuL");
      expect(async()=>{
        const txn = await Mina.transaction(senderAccount, async ()=>{
          await zkApp.verify(id, creditScore, signature);
        })
      }).rejects;
    });

    it('Throws error if score is above 700 but signature is invalid', async ()=>{
      const id = Field(1);
      const creditScore = Field(787);
      const signature = Signature.fromBase58(
        '7mXPv97hRN7AiUxBjuHgeWjzoSgL3z61a5QZacVgd1PEGain6FmyxQ8pbAYd5oycwLcAbqJLdezY7PRAUVtokFaQP8AJDEGX'
      );
      expect(async()=>{
        const txn = await Mina.transaction(senderAccount, async ()=>{
          await zkApp.verify(id, creditScore, signature);
        })
      }).rejects;
    })
  });
  describe('Oracle.verify()', () => {
    it('Emits an event containing the id if the credit score is greater than 700 with a valid signature', async () => {
      await localDeploy()
      const response = await fetch(
        'https://07-oracles.vercel.app/api/credit-score?user=1'
      )
      const data = await response.json()
      let { id, creditScore , signature } = data.data

      id = Field(id)
      creditScore = Field(creditScore)
      signature = Signature.fromBase58(signature)

      const txn =  await Mina.transaction(senderAccount, async()=>{
        await zkApp.verify(id, creditScore, signature)
      })
      await txn.prove()
      await txn.sign([senderKey]).send()

      const events = await zkApp.fetchEvents()
      const verifiedEventValue = events[0].event.data.toFields(null)[0]
      expect(verifiedEventValue).toEqual(id)
    });
    it('Throws an error if the credit score is less than 700 or the signature is invalid', async () => {
      await localDeploy()
      const response = await fetch(
        'https://07-oracles.vercel.app/api/credit-score?user=2'
      )
      const data = await response.json()
      let { id, creditScore , signature } = data.data

      id = Field(id)
      creditScore = Field(creditScore)
      signature = Signature.fromBase58(signature)

      expect(async()=>{
        const txn =  await Mina.transaction(senderAccount, async()=>{
          await zkApp.verify(id, creditScore, signature)
        })
        await txn.prove()
        await txn.sign([senderKey]).send()
      }).rejects
    });
  });
});
