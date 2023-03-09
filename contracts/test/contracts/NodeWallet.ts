import { Address, beginCell, Cell, Contract, contractAddress, ContractProvider, ContractState, Sender, storeStateInit, toNano } from "ton-core";
import { sign } from "ton-crypto";

import { hex as codeHex } from '../../build/node-wallet.compiled.json';

enum OPS {
    Withdraw = 0x3f6e74,
    CreateCall = 0xf3672d9,
    EndCall = 0x2c2c9c5e,
}

export type NodeWalletData = {
    contractBalance: bigint
    nodeHost: string
    owner: Address
    master: Address
}

export class NodeWallet implements Contract {
    static readonly code = Cell.fromBoc(Buffer.from(codeHex, "hex"))[0];

    constructor(readonly address: Address) {}

    async sendWithdraw(provider: ContractProvider, via: Sender, params: {
        value: bigint
        amount: bigint
    }) {
        await provider.internal(via, {
            value: params.value,
            body: beginCell()
                .storeUint(OPS.Withdraw, 32)
                .storeUint(0, 64) // query_id
                .storeCoins(params.amount)
                .endCell()
        })
    }

    async sendCreateCall(provider: ContractProvider, via: Sender, params: {
        value: bigint,
        userAddress: Address,
        userSecretKey: Buffer,
        callId: number,
    }) {
        const validUntil = Math.floor(Date.now() / 1000) + 60;
        const signingMessage = beginCell()
            .storeUint(params.callId, 64)
            .storeUint(validUntil, 32)
        const signature = sign(signingMessage.endCell().hash(), params.userSecretKey);
        await provider.internal(via, {
            value: params.value,
            body: beginCell()
                .storeUint(OPS.CreateCall, 32)
                .storeUint(0, 64) // query_id
                .storeAddress(params.userAddress)
                .storeRef(beginCell()
                            .storeBuffer(signature)
                            .storeBuilder(signingMessage)
                            .endCell()
                )
                .endCell()
        })
    }

    async sendEndCall(provider: ContractProvider, via: Sender, params: {
        value: bigint
        userAddress: Address,
        userSecretKey: Buffer,
        callId: number,
        spentMinutes: number,
    }) {
        const validUntil = Math.floor(Date.now() / 1000) + 60;
        const signingMessage = beginCell()
            .storeUint(params.callId, 64)
            .storeUint(validUntil, 32)
            .storeUint(params.spentMinutes, 32)
        const signature = sign(signingMessage.endCell().hash(), params.userSecretKey);
        await provider.internal(via, {
            value: params.value,
            body: beginCell()
                .storeUint(OPS.EndCall, 32)
                .storeUint(0, 64) // query_id
                .storeAddress(params.userAddress)
                .storeRef(beginCell()
                            .storeBuffer(signature)
                            .storeBuilder(signingMessage)
                            .endCell()
                )
                .endCell()
        })
    }

    async getData(provider: ContractProvider): Promise<NodeWalletData> {
        const { balance } = await provider.getState();
        const { stack } = await provider.get('get_wallet_data', [])

        return {
            contractBalance: balance,
            nodeHost: stack.readString(),
            owner: stack.readAddress(),
            master: stack.readAddress()
        }
    }
}