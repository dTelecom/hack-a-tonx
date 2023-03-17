import { SendTransactionRequest, TonConnect, UserRejectsError, WalletInfo, WalletInfoInjected } from '@tonconnect/sdk';
import {appStore} from './stores/appStore';

const dappMetadata = {
	manifestUrl:
		'https://raw.githubusercontent.com/dTelecom/hack-a-tonx/main/contracts/web/tonconnect-manifest.json',
};

export const connector = new TonConnect(dappMetadata);

connector.onStatusChange(walletInfo => {
    appStore.setCurrentUser(walletInfo);
    if (walletInfo) {
        // const tonweb = new TonWeb(new TonWeb.HttpProvider('https://testnet.toncenter.com/api/v2/jsonRPC', {apiKey: 'c1b5cb3616604568237afcd9148022fe13c4644ce9eec37045533d9bf888118f'}));
        // appStore.setTonweb(tonweb);

        const isTestnet = true;
        const client = new TonClient({ endpoint: `https://${isTestnet ? "testnet." : ""}toncenter.com/api/v2/jsonRPC` });
        appStore.setTonClient(client);
    } else {
        appStore.setTonClient(nul);
    }
}, console.error);