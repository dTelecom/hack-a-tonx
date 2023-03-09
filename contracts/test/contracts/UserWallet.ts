import { Address, beginCell, Cell, Contract, contractAddress, ContractProvider, ContractState, Sender, storeStateInit, toNano } from "ton-core";

import { hex as codeHex } from '../../build/user-wallet.compiled.json';

export type UserWalletData = {
    contractBalance: bigint
    publicKey: bigint
    owner: Address
    master: Address
}

export class UserWallet implements Contract {
    static readonly code = Cell.fromBoc(Buffer.from(codeHex, "hex"))[0];

    constructor(readonly address: Address) {}

    async getData(provider: ContractProvider): Promise<UserWalletData> {
        const { balance } = await provider.getState();
        const { stack } = await provider.get('get_wallet_data', [])
        return {
            contractBalance: balance,
            publicKey: stack.readBigNumber(),
            owner: stack.readAddress(),
            master: stack.readAddress(),
        }
    }

    async getCallIds(provider: ContractProvider): Promise<number[]> {
        const { stack } = await provider.get('get_call_ids_list', [])

        const callIds: number[] = [];
        let tuple = stack.readTupleOpt()
        while(tuple !== null) {
            callIds.push(tuple.readNumber())
            tuple = tuple.readTupleOpt()
        }

        return callIds;
    }
}