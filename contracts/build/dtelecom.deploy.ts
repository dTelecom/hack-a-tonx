import {Address, beginCell, Cell, Dictionary, OpenedContract, TonClient, WalletContractV4} from "ton"

import { hex as userWalletCodeHex } from './user-wallet.compiled.json'
import { hex as nodeWalletCodeHex } from './node-wallet.compiled.json'
import { Dtelecom } from "../test/contracts/Dtelecom"

const USER_WALLET_CODE = Cell.fromBoc(Buffer.from(userWalletCodeHex, 'hex'))[0]
const NODE_WALLET_CODE = Cell.fromBoc(Buffer.from(nodeWalletCodeHex, 'hex'))[0]

export function initData(deployingWalletAddress: Address) : Cell {
    return beginCell()
        .storeDict(Dictionary.empty())
        .storeAddress(deployingWalletAddress)
        .storeRef(USER_WALLET_CODE)
        .storeRef(NODE_WALLET_CODE)
        .endCell();
}

export function initMessage() : Cell | null {
    return null;
}

export async function postDeployTest(client: TonClient, walletContract: OpenedContract<WalletContractV4>, secretKey: Buffer, newContractAddress: Address) {
    // const dtelecom = client.open(new Dtelecom(newContractAddress));

    // console.log(await dtelecom.getNodeHosts());
    return null;
}