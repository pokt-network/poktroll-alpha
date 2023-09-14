// Generated by Ignite ignite.com/cli

import { StdFee } from "@cosmjs/launchpad";
import { SigningStargateClient, DeliverTxResponse } from "@cosmjs/stargate";
import { EncodeObject, GeneratedType, OfflineSigner, Registry } from "@cosmjs/proto-signing";
import { msgTypes } from './registry';
import { IgniteClient } from "../client"
import { MissingWalletError } from "../helpers"
import { Api } from "./rest";
import { MsgUnstakeServicer } from "./types/poktroll/servicer/tx";
import { MsgStakeServicer } from "./types/poktroll/servicer/tx";

import { Params as typeParams} from "./types"
import { Servicers as typeServicers} from "./types"

export { MsgUnstakeServicer, MsgStakeServicer };

type sendMsgUnstakeServicerParams = {
  value: MsgUnstakeServicer,
  fee?: StdFee,
  memo?: string
};

type sendMsgStakeServicerParams = {
  value: MsgStakeServicer,
  fee?: StdFee,
  memo?: string
};


type msgUnstakeServicerParams = {
  value: MsgUnstakeServicer,
};

type msgStakeServicerParams = {
  value: MsgStakeServicer,
};


export const registry = new Registry(msgTypes);

type Field = {
	name: string;
	type: unknown;
}
function getStructure(template) {
	const structure: {fields: Field[]} = { fields: [] }
	for (let [key, value] of Object.entries(template)) {
		let field = { name: key, type: typeof value }
		structure.fields.push(field)
	}
	return structure
}
const defaultFee = {
  amount: [],
  gas: "200000",
};

interface TxClientOptions {
  addr: string
	prefix: string
	signer?: OfflineSigner
}

export const txClient = ({ signer, prefix, addr }: TxClientOptions = { addr: "http://localhost:26657", prefix: "cosmos" }) => {

  return {
		
		async sendMsgUnstakeServicer({ value, fee, memo }: sendMsgUnstakeServicerParams): Promise<DeliverTxResponse> {
			if (!signer) {
					throw new Error('TxClient:sendMsgUnstakeServicer: Unable to sign Tx. Signer is not present.')
			}
			try {			
				const { address } = (await signer.getAccounts())[0]; 
				const signingClient = await SigningStargateClient.connectWithSigner(addr,signer,{registry, prefix});
				let msg = this.msgUnstakeServicer({ value: MsgUnstakeServicer.fromPartial(value) })
				return await signingClient.signAndBroadcast(address, [msg], fee ? fee : defaultFee, memo)
			} catch (e: any) {
				throw new Error('TxClient:sendMsgUnstakeServicer: Could not broadcast Tx: '+ e.message)
			}
		},
		
		async sendMsgStakeServicer({ value, fee, memo }: sendMsgStakeServicerParams): Promise<DeliverTxResponse> {
			if (!signer) {
					throw new Error('TxClient:sendMsgStakeServicer: Unable to sign Tx. Signer is not present.')
			}
			try {			
				const { address } = (await signer.getAccounts())[0]; 
				const signingClient = await SigningStargateClient.connectWithSigner(addr,signer,{registry, prefix});
				let msg = this.msgStakeServicer({ value: MsgStakeServicer.fromPartial(value) })
				return await signingClient.signAndBroadcast(address, [msg], fee ? fee : defaultFee, memo)
			} catch (e: any) {
				throw new Error('TxClient:sendMsgStakeServicer: Could not broadcast Tx: '+ e.message)
			}
		},
		
		
		msgUnstakeServicer({ value }: msgUnstakeServicerParams): EncodeObject {
			try {
				return { typeUrl: "/poktroll.servicer.MsgUnstakeServicer", value: MsgUnstakeServicer.fromPartial( value ) }  
			} catch (e: any) {
				throw new Error('TxClient:MsgUnstakeServicer: Could not create message: ' + e.message)
			}
		},
		
		msgStakeServicer({ value }: msgStakeServicerParams): EncodeObject {
			try {
				return { typeUrl: "/poktroll.servicer.MsgStakeServicer", value: MsgStakeServicer.fromPartial( value ) }  
			} catch (e: any) {
				throw new Error('TxClient:MsgStakeServicer: Could not create message: ' + e.message)
			}
		},
		
	}
};

interface QueryClientOptions {
  addr: string
}

export const queryClient = ({ addr: addr }: QueryClientOptions = { addr: "http://localhost:1317" }) => {
  return new Api({ baseURL: addr });
};

class SDKModule {
	public query: ReturnType<typeof queryClient>;
	public tx: ReturnType<typeof txClient>;
	public structure: Record<string,unknown>;
	public registry: Array<[string, GeneratedType]> = [];

	constructor(client: IgniteClient) {		
	
		this.query = queryClient({ addr: client.env.apiURL });		
		this.updateTX(client);
		this.structure =  {
						Params: getStructure(typeParams.fromPartial({})),
						Servicers: getStructure(typeServicers.fromPartial({})),
						
		};
		client.on('signer-changed',(signer) => {			
		 this.updateTX(client);
		})
	}
	updateTX(client: IgniteClient) {
    const methods = txClient({
        signer: client.signer,
        addr: client.env.rpcURL,
        prefix: client.env.prefix ?? "cosmos",
    })
	
    this.tx = methods;
    for (let m in methods) {
        this.tx[m] = methods[m].bind(this.tx);
    }
	}
};

const Module = (test: IgniteClient) => {
	return {
		module: {
			PoktrollServicer: new SDKModule(test)
		},
		registry: msgTypes
  }
}
export default Module;