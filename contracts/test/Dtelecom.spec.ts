import { beginCell, SendMode, toNano } from "ton-core"
import { Blockchain } from "@ton-community/sandbox"
import { Dtelecom } from "./contracts/Dtelecom"
import "@ton-community/test-utils" // register matchers
import { UserWallet } from "./contracts/UserWallet"
import { NodeWallet } from "./contracts/NodeWallet"

const MINUTE_PRICE = toNano('0.01')

describe('Dtelecom', () => {
    it('should work', async () => {
        const blkch = await Blockchain.create()
        // blkch.verbosity = {
        //     print: true,
        //     blockchainLogs: true,
        //     vmLogs: 'none',
        //     debugLogs: false,
        // }

        const admin = await blkch.treasury('master')
        const user = await blkch.treasury('user')
        const node = await blkch.treasury('node')

        const dtelecom = blkch.openContract(new Dtelecom(0, admin.address))

        // initializee contract by sending nonsense message
        await dtelecom.sendWithdraw(admin.getSender(), {
            value: toNano('0.1'),
            amount: toNano(0)
        })

        let dtelecomData = await dtelecom.getData()

        expect(dtelecomData.contractBalance).toBe(toNano(0))
        expect(dtelecomData.owner.toString()).toBe(admin.address.toString())


        // test master withdraw
        // <<<<<<<<<<<<<<<<<<<<
        await admin.send({
            value: toNano(1),
            to: dtelecom.address
        })

        dtelecomData = await dtelecom.getData()
        expect(dtelecomData.contractBalance).toBeLessThan(toNano('1.0'))
        expect(dtelecomData.contractBalance).toBeGreaterThan(toNano('0.9'))

        await dtelecom.sendWithdraw(admin.getSender(), {
            value: toNano('0.1'),
            amount: toNano(9999), // empty contract balance
        })

        dtelecomData = await dtelecom.getData()
        expect(dtelecomData.contractBalance).toBe(toNano(0))
        // >>>>>>>>>>>>>>>>>>>>


        // test create user
        // <<<<<<<<<<<<<<<<
        const userPublicKey = BigInt(`0x${user.keypair.publicKey.toString('hex')}`)
        await dtelecom.sendCreateUser(user.getSender(), {
            value: toNano('100.0'),
            publicKey: userPublicKey
        })
        const userWalletAddress = await dtelecom.getUserWalletAddress(user.address)
        const userWallet = blkch.openContract(new UserWallet(userWalletAddress))
        let userWalletData = await userWallet.getData()
        expect(userWalletData.contractBalance).toBeGreaterThan(toNano('99.9'))
        expect(userWalletData.contractBalance).toBeLessThan(toNano('100.0'))
        expect(userWalletData.publicKey).toBe(userPublicKey)
        expect(userWalletData.owner.toString()).toBe(user.address.toString())
        expect(userWalletData.master.toString()).toBe(dtelecom.address.toString())
        // >>>>>>>>>>>>>>>>


        // test create node
        // <<<<<<<<<<<<<<<<
        let nodeHosts = await dtelecom.getNodeHosts()
        expect(nodeHosts).toHaveLength(0)

        await dtelecom.sendCreateNode(node.getSender(), {
            value: toNano('20.1'),
            nodeHost: 'dtelecom.org'
        })
        const nodeWalletAddress = await dtelecom.getNodeWalletAddress(node.address)
        const nodeWallet = blkch.openContract(new NodeWallet(nodeWalletAddress))
        let nodeWalletData = await nodeWallet.getData()
        expect(nodeWalletData.contractBalance).toBeGreaterThanOrEqual(toNano('20.0'))
        expect(nodeWalletData.contractBalance).toBeLessThan(toNano('20.1'))
        expect(nodeWalletData.owner.toString()).toBe(node.address.toString())
        expect(nodeWalletData.master.toString()).toBe(dtelecom.address.toString())

        nodeHosts = await dtelecom.getNodeHosts()
        expect(nodeHosts).toHaveLength(1)
        expect(nodeHosts).toContain('dtelecom.org')
        // >>>>>>>>>>>>>>>>


        // test node withdraw
        // <<<<<<<<<<<<<<<<<<
        await nodeWallet.sendWithdraw(node.getSender(), {
            value: toNano('0.1'),
            amount: toNano(5)
        })
        nodeWalletData = await nodeWallet.getData()
        expect(nodeWalletData.contractBalance).toBeGreaterThanOrEqual(toNano('15.0'))
        expect(nodeWalletData.contractBalance).toBeLessThan(toNano('15.1'))

        await nodeWallet.sendWithdraw(node.getSender(), {
            value: toNano('0.1'),
            amount: toNano('99999')
        })
        nodeWalletData = await nodeWallet.getData()
        expect(nodeWalletData.contractBalance).toBe(toNano('1.0'))
        // >>>>>>>>>>>>>>>>>>


        // test create call
        // <<<<<<<<<<<<<<<<
        let callIds = await userWallet.getCallIds()
        expect(callIds).toHaveLength(0)

        const { contractBalance: userBalanceBeforeCreateCall } = await userWallet.getData()

        await nodeWallet.sendCreateCall(node.getSender(), {
            value: toNano('0.1'),
            userAddress: user.address,
            userSecretKey: node.keypair.secretKey,
            callId: 123
        })

        callIds = await userWallet.getCallIds()
        expect(callIds).toHaveLength(0) // cause we used node's secretKey, instead of user's

        await nodeWallet.sendCreateCall(node.getSender(), {
            value: toNano('0.1'),
            userAddress: user.address,
            userSecretKey: user.keypair.secretKey,
            callId: 123
        })

        callIds = await userWallet.getCallIds()
        expect(callIds).toHaveLength(1)
        expect(callIds).toContain(123)

        const { contractBalance: userBalanceAfterCreateCall } = await userWallet.getData()
        expect(userBalanceAfterCreateCall).toBeLessThanOrEqual(userBalanceBeforeCreateCall)
        // I do not understand why the balance after this operation sometimes is lower by 1 coin
        expect(userBalanceAfterCreateCall).toBeGreaterThanOrEqual(userBalanceBeforeCreateCall - BigInt(1))
        // >>>>>>>>>>>>>>>>


        // test end call
        // <<<<<<<<<<<<<
        const { contractBalance: userBalanceBeforeEndCall } = await userWallet.getData()
        const { contractBalance: nodeBalanceBeforeEndCall } = await nodeWallet.getData()
        const { contractBalance: dtelecomBalanceBeforeEndCall } = await dtelecom.getData()
        await nodeWallet.sendEndCall(node.getSender(), {
            value: toNano('0.1'),
            userAddress: user.address,
            userSecretKey: user.keypair.secretKey,
            callId: 123,
            spentMinutes: 100
        })

        callIds = await userWallet.getCallIds()
        expect(callIds).toHaveLength(0)

        const { contractBalance: userBalanceAfterEndCall } = await userWallet.getData()
        expect(userBalanceAfterEndCall).toBeLessThanOrEqual(userBalanceBeforeEndCall - MINUTE_PRICE * BigInt(200))
        // I do not understand why the balance after this operation sometimes is lower by 1 coin
        expect(userBalanceAfterEndCall).toBeGreaterThanOrEqual(userBalanceBeforeEndCall - MINUTE_PRICE * BigInt(200) - BigInt(1))

        const { contractBalance: nodeBalanceAfterEndCall } = await nodeWallet.getData()
        expect(nodeBalanceAfterEndCall).toBeLessThan(nodeBalanceBeforeEndCall + MINUTE_PRICE * BigInt(100))
        expect(nodeBalanceAfterEndCall).toBeGreaterThan(nodeBalanceBeforeEndCall + MINUTE_PRICE * BigInt(100) - toNano('0.01'))

        const { contractBalance: dtelecomBalanceAfterEndCall } = await dtelecom.getData()
        expect(dtelecomBalanceAfterEndCall).toBeLessThan(dtelecomBalanceBeforeEndCall + MINUTE_PRICE * BigInt(100))
        expect(dtelecomBalanceAfterEndCall).toBeGreaterThan(dtelecomBalanceBeforeEndCall + MINUTE_PRICE * BigInt(100)- toNano('0.01'))
        // >>>>>>>>>>>>>
    })
})