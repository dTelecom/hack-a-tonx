import { beginCell } from "ton";

export class MasterContract {

    constructor(address) {
        this.address = address;
    }

    async getUserWalletAddress(provider, userAddress) {
        const { stack } = await provider.get('get_user_wallet_address', [{ type: 'slice', cell: beginCell().storeAddress(userAddress).endCell() }])
        return stack.readAddress()
    }

    async getNodeWalletAddress(provider, nodeAddress) {
        const { stack } = await provider.get('get_node_wallet_address', [{ type: 'slice', cell: beginCell().storeAddress(nodeAddress).endCell() }])
        return stack.readAddress()
    }

}